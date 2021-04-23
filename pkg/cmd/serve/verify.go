/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package serve

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/meta"
)

func (sh *handler) verify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := strings.ToLower(vars["hash"])

	if sh.lcHost != "" && sh.lcPort != "" {
		// todo @Michele move getLcUser in handler sh constructor
		lcUser, err := getLcUser(r, sh.lcHost, sh.lcPort, sh.lcCert, sh.lcSkipTlsVerify, sh.lcNoTls)
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		if lcUser.Client.ApiKey == "" {
			writeError(w, http.StatusUnauthorized, fmt.Errorf("api key not provided"))
			return
		}
		err = lcUser.Client.Connect()
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}

		ar, verified, err := lcUser.LoadArtifact(hash, "", 0)
		if err != nil {
			if err == api.ErrNotVerified {
				writeError(w, http.StatusConflict, err)
				return
			}
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeLcResult(w, http.StatusOK, types.NewLcResult(ar, verified))
		return
	}

	var keys []string
	org := r.URL.Query().Get("org")
	if org != "" {
		bo, err := api.GetBlockChainOrganisation(org)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		keys = bo.MembersIDs()
	} else {
		ks := r.URL.Query().Get("signers")
		if ks != "" {
			keys = strings.Split(ks, ",")
			// add 0x if missing, lower case
			for i, k := range keys {
				if !strings.HasPrefix(k, "0x") {
					keys[i] = "0x" + k
				}
				keys[i] = strings.ToLower(keys[i])
			}
		}
	}

	var err error
	var verification *api.BlockchainVerification
	user, _, err := getCredential(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	// if keys have been passed, check for a verification matching them
	if len(keys) > 0 {
		verification, err = api.VerifyMatchingSignerIDs(hash, keys)
	} else {
		// if we have an user, check for verification matching user's key first
		userKey := ""
		if user != nil {
			userKey, err = user.SignerID()
			if err != nil {
				writeError(w, http.StatusConflict, err)
				return
			}
		}
		if userKey != "" {
			verification, err = api.VerifyMatchingSignerIDWithFallback(hash, userKey)
		} else {
			// if no passed keys nor user,
			// just get the last with highest level available verification
			verification, err = api.Verify(hash)
		}
	}

	if err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}

	name := ""
	var artifact *api.ArtifactResponse
	if !verification.Unknown() {
		artifact, _ = api.LoadArtifact(user, hash, verification.MetaHash())
		if artifact != nil {
			name = artifact.Name
		}
	}

	// todo(ameingast/leogr): remove reduntat event - need backend improvement
	api.TrackPublisher(user, meta.VcnVerifyEvent)
	api.TrackVerify(user, hash, name)

	writeResult(w, http.StatusOK, types.NewResult(nil, artifact, verification))
}
