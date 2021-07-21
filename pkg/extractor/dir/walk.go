/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
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

	"github.com/vchain-us/vcn/pkg/bundle"
)

func Walk(root string) (files []bundle.Descriptor, err error) {
	files = make([]bundle.Descriptor, 0)
	ignore, err := newIgnoreFileMatcher(root)
	if err != nil {
		return
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// skip irregular files (e.g. dir, symlink, pipe, socket, device...)
		if !info.Mode().IsRegular() {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		// descriptor's path must be OS agnostic
		relPath = filepath.ToSlash(relPath)

		// skip manifest and files matching the ignore patterns
		if relPath == bundle.ManifestFilename || ignore.Match(strings.Split(relPath, "/"), false) {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		d, err := bundle.NewDescriptor(relPath, file)
		file.Close()
		if err != nil {
			return err
		}
		files = append(files, *d)

		return nil
	})
	return
}
