/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom

import (
	"github.com/vchain-us/vcn/pkg/bom/artifact"
	"github.com/vchain-us/vcn/pkg/bom/dotnet"
	"github.com/vchain-us/vcn/pkg/bom/golang"
	"github.com/vchain-us/vcn/pkg/bom/java"
	"github.com/vchain-us/vcn/pkg/bom/python"
)

// New returns Artifact implementation of type, matching the artifact language/environment
func New(filename string) artifact.Artifact {
	// try all language options, return the one that matches
	if pkg := golang.New(filename); pkg != nil {
		return pkg
	}
	if pkg := python.New(filename); pkg != nil {
		return pkg
	}
	if a := dotnet.New(filename); a != nil {
		return a
	}
	if pkg := java.New(filename); pkg != nil {
		return pkg
	}

	return nil
}
