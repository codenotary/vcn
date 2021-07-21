/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom_package

import (
	"github.com/vchain-us/vcn/pkg/bom_component"
	"github.com/vchain-us/vcn/pkg/bom_java"
	"github.com/vchain-us/vcn/pkg/bom_node"

	// add language-specific packages below
	"github.com/vchain-us/vcn/pkg/bom_dotnet"
	"github.com/vchain-us/vcn/pkg/bom_go"
	"github.com/vchain-us/vcn/pkg/bom_python"
)

// Package is a language- and/or environment-specific interface which finds components
type Package interface {
	Path() string
	Type() string
	Components() ([]bom_component.Component, error)
	Close()
}

func New(filename string) Package {
	// try all language options, return the one that matches
	if pkg := bom_go.New(filename); pkg != nil {
		return pkg
	}
	if pkg := bom_python.New(filename); pkg != nil {
		return pkg
	}
	if pkg := bom_dotnet.New(filename); pkg != nil {
		return pkg
	}
	if pkg := bom_java.New(filename); pkg != nil {
		return pkg
	}
	if pkg := bom_node.New(filename); pkg != nil {
		return pkg
	}
	return nil
}
