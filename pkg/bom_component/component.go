package bom_component

// Component is a single building block, used for building the package
type Component struct {
	Name    string
	Version string
	Hash    string
}

func (c Component) IsTrusted() bool {
	// TODO query ledger
	return true
}

func (c Component) Trust() error {
	// TODO update ledger
	return nil
}

func (c Component) Untrust() error {
	// TODO update ledger
	return nil
}
