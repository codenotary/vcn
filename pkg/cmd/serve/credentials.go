/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package serve

import (
	"net/http"

	"github.com/vchain-us/vcn/pkg/api"
)

func getCredential(r *http.Request) (user *api.User, passphrase string, err error) {
	if email, password, ok := r.BasicAuth(); ok {
		user = api.NewUser(email)
		// we don't support otp from serve
		err = user.Authenticate(password, "")
		if err == nil {
			if empty := r.Header.Get("x-notarization-password-empty"); empty == "" {
				passphrase = r.Header.Get("x-notarization-password")
				if passphrase == "" {
					passphrase = password
				}
			}
			// else use empty passphrase
		}
	}
	return
}

func getLcUser(r *http.Request, lcHost, lcPort, lcCert string) (*api.LcUser, error) {
	apikey := r.Header.Get("x-notarization-lc-api-key")
	return api.NewLcUser(apikey, lcHost, lcPort, lcCert)
}
