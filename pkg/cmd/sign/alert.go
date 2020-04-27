/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package sign

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vchain-us/vcn/pkg/store"

	"github.com/vchain-us/vcn/pkg/extractor/dir"
	"github.com/vchain-us/vcn/pkg/extractor/file"
	"github.com/vchain-us/vcn/pkg/extractor/git"
	"github.com/vchain-us/vcn/pkg/uri"

	"github.com/vchain-us/vcn/pkg/api"
)

type alertOptions struct {
	arg   string
	name  string
	email string
}

func handleAlert(opts *alertOptions, u api.User, a *api.Artifact, v api.BlockchainVerification, output string) error {
	if opts == nil {
		return nil
	}

	m := api.Metadata{}

	// make path absolute
	aURI, err := uri.Parse(opts.arg)
	if err != nil {
		return fmt.Errorf("invalid argument for alert: %s", opts.arg)
	}

	switch a.Kind {
	case file.Scheme:
		fallthrough
	case dir.Scheme:
		fallthrough
	case git.Scheme:
		absPath, err := filepath.Abs(strings.TrimPrefix(aURI.Opaque, "//"))
		if err != nil {
			return err
		}
		if aURI.Scheme == "" {
			aURI.Opaque = absPath
		} else {
			aURI.Opaque = "//" + absPath
		}
		opts.arg = aURI.String()
		m["path"] = absPath
	}

	hostname, _ := os.Hostname()
	if hostname != "" {
		m["hostname"] = hostname
	}

	if opts.name == "" {
		opts.name = hostname
	}

	alertConfig, err := u.CreateAlert(opts.name, opts.email, a, v, m)
	if err != nil {
		return fmt.Errorf("cannot create alert: %s", err)
	}

	if output == "" {
		fmt.Printf("\nAlert %s has been created.\n", alertConfig.AlertUUID)
	}

	if err := store.SaveAlert(u.Email(), alertConfig.AlertUUID, store.Alert{
		Name:   opts.name,
		Arg:    opts.arg,
		Config: alertConfig,
	}); err != nil {
		return fmt.Errorf("cannot save alert: %s", err)
	}

	if output == "" {
		configPath, _ := store.AlertFilepath(u.Email())
		fmt.Printf("\nThe alert configuration has been added to %s.\n", configPath)
	}

	return nil
}
