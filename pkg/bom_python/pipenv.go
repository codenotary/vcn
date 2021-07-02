package bom_python

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/vchain-us/vcn/pkg/bom_component"
)

// for pipenv Pipfile.lock JSON file contains all needed information in "default" section
func procPipenv(dir string) ([]bom_component.Component, error) {
	file, err := os.Open(filepath.Join(dir, pipenvFileName))
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var content map[string]interface{}
	err = json.Unmarshal(buf, &content)
	if err != nil {
		return nil, err
	}

	deflt, ok := content["default"]
	if !ok {
		return nil, errors.New("malformed Pipfile.lock - missing \"default\"")
	}
	packages, ok := deflt.(map[string]interface{})
	if !ok {
		return nil, errors.New("malformed Pipfile.lock - \"default\" has a wrong data type")
	}

	res := make([]bom_component.Component, 0, len(packages))
	for name, pkg := range packages {
		pkgContent, ok := pkg.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed \"%s\" section", name)
		}

		var comp bom_component.Component
		comp.Name = name
		field, ok := pkgContent["hashes"]
		if !ok {
			return nil, fmt.Errorf("malformed \"%s\" section", name)
		}
		hashArray, ok := field.([]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed \"%s\" section", name)
		}

		hashes := make([]string, len(hashArray))
		for i := range hashArray {
			hashes[i] = hashArray[i].(string)
		}

		comp.Hash, comp.HashType, err = combineHashes(hashes)
		if err != nil {
			return nil, err
		}

		field, ok = pkgContent["version"]
		if !ok {
			return nil, fmt.Errorf("malformed \"%s\" section", name)
		}
		version, ok := field.(string)
		if !ok {
			return nil, fmt.Errorf("malformed \"%s\" section", name)
		}
		comp.Version = strings.TrimPrefix(version, "==")

		res = append(res, comp)
	}

	return res, nil
}
