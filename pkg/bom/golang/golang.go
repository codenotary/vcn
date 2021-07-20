/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package golang

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"github.com/vchain-us/vcn/pkg/bom/artifact"
)

const AssetType = "go"

// goArtifact implements Artifact interface
type goArtifact struct {
	artifact.GenericArtifact
	path string
}

// New returns new GoArtifact object, or nil if filename doesn't referer to ELF, built from Go source, or isn't
// a directory containing Go files
func New(path string) artifact.Artifact {
	fi, err := os.Stat(path)
	if err != nil {
		return nil
	}
	if fi.IsDir() {
		// a directory - checks the presence of Go files, if no Go files present, look for go.sum
		files, err := filepath.Glob(filepath.Join(path, "*.go"))
		if err == nil && len(files) > 0 {
			return &goArtifactFromGoList{goArtifact: goArtifact{ path: path}}
		}
		_, err = os.Stat(filepath.Join(path, "go.sum"))
		if err != nil {
			return nil
		}
		return &goArtifactFromSum{goArtifact: goArtifact{path: path}}
	} else {
		// not a directory - check if file is executable and contains Go build info section
		file, err := openExe(path)
		if err != nil {
			return nil // not a ELF binary
		}
		if file.DataStart() == 0 {
			file.Close()
			return nil // cannot find build info
		}
		return &goArtifactFromExe{goArtifact: goArtifact{path: path}, file: file}
	}
}

func (p goArtifact) Type() string {
	return AssetType
}

func (p goArtifact) Path() string {
	return p.path
}

func ModHash(encoded string) (string, artifact.HashType, error) {
	hashType := artifact.HashInvalid
	fields := strings.SplitN(encoded, ":", 2)
	if len(fields) != 2 {
		return "", artifact.HashInvalid, errors.New("malformed hash value")
	}
	// At the time of writing "h1" (SHA256) is the only hash type, supported by Go
	if fields[0] == "h1" {
		hashType = artifact.HashSHA256
	} else {
		return "", artifact.HashInvalid, errors.New("unsupported hash type")
	}
	hash, err := base64.StdEncoding.DecodeString(fields[1])
	if err != nil {
		return "", artifact.HashInvalid, fmt.Errorf("cannot decode base64 hash: %w", err)
	}

	return hex.EncodeToString(hash), hashType, nil
}
