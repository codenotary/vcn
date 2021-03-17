/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"github.com/vchain-us/vcn/pkg/meta"
)

// SignOption is a functional option for signing operations
type LcSignOption func(*lcSignOpts) error

type lcSignOpts struct {
	status     meta.Status
	visibility meta.Visibility
	attach     []string
}

func makeLcSignOpts(opts ...LcSignOption) (o *lcSignOpts, err error) {
	o = &lcSignOpts{
		status:     meta.StatusTrusted,
		visibility: meta.VisibilityPrivate,
		attach:     nil,
	}

	for _, option := range opts {
		if option == nil {
			continue
		}
		if err := option(o); err != nil {
			return nil, err
		}
	}

	return
}

// SignWithStatus returns the functional option for the given status.
func LcSignWithStatus(status meta.Status) LcSignOption {
	return func(o *lcSignOpts) error {
		o.status = status
		return nil
	}
}

// SignWithVisibility returns the functional option for the given visibility.
func LcSignWithVisibility(visibility meta.Visibility) LcSignOption {
	return func(o *lcSignOpts) error {
		o.visibility = visibility
		return nil
	}
}

// LcSignWithAttachments returns the functional option for the given status.
func LcSignWithAttachments(attach []string) LcSignOption {
	return func(o *lcSignOpts) error {
		o.attach = attach
		return nil
	}
}
