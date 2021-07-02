package bom_python

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vchain-us/vcn/pkg/bom_component"
)

type pypiUrls struct {
	Files []pypiFile `json:"urls"`
}

type pypiFile struct {
	Digests pypiDigests `json:"digests"`
}

type pypiDigests struct {
	Md5    string `json:"md5,omitempty"`
	Sha256 string `json:"sha256,omitempty"`
}

type module struct {
	needed  bool
	version string
}

type task struct {
	name    string
	version string
}
type result struct {
	name     string
	hash     string
	hashType int
	deps     []string
	err      error
}

const maxGoroutines = 8

var moduleListArgs = []string{"-m", "pip", "list", "-v"}
var moduleDetailsArgs = []string{"-m", "pip", "show"}

// collect info about all installed modules, find module relations, populate module graph and then recursively
// select only the needed modules, using content of 'requirements.txt' as a starting point
func procPip(dir string) ([]bom_component.Component, error) {
	// first try "python", if it fails, try "python3"
	pythonExe := "python"
	buf, err := exec.Command(pythonExe, moduleListArgs...).Output()
	if err != nil || len(buf) == 0 {
		pythonExe = "python3"
		buf, err = exec.Command(pythonExe, moduleListArgs...).Output()
		if err != nil {
			return nil, fmt.Errorf("cannot get python module list: %w", err)
		}
	}

	// output has two header lines, and then entries of the format "<package> <version> <location> <installer>"
	scanner := bufio.NewScanner(bytes.NewReader(buf))
	scanner.Split(bufio.ScanLines)
	for i := 0; i < 2; {
		if !scanner.Scan() {
			return nil, errors.New("got unexpected result to pip module list request")
		}
		text := scanner.Text()
		if text[0] == '#' {
			// skip possible Python warnings - lines starting with #
			continue
		}
		i++
	}

	// store all known modules
	moduleGraph := make(map[string]*module)
	for scanner.Scan() {
		text := scanner.Text()
		if text[0] == '#' {
			// skip possible Python warnings - lines starting with #
			continue
		}
		fields := strings.Fields(text)
		moduleGraph[fields[0]] = &module{version: fields[1]}
	}

	// root dependencies from requirements.txt
	buf, err = ioutil.ReadFile(filepath.Join(dir, "requirements.txt"))
	if err != nil {
		return nil, err
	}

	// init goroutine throttling - channels, start goroutines.
	// We can be sure that there will be no more in-flight messages in channels than known modules
	tasks := make(chan task, len(moduleGraph))
	results := make(chan result, len(moduleGraph))
	for i := 0; i < maxGoroutines; i++ {
		go worker(tasks, results, pythonExe)
	}

	taskCount := 0

	// initial tasks - content of requirements.txt
	scanner = bufio.NewScanner(bytes.NewReader(buf))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		endPos := strings.IndexAny(line, "=>")
		if endPos > 0 {
			line = line[:endPos]
		}

		mod, ok := moduleGraph[line]
		if !ok {
			// not found (should never happen) or already processed
			log.Printf("Unknown module %s - ignoring", line)
			continue
		}
		mod.needed = true

		tasks <- task{name: line, version: mod.version}
		taskCount++
	}

	// get dependencies, run tasks for dependencies, collect info about all used modules
	res := make([]bom_component.Component, 0)
	for done := 0; taskCount == 0 || done < taskCount; done++ {
		result := <-results
		if result.err != nil {
			close(tasks) // signal workers to stop
			return nil, err
		}
		res = append(res, bom_component.Component{
			Name:     result.name,
			Version:  moduleGraph[result.name].version,
			Hash:     result.hash,
			HashType: result.hashType})
		for _, v := range result.deps {
			if v == "" {
				continue
			}
			mod, ok := moduleGraph[v]
			if !ok {
				// not found (should never happen) or already processed
				log.Printf("Unknown module %s - ignoring", v)
				continue
			}

			if mod.needed {
				continue // already being processed
			}
			mod.needed = true
			tasks <- task{name: v, version: mod.version}
			taskCount++
		}
	}
	close(tasks)   // signal workers to stop
	close(results) // it is safe to close result channel because workers do nothing at this point

	return res, nil
}

func worker(tasks <-chan task, results chan<- result, pythonExe string) {
	for task := range tasks {
		hash, hashType, err := queryHash(task.name, task.version)
		if err != nil {
			results <- result{err: err}
			continue
		}
		deps, err := preRequisites(pythonExe, task.name)
		if err != nil {
			results <- result{err: err}
			continue
		}

		results <- result{name: task.name, hash: hash, deps: deps, hashType: hashType, err: nil}
	}
}

func preRequisites(pythonExe string, module string) ([]string, error) {
	output, err := exec.Command(pythonExe, append(moduleDetailsArgs, module)...).Output()
	if err != nil {
		return nil, err
	}

	// output has two header lines, and then entries of the format "<package> <version> <location> <installer>"
	scanner := bufio.NewScanner(bytes.NewReader(output))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), ": ", 2)
		if fields[0] == "Requires" {
			return strings.Split(fields[1], ", "), nil
		}
	}
	return nil, errors.New("malformed output from pip module query")
}

// query PyPI.org for module hash, combine all available hashes using XOR
func queryHash(name, version string) (string, int, error) {
	resp, err := http.Get("https://pypi.org/pypi/" + name + "/" + version + "/json")
	if err != nil {
		return "", bom_component.HashInvalid, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", bom_component.HashInvalid, errors.New("cannot query PyPI for package details")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", bom_component.HashInvalid, err
	}

	var urls pypiUrls
	err = json.Unmarshal(body, &urls)
	if err != nil {
		return "", bom_component.HashInvalid, err
	}

	// assuming that all files have the same type of hash, with priority for SHA-256
	hashType := bom_component.HashMD5
	if urls.Files[0].Digests.Sha256 != "" {
		hashType = bom_component.HashSHA256
	}
	hashes := make([]string, len(urls.Files))
	for i, file := range urls.Files {
		if hashType == bom_component.HashSHA256 {
			hashes[i] = file.Digests.Sha256
		} else {
			hashes[i] = file.Digests.Md5
		}
	}

	hash, _, err := combineHashes(hashes)
	if err != nil {
		return "", bom_component.HashInvalid, errors.New("malformed hash value")
	}

	return hash, hashType, nil
}
