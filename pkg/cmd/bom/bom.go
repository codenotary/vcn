/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vchain-us/vcn/pkg/bom_package"
)

// NewCommand returns the cobra command for `vcn info`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bom",
		Example: "  vcn bom /bin/vcn",
		Short:   "Display BoM information",
		Long:    ``,
		RunE:    runBom,
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.ExactValidArgs(1)(cmd, args)
		},
	}

	return cmd
}

func runBom(cmd *cobra.Command, args []string) error {
	pkg := bom_package.New(args[0])

	if pkg == nil {
		return fmt.Errorf("unsupported package format/language")
	}

	spdx, err := bom_package.OutputSpdxText(pkg)
	if err != nil {
		return err
	}
	fmt.Printf(spdx)

	return nil
}
