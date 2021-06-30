package bom_python

import (
	"errors"
	"io/ioutil"

	component "github.com/codenotary/vcn/pkg/bom_component"
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
			break loop		// don't look further
		case "poetry.lock":
			packageType = poetry
			break loop		// don't look further
		case "requirements.txt":
			packageType = reqs
			// keep looking for other files
		}
	}

	if packageType == unknown {
		return nil	// not Python signature files found
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
		return procReqs(p.dirName)
	default:
		// should never happen
		return nil, errors.New("unknown package type")
	}
}



func procPoetry(dir string) ([]component.Component, error) {
	return nil, nil
}

func procReqs(dir string) ([]component.Component, error) {
	return nil, nil
}