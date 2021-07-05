package bom_component

import (
	"fmt"
)

const (
	HashInvalid = iota
	HashSHA1
	HashSHA224
	HashSHA256
	HashSHA384
	HashSHA512
	HashMD2
	HashMD4
	HashMD5
	HashMD6
	hashMax = HashMD6
)

const MaxGoroutines = 8		// used by other packages that query components from external sources

var hashText = [hashMax]string{"SHA1", "SHA224", "SHA256", "SHA384", "SHA512", "MD2", "MD4", "MD5", "MD6"}

// Component is a single building block, used for building the package
type Component struct {
	Name     string
	Version  string
	Hash     string
	HashType int
}

func (c Component) HashTypeName() (string, error) {
	if c.HashType <= HashInvalid || c.HashType >= hashMax {
		return "", fmt.Errorf("invalid hash type %d", c.HashType)
	}
	return hashText[c.HashType-1], nil
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
