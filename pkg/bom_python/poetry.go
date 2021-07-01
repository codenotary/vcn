package bom_python

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	component "github.com/vchain-us/vcn/pkg/bom_component"
)

// poetry.lock TOML structure
type poetryFile struct {
	Packages []pkg    `toml:"package"`
	Meta     metadata `toml:"metadata"`
}
type pkg struct {
	Name    string
	Version string
}
type metadata struct {
	Files map[string][]file
}
type file struct {
	File string
	Hash string
}

// poetry.lock file contains list of all dependencies with hashes
func procPoetry(dir string) ([]component.Component, error) {
	var poetry poetryFile
	_, err := toml.DecodeFile(filepath.Join(dir, "poetry.lock"), &poetry)
	if err != nil {
		return nil, err
	}
	res := make([]component.Component, 0, len(poetry.Packages))
	for _, pkg := range poetry.Packages {
		var hash string
		meta, ok := poetry.Meta.Files[pkg.Name]
		if ok {
			hashes := make([]string, len(meta))
			for i, f := range meta {
				hashes[i] = f.Hash
			}
			hash, err = combineHashes(hashes)
			if err != nil {
				return nil, err
			}
		}
		res = append(res, component.Component{Name: pkg.Name, Version: pkg.Version, Hash: hash})
	}

	return res, nil
}
