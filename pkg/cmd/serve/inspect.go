/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package serve

import (
	i "github.com/vchain-us/vcn/pkg/cmd/inspect"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (sh *handler) inspectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := strings.ToLower(vars["hash"])
	var signerID string
	keys, ok := r.URL.Query()["signerid"]

	if ok && len(keys[0]) > 0 {
		signerID = keys[0]
	}

	if sh.lcHost != "" && sh.lcPort != "" {
		// todo @Michele move getLcUser in handler sh constructor
		user, err := getLcUser(r, sh.lcHost, sh.lcPort, sh.lcCert, sh.lcSkipTlsVerify, sh.lcNoTls)
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		lcInspect(user, hash, signerID, w)
		return
	}

	user, _, err := getCredential(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	results, err := i.GetResults(hash, user)

	writeResults(w, http.StatusOK, results)
}
