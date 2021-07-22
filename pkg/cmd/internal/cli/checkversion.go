/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/vchain-us/vcn/pkg/meta"

	"github.com/blang/semver"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/store"
)

// CheckVersion searches for newer vcn version and, if any, prints a message.
func CheckVersion() {
	if lastCheck := store.VersionCheckTime(); lastCheck != nil && lastCheck.After(time.Now().AddDate(0, 0, -1)) {
		return
	}
	store.SetVersionCheckTime()

	v, m, err := api.LatestCLIVersion()
	if v == "" || err != nil {
		return
	}

	curV, err := semver.Parse(strings.TrimPrefix(meta.Version(), "v"))
	if err != nil {
		return
	}
	latestV, err := semver.Parse(strings.TrimPrefix(v, "v"))
	if err != nil {
		return
	}

	if latestV.GT(curV) {

		downloadURL, err := api.LatestCLIDownloadURL()
		if err != nil {
			return
		}

		fmt.Println()
		color.Set(meta.StyleAffordance())
		fmt.Println("A newer version of vcn is available to download.")
		color.Unset()
		fmt.Printf(`		
Your version: %s
Latest version: %s
Download URL: %s
%s
`, meta.Version(), v, downloadURL, m)
	}
}
