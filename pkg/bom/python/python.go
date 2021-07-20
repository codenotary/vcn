/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package python

import (
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
	"encoding/json"
	"net/http"

	"github.com/vchain-us/vcn/pkg/bom/artifact"
)

const AssetType = "python"

// pythonArtifact implements Artifact interface
type pythonArtifact struct {
	artifact.GenericArtifact
	path string
}

type pypiUrls struct {
	Files []pypiFile `json:"urls"`
}

type pypiFile struct {
	Digests pypiDigests `json:"digests"`
}

type pypiDigests struct {
	Md5    string `json:"md5,omitempty"`
	Sha256 string `json:"sha256,omitempty"`
}

const (
	pipenvFileName = "Pipfile.lock"
	poetryFileName = "poetry.lock"
	pipFileName    = "requirements.txt"
	pypiApiPrefix = "https://pypi.org/pypi/"
)

// New returns new Artifact, or nil if the path doesn't point to directory with Python package
func New(path string) artifact.Artifact {
	_, err := os.Stat(filepath.Join(path, pipenvFileName))
	if err == nil {
		return &pythonArtifactFromPipEnv{pythonArtifact: pythonArtifact{path: path}}
	}

	_, err = os.Stat(filepath.Join(path, poetryFileName))
	if err == nil {
		return &pythonArtifactFromPoetry{pythonArtifact: pythonArtifact{path: path}}
	}

	_, err = os.Stat(filepath.Join(path, pipFileName))
	if err == nil {
		return &pythonArtifactFromPip{pythonArtifact: pythonArtifact{path: path}}
	}

	return nil
}

func (p pythonArtifact) Type() string {
	return AssetType
}

func (p pythonArtifact) Path() string {
	return p.path
}

// combine multiple hashes into single hash by XORing them. Return Base16-encoded hash
// hash entry has a form of "<hash_type>:<base16-encoded hash>", where "<hash_type>:" is optional
func combineHashes(hashes []string) (string, artifact.HashType, error) {
	if len(hashes) == 0 {
		return "", artifact.HashInvalid, nil
	}
	hashType := artifact.HashInvalid
	var res []byte
	for _, v := range hashes {
		fields := strings.SplitN(v, ":", 2)
		hash := fields[0]
		if len(fields) >= 2 {
			hash = fields[1]
			switch fields[0] {
			case "sha256":
				hashType = artifact.HashSHA256
			case "md5":
				hashType = artifact.HashMD5
			}
		}
		comp, err := hex.DecodeString(hash)
		if err != nil {
			return "", artifact.HashInvalid, errors.New("malformed hash value")
		}
		if res == nil {
			res = comp
		} else {
			if len(comp) != len(res) {
				// should never happen - all hashes must be of the same length
				return "", artifact.HashInvalid, errors.New("malformed hash value")
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
func QueryHash(name, version string) (string, artifact.HashType, error) {
	resp, err := http.Get(pypiApiPrefix + name + "/" + version + "/json")
	if err != nil {
		return "", artifact.HashInvalid, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", artifact.HashInvalid, errors.New("cannot query PyPI for package details")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", artifact.HashInvalid, err
	}

	var urls pypiUrls
	err = json.Unmarshal(body, &urls)
	if err != nil {
		return "", artifact.HashInvalid, err
	}

	// assuming that all files have the same type of hash, with priority for SHA-256
	hashType := artifact.HashMD5
	if urls.Files[0].Digests.Sha256 != "" {
		hashType = artifact.HashSHA256
	}
	hashes := make([]string, len(urls.Files))
	for i, file := range urls.Files {
		if hashType == artifact.HashSHA256 {
			hashes[i] = file.Digests.Sha256
		} else {
			hashes[i] = file.Digests.Md5
		}
	}

	hash, _, err := combineHashes(hashes)
	if err != nil {
		return "", artifact.HashInvalid, errors.New("malformed hash value")
	}

	return hash, hashType, nil
}
