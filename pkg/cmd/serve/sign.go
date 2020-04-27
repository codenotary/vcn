/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package serve

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/meta"
)

func signHander(state meta.Status) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s := state
		k := make(map[string]bool)
		for _, scheme := range extractor.Schemes() {
			k[scheme] = true
		}
		sign(s, k, w, r)
	}
}

func sign(status meta.Status, kinds map[string]bool, w http.ResponseWriter, r *http.Request) {
	user, passphrase, err := getCredential(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	if user == nil {
		writeError(w, http.StatusUnauthorized, fmt.Errorf("bad or missing credentials"))
		return
	}

	keyin, _, offline, err := user.Secret()
	if err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	if offline {
		writeError(w, http.StatusConflict, fmt.Errorf("offline secret is not yet supported"))
		return
	}

	opts := []api.SignOption{
		api.SignWithKey(keyin, passphrase),
		api.SignWithStatus(status),
	}

	if _, public := r.URL.Query()["public"]; public {
		opts = append(opts, api.SignWithVisibility(meta.VisibilityPublic))
	}

	decoder := json.NewDecoder(r.Body)
	artifact := new(api.Artifact)
	err = decoder.Decode(artifact)

	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if artifact.Name == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("name cannot be empty"))
		return
	}

	if !kinds[artifact.Kind] {
		writeError(w, http.StatusBadRequest, fmt.Errorf(`"%s" is not a valid value for kind`, artifact.Kind))
		return
	}

	artifact.Hash = strings.ToLower(artifact.Hash)

	verification, err := user.Sign(
		artifact,
		opts...,
	)

	// todo(ameingast/leogr): remove reduntat event - need backend improvement
	api.TrackPublisher(user, meta.VcnSignEvent)
	api.TrackSign(user, artifact.Hash, artifact.Name, status)

	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var ar *api.ArtifactResponse
	if !verification.Unknown() {
		ar, _ = api.LoadArtifact(user, artifact.Hash, verification.MetaHash())
	}

	writeResult(w, http.StatusOK, types.NewResult(artifact, ar, verification))
}
