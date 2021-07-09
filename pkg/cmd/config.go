/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package cmd

import (
	"fmt"
	"github.com/vchain-us/vcn/pkg/extractor/javacom"
	"github.com/vchain-us/vcn/pkg/extractor/wildcard"
	"os"

	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/extractor/dir"
	"github.com/vchain-us/vcn/pkg/extractor/docker"
	"github.com/vchain-us/vcn/pkg/extractor/file"
	"github.com/vchain-us/vcn/pkg/extractor/git"
	"github.com/vchain-us/vcn/pkg/extractor/gocom"

	"github.com/vchain-us/vcn/pkg/store"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Register metadata extractors
	extractor.Register("", wildcard.Artifact)
	extractor.Register(file.Scheme, file.Artifact)
	extractor.Register(dir.Scheme, dir.Artifact)
	extractor.Register(docker.Scheme, docker.Artifact)
	extractor.Register(docker.SchemePodman, docker.Artifact)
	extractor.Register(git.Scheme, git.Artifact)
	extractor.Register(wildcard.Scheme, wildcard.Artifact)
	extractor.Register(javacom.Scheme, javacom.Artifact)
	extractor.Register(gocom.Scheme, gocom.Artifact)

	// Load config
	if cfgFile != "" {
		store.SetConfigFile(cfgFile)
		if output, _ := rootCmd.PersistentFlags().GetString("output"); output == "" {
			fmt.Println("Using config file: ", store.ConfigFile())
		}
	}
	if err := store.LoadConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
