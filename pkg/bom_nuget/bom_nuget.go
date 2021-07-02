package bom_nuget

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/vchain-us/vcn/pkg/bom_component"
)

type PackageFile struct {
	Targets map[string]Deps `json:"dependencies"`
}

type Deps map[string]Dep

type Dep struct {
	Type    string
	Version string `json:"resolved"`
	Hash    string `json:"contentHash"`
}

const nugetFileName = "packages.lock.json"

// GoPackage implements Package interface
type NugetPackage struct {
	dirName string
}

type mapKey struct {
	name    string
	version string
}

// New returns new NugetPackage object if Nuget packages file is found in directory
func New(dir string) *NugetPackage {
	_, err := os.Stat(filepath.Join(dir, nugetFileName))
	if err != nil {
		return nil
	}

	return &NugetPackage{dirName: dir}
}

func (p *NugetPackage) Type() string {
	return ".Net"
}

func (p *NugetPackage) Path() string {
	return p.dirName
}

func (p *NugetPackage) Close() {
	// do nothing, function needed to comply with Package interface
}

func (p *NugetPackage) Components() ([]bom_component.Component, error) {
	file, err := os.Open(filepath.Join(p.dirName, nugetFileName))
	if err != nil {
		return nil, fmt.Errorf("cannot open Nuget dependencies file: %w", err)
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read Nuget dependencies file: %w", err)
	}

	var lockFile PackageFile
	err = json.Unmarshal(buf, &lockFile)
	if err != nil {
		return nil, fmt.Errorf("cannot parse Nuget dependencies file: %w", err)
	}

	// file may contain duplicate modules for different .Net version, so process all .Net versions and select only unique packages
	res := make([]bom_component.Component, 0)
	seen := map[mapKey]struct{}{}
	for _, target := range lockFile.Targets {
		for name, details := range target {
			if details.Type == "Project" {
				continue // skip entry for the parent project
			}
			_, ok := seen[mapKey{name: name, version: details.Version}]
			if ok {
				continue // module already processed
			}
			seen[mapKey{name: name, version: details.Version}] = struct{}{} // mark as seen

			hash, err := base64.StdEncoding.DecodeString(details.Hash)
			if err != nil {
				return nil, fmt.Errorf("cannot decode hash for %s: %w", name, err)
			}

			res = append(res, bom_component.Component{
				Name:     name,
				Version:  details.Version,
				Hash:     hex.EncodeToString(hash),
				HashType: bom_component.HashSHA512})
		}
	}

	return res, nil
}
