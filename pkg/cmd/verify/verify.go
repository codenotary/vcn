/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package verify

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/vchain-us/vcn/pkg/extractor"

	"github.com/vchain-us/vcn/pkg/store"

	"github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vchain-us/vcn/internal/cli"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/meta"

	"github.com/spf13/cobra"
)

// NewCmdVerify returns the cobra command for `vcn verify`
func NewCmdVerify() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "verify",
		Example: "  vcn verify /bin/vcn",
		Aliases: []string{"v"},
		Short:   "Verify digital artifact against blockchain",
		Long:    ``,
		RunE:    runVerify,
		Args: func(cmd *cobra.Command, args []string) error {
			if hash, _ := cmd.Flags().GetString("hash"); hash != "" {
				if len(args) > 0 {
					return fmt.Errorf("cannot use arg(s) with --hash")
				}
				return nil
			}
			return cobra.MinimumNArgs(1)(cmd, args)
		},
	}

	cmd.SetUsageTemplate(
		strings.Replace(cmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}} ...ARG(s)", 1),
	)

	cmd.Flags().StringP("key", "k", "", "specify the public key <vcn> should use, if not set the last available is used")
	cmd.Flags().String("hash", "", "specify a hash to verify, if set no arg(s) can be used")

	return cmd
}

func runVerify(cmd *cobra.Command, args []string) error {
	hash, err := cmd.Flags().GetString("hash")
	if err != nil {
		return err
	}
	pubKey, err := cmd.Flags().GetString("key")
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true

	user := api.NewUser(store.Config().CurrentContext)

	// by hash
	if hash != "" {
		a := &api.Artifact{
			Hash: hash,
		}
		if err := verify(a, pubKey, user); err != nil {
			return err
		}
		return nil
	}

	// else by args
	for _, arg := range args {
		a, err := extractor.Extract(arg)
		if err != nil {
			return err
		}
		if err := verify(a, pubKey, user); err != nil {
			return err
		}
	}

	return nil
}

func verify(a *api.Artifact, pubKey string, user *api.User) (err error) {
	var verification *api.BlockchainVerification
	if pubKey != "" {
		verification, err = api.BlockChainVerifyMatchingPublicKey(a.Hash, pubKey)
	} else if pubKeys := user.Keys(); len(pubKeys) > 0 {
		verification, err = api.BlockChainVerifyMatchingPublicKeys(a.Hash, pubKeys)
	} else {
		verification, err = api.BlockChainVerify(a.Hash)
	}
	if err != nil {
		return fmt.Errorf("unable to verify hash: %s", err)
	}

	var artifact *api.ArtifactResponse
	if verification.Owner != common.BigToAddress(big.NewInt(0)) {
		artifact, _ = api.LoadArtifactForHash(user, a.Hash, verification.MetaHash())
	}
	if artifact != nil {
		cli.PrintColumn("Asset", artifact.Filename, a.Name)
		cli.PrintColumn("Kind", artifact.Kind, "NA")
		cli.PrintColumn("Hash", a.Hash, "NA")
		cli.PrintColumn("Date", verification.Timestamp.String(), "NA")
		cli.PrintColumn("Signer", artifact.Publisher, "NA")
		cli.PrintColumn("Key", strings.ToLower(verification.Owner.Hex()), "NA")
		cli.PrintColumn("Name", artifact.Name, "NA")
		if artifact.FileSize > 0 {
			cli.PrintColumn("Size", humanize.Bytes(artifact.FileSize), "NA")
		} else {
			cli.PrintColumn("Size", "NA", "NA")
		}
		cli.PrintColumn("MimeType", artifact.MimeType, "NA")
		cli.PrintColumn("Platform", artifact.Platform, "NA")
		cli.PrintColumn("Arch", artifact.Architecture, "NA")
		cli.PrintColumn("Url", artifact.Url, "NA")
		cli.PrintColumn("License", artifact.License, "NA")
		metadata := ""
		for k, v := range artifact.Metadata {
			if vv, err := json.Marshal(v); err == nil {
				metadata += fmt.Sprintf("%s=%s\t", k, string(vv))
			}
		}
		cli.PrintColumn("Metadata", metadata, "NA")
		cli.PrintColumn("Company", artifact.PublisherCompany, "NA")
		cli.PrintColumn("Website", artifact.PublisherWebsiteUrl, "NA")
		cli.PrintColumn("Level", meta.LevelName(verification.Level), "NA")
	} else {
		cli.PrintColumn("Asset", a.Name, "NA")
		cli.PrintColumn("Kind", a.Kind, "NA")
		cli.PrintColumn("Hash", a.Hash, "NA")
		if verification.Timestamp != time.Unix(0, 0) {
			cli.PrintColumn("Date", verification.Timestamp.String(), "NA")
		} else {
			cli.PrintColumn("Date", "NA", "NA")
		}
		cli.PrintColumn("Signer", "NA", "NA")
		if verification.Owner != common.BigToAddress(big.NewInt(0)) {
			cli.PrintColumn("Key", strings.ToLower(verification.Owner.Hex()), "NA")
		} else {
			cli.PrintColumn("Key", "NA", "NA")
		}
	}

	c, s := meta.StatusColor(verification.Status)
	cli.PrintColumn("Status", meta.StatusName(verification.Status), "NA", c, s)

	// todo(ameingast): redundant tracking events?
	_ = api.TrackPublisher(user, meta.VcnVerifyEvent)
	_ = api.TrackVerify(user, a.Hash, a.Name)

	if verification.Status != meta.StatusTrusted {
		if pubKey != "" {
			err = fmt.Errorf("%s is not verified by %s", a.Hash, pubKey)
		} else if email := user.Email(); email != "" {
			err = fmt.Errorf("%s is not verified by %s", a.Hash, email)
		} else {
			err = fmt.Errorf("%s is not verified", a.Hash)
		}
	}

	return
}
