/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"io"

	"github.com/vchain-us/vcn/pkg/meta"
)

// SignOption is a functional option for signing operations
type SignOption func(*signOpts) error

type signOpts struct {
	status     meta.Status
	visibility meta.Visibility
	keyin      io.Reader
	passphrase string
}

func makeSignOpts(opts ...SignOption) (o *signOpts, err error) {
	o = &signOpts{
		status:     meta.StatusTrusted,
		visibility: meta.VisibilityPrivate,
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
func SignWithStatus(status meta.Status) SignOption {
	return func(o *signOpts) error {
		o.status = status
		return nil
	}
}

// SignWithVisibility returns the functional option for the given visibility.
func SignWithVisibility(visibility meta.Visibility) SignOption {
	return func(o *signOpts) error {
		o.visibility = visibility
		return nil
	}
}

// SignWithKey returns the functional option for the given keyin and passphrase.
func SignWithKey(keyin io.Reader, passphrase string) SignOption {
	return func(o *signOpts) error {
		o.keyin = keyin
		o.passphrase = passphrase
		return nil
	}
}
