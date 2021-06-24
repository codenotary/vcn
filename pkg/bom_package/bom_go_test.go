package bom_package

import (
	"fmt"
	"testing"
)

func TestGo(t *testing.T) {
	binary := "../../vcn"
	pkg := New(binary)
	if pkg == nil {
		t.Fatalf("Cannot process '%s' - unsupported file format/language", binary)
	}

	comp, err := pkg.Components()
	if err != nil {
		t.Fatalf("Cannot read component for '%s': %v", binary, err)
	}

	fmt.Printf("'%s' as a %s package\n", binary, pkg.Type())

	for _, c := range comp {
		fmt.Printf("%s\t%s\t%s\n", c.Name, c.Version, c.Hash)
	}
	pkg.Close()
}
