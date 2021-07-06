/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"

	sdk "github.com/vchain-us/ledger-compliance-go/grpcclient"
	"github.com/vchain-us/vcn/pkg/meta"
	"github.com/vchain-us/vcn/pkg/store"
)

// User represent a CodeNotary platform user.
type LcUser struct {
	Client *sdk.LcClient
}

// NewLcUser returns a new User instance configured with provided parameters.
// LcLedger parameter is used when a cross-ledger key is provided in order to specify the ledger on which future operations will be directed. Empty string is accepted
func NewLcUser(lcApiKey, lcLedger, host, port, lcCert string, skipTlsVerify bool, noTls bool) (*LcUser, error) {
	client, err := NewLcClient(lcApiKey, lcLedger, host, port, lcCert, skipTlsVerify, noTls)
	if err != nil {
		return nil, err
	}
	store.Config().NewLcUser(host, port, lcCert, skipTlsVerify, noTls)

	return &LcUser{
		Client: client,
	}, nil
}

// NewLcUserVolatile returns a new User instance without a backing cfg file.
func NewLcUserVolatile(lcApiKey, lcLedger string, host string, port string) *LcUser {
	p, _ := strconv.Atoi(port)
	return &LcUser{
		Client: sdk.NewLcClient(
			sdk.ApiKey(lcApiKey),
			sdk.MetadataPairs([]string{
				meta.VcnLCLedgerHeaderName, lcLedger,
				meta.VcnLCVersionHeaderName, meta.Version(),
			}),
			sdk.Host(host),
			sdk.Port(p),
			sdk.Dir(store.CurrentConfigFilePath())),
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

func GetSignerIDByApiKey(lcApiKey string) string {
	ris := strings.Split(lcApiKey, ".")
	// new apikey format {friendlySignerID}.{secret}
	if len(ris) > 1 {
		return strings.Join(ris[:len(ris)-1], ".")
	}
	// old apikey format {secret}
	hash := sha256.Sum256([]byte(lcApiKey))
	return base64.URLEncoding.EncodeToString(hash[:])
}
