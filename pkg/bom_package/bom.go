package bom_package

import (
	component "github.com/codenotary/vcn/pkg/bom_component"

	// add language-specific packages below
	"github.com/codenotary/vcn/pkg/bom_go"
)

// Package is a language- and/or environment-specific interface which finds components
type Package interface {
	Type() string
	Components() ([]component.Component, error)
	Close()
}

func New(filename string) Package {
	// try all language options, return the one that match
	if pkg := bom_go.New(filename); pkg != nil {
		return pkg
	}

	return nil
}
