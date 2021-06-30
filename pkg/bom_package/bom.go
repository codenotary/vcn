package bom_package

import (

	// add language-specific packages below
	"github.com/vchain-us/vcn/pkg/bom_component"
	"github.com/vchain-us/vcn/pkg/bom_go"
)

// Package is a language- and/or environment-specific interface which finds components
type Package interface {
	Type() string
	Components() ([]bom_component.Component, error)
	Close()
}

func New(filename string) Package {
	// try all language options, return the one that matches
	if pkg := bom_go.New(filename); pkg != nil {
		return pkg
	}

	return nil
}
