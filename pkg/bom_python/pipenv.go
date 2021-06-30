package bom_python

import (
	"os"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"encoding/hex"

	component "github.com/vchain-us/vcn/pkg/bom_component"
)

// for pipenv Pipfile.lock JSON file contains all needed information in "default" section
func procPipenv(dir string) ([]component.Component, error) {
	file, err := os.Open(filepath.Join(dir, "Pipfile.lock"))
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

	res := make([]component.Component, 0, len(packages))
	for name, pkg := range packages {
		pkgContent, ok := pkg.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed \"%s\" section", name)
		}
	
		var comp component.Component
		comp.Name = name
		field, ok := pkgContent["hashes"]
		if !ok {
			return nil, fmt.Errorf("malformed \"%s\" section", name)
		}
		hashes, ok := field.([]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed \"%s\" section", name)
		}

		comp.Hash, err = combineHashes(hashes)
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

// combine multiple hashes into single hash by XORing them
// hash entry has a form of "<hash_type>:<base16-encoded hash>"
func combineHashes(hashes []interface{}) (string, error) {
	if len(hashes) == 0 {
		return "", nil
	}
	var res []byte
	for _, v := range hashes {
		hash, ok := v.(string)
		if !ok { 
			return "", errors.New("malformed hash value")
		}
		split := strings.SplitN(hash, ":", 2)
		if len(split) < 2 {
			return "", errors.New("malformed hash value")
		}
		comp, err := hex.DecodeString(split[1])
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
