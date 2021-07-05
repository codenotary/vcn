package bom_go

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"os"
	"io/ioutil"

	"github.com/vchain-us/vcn/pkg/bom_component"
)

// GoPackage implements Package interface
type GoPackage struct {
	path string
	file exe		// file is nil when getting info from go.mod, not binary
}

// New returns new GoPackage object, or nil if filename doesn't referer to ELF, built from Go source, or directory
// containing Go files
func New(path string) *GoPackage {
	// if path is directory, checks the presence of Go files, otherwise check if file is executable and contains
	// Go build info section
	fi, err := os.Stat(path)
	if err != nil {
		return nil
	}
	if fi.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return nil
		}
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".go") {
				return &GoPackage{path: path}
			}
		}
		return nil
	} else {
		file, err := openExe(path)
		if err != nil {
			return nil // not a ELF binary
		}
		if file.DataStart() == 0 {
			file.Close()
			return nil // cannot find build info
		}
		return &GoPackage{path: path, file: file}
	}
}

func (p *GoPackage) Components() ([]bom_component.Component, error) {
	if p.file != nil {
		return exeComponents(p.file)
	}
	return goListComponents(p.path)
}

func (p *GoPackage) Type() string {
	return "Go"
}

func (p *GoPackage) Path() string {
	return p.path
}

func (p *GoPackage) Close() {
	if p.file != nil {
		p.file.Close()
	}
}

func goModHash(encoded string) (string, int, error) {
	hashType := bom_component.HashInvalid
	fields := strings.SplitN(encoded, ":", 2)
	if len(fields) != 2 {
		return "", bom_component.HashInvalid, errors.New("malformed hash value")
	}
	// At the time of writing "h1" (SHA256) is the only hash type, supported by Go
	if fields[0] == "h1" {
		hashType = bom_component.HashSHA256
	} else {
		return "", bom_component.HashInvalid, errors.New("unsupported hash type")
	}
	hash, err := base64.StdEncoding.DecodeString(fields[1])
	if err != nil {
		return "", bom_component.HashInvalid, fmt.Errorf("cannot decode base64 hash: %w", err)
	}

	return hex.EncodeToString(hash), hashType, nil
}
