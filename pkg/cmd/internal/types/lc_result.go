/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package types

import (
	"github.com/vchain-us/vcn/pkg/api"
)

type LcResult struct {
	api.LcArtifact `yaml:",inline"`
	Verified       bool           `json:"verified" yaml:"verified" vcn:"Verified"`
	Verbose        *LcVerboseInfo `yaml:",inline" vcn:"Verbose"`
	Errors         []error        `json:"error,omitempty" yaml:"error,omitempty"`
}

type LcVerboseInfo struct {
	LedgerName string `json:"ledgerName" yaml:"ledgerName" vcn:"LedgerName"`
	LocalSID   string `json:"localSID" yaml:"localSID" vcn:"LocalSID"`
	ApiKey     string `json:"apiKey" yaml:"apiKey" vcn:"ApiKey"`
}

func (r *LcResult) AddError(err error) {
	r.Errors = append(r.Errors, err)
}

func NewLcResult(lca *api.LcArtifact, verified bool, verbose *LcVerboseInfo) *LcResult {

	var r LcResult

	switch true {
	case lca != nil:
		r = LcResult{LcArtifact: *lca, Verified: verified, Verbose: verbose}
	default:
		r = LcResult{}
	}

	// Do not show status and level from platform
	//r.Status =

	return &r
}
