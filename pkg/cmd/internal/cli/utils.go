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
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func readPassword(msg string) (string, error) {
	fmt.Print(msg)
	password, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println(".")
	if err != nil {
		return "", err
	}
	return string(password), nil
}
