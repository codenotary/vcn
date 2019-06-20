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

// NewCmdUnsupport returns the cobra command for `vcn unsupport`
func NewCmdUnsupport() *cobra.Command {
	cmd := NewCmdSign()
	cmd.Use = "unsupport"
	cmd.Aliases = []string{"us"}
	cmd.Short = "Unsupport a digital asset"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runSignWithState(cmd, args, meta.StatusUnsupported)
	}
	return cmd
}
