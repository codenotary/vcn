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
	"github.com/vchain-us/vcn/pkg/meta"
)

func lcSign(user *api.LcUser, status meta.Status, kinds map[string]bool, w http.ResponseWriter, r *http.Request) {
	if user.Client.ApiKey == "" {
		writeError(w, http.StatusUnauthorized, fmt.Errorf("api key not provided"))
		return
	}
	err := user.Client.Connect()
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	opts := []api.LcSignOption{
		api.LcSignWithStatus(status),
	}

	if r.Body == http.NoBody {
		writeError(w, http.StatusBadRequest, fmt.Errorf("no artifact submitted"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var artifact api.Artifact
	err = decoder.Decode(&artifact)
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

	verified, tx, err := user.Sign(
		artifact,
		opts...,
	)

	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ar, verified, err := user.LoadArtifact(artifact.Hash, "", "", tx)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeLcResult(w, http.StatusOK, types.NewLcResult(ar, verified, nil))
}
