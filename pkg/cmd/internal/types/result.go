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

type Result struct {
	api.ArtifactResponse `yaml:",inline"`
	Verification         *api.BlockchainVerification `json:"verification" yaml:"verification"`
	Errors               []error                     `json:"error,omitempty" yaml:"error,omitempty"`
}

func (r *Result) AddError(err error) {
	r.Errors = append(r.Errors, err)
}

func NewResult(a *api.Artifact, ar *api.ArtifactResponse, v *api.BlockchainVerification) *Result {

	var vv *api.BlockchainVerification
	if v != nil {
		vCopy := *v
		vv = &vCopy
	}

	var r Result

	switch {
	case ar != nil:
		r = Result{*ar, vv, nil}
	case a != nil:
		r = Result{api.ArtifactResponse{
			Name:        a.Name,
			Kind:        a.Kind,
			Hash:        a.Hash,
			Size:        a.Size,
			ContentType: a.ContentType,
			Metadata:    a.Metadata,
		}, vv, nil}
	default:
		r = Result{}
		r.Verification = vv
	}

	// Do not show status and level from platform
	r.Status = ""
	r.Level = 0

	return &r
}
