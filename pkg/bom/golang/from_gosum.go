package golang

import (
	"os"
	"bufio"
	"path/filepath"
	"strings"

	"github.com/vchain-us/vcn/pkg/bom/artifact"
)

// goArtifactFromSum implements Artifact interface
type goArtifactFromSum struct {
	goArtifact
}

// Dependencies returns list of Go dependencies used during the build
func (a *goArtifactFromSum) Dependencies() ([]artifact.Dependency, error) {
	if a.Deps != nil {
		return a.Deps, nil
	}
	// cache content of go.sum to avoid dupes
	sumFile, err := os.Open(filepath.Join(a.path, "go.sum"))
	if err != nil {
		return nil, err
	}
	defer sumFile.Close()

	scanner := bufio.NewScanner(sumFile)
	hashes := make(map[mapKey]string)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 3 {
			continue	// skip malformed lines
		}
		if strings.HasSuffix(fields[1], "/go.mod") {
			continue	// skip go.mod hashes
		}
		key := mapKey{name: fields[0], version: fields[1]}
		if _, ok := hashes[key]; ok {
			continue
		}
		hashes[key] = fields[2]
	}
	

	res := make([]artifact.Dependency, 0, len(hashes))
	for k, v := range hashes {
		hash, hashType, err := ModHash(v)
		if err != nil {
			return nil, err
		}
		res = append(res, artifact.Dependency{
			Name: k.name,
			Version: k.version,
			Hash: hash,
			HashType: hashType})
	}

	a.Deps = res
	return res, nil
}
