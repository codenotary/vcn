/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package meta

import (
	"math/big"
	"os"
	"time"
)

const DefaultGasLimit = 20000000

// DashboardURL returns the CodeNotary's dashboard URL.
func DashboardURL() string {
	switch StageEnvironment() {
	case StageTest:
		return os.Getenv("VCN_TEST_DASHBOARD")
	case StageStaging:
		return "https://dashboard.staging.codenotary.io"
	case StageProduction:
		fallthrough
	default:
		return "https://dashboard.codenotary.io"
	}
}

// MainNet returns the CodeNotary mainnet URL.
func MainNet() string {
	switch StageEnvironment() {
	case StageTest:
		return os.Getenv("VCN_TEST_NET")
	case StageStaging:
		return "https://main.staging.codenotary.io"
	case StageProduction:
		fallthrough
	default:
		return "https://main.codenotary.io"
	}
}

// APIEndpoint returns the API's endpoint URL for a given resource.
func APIEndpoint(resource string) string {
	base := ""
	switch StageEnvironment() {
	case StageTest:
		base = os.Getenv("VCN_TEST_API")
	case StageStaging:
		base = "https://api.staging.codenotary.io/foundation"
	case StageProduction:
		fallthrough
	default:
		base = "https://api.codenotary.io/foundation"
	}

	return base + "/v1/" + resource
}

// AssetsRelayContractAddress returns the AssetsRelay smart contract public address.
func AssetsRelayContractAddress() string {
	switch StageEnvironment() {
	case StageTest:
		return os.Getenv("VCN_TEST_CONTRACT")
	case StageStaging:
		return "0x4eb8d2866da4341796ce64a983786a01b1072939"
	case StageProduction:
		fallthrough
	default:
		return "0x41a749a79a78b388607df06c25adbc73dbbf1e87"
	}
}

// OrganisationsRelayContractAddress returns the OrganisationsRelay smart contract public address.
func OrganisationsRelayContractAddress() string {
	switch StageEnvironment() {
	case StageTest:
		return os.Getenv("VCN_TEST_CONTRACT_ORG")
	case StageStaging:
		return "0x4a9a0547949ec55ecbf06738e8c2bad747f410bb"
	case StageProduction:
		fallthrough
	default:
		return "0x258e39ff07e6e3a2430aa951f387cfbd808835bc"
	}
}

// TxVerificationRounds returns the maximum number of rounds to try before considering a pending transaction failed.
// The duration of each round is returned by PollInterval()
func TxVerificationRounds() uint64 {
	return 30
}

// PollInterval returns the waiting time between each round.
// See TxVerificationRounds().
func PollInterval() time.Duration {
	return 2 * time.Second
}

// GasPrice for transactions.
func GasPrice() *big.Int {
	return big.NewInt(0)
}

// GasLimit for transactions.
func GasLimit() uint64 {
	return DefaultGasLimit
}
