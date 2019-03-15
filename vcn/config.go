/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 * Built on top of CLI (MIT license)
 * https://github.com/urfave/cli#overview
 */

package main

import (
	"math/big"
	"os"
	"time"

	"github.com/fatih/color"
)

type Level int

const (
	DISABLED          Level = -1
	UNKNOWN           Level = 0
	EMAIL_VERIFIED    Level = 1
	SOCIAL_VERIFIED   Level = 2
	ID_VERIFIED       Level = 3
	LOCATION_VERIFIED Level = 4
	VCHAIN            Level = 99
)

type Status int

const (
	TRUSTED        Status = 0
	UNTRUSTED      Status = 1
	STATUS_UNKNOWN Status = 2
	UNSUPPORTED    Status = 3
)

type Visibility int

const (
	PUBLIC  Visibility = 0
	PRIVATE Visibility = 1
)

func DashboardURL() string {
	// production:  https://dashboard.codenotary.io
	// staging: 	https://dashboard.staging.codenotary.io
	return "https://dashboard.codenotary.io"
}
func MainNetEndpoint() string {
	// production:  https://main.codenotary.io
	// staging: 	https://main.staging.codenotary.io
	return "https://main.codenotary.io"
}

func FoundationEndpoint() string {
	// production: 	https://api.codenotary.io/foundation
	// staging: 	https://api.staging.codenotary.io/foundation
	return "https://api.codenotary.io/foundation"
}

func AssetsRelayContractAddress() string {
	// production:  0x495021fe1a48a5b0c85ef1abd68c42cdfc7cda08
	// staging: 	0xf1d4b9fe8290bb5718db5d46c313e7b266570c21
	return "0x495021fe1a48a5b0c85ef1abd68c42cdfc7cda08"
}

func StyleAffordance() (color.Attribute, color.Attribute, color.Attribute) {
	return color.FgHiBlue, color.BgWhite, color.Bold
}

func StyleError() (color.Attribute, color.Attribute, color.Attribute) {
	return color.FgHiRed, color.BgHiWhite, color.Bold
}

func StyleSuccess() (color.Attribute, color.Attribute, color.Attribute) {
	return color.FgHiWhite, color.BgHiCyan, color.Bold
}

func ErrorWikiURL() string {
	return "https://github.com/vchain-us/vcn/wiki/Errors#"
}

func TrackingEvent() string {
	return FoundationEndpoint() + "/v1/tracking-event"
}

func TokenCheckEndpoint() string {
	return PublisherEndpoint() + "/auth/check"
}

func PublisherEndpoint() string {
	return FoundationEndpoint() + "/v1/publisher"
}

func WalletEndpoint() string {
	return FoundationEndpoint() + "/v1/wallet"
}

func ArtifactEndpoint(walletAddress string) string {
	return FoundationEndpoint() + "/v1/artifact?wallet-address=" + walletAddress
}

func TxVerificationRounds() uint64 {
	return 10
}

func PollInterval() time.Duration {
	return 1 * time.Second
}

func VcnDirectory() string {
	return os.Getenv("HOME") + "/.vcn"
}

func WalletDirectory() string {
	return VcnDirectory() + "/wallets"
}

func TokenFile() string {
	return VcnDirectory() + "/t"
}

func GasPrice() *big.Int {
	return big.NewInt(0)
}

func GasLimit() uint64 {
	return 20000000
}
