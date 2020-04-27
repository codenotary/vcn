/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bundle

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/opencontainers/go-digest"

	"github.com/stretchr/testify/assert"

	"testing"
)

func getTestManifest(t *testing.T) *Manifest {
	items := make([]Descriptor, 0)
	for path, src := range map[string]io.Reader{
		"digits.txt":     strings.NewReader("1234567890"),
		"dup-digits.txt": strings.NewReader("1234567890"),
		"letters.txt":    strings.NewReader("abcdef"),
	} {
		d, err := NewDescriptor(path, src)
		assert.NoError(t, err)
		items = append(items, *d)
	}

	return NewManifest(items...)
}

func TestManifest(t *testing.T) {

	m := getTestManifest(t)

	assert.NotNil(t, m)

	d, err := m.Digest()
	assert.NoError(t, err)
	assert.Equal(t, "sha256:2cc48ce16beff9987ad31f32a7623bda48317d9232dbb75ba9a83f6d85ed073e", d.String())

	j, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.Equal(
		t,
		`{"schemaVersion":1,"items":[{"digest":"sha256:bef57ec7f53a6d40beb640a780a639c83bc29ac8a9816f1fc6c5c6dcd93c4721","size":6,"paths":["letters.txt"]},{"digest":"sha256:c775e7b757ede630cd0aa1113bd102661ab38829ca52a6422ab782862f268646","size":10,"paths":["digits.txt","dup-digits.txt"]}]}`,
		string(j),
	)
}

func TestWriteReadManifest(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "vcn-manifest")
	if err != nil {
		log.Fatal(err)
	}
	filename := tmpfile.Name()
	defer os.Remove(filename) // clean up
	m := getTestManifest(t)

	assert.NoError(t, WriteManifest(*m, filename))

	assert.NoError(t, m.Normalize())

	mm, err := ReadManifest(filename)
	assert.NoError(t, err)
	assert.Equal(t, m, mm)
}

func TestManifestWrongVersion(t *testing.T) {
	m := Manifest{
		SchemaVersion: 11,
	}
	err := m.Normalize()
	assert.Error(t, err)
}

func TestManifestWrongAlgo(t *testing.T) {
	m := Manifest{
		SchemaVersion: 1,
		Items: []Descriptor{
			{
				Digest: digest.Digest("sha512:cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"),
				Size:   0,
				Paths:  []string{"file"},
			},
		},
	}
	assert.Error(t, m.Normalize())

	m.Items[0].Digest = digest.Digest("sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	assert.NoError(t, m.Normalize())
}

func TestManifestSameDigestMultipleSizes(t *testing.T) {
	m := Manifest{
		SchemaVersion: 1,
		Items: []Descriptor{
			{
				Digest: digest.Digest("sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
				Size:   0,
				Paths:  []string{"file"},
			},
			{
				Digest: digest.Digest("sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
				Size:   1,
				Paths:  []string{"file1"},
			},
		},
	}
	assert.Error(t, m.Normalize())
}
