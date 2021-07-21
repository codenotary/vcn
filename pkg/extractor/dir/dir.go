/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/bundle"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/uri"
)

// Scheme for dir
const Scheme = "dir"

// ManifestKey is the metadata's key for storing the manifest
const ManifestKey = "manifest"

// PathKey is the metadata's key for the directory path
const PathKey = "path"

type opts struct {
	initIgnoreFile    bool
	skipIgnoreFileErr bool
}

// Artifact returns a file *api.Artifact from a given u
func Artifact(u *uri.URI, options ...extractor.Option) ([]*api.Artifact, error) {

	if u.Scheme != Scheme {
		return nil, nil
	}

	opts := &opts{}
	if err := extractor.Options(options).Apply(opts); err != nil {
		return nil, err
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
		return nil, fmt.Errorf("read %s: is not a directory", path)
	}

	if opts.initIgnoreFile {
		if err := InitIgnoreFile(path); err != nil {
			if !opts.skipIgnoreFileErr {
				return nil, err
			}
		}
	}

	files, err := Walk(path)
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

	return []*api.Artifact{{
		Kind:     Scheme,
		Hash:     digest.Encoded(),
		Name:     stat.Name(),
		Metadata: m,
	}}, nil
}

// WithIgnoreFileInit returns a functional option to instruct the dir's extractor to create the defualt ignore file
// when not yet present into the targeted directory.
func WithIgnoreFileInit() extractor.Option {
	return func(o interface{}) error {
		if o, ok := o.(*opts); ok {
			o.initIgnoreFile = true
		}
		return nil
	}
}

// WithSkipIgnoreFileErr returns a functional option to instruct the dir's extractor to skip errors during
// ignore file initialization.
func WithSkipIgnoreFileErr() extractor.Option {
	return func(o interface{}) error {
		if o, ok := o.(*opts); ok {
			o.skipIgnoreFileErr = true
		}
		return nil
	}
}
