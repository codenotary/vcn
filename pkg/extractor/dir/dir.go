/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package dir

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/vchain-us/vcn/pkg/meta"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/bundle"
	"github.com/vchain-us/vcn/pkg/uri"
)

// Scheme for dir
const Scheme = "dir"

// ManifestKey is the metadata's key for storing the manifest
const ManifestKey = "manifest"

// PathKey is the metadata's key for the directory path
const PathKey = "path"

// Artifact returns a file *api.Artifact from a given u
func Artifact(u *uri.URI) (*api.Artifact, error) {

	if u.Scheme != "" && u.Scheme != Scheme {
		return nil, nil
	}

	path := strings.TrimPrefix(u.Opaque, "//")
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	d, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	// get file info and check if is a directory
	stat, err := d.Stat()
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, nil
	}
    _, isset := os.LookupEnv(meta.VcnWriteIgnoreFile);
    if isset {
	    if err := initIgnoreFile(path); err != nil {
		    return nil, err
        }
	}

	files, err := walk(path)
	if err != nil {
		return nil, err
	}

	manifest := bundle.NewManifest(files...)
	digest, err := manifest.Digest()
	if err != nil {
		return nil, err
	}

	// Metadata container
	m := api.Metadata{
		ManifestKey: manifest,
		PathKey:     path,
	}

	return &api.Artifact{
		Kind:     Scheme,
		Hash:     digest.Encoded(),
		Name:     stat.Name(),
		Metadata: m,
	}, nil
}
