/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package meta

import (
	"github.com/fatih/color"
)

func StatusColor(status Status) (color.Attribute, color.Attribute) {
	switch status {
	case StatusTrusted:
		return StyleSuccess()
	case StatusUnknown:
		return StyleWarning()
	default:
		return StyleError()
	}
}

func StyleAffordance() (color.Attribute, color.Attribute) {
	return color.FgHiBlue, color.Bold
}

func StyleError() (color.Attribute, color.Attribute) {
	return color.FgRed, color.Bold
}

func StyleWarning() (color.Attribute, color.Attribute) {
	return color.FgYellow, color.Bold
}

func StyleSuccess() (color.Attribute, color.Attribute) {
	return color.FgGreen, color.Bold
}
