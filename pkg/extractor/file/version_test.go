/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferVer(t *testing.T) {
	testCases := map[string]string{
		// Supported
		"vcn-v0.4.0-darwin-10.6-amd64":     "0.4.0",
		"vcn-v0.4.0-linux-amd64":           "0.4.0",
		"vcn-v0.4.0-windows-4.0-amd64.exe": "0.4.0",

		// Unsupported
		"codenotary_vcn_0.4.0_setup.exe": "",
	}

	for filename, ver := range testCases {
		assert.Equal(t, ver, inferVer(filename), "wrong version for %s", filename)
	}

}
