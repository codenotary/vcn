module github.com/codenotary/vcn/pkg/bom_package

go 1.16

replace github.com/codenotary/vcn/pkg/bom_go => ../bom_go

replace github.com/codenotary/vcn/pkg/bom_python => ../bom_python

replace github.com/codenotary/vcn/pkg/bom_component => ../bom_component

require (
	github.com/codenotary/vcn/pkg/bom_component v0.0.0-00010101000000-000000000000
	github.com/codenotary/vcn/pkg/bom_go v0.0.0-00010101000000-000000000000
	github.com/codenotary/vcn/pkg/bom_python v0.0.0-00010101000000-000000000000
)
