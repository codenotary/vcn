/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package python

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/vchain-us/vcn/pkg/bom/artifact"
)

// pythonArtifactFromPipEnv implements Artifact interface
type pythonArtifactFromPipEnv struct {
	pythonArtifact
}

// Dependencies returns list of Python dependencies for the artifact
// for pipenv Pipfile.lock JSON file contains all needed information in "default" section
func (a *pythonArtifactFromPipEnv) Dependencies() ([]artifact.Dependency, error) {
	if a.Deps != nil {
		return a.Deps, nil
	}
	file, err := os.Open(filepath.Join(a.path, pipenvFileName))
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

	res := make([]artifact.Dependency, 0, len(packages))
	for name, pkg := range packages {
		pkgContent, ok := pkg.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed \"%s\" section", name)
		}

		var dep artifact.Dependency
		dep.Name = name
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

		dep.Hash, dep.HashType, err = combineHashes(hashes)
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
		dep.Version = strings.TrimPrefix(version, "==")

		res = append(res, dep)
	}

	a.Deps = res
	return res, nil
}
