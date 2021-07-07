/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom_python

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/vchain-us/vcn/pkg/bom_component"
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
func procPoetry(dir string) ([]bom_component.Component, error) {
	var poetry poetryFile
	_, err := toml.DecodeFile(filepath.Join(dir, poetryFileName), &poetry)
	if err != nil {
		return nil, err
	}
	res := make([]bom_component.Component, 0, len(poetry.Packages))
	for _, pkg := range poetry.Packages {
		var hash string
		var hashType int
		meta, ok := poetry.Meta.Files[pkg.Name]
		if ok {
			hashes := make([]string, len(meta))
			for i, f := range meta {
				hashes[i] = f.Hash
			}
			hash, hashType, err = combineHashes(hashes)
			if err != nil {
				return nil, err
			}
		}
		res = append(res, bom_component.Component{
			Name:     pkg.Name,
			Version:  pkg.Version,
			Hash:     hash,
			HashType: hashType})
	}

	return res, nil
}
