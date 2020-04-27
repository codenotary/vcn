/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package store

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"

	"github.com/vchain-us/vcn/pkg/bundle"
)

func ManifestFilepath(kind, target string) (string, error) {
	target, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, defaultManifestsDir)
	if err := ensureDir(path); err != nil {
		return "", err
	}
	id := sha256.Sum256([]byte(target))

	return filepath.Join(path, fmt.Sprintf("%s_%x.json", kind, id)), nil
}

func SaveManifest(kind, target string, manifest bundle.Manifest) error {
	path, err := ManifestFilepath(kind, target)
	if err != nil {
		return err
	}
	return bundle.WriteManifest(manifest, path)
}

func ReadManifest(kind, target string) (*bundle.Manifest, error) {
	path, err := ManifestFilepath(kind, target)
	if err != nil {
		return nil, err
	}
	return bundle.ReadManifest(path)
}
