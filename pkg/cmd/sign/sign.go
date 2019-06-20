/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package sign

import (
	"fmt"
	"strings"

	"github.com/caarlos0/spin"
	"github.com/spf13/cobra"
	"github.com/vchain-us/vcn/internal/assert"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/meta"
	"github.com/vchain-us/vcn/pkg/store"
)

// NewCmdSign returns the cobra command for `vcn sign`
func NewCmdSign() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sign",
		Aliases: []string{"s"},
		Short:   "Sign digital assets' hashes onto the blockchain",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSignWithState(cmd, args, meta.StatusTrusted)
		},

		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().VarP(make(mapOpts), "attr", "a", "add user defined attributes (format: --attr key=value)")
	cmd.Flags().StringP("key", "k", "", "specify which user's key to use for signing, if not set the last available is used")
	cmd.Flags().StringP("name", "n", "", "set the asset's name")
	cmd.Flags().BoolP("public", "p", false, "when signed as public, the asset name and the signer's identity will be visible to everyone")

	cmd.SetUsageTemplate(
		strings.Replace(cmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}} ARG", 1),
	)

	return cmd
}

func runSignWithState(cmd *cobra.Command, args []string, state meta.Status) error {

	var hash string
	if hashFlag := cmd.Flags().Lookup("hash"); hashFlag != nil {
		var err error
		hash, err = cmd.Flags().GetString("hash")
		if err != nil {
			return err
		}
	}

	public, err := cmd.Flags().GetBool("public")
	if err != nil {
		return err
	}

	pubKey, err := cmd.Flags().GetString("key")
	if err != nil {
		return err
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	metadata := cmd.Flags().Lookup("attr").Value.(mapOpts).StringToInterface()

	cmd.SilenceUsage = true

	// User
	if err := assert.UserLogin(); err != nil {
		return err
	}
	u := api.NewUser(store.Config().CurrentContext)

	if err := assert.UserKeystore(); err != nil {
		return err
	}

	// Make the artifact to be signed
	var a *api.Artifact
	if hash != "" {
		// Load existing artifact
		if ar, err := u.LoadArtifact(hash); err == nil && ar != nil {
			a = ar.Artifact()
		} else {
			if err == nil {
				return fmt.Errorf("no asset found for %s", hash)
			}
			return err
		}
	} else {
		// Extract artifact from arg
		a, err = extractor.Extract(args[0])
		if err != nil {
			return err
		}
	}

	// Override the asset's name, if provided by --name
	if name != "" {
		a.Name = name
	}

	// Copy user provided custom attributes
	a.Metadata.SetValues(metadata)

	return sign(u, a, pubKey, state, meta.VisibilityForFlag(public), output)
}

func sign(u *api.User, a *api.Artifact, pubKey string, state meta.Status, visibility meta.Visibility, output string) error {

	if pubKey == "" {
		pubKey = u.DefaultKey()
	}
	if output == "" {
		fmt.Println("Signer:\t" + u.Email())
		fmt.Println("Key:\t" + pubKey)
	}
	passphrase, err := cli.ProvidePassphrase()
	if err != nil {
		return err
	}
	s := spin.New("%s Signing asset...")
	if output == "" {
		s.Set(spin.Spin1)
		s.Start()
	}
	// TODO: return and display: block #, trx #
	verification, err := u.Sign(*a, pubKey, passphrase, state, visibility)

	if output == "" {
		s.Stop()
	}
	if err != nil {
		return err
	}

	if output == "" {
		fmt.Println()
	}
	cli.Print(output, types.NewResult(a, nil, verification))
	return nil
}
