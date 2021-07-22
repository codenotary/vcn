/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/bom"
	"github.com/vchain-us/vcn/pkg/bom/artifact"
	"github.com/vchain-us/vcn/pkg/meta"
	"github.com/vchain-us/vcn/pkg/store"
)

// NewCommand returns the cobra command for `vcn info`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bom",
		Example: "  vcn bom /bin/vcn",
		Short:   "Collect BoM information",
		Long: `
Collect BoM (Bill of Material) information

It identifies dependencies of build artifact and produces the BoM, which can be
later notorized using 'vcn notorize' command with '--bom' option.
It can optionally notorize unsupported dependencies, if '--nd' option is specified.
`,
		RunE: runBom,
		PreRun: func(cmd *cobra.Command, args []string) {
			// Bind to all flags to env vars (after flags were parsed),
			// but only ones retrivied by using viper will be used.
			viper.BindPFlags(cmd.Flags())
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.ExactValidArgs(1)(cmd, args)
		},
	}

	cmd.Flags().String("signerID", "", "use specified ID for authentication and notarizing dependencies")
	cmd.Flags().String("lc-host", "", meta.VcnLcHostFlagDesc)
	cmd.Flags().String("lc-port", "443", meta.VcnLcPortFlagDesc)
	cmd.Flags().String("lc-cert", "", meta.VcnLcCertPathDesc)
	cmd.Flags().Bool("lc-skip-tls-verify", false, meta.VcnLcSkipTlsVerifyDesc)
	cmd.Flags().Bool("lc-no-tls", false, meta.VcnLcNoTlsDesc)
	cmd.Flags().String("lc-api-key", "", meta.VcnLcApiKeyDesc)
	cmd.Flags().String("lc-ledger", "", meta.VcnLcLedgerDesc)
	cmd.Flags().Bool("nd", false, "enables automatic notarization of unsupported dependencies")
	cmd.Flags().Int("max-unsupported", 0, "max number (in %) of unsupported dependencies allowed")
	cmd.Flags().String("spdx", "", "name of the file to output BoM in SPDX format")
	cmd.Flags().Int("trust-level", int(artifact.Trusted), "min trust level: 0 (untrusted) / 1 (unsupported) / 2 (trusted)")

	return cmd
}

func runBom(cmd *cobra.Command, args []string) error {

	lcHost := viper.GetString("lc-host")
	lcPort := viper.GetString("lc-port")
	lcCert := viper.GetString("lc-cert")
	skipTlsVerify := viper.GetBool("lc-skip-tls-verify")
	noTls := viper.GetBool("lc-no-tls")
	lcApiKey := viper.GetString("lc-api-key")
	lcLedger := viper.GetString("lc-ledger")
	autoNotarize := viper.GetBool("nd")
	threshold := viper.GetInt("max-unsupported")
	signerID := viper.GetString("signerID")
	if threshold < 0 || threshold > 100 {
		return errors.New("max-unsupported must be a valid percentage value (0-100)")
	}
	spdxFilename := viper.GetString("spdx")
	trustLevel := viper.GetInt("trust-level")
	if trustLevel < int(artifact.MinTrustLevel) || trustLevel > int(artifact.MaxTrustLevel) {
		return fmt.Errorf("invalid trust level, expected values %d-%d", artifact.MinTrustLevel, artifact.MaxTrustLevel)
	}

	cmd.SilenceUsage = true

	buildArtifact := bom.New(args[0])
	if buildArtifact == nil {
		return fmt.Errorf("unsupported artifact format/language")
	}

	//check if an lcUser is present inside the context
	var lcUser *api.LcUser
	uif, err := api.GetUserFromContext(store.Config().CurrentContext, lcApiKey, lcLedger)
	if err != nil {
		return err
	}
	if lctmp, ok := uif.(*api.LcUser); ok {
		lcUser = lctmp
	}

	// use credentials if at least ledger compliance host is provided
	if lcHost != "" && lcApiKey != "" {
		lcUser, err = api.NewLcUser(lcApiKey, "", lcHost, lcPort, lcCert, skipTlsVerify, noTls)
		if err != nil {
			return err
		} // Store the new config
		if err := store.SaveConfig(); err != nil {
			return err
		}
	}

	if lcUser != nil {
		err = lcUser.Client.Connect()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("cannot load the current user")
	}

	if signerID == "" {
		signerID = api.GetSignerIDByApiKey(lcUser.Client.ApiKey)
	}

	fmt.Printf("Resolving dependencies...\n")
	deps, err := buildArtifact.Dependencies()
	if err != nil {
		return err
	}

	err = artifact.AuthDependencies(lcUser, buildArtifact, signerID, artifact.TrustLevel(trustLevel), autoNotarize, uint(threshold))
	if err != nil {
		return err
	}

	// write list of dependencies into .bom file for later processing by 'vcn a --bom ...'
	f, err := os.OpenFile(".bom", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, dep := range deps {
		_, err = f.WriteString("vcn." + dep.SignerID + "." + dep.Hash + "\n")
		if err != nil {
			return err
		}
	}

	if spdxFilename != "" {
		fmt.Printf("Outputting SPDX...\n")
		spdx, err := bom.OutputSpdxText(buildArtifact)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(spdxFilename, []byte(spdx), 0644)
		if err != nil {
			// show warning, but not error, because notarization succeeded
			fmt.Printf("Cannot output SPDX: %v", err)
		}
	}

	return nil
}
