/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package verify

import (
	"fmt"
	"math/big"
	"path/filepath"
	"strings"
	"time"

	"github.com/vchain-us/vcn/pkg/store"

	"github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vchain-us/vcn/internal/cli"
	"github.com/vchain-us/vcn/internal/docker"
	"github.com/vchain-us/vcn/internal/utils"
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
		Args:    cobra.MinimumNArgs(1),
	}

	cmd.SetUsageTemplate(
		strings.Replace(cmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}} ...ARG(s)", 1),
	)

	cmd.Flags().StringP("key", "k", "", "specify the public key <vcn> should use, if not set the last available is used")

	return cmd
}

func runVerify(cmd *cobra.Command, args []string) error {
	pubKey, err := cmd.Flags().GetString("key")
	if err != nil {
		return err
	}
	user := api.NewUser(store.Config().CurrentContext)
	for _, spec := range args {
		if ok, err := verify(spec, pubKey, user); !ok {
			cmd.SilenceUsage = true
			if err != nil {
				return err
			}
			if pubKey == "" {
				return fmt.Errorf("%s is not verified", spec)
			}
			return fmt.Errorf("%s is not verified by %s", spec, pubKey)
		}
	}

	return nil
}

func verify(filename string, pubKey string, user *api.User) (success bool, err error) {
	var artifactHash string
	// fixme(leogr): refactor to spec
	if strings.HasPrefix(filename, "docker:") {
		artifactHash, err = docker.GetHash(filename)
		if err != nil {
			return false, fmt.Errorf("failed to get hash for docker image: %s", err)
		}
	} else {
		hash, err := utils.HashFile(filename)
		if err != nil {
			return false, err
		}
		artifactHash = strings.TrimSpace(hash)
	}

	var verification *api.BlockchainVerification
	if pubKey == "" {
		verification, err = api.BlockChainVerify(artifactHash)
	} else {
		verification, err = api.BlockChainVerifyMatchingPublicKey(artifactHash, pubKey)
	}
	if err != nil {
		return false, fmt.Errorf("unable to verify hash: %s", err)
	}

	var artifact *api.ArtifactResponse
	if verification.Owner != common.BigToAddress(big.NewInt(0)) {
		artifact, _ = api.LoadArtifactForHash(user, artifactHash, verification.MetaHash())
	}
	if artifact != nil {
		cli.PrintColumn("Asset", artifact.Filename, filepath.Base(filename))
		cli.PrintColumn("Hash", artifactHash, "NA")
		cli.PrintColumn("Date", verification.Timestamp.String(), "NA")
		cli.PrintColumn("Signer", artifact.Publisher, "NA")
		cli.PrintColumn("Key", strings.ToLower(verification.Owner.Hex()), "NA")
		cli.PrintColumn("Name", artifact.Name, "NA")
		if artifact.FileSize > 0 {
			cli.PrintColumn("Size", humanize.Bytes(artifact.FileSize), "NA")
		} else {
			cli.PrintColumn("Size", "NA", "NA")
		}
		cli.PrintColumn("Company", artifact.PublisherCompany, "NA")
		cli.PrintColumn("Website", artifact.PublisherWebsiteUrl, "NA")
		cli.PrintColumn("Level", meta.LevelName(verification.Level), "NA")
	} else {
		cli.PrintColumn("Asset", filepath.Base(filename), "NA")
		cli.PrintColumn("Hash", artifactHash, "NA")
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
		cli.PrintColumn("Name", "NA", "NA")
		cli.PrintColumn("Company", "NA", "NA")
		cli.PrintColumn("Website", "NA", "NA")
		cli.PrintColumn("Size", "NA", "NA")
		cli.PrintColumn("Level", "NA", "NA")
	}

	c, s := meta.StatusColor(verification.Status)
	cli.PrintColumn("Status", meta.StatusName(verification.Status), "NA", c, s)
	success = verification.Status == meta.StatusTrusted

	// todo(ameingast): redundant tracking events?
	_ = api.TrackPublisher(user, meta.VcnVerifyEvent)
	_ = api.TrackVerify(user, artifactHash, filepath.Base(filename))
	return
}
