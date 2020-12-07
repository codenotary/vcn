/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"github.com/vchain-us/vcn/pkg/store"
)

// GetUserFromContext returns a new the correct user based on the context
func GetUserFromContext(context store.CurrentContext) (interface{}, error) {
	if context.Email != "" {
		return &User{
			cfg: store.Config().UserByMail(context.Email),
		}, nil
	}
	if context.LcApiKey != "" {
		client, err := NewLcClient(context.LcApiKey, context.LcHost, context.LcPort, context.LcCert)
		if err != nil {
			return nil, err
		}
		return &LcUser{
			Client: client,
			cfg:    store.Config().UserByLcApiKey(context.LcApiKey),
		}, nil
	}
	return nil, nil
}
