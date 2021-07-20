/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package gocom

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/uri"
	"github.com/vchain-us/vcn/pkg/bom/artifact"
	"github.com/vchain-us/vcn/pkg/bom/golang"
)

// Scheme for Go component
const Scheme = "gocom"

type modInfo struct {
	Path    string
	Version string
	Sum     string
}

var modInfoArgs = []string{"mod", "download", "-json"}

// Artifact returns a file *api.Artifact from a given uri
func Artifact(u *uri.URI, options ...extractor.Option) ([]*api.Artifact, error) {
	if u.Scheme != Scheme {
		return nil, nil
	}

	path := strings.TrimPrefix(u.Opaque, "//")

	buf, err := exec.Command("go", append(modInfoArgs, path)...).Output()
	if err != nil {
		if len(buf) > 0 {
			// Error field in json may contain error message
			var fields map[string]string
			if nil != json.Unmarshal(buf, &fields) {
				// output isn't a valid JSON - return original execution error
				return nil, fmt.Errorf("cannot get Go module: %w", err)
			}
			errMsg, ok := fields["Error"]
			if !ok {
				return nil, fmt.Errorf("cannot get Go module: %w", err)
			}
			return nil, fmt.Errorf("cannot get Go module: %s", errMsg)
		}
	}

	var info modInfo
	err = json.Unmarshal(buf, &info)
	if err != nil {
		return nil, fmt.Errorf("cannot parse Go command output: %w", err)
	}

	hash, hashType, err := golang.ModHash(info.Sum)
	if err != nil {
		return nil, fmt.Errorf("cannot decode module hash: %w", err)
	}

	return []*api.Artifact{{
		Kind:        golang.AssetType,
		Name:        info.Path,
		Hash:        hash,
		Size:        uint64(len(buf)),
		ContentType: "text/json; charset=utf-8",
		Metadata:    api.Metadata{
			"path": info.Path,
			"version": info.Version,
			"hashType": artifact.HashTypeName(hashType)},
	}}, nil
}
