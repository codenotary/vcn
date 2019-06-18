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
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/meta"
	"github.com/vchain-us/vcn/pkg/store"
)

type errorResponse struct {
	Message string `json:"message"`
	Code    uint64 `json:"code"`
	Error   error  `json:"error"`
}

// NewCmdServe returns the cobra command for `vcn serve`
func NewCmdServe() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start a web server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runServe(cmd)
		},
		Args: cobra.NoArgs,
	}
	cmd.Flags().String("host", "localhost", "host address to serve the application")
	cmd.Flags().String("port", "8080", "port to serve the application")
	return cmd
}

func runServe(cmd *cobra.Command) error {
	//CHECK PASSphrase
	passphrase := os.Getenv(meta.KeyStorePasswordEnv)
	if passphrase == "" {
		return fmt.Errorf("Server needs KEYSTORE_PASSWORD env")
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/sign", signHander(meta.StatusTrusted)).Methods("POST")
	router.HandleFunc("/untrust", signHander(meta.StatusUntrusted)).Methods("POST")
	router.HandleFunc("/unsupport", signHander(meta.StatusUnsupported)).Methods("POST")
	router.HandleFunc("/verify/{hash}", verify).Methods("GET")

	host, err := cmd.Flags().GetString("host")
	if err != nil {
		return nil
	}

	port, err := cmd.Flags().GetString("port")
	if err != nil {
		return nil
	}

	host += ":" + port

	fmt.Println("Starting server http://" + host)
	log.Fatal(http.ListenAndServe(host, router))

	return nil
}

func currentUser() (*api.User, error) {
	email := store.Config().CurrentContext
	if email == "" {
		return nil, fmt.Errorf("no user has been set for current context")
	}
	u := api.NewUser(email)
	hasAuth, err := u.IsAuthenticated()
	if err != nil {
		return u, fmt.Errorf("current user is not authenticated")
	}
	if !hasAuth {
		return u, fmt.Errorf("current user is not authenticated")
	}
	return u, nil
}

func writeErrorResponse(w http.ResponseWriter, message string, err error, code uint64) {
	w.Header().Set("Content-Type", "application/json")
	var errResponse errorResponse
	errResponse.Message = message
	errResponse.Code = code
	if err != nil {
		errResponse.Error = err
	}
	b, jerr := json.Marshal(errResponse)
	w.WriteHeader(http.StatusBadRequest)
	if jerr == nil {
		fmt.Fprintln(w, string(b))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func writeResponse(w http.ResponseWriter, r *types.Result) {
	b, err := json.Marshal(r)
	if err != nil {
		writeErrorResponse(w, "", err, 400)
		return
	}
	if b == nil {
		writeErrorResponse(w, "", err, 400)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func signHander(state meta.Status) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s := state
		sign(s, w, r)
	}
}

func sign(state meta.Status, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var jsonRequest api.ArtifactRequest
	err := decoder.Decode(&jsonRequest)

	if err != nil {
		writeErrorResponse(w, "invalid Request Body", err, http.StatusBadRequest)
		return
	}

	// check hash in request
	if jsonRequest.Hash == "" {
		writeErrorResponse(w, "invalid hash", nil, http.StatusBadRequest)
		return
	}
	if jsonRequest.Size <= 0 {
		writeErrorResponse(w, "invalid size", nil, http.StatusBadRequest)
		return
	}
	if jsonRequest.Name == "" {
		writeErrorResponse(w, "invalid name", nil, http.StatusBadRequest)
		return
	}
	if jsonRequest.ContentType == "" {
		writeErrorResponse(w, "invalid contentType", nil, http.StatusBadRequest)
		return
	}
	if jsonRequest.Kind == "" {
		writeErrorResponse(w, "invalid kind", nil, http.StatusBadRequest)
		return
	}

	user, err := currentUser()
	if user == nil || err != nil {
		writeErrorResponse(w, "no sourch user", err, http.StatusBadRequest)
		return
	}

	pubKey := user.DefaultKey()
	if pubKey == "" {
		writeErrorResponse(w, "invalid pubKey", nil, http.StatusBadRequest)
		return
	}

	// Make the artifact to be signed
	var a api.Artifact
	m := api.Metadata{}
	a.Hash = jsonRequest.Hash
	a.Name = jsonRequest.Name
	a.Size = jsonRequest.Size
	a.Kind = jsonRequest.Kind
	a.ContentType = jsonRequest.ContentType
	a.Metadata = m

	verification, err := user.Sign(a, pubKey, os.Getenv(meta.KeyStorePasswordEnv), state, meta.VisibilityFromString(jsonRequest.Visibility))
	if err != nil {
		writeErrorResponse(w, "sign error", err, http.StatusBadRequest)
		return
	}

	writeResponse(w, types.NewResult(&a, nil, verification))
}

func verify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	user, _ := currentUser()

	verification, err := api.BlockChainVerify(hash)
	if err != nil {
		panic(err)
	}

	var artifact *api.ArtifactResponse
	if !verification.Unknown() {
		artifact, _ = api.LoadArtifactForHash(user, hash, verification.MetaHash())
	}

	writeResponse(w, types.NewResult(nil, artifact, verification))
}

func index(w http.ResponseWriter, r *http.Request) {
	// can be used for healthcheck
	fmt.Fprintln(w, "OK")
}
