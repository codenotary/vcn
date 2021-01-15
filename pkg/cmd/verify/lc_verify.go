package verify

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/meta"
)

func lcVerify(a *api.Artifact, user *api.LcUser, signerID string, output string) (err error) {

	ar, verified, err := user.LoadArtifact(a.Hash, signerID, 0)
	if err != nil {
		return cli.PrintWarning(output, err.Error())
	}

	if !verified {
		color.Set(meta.StyleError())
		fmt.Println("the ledger is compromised. Please contact the CodeNotary Ledger Compliance administrators")
		color.Unset()
		fmt.Println()
		ar.Status = meta.StatusUnknown
	}

	cli.PrintLc(output, types.NewLcResult(ar, verified))

	return
}
