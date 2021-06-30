package bom

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/codenotary/vcn/pkg/bom_package"
)

// NewCommand returns the cobra command for `vcn info`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bom",
		Example: "  vcn bom /bin/vcn",
		Short: "Display BoM information",
		Long:  ``,
		RunE:  runBom,
		Args:  func(cmd *cobra.Command, args []string) error {
			return cobra.ExactValidArgs(1)(cmd, args)
		},
	}

	return cmd
}

func runBom(cmd *cobra.Command, args []string) error {
	pkg := bom_package.New(args[0])

	if pkg == nil {
		return fmt.Errorf("unsupported package format/language")
	}

	comp, err := pkg.Components()
	if err != nil {
		return fmt.Errorf("cannot extract components: %w", err)
	}

	for _, c := range comp {
		fmt.Printf("%s\t%s\t%s\n", c.Name, c.Version, c.Hash)
	}

	return nil
}
