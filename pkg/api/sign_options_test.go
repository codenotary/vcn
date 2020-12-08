/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"testing"

	"github.com/vchain-us/vcn/pkg/meta"

	"github.com/stretchr/testify/assert"
)

func testUser() User {
	return *NewUser("example@example.net")
}

func TestMakeSignOpts(t *testing.T) {
	reader := "Hi!"
	pass := "word"
	o, err := makeSignOpts(
		SignWithKey(reader, pass),
	)
	assert.NoError(t, err)

	assert.Equal(t, reader, o.keyin)
	assert.Equal(t, pass, o.passphrase)

	// test defaults
	assert.Equal(t, meta.StatusTrusted, o.status)
	assert.Equal(t, meta.VisibilityPrivate, o.visibility)
}

func TestSignWithStatus(t *testing.T) {
	o := &signOpts{}
	SignWithStatus(meta.StatusUnsupported)(o)

	assert.Equal(t, meta.StatusUnsupported, o.status)
}

func TestSignWithVisibility(t *testing.T) {
	o := &signOpts{}
	SignWithVisibility(meta.VisibilityPublic)(o)

	assert.Equal(t, meta.VisibilityPublic, o.visibility)
}

func TestSignWithKey(t *testing.T) {
	reader := "Hi!"
	pass := "word"

	o := &signOpts{}
	SignWithKey(reader, pass)(o)

	assert.Equal(t, reader, o.keyin)
	assert.Equal(t, pass, o.passphrase)
}
