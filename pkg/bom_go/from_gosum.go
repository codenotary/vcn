package bom_go

import (
	"os"
	"bufio"
	"path/filepath"
	"strings"

	"github.com/vchain-us/vcn/pkg/bom_component"
)

func goModComponents(path string) ([]bom_component.Component, error) {
	// cache content of go.sum to avoid dupes
	sumFile, err := os.Open(filepath.Join(path, "go.sum"))
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
	

	res := make([]bom_component.Component, 0, len(hashes))
	for k, v := range hashes {
		hash, hashType, err := ModHash(v)
		if err != nil {
			return nil, err
		}
		res = append(res, bom_component.Component{
			Name: k.name,
			Version: k.version,
			Hash: hash,
			HashType: hashType})
	}

	return res, nil
}
