package bom_python

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"strings"

	component "github.com/vchain-us/vcn/pkg/bom_component"
)

const (
	unknown = iota
	pipenv
	poetry
	reqs
)

// PythonPackage implements Package interface
type PythonPackage struct {
	dirName string
	pkgType int
}

// New returns new PythonPackage object, or nil if the path doesn't point to directory with Python package
func New(dir string) *PythonPackage {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil // not a directory or cannot be accessed
	}

	packageType := unknown
	// check what type of package is it
loop:
	for _, file := range files {
		switch file.Name() {
		case "Pipfile.lock":
			packageType = pipenv
			break loop // don't look further
		case "poetry.lock":
			packageType = poetry
			break loop // don't look further
		case "requirements.txt":
			packageType = reqs
			// keep looking for other files
		}
	}

	if packageType == unknown {
		return nil // not Python signature files found
	}

	return &PythonPackage{pkgType: packageType, dirName: dir}
}

func (p *PythonPackage) Type() string {
	return "Python"
}

func (p *PythonPackage) Close() {
	// do nothing, function needed to comply with Package interface
}

func (p *PythonPackage) Components() ([]component.Component, error) {
	switch p.pkgType {
	case pipenv:
		return procPipenv(p.dirName)
	case poetry:
		return procPoetry(p.dirName)
	case reqs:
		return procPip(p.dirName)
	default:
		// should never happen
		return nil, errors.New("unknown package type")
	}
}

// combine multiple hashes into single hash by XORing them
// hash entry has a form of "<hash_type>:<base16-encoded hash>", where "<hash_type>:" is optional
func combineHashes(hashes []string) (string, error) {
	if len(hashes) == 0 {
		return "", nil
	}
	var res []byte
	for _, v := range hashes {
		split := strings.SplitN(v, ":", 2)
		hash := split[0]
		if len(split) >= 2 {
			hash = split[1]
		}
		comp, err := hex.DecodeString(hash)
		if err != nil {
			return "", errors.New("malformed hash value")
		}
		if res == nil {
			res = comp
		} else {
			if len(comp) != len(res) {
				// should never happen - all hashes must be of the same length
				return "", errors.New("malformed hash value")
			}
			// XOR hash
			for i := 0; i < len(res); i++ {
				res[i] ^= comp[i]
			}
		}
	}

	return hex.EncodeToString(res), nil
}
