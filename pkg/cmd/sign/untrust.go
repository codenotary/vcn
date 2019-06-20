/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package sign

import (
	"github.com/spf13/cobra"
	"github.com/vchain-us/vcn/pkg/meta"
)

// NewCmdUntrust returns the cobra command for `vcn untrust`
func NewCmdUntrust() *cobra.Command {
	cmd := NewCmdSign()
	cmd.Use = "untrust"
	cmd.Aliases = []string{"ut"}
	cmd.Short = "Untrust a digital asset"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runSignWithState(cmd, args, meta.StatusUntrusted)
	}
	return cmd
}
