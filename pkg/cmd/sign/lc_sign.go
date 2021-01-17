package sign

import (
	"fmt"
	"github.com/caarlos0/spin"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/meta"
)

func LcSign(u *api.LcUser, artifacts []*api.Artifact, state meta.Status, output string) error {

	if output == "" {
		color.Set(meta.StyleAffordance())
		fmt.Println("Your assets will not be uploaded. They will be processed locally.")
		color.Unset()
		fmt.Println()
	}

	s := spin.New("%s Notarization in progress...")
	s.Set(spin.Spin1)

	var bar *progressbar.ProgressBar
	lenArtifacts := len(artifacts)
	if lenArtifacts > 1 && output == "" {
		bar = progressbar.Default(int64(lenArtifacts))
	}

	for _, a := range artifacts {
		// @todo mmeloni use verified sign
		verified, tx, err := u.Sign(
			*a,
			api.LcSignWithStatus(state),
		)
		if err != nil {
			return err
		}
		if !verified {
			color.Set(meta.StyleError())
			fmt.Println("the ledger is compromised. Please contact the CodeNotary Ledger Compliance administrators")
			color.Unset()
			fmt.Println()
			return nil
		}

		if err != nil {
			return cli.PrintWarning(output, err.Error())
		}
		if output == "" && lenArtifacts == 0 {
			fmt.Println()
		}

		artifact, verified, err := u.LoadArtifact(a.Hash, "", tx)
		if err != nil {
			return cli.PrintWarning(output, err.Error())
		}

		if !verified {
			color.Set(meta.StyleError())
			fmt.Println("the ledger is compromised. Please contact the CodeNotary Ledger Compliance administrators")
			color.Unset()
			fmt.Println()
			artifact.Status = meta.StatusUnknown
			return nil
		}

		if lenArtifacts > 1 && output == "" {
			if err := bar.Add(1); err != nil {
				return err
			}
		} else {
			cli.PrintLc(output, types.NewLcResult(artifact, verified))
		}
	}
	if lenArtifacts > 1 && output == "" {
		color.Set(meta.StyleSuccess())
		fmt.Printf("notarized %d items", lenArtifacts)
		color.Unset()
		fmt.Println()
	}
	return nil
}
