package verify

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/meta"
)

func lcVerify(cmd *cobra.Command, a *api.Artifact, user *api.LcUser, signerID string, uid string, attach string, lcAttachForce bool, verbose bool, output string) (err error) {
	hook := newHook(cmd, a)
	err = hook.lcFinalizeWithoutAlert(user, output, 0)
	if err != nil {
		return err
	}
	var attachmentList []api.Attachment

	if attach != "" {
		attachmentList, uid, err = user.GetArtifactAttachmentListByLabel(a.Hash, signerID, attach)
		if err != nil {
			return err
		}
	}

	ar, verified, err := user.LoadArtifact(
		a.Hash,
		signerID,
		uid,
		0,
		map[string][]string{meta.VcnLCCmdHeaderName: {meta.VcnLCVerifyCmdHeaderValue}})
	if err != nil {
		if err == api.ErrNotFound {
			err = fmt.Errorf("%s was not notarized", a.Hash)
			viper.Set("exit-code", strconv.Itoa(meta.StatusUnknown.Int()))
		}
		if err == api.ErrNotVerified {
			color.Set(meta.StyleError())
			fmt.Println("the ledger is compromised. Please contact the CodeNotary Immutable Ledger administrators")
			color.Unset()
			fmt.Println()
			viper.Set("exit-code", strconv.Itoa(meta.StatusUnknown.Int()))
		}
		return cli.PrintWarning(output, err.Error())
	}
	if ar.Revoked != nil && !ar.Revoked.IsZero() {
		viper.Set("exit-code", strconv.Itoa(meta.StatusApikeyRevoked.Int()))
		ar.Status = meta.StatusApikeyRevoked
	}

	if len(attachmentList) == 0 && ar.Attachments != nil {
		attachmentList = ar.Attachments
	}

	if output == "attachments" {
		cli.PrintAttachmentList(attach, attachmentList)

		var bar *progressbar.ProgressBar
		lenAttachments := len(attachmentList)
		if lenAttachments >= 1 {
			bar = progressbar.Default(int64(lenAttachments))
		}

		for _, a := range attachmentList {
			_ = bar.Add(1)
			err := user.DownloadAttachment(&a, ar, 0, lcAttachForce)
			if err != nil {
				return err
			}
		}
		fmt.Println()
	}
	if !verified {
		color.Set(meta.StyleError())
		fmt.Println("the ledger is compromised. Please contact the CodeNotary Immutable Ledger administrators")
		color.Unset()
		fmt.Println()
		viper.Set("exit-code", strconv.Itoa(meta.StatusUnknown.Int()))
		ar.Status = meta.StatusUnknown
	}

	exitCode, err := cmd.Flags().GetInt("exit-code")
	if err != nil {
		return err
	}
	// if exitCode == VcnDefaultExitCode user didn't specify to use a custom exit code in case of success.
	// In that case we return the ar.Status as exit code.
	// User defined exit code is returned only if the viper exit-code status is == 0 (status trusted)
	if exitCode == meta.VcnDefaultExitCode && viper.GetInt("exit-code") == 0 {
		viper.Set("exit-code", strconv.Itoa(ar.Status.Int()))
	}
	var verbInfos *types.LcVerboseInfo
	if verbose {
		verbInfos = &types.LcVerboseInfo{
			LedgerName: ar.Ledger,
			LocalSID:   api.GetSignerIDByApiKey(user.Client.ApiKey),
			ApiKey:     user.Client.ApiKey,
		}
	}
	cli.PrintLc(output, types.NewLcResult(ar, verified, verbInfos))

	return
}
