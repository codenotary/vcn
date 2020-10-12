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
	"github.com/vchain-us/vcn/pkg/meta"
	"github.com/vchain-us/vcn/pkg/store"
	"strconv"
)

// GetUserFromContext returns a new the correct user based on the context
func GetUserFromContext(context store.CurrentContext) (interface{}, error) {
	if context.Email != "" {
		return &User{
			cfg: store.Config().UserByMail(context.Email),
		}, nil
	}
	if context.LcApiKey != "" {
		p, err := strconv.Atoi(context.LcPort)
		if err != nil {
			return nil, err
		}
		return &LcUser{
			Client: sdk.NewLcClient(sdk.Host(context.LcHost), sdk.Port(p), sdk.ApiKey(context.LcApiKey), sdk.Dir(store.CurrentConfigFilePath()), sdk.PluginPrefix(meta.VcnLCPrefix)),
			cfg:    store.Config().UserByLcApiKey(context.LcApiKey),
		}, nil
	}
	return nil, nil
}
