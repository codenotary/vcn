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
)

type signRequest struct {
	Dataa string `json:"dataa"`
	Datab string `json:"datab"`
	Datac string `json:"datac"`
	Datad string `json:"datad"`
	Datae string `json:"datae"`
}

type verifyResult struct {
	Artifact     *api.ArtifactResponse       `json:"artifact"`
	Verification *api.BlockchainVerification `json:"verification"`
	Hash         string                      `json:"hash"`
}

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

func sign(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var t signRequest
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}

	user := api.NewUser("user@cloud-r.eu")
	isExist, err := user.IsExist()
	if err != nil {
		panic(err)
	}
	if !isExist {
		panic(fmt.Errorf("no such user"))
	}
	// Make the artifact to be signed
	var a api.Artifact
	a.Hash = "myhash"
	ver, err := user.Sign(a, "<key>", "<passphrase>", 0, 0)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, ver)
	fmt.Println("Starting server")

}

func verify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	w.WriteHeader(http.StatusOK)

	var user *api.User // todo: get current user

	verification, err := api.BlockChainVerify(hash)
	if err != nil {
		panic(err)
	}

	var artifact *api.ArtifactResponse
	if !verification.Unknown() {
		artifact, _ = api.LoadArtifactForHash(user, hash, verification.MetaHash())
	}

	res := verifyResult{
		Verification: verification,
		Artifact:     artifact,
		Hash:         hash,
	}

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, string(b))

}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Index")
}
