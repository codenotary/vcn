package verify

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/meta"
	"strconv"
)

func lcVerify(cmd *cobra.Command, a *api.Artifact, user *api.LcUser, signerID string, uid string, attach string, lcAttachForce bool, verbose bool, output string) (err error) {
	hook := newHook(cmd, a)
	err = hook.lcFinalizeWithoutAlert(user, output, 0)
	if err != nil {
		return err
	}
	var ar *api.LcArtifact
	var verified bool
	var attachmentList []api.Attachment

	var attachmentMap = make(map[string][]*api.Attachment)
	if attach != "" {
		attachmentMap, err = user.GetArtifactUIDAndAttachmentsListByAttachmentLabel(a.Hash, signerID, attach)
		if err != nil {
			return err
		}
		// attachmentFileNameMap is used internally to produce a map to handle equal file name attachments
		attachmentFileNameMap := make(map[string][]*api.Attachment)
		for k, attachMapEntry := range attachmentMap {
			// latest uid, needed to authenticate the latest notarized artifact
			if uid == "" {
				uid = k
			}
			for _, att := range attachMapEntry {
				fn := att.Filename
				if _, ok := attachmentFileNameMap[fn]; ok && len(attachmentFileNameMap[fn]) > 0 {
					// if there is a newer filename here a postfix is added. ~1,~2 ... ~N
					att.Filename = fn + "~" + strconv.Itoa(len(attachmentFileNameMap[fn]))
				}
				attachmentFileNameMap[fn] = append(attachmentFileNameMap[fn], att)
				// attachmentList contains all attachments with latest first order
				attachmentList = append(attachmentList, *att)
			}
		}
	}

	ar, verified, err = user.LoadArtifact(a.Hash, signerID, uid, 0)

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
		color.Set(meta.StyleAffordance())
		attachDownloadMessage := "downloading attachments ..."
		if attach != "" {
			attachDownloadMessage = "downloading all notarizations attachments associated to the provided label ..."
		}
		fmt.Println(attachDownloadMessage)
		color.Unset()

		fmt.Printf("Download list:\n")
		for _, a := range attachmentList {
			fmt.Printf("\t\t- Filename:\t%s\n", a.Filename)
			fmt.Printf("\t\t  Hash:\t\t%s\n", a.Hash)
			fmt.Printf("\t\t  Mime:\t\t%s\n", a.Mime)
		}
		fmt.Printf("\n")
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
