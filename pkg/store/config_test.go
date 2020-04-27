/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package store

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mkTmpForConfig(t *testing.T) string {
	tdir, err := ioutil.TempDir("", "vcn-test-store-config")
	if err != nil {
		t.Fatal(err)
	}
	return tdir
}

func TestLoadConfig(t *testing.T) {
	tdir := mkTmpForConfig(t)
	SetDir(tdir + "/" + DefaultDirName)

	assert.Nil(t, Config())

	err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, Config())
	assert.Equal(t, filepath.Join(tdir, DefaultDirName, configFilename), ConfigFile())
	assert.FileExists(t, ConfigFile())
	assert.NotNil(t, Config())
}

func TestSaveConfig(t *testing.T) {
	tdir := mkTmpForConfig(t)
	SetDir(tdir + "/" + DefaultDirName)

	email := "example@example.net"

	cfg = &ConfigRoot{
		CurrentContext: email,
		Users: []*User{
			{
				Email:    email,
				KeyStore: filepath.Join(tdir, "u", email, "k"),
			},
		},
	}

	err := SaveConfig()
	assert.NoError(t, err)

	err = LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, email, Config().CurrentContext)
}

func TestConfigClearContext(t *testing.T) {
	email := "example@example.net"

	cfg = &ConfigRoot{
		CurrentContext: email,
		Users: []*User{
			{
				Email: email,
				Token: "dummy",
			},
		},
	}

	cfg.ClearContext()
	assert.Empty(t, cfg.CurrentContext)

	for _, u := range cfg.Users {
		assert.Empty(t, u.Token)
	}
}
