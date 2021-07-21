/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom_node

import (
	"encoding/json"
	"fmt"
	"github.com/opencontainers/go-digest"
	"github.com/vchain-us/vcn/pkg/bom_component"
	"github.com/vchain-us/vcn/pkg/bundle"
	"github.com/vchain-us/vcn/pkg/extractor/dir"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

const packageLockJsonFN = "package-lock.json"
const packageJsonFN = "package.json"

type NodePackage struct {
	path string
}

func New(path string) *NodePackage {
	plock, err := getPackageLock(path)
	if err != nil {
		return nil
	}
	return &NodePackage{path: plock}
}

func (p *NodePackage) Components() ([]bom_component.Component, error) {
	comps := make([]bom_component.Component, 0)
	fname, err := exec.LookPath("npm")
	if err != nil {
		return nil, fmt.Errorf("please install npm tool follwing this link: https://docs.npmjs.com/getting-started. Error reported: %w", err)
	}

	command := exec.Command(fname, "ls", "-a", "-l", "-p", "--json")
	o, err := command.Output()
	if err != nil {
		return nil, err
	}

	var plJson map[string]interface{}

	err = json.Unmarshal(o, &plJson)
	if err != nil {
		return nil, err
	}

	for _, dep := range plJson["dependencies"].(map[string]interface{}) {
		comp := bom_component.Component{}
		d, ok := dep.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("impossible to calculate digest of %v dependecy", dep)
		}
		name, ok := d["name"].(string)
		if !ok {
			return nil, fmt.Errorf("impossible to calculate digest of %v dependecy", d)
		}
		version, ok := d["version"].(string)
		if !ok {
			return nil, fmt.Errorf("impossible to calculate digest of %s dependecy", name)
		}
		nodeComPath, ok := d["path"].(string)
		if !ok {
			return nil, fmt.Errorf("impossible to calculate digest of %s-%s dependecy", name, version)
		}
		hash, err := GetNodeComDigest(nodeComPath)
		if err != nil {
			return nil, err
		}
		comp.Hash = hash.Encoded()
		comp.Version = version
		comp.Name = name
		comp.HashType = bom_component.HashSHA256
		comps = append(comps, comp)
	}

	return comps, nil
}

func getPackageLock(p string) (string, error) {
	fp := path.Join(p, packageLockJsonFN)
	_, err := os.Stat(fp)
	if os.IsNotExist(err) {
		return "", err
	}
	return fp, nil
}

func (p *NodePackage) Type() string {
	return "node"
}

func (p *NodePackage) Path() string {
	return p.path
}

func (p *NodePackage) Close() {}

func GetPackageJsonPath(p string) (string, error) {
	fp := path.Join(p, packageJsonFN)
	_, err := os.Stat(fp)
	if os.IsNotExist(err) {
		return "", err
	}
	return fp, nil
}

func GetNodeComDigest(componentFolder string) (digest.Digest, error) {
	path, err := filepath.Abs(componentFolder)
	if err != nil {
		return "", err
	}

	d, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer d.Close()

	stat, err := d.Stat()
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		return "", fmt.Errorf("read node component %s: is not a directory", path)
	}

	files, err := dir.Walk(path)
	if err != nil {
		return "", err
	}

	manifest := bundle.NewManifest(files...)
	return manifest.Digest()
}
