/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package python

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/vchain-us/vcn/pkg/bom/artifact"
)

// pythonArtifactFromPoetry implements Artifact interface
type pythonArtifactFromPoetry struct {
	pythonArtifact
}

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
func (a *pythonArtifactFromPoetry) Dependencies() ([]artifact.Dependency, error) {
	if a.Deps != nil {
		return a.Deps, nil
	}
	var poetry poetryFile
	_, err := toml.DecodeFile(filepath.Join(a.path, poetryFileName), &poetry)
	if err != nil {
		return nil, err
	}
	res := make([]artifact.Dependency, 0, len(poetry.Packages))
	for _, pkg := range poetry.Packages {
		var hash string
		var hashType artifact.HashType
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
		res = append(res, artifact.Dependency{
			Name:     pkg.Name,
			Version:  pkg.Version,
			Hash:     hash,
			HashType: hashType})
	}

	a.Deps = res
	return res, nil
}
