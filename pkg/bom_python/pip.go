package bom_python

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
//	"sync"
	"net/http"
	"encoding/json"
	"io"
	"encoding/hex"

	component "github.com/vchain-us/vcn/pkg/bom_component"
)

type pypiUrls struct {
	Files []pypiFile	`json:"urls"`
}

type pypiFile struct {
	Digests pypiDigests `json:"digests"`
}

type pypiDigests struct {
	Md5 string			`json:"md5,omitempty"`
	Sha256 string		`json:"sha256,omitempty"`
}

type module struct {
	needed	bool
	version string
}

const (
	pythonExe = "python"
	maxGoroutines = 8
)

var moduleListArgs = []string{"-m", "pip", "list", "-v"}
var moduleDetailsArgs = []string{"-m", "pip", "show"}

// collect info about all installed modules, find module relations, populate module graph and then recursively
// select only the needed modules, using content of 'requirements.txt' as a starting point
func procPip(dir string) ([]component.Component, error) {
	buf, err := exec.Command(pythonExe, moduleListArgs...).Output()
	if err != nil {
		return nil, err
	}

	// output has two header lines, and then entries of the format "<package> <version> <location> <installer>"
	scanner := bufio.NewScanner(bytes.NewReader(buf))
	scanner.Split(bufio.ScanLines)
	if !scanner.Scan() || !scanner.Scan() { // not an error - need to skip two lines
		return nil, errors.New("got unexpected result to pip module list request")
	}

	// collect module dependencies
	moduleGraph := make(map[string]*module)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		moduleGraph[fields[0]] = &module{version: fields[1]}
	}

	// read root dependencies from requirements.txt and then traverse dependency graph, identifying needed modules
	buf, err = os.ReadFile(filepath.Join(dir, "requirements.txt"))
	if err != nil {
		return nil, err
	}
	scanner = bufio.NewScanner(bytes.NewReader(buf))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		endPos := strings.IndexAny(line, "=>")
		if endPos > 0 {
			err = traverseDependencies(moduleGraph, line[:endPos])
		} else {
			err = traverseDependencies(moduleGraph, line)
		}
		if err != nil {
			return nil, err
		}
	}

	// select only needed modules
	res := make([]component.Component, 0)
	for k, v := range moduleGraph {
		if v.needed {
			hash, err := queryHash(k, v.version)
			if err != nil {
				return nil, err
			}
			res = append(res, component.Component{Name: k, Version: v.version, Hash: hash})
		}
	}

	return res, nil
}

func preRequisites(module string) ([]string, error) {
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

func traverseDependencies(graph map[string]*module, line string) error {
	mod, ok := graph[line]
	if !ok || mod.needed {
		return nil // not found (should never happen) or already processed
	}
	mod.needed = true

	preReqs, err := preRequisites(line)
	if err != nil {
		return err
	}

	for _, dep := range preReqs {
		err = traverseDependencies(graph, dep)
		if err != nil {
			return err
		}	
	}

	return nil
}

// query PyPI.org for module hash, combine all available hashes using XOR
func queryHash(name, version string) (string, error) {
	resp, err := http.Get("https://pypi.org/pypi/" + name + "/" + version + "/json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if  resp.StatusCode != http.StatusOK {
		return "", errors.New("cannot query PyPI for package details")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var urls pypiUrls
	err = json.Unmarshal(body, &urls)
	if err != nil {
		return "", err
	}

	var res []byte
	for _, file := range urls.Files {
		var raw []byte
		if file.Digests.Sha256 != "" {
			raw, err = hex.DecodeString(file.Digests.Sha256)
		} else {
			raw, err = hex.DecodeString(file.Digests.Md5)
		}
		if err != nil {
			return "", errors.New("malformed hash value")
		}
		if res == nil {
			res = raw
		} else {
			if len(raw) != len(res) {
				// should never happen - all hashes must be of the same length
				return "", errors.New("malformed hash value")
			}
			// XOR hash
			for i := 0; i < len(res); i++ {
				res[i] ^= raw[i]
			}
		}
	}
	return hex.EncodeToString(res), nil
}
