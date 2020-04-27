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
)

const verifiedUser = "leonardo@vchain.us"

func TestCheckTokenNoInput(t *testing.T) {
	token := ""
	ret, _ := checkToken(token)

	if ret != false {
		t.Error("CheckToken() with empty string input must return false")
	}
}

func TestPublisherExists(t *testing.T) {
	ret, err := checkUserExists(verifiedUser)

	if ret == false || err != nil {
		t.Errorf("checkUserExists() must return true for infamous <%s> \n", verifiedUser)
	}
}
