/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package cli

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintColumn(field string, value string, fallback string, p ...color.Attribute) {
	var spaces string
	for i := len(field); i < 8; i++ {
		spaces += " "
	}
	fmt.Print(field + ":" + spaces)
	if p != nil {
		c := color.New(p...)
		c.Set()
	}
	if value != "" {
		fmt.Print(value)
	} else {
		fmt.Print(fallback)
	}
	if p != nil {
		color.Unset()
	}
	fmt.Println()
}
