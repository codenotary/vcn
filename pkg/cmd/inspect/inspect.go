/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package inspect

import (
	"fmt"
	"strings"

	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"

	"github.com/spf13/cobra"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/store"
)

// NewCommand returns the cobra command for `vcn inspect`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect",
		Aliases: []string{"i"},
		Short:   "Return the asset history with low-level information",
		Long:    ``,
		RunE:    runInspect,
		Args: func(cmd *cobra.Command, args []string) error {
			if hash, _ := cmd.Flags().GetString("hash"); hash != "" {
				if len(args) > 0 {
					return fmt.Errorf("cannot use ARG(s) with --hash")
				}
				return nil
			}
			return cobra.MinimumNArgs(1)(cmd, args)
		},
	}

	cmd.SetUsageTemplate(
		strings.Replace(cmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}} ARG", 1),
	)

	cmd.Flags().String("hash", "", "specify a hash to inspect, if set no ARG can be used")
	cmd.Flags().Bool("extract-only", false, "if set, print only locally extracted info")

	return cmd
}

func runInspect(cmd *cobra.Command, args []string) error {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	hash, err := cmd.Flags().GetString("hash")
	if err != nil {
		return err
	}
	hash = strings.ToLower(hash)

	extractOnly, err := cmd.Flags().GetBool("extract-only")
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true

	if hash == "" {
		if len(args) < 1 {
			return fmt.Errorf("no argument")
		}
		if hash, err = extractInfo(args[0], output); err != nil {
			return err
		}
		if output == "" {
			fmt.Print("\n\n")
		}
	}

	if extractOnly {
		return nil
	}

	u := api.NewUser(store.Config().CurrentContext)
	if hasAuth, _ := u.IsAuthenticated(); hasAuth && output == "" {
		fmt.Printf("Current user: %s\n", u.Email())
	}

	return inspect(hash, u, output)
}

func extractInfo(arg, output string) (hash string, err error) {
	a, err := extractor.Extract(arg)
	if err != nil {
		return "", err
	}
	if a == nil {
		return "", fmt.Errorf("unable to process the input asset provided: %s", arg)
	}

	hash = a.Hash

	if output == "" {
		fmt.Printf("Extracted info from: %s\n\n", arg)
	}
	cli.Print(output, types.NewResult(a, nil, nil))
	return
}

func inspect(hash string, u *api.User, output string) error {
	verifications, err := api.BlockChainInspect(hash)
	if err != nil {
		return err
	}
	l := len(verifications)

	if output == "" {
		fmt.Printf(
			`%d notarizations found for "%s"

`,
			l, hash)
	}

	results := make([]types.Result, l)
	for i, v := range verifications {
		ar, err := api.LoadArtifact(u, hash, v.MetaHash())
		results[i] = *types.NewResult(nil, ar, &v)
		if err != nil {
			results[i].AddError(err)
		}
		// check if artifact is synced, if any
		if ar != nil {
			if v.Status.String() != ar.Status {
				results[i].AddError(fmt.Errorf(
					"status not in sync (blockchain: %s, platform: %s)", v.Status.String(), ar.Status,
				))
			}
			if int64(v.Level) != ar.Level {
				results[i].AddError(fmt.Errorf(
					"level not in sync (blockchain: %d, platform: %d)", v.Level, ar.Level,
				))
			}
		}
	}

	return cli.PrintSlice(output, results)
}
