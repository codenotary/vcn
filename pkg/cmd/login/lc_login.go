/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/fatih/color"
	vcnerr "github.com/vchain-us/vcn/internal/errors"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/meta"
	"github.com/vchain-us/vcn/pkg/store"
	"google.golang.org/grpc/metadata"
)

// Execute the login action
func ExecuteLC(host, port, lcCert, lcApiKey, lcLedger string, skipTlsVerify, lcNoTls bool) error {
	if store.CNioContext() == true {
		return errors.New("Already logged on CodeNotary.io. Please logout first.")
	}

	color.Set(meta.StyleAffordance())
	fmt.Println("Logging into CodeNotary Ledger Compliance.")
	color.Unset()

	if lcApiKey != "" {
		u, err := api.NewLcUser(lcApiKey, lcLedger, host, port, lcCert, skipTlsVerify, lcNoTls)
		if err != nil {
			return err
		}
		if u != nil {
			err = u.Client.Connect()
			if err != nil {
				return err
			}
			md := metadata.Pairs(meta.VcnLCPluginTypeHeaderName, meta.VcnLCPluginTypeHeaderValue)
			ctx := metadata.NewOutgoingContext(context.Background(), md)
			_, err = u.Client.Health(ctx)
			if err != nil {
				return err
			}
			// Store the new config
			if err := store.SaveConfig(); err != nil {
				return err
			}
		}
	}
	if lcApiKey == "" {
		return vcnerr.ErrNoLcApiKeyEnv
	}
	// shouldn't happen
	return nil
}
