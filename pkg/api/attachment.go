/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

// Attachment holds Attachment attributes
type Attachment struct {
	Filename string `json:"filename" yaml:"filename" vcn:"filename"`
	Hash     string `json:"hash" yaml:"hash" vcn:"hash"`
	Mime     string `json:"mime" yaml:"mime" vcn:"mime"`
}
