module github.com/codenotary/vcn/pkg/cmd/bom

go 1.16

replace github.com/codenotary/vcn/pkg/bom_component => ../../bom_component

replace github.com/codenotary/vcn/pkg/bom_package => ../../bom_package

replace github.com/codenotary/vcn/pkg/bom_go => ../../bom_go

replace github.com/codenotary/vcn/pkg/bom_python => ../../bom_python

require (
	github.com/codenotary/vcn/pkg/bom_package v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.1.3
)
