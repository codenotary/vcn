/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package dashboard

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/vchain-us/vcn/pkg/meta"

	"github.com/spf13/cobra"
)

// NewCommand returns the cobra command for `vcn dashboard`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dashboard",
		Aliases: []string{"d"},
		Short:   "Open " + meta.DashboardURL() + " in browser",
		Long:    ``,
		Run:     runDashboard,
		Args:    cobra.NoArgs,
	}

	return cmd
}

func runDashboard(cmd *cobra.Command, args []string) {
	// We intentionally do not read the customer's token from disk
	// and GET the dashboard => this would be insecure as tokens would
	// be visible in server logs. In case the anyhow long-running web session
	// has expired the customer will have to log in.
	url := meta.DashboardURL()
	if output, _ := cmd.Flags().GetString("output"); output == "" {
		fmt.Printf("Taking you to <%s> \n", url)
	}
	browser.OpenURL(url)
}
