/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package serve

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/store"
)

// NewCmdLogin returns the cobra command for `vcn login`
func NewCmdServe() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start Web Server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return Execute()
		},
		Args: cobra.NoArgs,
	}

	return cmd
}

// Execute the login action
func Execute() error {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/sign", sign).Methods("POST")
	router.HandleFunc("/verify/{hash}", verify).Methods("GET")
	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8080", router))
	return nil
}

func currentUser() (*api.User, error) {
	email := store.Config().CurrentContext
	if email != "" {
		return nil, fmt.Errorf("No user has been set for current context")
	}
	u := api.NewUser(email)
	hasAuth, err := u.IsAuthenticated()
	if err != nil {
		return nil, err
	}
	if !hasAuth {
		return nil, fmt.Errorf("Current user is not authenticated")
	}
	return u, nil
}

func writeErr(w http.ResponseWriter, err error) {
	b, _ := json.MarshalIndent(types.NewError(err), "", "  ")
	fmt.Fprintln(w, string(b))
}

func writeResult(w http.ResponseWriter, r *types.Result) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		writeErr(w, err)
		return
	}
	fmt.Fprintln(w, string(b))
}

func sign(w http.ResponseWriter, r *http.Request) {
	user, err := currentUser()
	if err != nil {
		writeErr(w, err)
		return
	}
	isExist, err := user.IsExist()
	if err != nil {
		writeErr(w, err)
		return
	}
	if !isExist {
		writeErr(w, fmt.Errorf("no such user"))
		return
	}

	var a api.Artifact
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&a)
	if err != nil {
		writeErr(w, err)
		return
	}

	verification, err := user.Sign(a, "<key>", "<passphrase>", 0, 0)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeResult(w, types.NewResult(&a, nil, verification))
}

func verify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	w.WriteHeader(http.StatusOK)

	user, _ := currentUser()

	verification, err := api.BlockChainVerify(hash)
	if err != nil {
		writeErr(w, err)
		return
	}

	var artifact *api.ArtifactResponse
	if !verification.Unknown() {
		artifact, _ = api.LoadArtifactForHash(user, hash, verification.MetaHash())
	}

	writeResult(w, types.NewResult(nil, artifact, verification))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Index")
}
