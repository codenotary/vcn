/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package info

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/vchain-us/vcn/internal/logs"
	"github.com/vchain-us/vcn/pkg/meta"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/store"

	"github.com/spf13/cobra"
)

// NewCommand returns the cobra command for `vcn info`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Display vcn information",
		Long:  ``,
		RunE:  runInfo,
		Args:  cobra.NoArgs,
	}

	return cmd
}

func runInfo(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	context := store.Config().CurrentContext
	if context.Email == "" && context.LcApiKey == "" {
		fmt.Println("\nNo user logged in.")
		return nil
	}

	uif, err := api.GetUserFromContext(context)
	if err != nil {
		return err
	}
	if lctmp, ok := uif.(*api.LcUser); ok {
		hasher := sha256.New()
		hasher.Write([]byte(lctmp.LcApiKey()))
		contextSignerID := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		fmt.Printf(`
Version:		%s
Git Rev:		%s
UserAgent:		%s
Config file:		%s
Log level:		%s
Api key:		%s
Host:			%s
Port:			%s
No-tls:			%t
Skip-verify-tls:	%t
Certificate:	%s
Current signerID:	%s
`,
			meta.Version(),
			meta.GitRevision(),
			meta.UserAgent(),
			store.ConfigFile(),
			logs.LOG.GetLevel().String(),
			lctmp.LcApiKey(),
			context.LcHost,
			context.LcPort,
			context.LcNoTls,
			context.LcSkipTlsVerify,
			context.LcCert,
			contextSignerID,
		)
		return nil
	}

	fmt.Printf(`
Version:        %s
Git Rev:        %s
UserAgent:      %s
Config file:    %s
Stage:          %s
Log level:      %s
API endpoint:   %s
MainNet:        %s
Contract Addr.: %s
`,
		meta.Version(),
		meta.GitRevision(),
		meta.UserAgent(),
		store.ConfigFile(),
		meta.StageEnvironment().String(),
		logs.LOG.GetLevel().String(),
		meta.APIEndpoint(""),
		meta.MainNet(),
		meta.AssetsRelayContractAddress(),
	)

	u := api.NewUser(context.Email)
	fmt.Printf("\nUser:		%s\n", u.Email())
	hasAuth, err := u.IsAuthenticated()
	if err != nil {
		return err
	}
	if !hasAuth {
		fmt.Println("\nUser not authenticated (token expired).")
		return nil
	}
	id, err := u.SignerID()
	if err != nil {
		return err
	}
	fmt.Printf("SignerID:	%s\n", id)
	return nil
}
