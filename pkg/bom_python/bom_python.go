/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom_python

import (
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
	"encoding/json"
	"net/http"

	"github.com/vchain-us/vcn/pkg/bom_component"
)

const (
	unknown = iota
	pipenv
	poetry
	pip
)

const (
	pipenvFileName = "Pipfile.lock"
	poetryFileName = "poetry.lock"
	pipFileName    = "requirements.txt"
	pypiApiPrefix = "https://pypi.org/pypi/"
)

// PythonPackage implements Package interface
type PythonPackage struct {
	dirName string
	pkgType int
}

// New returns new PythonPackage object, or nil if the path doesn't point to directory with Python package
func New(dir string) *PythonPackage {
	_, err := os.Stat(filepath.Join(dir, pipenvFileName))
	if err == nil {
		return &PythonPackage{pkgType: pipenv, dirName: dir}
	}

	_, err = os.Stat(filepath.Join(dir, poetryFileName))
	if err == nil {
		return &PythonPackage{pkgType: poetry, dirName: dir}
	}

	_, err = os.Stat(filepath.Join(dir, pipFileName))
	if err == nil {
		return &PythonPackage{pkgType: pip, dirName: dir}
	}

	return nil
}

func (p *PythonPackage) Type() string {
	return "Python"
}

func (p *PythonPackage) Path() string {
	return p.dirName
}

func (p *PythonPackage) Close() {
	// do nothing, function needed to comply with Package interface
}

func (p *PythonPackage) Components() ([]bom_component.Component, error) {
	switch p.pkgType {
	case pipenv:
		return procPipenv(p.dirName)
	case poetry:
		return procPoetry(p.dirName)
	case pip:
		return procPip(p.dirName)
	default:
		// should never happen
		return nil, errors.New("unknown package type")
	}
}

// combine multiple hashes into single hash by XORing them. Return Base16-encoded hash
// hash entry has a form of "<hash_type>:<base16-encoded hash>", where "<hash_type>:" is optional
func combineHashes(hashes []string) (string, int, error) {
	if len(hashes) == 0 {
		return "", bom_component.HashInvalid, nil
	}
	hashType := bom_component.HashInvalid
	var res []byte
	for _, v := range hashes {
		fields := strings.SplitN(v, ":", 2)
		hash := fields[0]
		if len(fields) >= 2 {
			hash = fields[1]
			switch fields[0] {
			case "sha256":
				hashType = bom_component.HashSHA256
			case "md5":
				hashType = bom_component.HashMD5
			}
		}
		comp, err := hex.DecodeString(hash)
		if err != nil {
			return "", bom_component.HashInvalid, errors.New("malformed hash value")
		}
		if res == nil {
			res = comp
		} else {
			if len(comp) != len(res) {
				// should never happen - all hashes must be of the same length
				return "", bom_component.HashInvalid, errors.New("malformed hash value")
			}
			// XOR hash
			for i := 0; i < len(res); i++ {
				res[i] ^= comp[i]
			}
		}
	}

	return hex.EncodeToString(res), hashType, nil
}

// query PyPI.org for module hash, combine all available hashes using XOR
func QueryHash(name, version string) (string, int, error) {
	resp, err := http.Get(pypiApiPrefix + name + "/" + version + "/json")
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