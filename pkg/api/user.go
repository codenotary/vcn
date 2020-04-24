/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/sirupsen/logrus"
	"github.com/vchain-us/vcn/internal/errors"
	"github.com/vchain-us/vcn/pkg/meta"
	"github.com/vchain-us/vcn/pkg/store"
)

// User represent a CodeNotary platform user.
type User struct {
	cfg *store.User
}

// NewUser returns a new User instance for the given email.
func NewUser(email string) *User {
	return &User{
		cfg: store.Config().User(email),
	}
}

// Email returns the User's email, if any, otherwise an empty string.
func (u User) Email() string {
	if u.cfg != nil {
		return u.cfg.Email
	}
	return ""
}

// Authenticate the User against the CodeNotary platform.
// If successful the auth token in stored within the User's config and used for subsequent API call.
func (u *User) Authenticate(password string, otp string) (err error) {
	if u == nil || u.Email() == "" {
		return makeFatal("user not initialized", nil)
	}

	email := u.Email()

	ok, err := checkUserExists(email)
	if err != nil {
		return
	}
	if !ok {
		return fmt.Errorf(`no such user "%s", please create an account at %s`, email, meta.DashboardURL())
	}

	token, err := authenticateUser(email, password, otp)
	if err != nil {
		return
	}

	u.cfg.Token = token
	return nil
}

// ClearAuth deletes the stored authentication token.
func (u *User) ClearAuth() {
	if u != nil && u.cfg != nil {
		u.cfg.Token = ""
	}
}

func (u *User) token() string {
	if u != nil && u.cfg != nil {
		return u.cfg.Token
	}
	return ""
}

// IsAuthenticated returns true if the stored auth token is still valid.
func (u User) IsAuthenticated() (bool, error) {
	if u.cfg == nil || u.cfg.Token == "" {
		return false, nil
	}

	return checkToken(u.cfg.Token)
}

// IsExist returns true if the User's was registered on the CodeNotary platform.
func (u User) IsExist() (bool, error) {
	email := u.Email()
	if email != "" {
		return checkUserExists(email)
	}
	return false, nil
}

// Config returns the User configuration object (see store.User), if any.
// It returns nil if the User is not properly initialized.
func (u User) Config() *store.User {
	if u.cfg != nil {
		return u.cfg
	}
	return nil
}

func (u User) getWallet() (address, keystore string, offline bool, err error) {
	authError := new(Error)
	pagedWalletResponse := new(struct {
		Content []struct {
			Address             string `json:"address"`
			KeyStore            string `json:"keyStore"`
			CreatedAt           string `json:"createdAt"`
			Name                string `json:"name"`
			PermissionSyncState string `json:"permissionSyncState"`
			LevelSyncState      string `json:"levelSyncState"`
		} `json:"content"`
	})
	r, err := newSling(u.token()).
		Get(meta.APIEndpoint("wallet")).
		Receive(pagedWalletResponse, authError)
	logger().WithFields(logrus.Fields{
		"response":  "HIDDEN",
		"err":       err,
		"authError": authError,
	}).Trace("getWallet")
	if err != nil {
		return
	}
	if r.StatusCode != 200 {
		err = fmt.Errorf("request failed: %s (%d)", authError.Message, authError.Status)
		return
	}

	wallets := pagedWalletResponse.Content
	for _, wallet := range wallets {
		if wallet.Address == strings.ToLower(address) {
			if wallet.PermissionSyncState == "SYNCED" && wallet.LevelSyncState == "SYNCED" {
				// everything is ok
				break // currently, the user can have just one wallet
			}
			err = fmt.Errorf(errors.AccountNotSynced)
			return
		}
	}

	// currently, the user can have just one wallet
	if len(wallets) > 0 {
		address = wallets[0].Address
		keystore = wallets[0].KeyStore
	}
	return
}

// Secret fetches the User's secret and returns an io.Reader for reading it.
func (u User) Secret() (reader io.Reader, id string, offline bool, err error) {
	id, keystore, offline, err := u.getWallet()

	switch true {
	case err != nil:
		// passthru error
	case id == "":
		err = fmt.Errorf("no secret found for %s", u.Email())
	case offline:
		// passthru offline
	case keystore == "":
		err = fmt.Errorf("no secret found for %s", u.Email())
	default:
		reader = bytes.NewReader([]byte(keystore))
	}

	return
}

// SignerID retrives the User's SignerID (the public address derived from the secret) from the platform.
func (u User) SignerID() (id string, err error) {
	id, _, _, err = u.getWallet()
	if err == nil && id == "" {
		err = fmt.Errorf("no SignerID found for %s", u.Email())
	}
	return
}

// UploadSecret uploads the User's secret to the platform.
func (u User) UploadSecret(secret io.Reader, passphrase string) (err error) {
	var buf bytes.Buffer
	tee := io.TeeReader(secret, &buf)
	transactor, err := bind.NewTransactor(tee, passphrase)
	if err != nil {
		return
	}
	address := strings.ToLower(transactor.From.Hex())

	b, err := ioutil.ReadAll(&buf)
	if err != nil {
		return err
	}
	content := &struct {
		KeyStore string `json:"keyStore"`
	}{string(b)}
	b, err = json.Marshal(content)
	if err != nil {
		return
	}

	authError := new(Error)
	r, err := newSling(u.token()).
		Patch(meta.APIEndpoint("wallet")+"?wallet-address="+address).
		Body(bytes.NewReader(b)).
		Add("Content-Type", "application/json").
		Receive(nil, authError)
	logger().WithFields(logrus.Fields{
		"err":       err,
		"authError": authError,
	}).Trace("UploadSecret")

	if err != nil {
		return
	}
	if r.StatusCode != 200 {
		err = fmt.Errorf("request failed: %s (%d)", authError.Message, authError.Status)
		return
	}
	return
}

// RemainingSignOps returns the number of remaining notarizations in the User's account subscription.
func (u User) RemainingSignOps() (uint64, error) {
	response := new(struct {
		Count uint64 `json:"count"`
	})
	restError := new(Error)
	r, err := newSling(u.token()).
		Get(meta.APIEndpoint("artifact/remaining-sign-operations")).
		Receive(&response, restError)
	logger().WithFields(logrus.Fields{
		"response":  response,
		"err":       err,
		"restError": restError,
	}).Trace("RemainingSignOps")
	if err != nil {
		return 0, err
	}
	switch r.StatusCode {
	case 200:
		return response.Count, nil
	}
	return 0, fmt.Errorf("count remaining sign operations failed: %+v", restError)
}

func (u User) trialExpired() (bool, error) {
	response := new(struct {
		TrialExpired bool `json:"trialExpired"`
	})
	restError := new(Error)
	r, err := newSling(u.token()).
		Get(publisherEndpoint()).
		Receive(&response, restError)
	logger().WithFields(logrus.Fields{
		"response":  response,
		"err":       err,
		"restError": restError,
	}).Trace("trialExpired")
	if err != nil {
		return false, err
	}
	switch r.StatusCode {
	case 200:
		return response.TrialExpired, nil
	}
	return false, fmt.Errorf("cannot get user info from platform: %+v", restError)
}
