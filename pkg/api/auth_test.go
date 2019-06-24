/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"fmt"
	"testing"
)

const VERIFIED_USER = "mathias@vchain.us"

func TestCheckTokenNoInput(t *testing.T) {
	token := ""
	ret, _ := checkToken(token)

	if ret != false {
		t.Error(fmt.Sprintf(`CheckToken() with empty string input must return false`))
	}

}

func TestPublisherExists(t *testing.T) {
	ret, err := checkUserExists(VERIFIED_USER)

	if ret == false || err != nil {
		t.Error(fmt.Sprintf(`checkUserExists() must return true for infamous <%s>`, VERIFIED_USER))
	}

}
