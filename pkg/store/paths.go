/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package store

import (
	"os"
	"path/filepath"

	"github.com/vchain-us/vcn/pkg/meta"
)

var dir = DefaultDirName
var configFilepath string

func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, DirPerm); err != nil {
			return err
		}
	}
	return nil
}

func defaultConfigFilepath() string {
	return filepath.Join(dir, configFilename)
}

// SetDefaultDir sets the default store working directory (eg. /tmp/.vcn)
func SetDefaultDir() error {
	// Find home directory
	tmpDir := os.TempDir()
	var vcn string
	switch meta.StageEnvironment() {
	case meta.StageStaging:
		vcn = DefaultDirName + ".staging"
	case meta.StageTest:
		vcn = DefaultDirName + ".test"
	default:
		vcn = DefaultDirName
	}
	SetDir(filepath.Join(tmpDir, vcn))
	return nil
}

// SetDir sets the store working directory (eg. /tmp/.vcn)
func SetDir(p string) {
	dir = p
}

// ConfigFile returns the config file path
func ConfigFile() string {
	if configFilepath == "" {
		return defaultConfigFilepath()
	}
	return configFilepath
}

// SetConfigFile sets the config file path (e.g. /tmp/.vcn/config.json)
func SetConfigFile(filepath string) {
	configFilepath = filepath
}

// CurrentConfigFilePath returns the current config file path (e.g. /tmp/.vcn/config.json)
func CurrentConfigFilePath() string {
	return dir
}
