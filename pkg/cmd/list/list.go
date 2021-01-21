/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package list

import (
	"fmt"

	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"

	"github.com/spf13/cobra"
	"github.com/vchain-us/vcn/internal/assert"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/store"
)

// NewCommand returns the cobra command for `vcn list`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List your notarized assets",
		Long:    ``,
		RunE:    runList,
		Args:    cobra.NoArgs,
	}

	cmd.Flags().UintP("page", "p", 0, "page number")

	return cmd
}

const (
	listFooter = `
Shown %d-%d of %d assets (current page %d)

To list next page, run:
vcn list --page %d
`
)

func runList(cmd *cobra.Command, args []string) error {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	page, err := cmd.Flags().GetUint("page")
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true

	if store.Config().CurrentContext.LcApiKey != "" {
		fmt.Printf("Not supported with CodeNotary Ledger Compliance credentials\n")
		return nil
	}

	if err := assert.UserLogin(); err != nil {
		return err
	}
	u := api.NewUser(store.Config().CurrentContext.Email)

	artifacts, err := u.ListArtifacts(page)
	if err != nil {
		return err
	}
	if output == "" {
		fmt.Printf("Listing assets for %s...\n\n", u.Email())
	}
	if err = cli.PrintList(output, artifacts.Content); err != nil {
		return err
	}
	if output == "" {
		if l := uint64(len(artifacts.Content)); l > 0 {
			offset := artifacts.Pageable.PageSize * artifacts.Pageable.PageNumber
			fmt.Printf(
				"%s's assets: %d-%d of %d (current page %d)\n\n",
				u.Email(),
				1+offset,
				uint64(len(artifacts.Content))+offset,
				artifacts.TotalElements,
				artifacts.Pageable.PageNumber,
			)

			if (offset + artifacts.Pageable.PageSize) <= artifacts.TotalElements {
				fmt.Printf(
					"To list next page, run:\nvcn list --page %d\n\n",
					artifacts.Pageable.PageNumber+1,
				)
			}
		} else {
			fmt.Printf("No results.\n\n")
		}
	}
	return nil
}
