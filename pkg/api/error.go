/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import "fmt"

// ErrNotVerified is returned when an artifact is not verified on CNLC
var ErrNotVerified = fmt.Errorf("artifact is not verified")

// ErrNotVerified is returned when an artifact is not found on CNLC
var ErrNotFound = fmt.Errorf("artifact is not found")

// Error represents a CodeNotary platform's API returned error.
type Error struct {
	Description string   `json:"description"`
	Status      int      `json:"status"`
	Message     string   `json:"message"`
	Path        string   `json:"path"`
	Timestamp   string   `json:"timestamp"`
	Error       string   `json:"error"`
	FieldErrors []string `json:"fieldErrors"`
}
