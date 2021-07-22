/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package golang

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/mod/sumdb"

	"github.com/vchain-us/vcn/pkg/bom/artifact"
)

// GoArtifactFromExe implements Artifact interface
type goArtifactFromGoList struct {
	goArtifact
}

type mapKey struct {
	name    string
	version string
}

type clientOps struct{}

var goListArgs = []string{"list", "--deps", "-f", "{{if not .Standard}}{{.Module.Path}} {{.Module.Version}}{{end}}"}
var sumDb = "sum.golang.org+033de0ae+Ac4zctda0e5eza+HJyk9SxEdh+s3Ux18htTTAD8OuAn8" // default sumdb server and it's public key

// Dependencies returns list of Go dependencies used during the build
// run 'go list' to get the list of used modules, and then get hashes from sumdb
func (a *goArtifactFromGoList) Dependencies() ([]artifact.Dependency, error) {
	if a.Deps != nil {
		return a.Deps, nil
	}
	absPath, err := filepath.Abs(a.path)
	if err != nil {
		return nil, err
	}
	// go won't work if cur directory is outside Go module root
	cmd := exec.Command("go", append(goListArgs, absPath)...)
	cmd.Dir = absPath
	buf, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// override default sumdb path with the value from env, if set
	tmp, ok := os.LookupEnv("GOSUMDB")
	if ok {
		sumDb = tmp
	}

	res := make([]artifact.Dependency, 0)
	client := sumdb.NewClient(new(clientOps))

	// start workers and response processor
	tasks := make(chan mapKey)
	results := make(chan string)
	var wg sync.WaitGroup
	for i := 0; i < artifact.MaxGoroutines; i++ {
		go func() {
			for tasks := range tasks {
				lines, err := client.Lookup(tasks.name, tasks.version)
				if err != nil {
					log.Printf("Cannot lookup package %s/%s: %v", tasks.name, tasks.version, err)
					results <- ""
				} else {
					results <- lines[0]
				}
			}
		}()
	}
	go func() {
		for details := range results {
			fields := strings.Fields(details)
			if len(fields) != 3 {
				continue
			}
			hash, hashType, err := ModHash(fields[2])
			if err == nil {
				res = append(res, artifact.Dependency{
					Name:     fields[0],
					Version:  fields[1],
					Hash:     hash,
					HashType: hashType})
			}
			wg.Done()
		}
	}()

	seen := make(map[mapKey]struct{})
	scanner := bufio.NewScanner(bytes.NewReader(buf))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			continue // skip malformed lines
		}

		_, ok := seen[mapKey{fields[0], fields[1]}]
		if ok {
			continue // already processed
		}
		seen[mapKey{fields[0], fields[1]}] = struct{}{}

		tasks <- mapKey{fields[0], fields[1]}
		wg.Add(1)
	}

	wg.Wait()
	close(tasks)
	close(results)

	a.Deps = res
	return res, nil
}

func (*clientOps) ReadConfig(file string) ([]byte, error) {
	if file == "key" {
		return []byte(sumDb), nil
	}
	if strings.HasSuffix(file, "/latest") {
		// Looking for cached latest tree head.
		// Empty result means empty tree.
		return []byte{}, nil
	}
	return nil, fmt.Errorf("unknown config %s", file)
}

func (*clientOps) WriteConfig(file string, old, new []byte) error {
	// Ignore writes.
	return nil
}

func (*clientOps) ReadCache(file string) ([]byte, error) {
	return nil, fmt.Errorf("no cache")
}

func (*clientOps) WriteCache(file string, data []byte) {
	// Ignore writes.
}

func (*clientOps) Log(msg string) {
	log.Print(msg)
}

func (*clientOps) SecurityError(msg string) {
	log.Fatal(msg)
}

func init() {
	http.DefaultClient.Timeout = 1 * time.Minute
}

func (*clientOps) ReadRemote(path string) ([]byte, error) {
	name := sumDb
	if i := strings.Index(name, "+"); i >= 0 {
		name = name[:i]
	}
	target := "https://" + name + path
	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GET %v: %v", target, resp.Status)
	}
	data, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	return data, nil
}
