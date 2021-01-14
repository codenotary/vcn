/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	sdk "github.com/vchain-us/ledger-compliance-go/grpcclient"
	"github.com/vchain-us/vcn/pkg/store"
	"strconv"
)

// User represent a CodeNotary platform user.
type LcUser struct {
	Client *sdk.LcClient
	cfg    *store.User
}

// NewUser returns a new User instance for the given email.
func NewLcUser(lcApiKey, host, port, lcCert string, skipTlsVerify bool) (*LcUser, error) {
	client, err := NewLcClient(lcApiKey, host, port, lcCert, skipTlsVerify)
	if err != nil {
		return nil, err
	}
	return &LcUser{
		Client: client,
		cfg:    store.Config().NewLcUser(lcApiKey, host, port, lcCert, skipTlsVerify),
	}, nil
}

// NewLcUserVolatile returns a new User instance without a backing cfg file.
func NewLcUserVolatile(lcApiKey string, host string, port string) *LcUser {
	p, _ := strconv.Atoi(port)
	return &LcUser{
		Client: sdk.NewLcClient(sdk.ApiKey(lcApiKey), sdk.Host(host), sdk.Port(p), sdk.Dir(store.CurrentConfigFilePath())),
		cfg:    &store.User{LcApiKey: lcApiKey},
	}
}

// Email returns the User's email, if any, otherwise an empty string.
func (u LcUser) LcApiKey() string {
	if u.cfg != nil {
		return u.cfg.LcApiKey
	}
	return ""
}

// ClearAuth deletes the stored authentication token.
func (u *LcUser) ClearAuth() {
	if u != nil && u.cfg != nil {
		u.cfg.LcApiKey = ""
	}
}

// Config returns the User configuration object (see store.User), if any.
// It returns nil if the User is not properly initialized.
func (u User) User() *store.User {
	if u.cfg != nil {
		return u.cfg
	}
	return nil
}
