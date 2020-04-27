/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"context"
	goErr "errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/vchain-us/vcn/internal/blockchain"
	"github.com/vchain-us/vcn/internal/errors"
	"github.com/vchain-us/vcn/pkg/meta"
)

// ErrWrongPassphrase is error for incorrect password notarization
var ErrWrongPassphrase = goErr.New("incorrect notarization password")

// Sign is invoked by the User to notarize an artifact using the given functional options,
// if successful a BlockchainVerification is returned.
// By default, the artifact is notarized using status = meta.StatusTrusted, visibility meta.VisibilityPrivate.
// At least the key (secret) must be provided using SignWithKey().
func (u *User) Sign(artifact *Artifact, options ...SignOption) (*BlockchainVerification, error) {
	if artifact.Hash == "" {
		return nil, makeError("hash is missing", nil)
	}

	hasAuth, err := u.IsAuthenticated()
	if err != nil {
		return nil, err
	}
	if !hasAuth {
		return nil, makeAuthRequiredError()
	}

	trialExpired, err := u.trialExpired()
	if err != nil {
		return nil, err
	}
	if trialExpired {
		return nil, fmt.Errorf(errors.TrialExpired)
	}

	opsLeft, err := u.RemainingSignOps()
	if err != nil {
		return nil, err
	}

	if opsLeft < 1 {
		return nil, fmt.Errorf(errors.NoRemainingSignOps)
	}

	return u.commitTransaction(
		artifact,
		options...,
	)
}

func (u *User) commitTransaction(
	artifact *Artifact,
	opts ...SignOption,
) (verification *BlockchainVerification, err error) {

	o, err := makeSignOpts(opts...)
	if err != nil {
		return
	}
	transactor, err := bind.NewTransactor(o.keyin, o.passphrase)
	if err != nil {
		if err.Error() == "could not decrypt key with given passphrase" {
			err = ErrWrongPassphrase
		}
		return
	}

	transactor.GasLimit = meta.GasLimit()
	transactor.GasPrice = meta.GasPrice()
	client, err := ethclient.Dial(meta.MainNet())
	if err != nil {
		err = makeError(
			errors.BlockchainCannotConnect,
			logrus.Fields{
				"error":   err,
				"network": meta.MainNet(),
			})
		return
	}
	address := common.HexToAddress(meta.AssetsRelayContractAddress())
	instance, err := blockchain.NewAssetsRelay(address, client)
	if err != nil {
		err = makeFatal(
			errors.BlockchainContractErr,
			logrus.Fields{
				"error":    err,
				"contract": meta.AssetsRelayContractAddress(),
			},
		)
		return
	}
	tx, err := instance.Sign(transactor, artifact.Hash, big.NewInt(int64(o.status)))
	if err != nil {
		err = makeFatal(
			errors.SignFailed,
			logrus.Fields{
				"error": err,
				"hash":  artifact.Hash,
			},
		)
		return
	}
	timeout, err := waitForTx(tx.Hash(), meta.TxVerificationRounds(), meta.PollInterval())
	if err != nil {
		err = makeFatal(
			errors.BlockchainPermission,
			logrus.Fields{
				"error": err,
			},
		)
		return
	}
	if timeout {
		err = makeFatal(
			errors.BlockchainTimeout,
			logrus.Fields{
				"error": err,
			},
		)
		return
	}

	signerID := transactor.From.Hex()
	verification, err = VerifyMatchingSignerID(artifact.Hash, signerID)
	if err != nil {
		return
	}

	err = u.createArtifact(verification, strings.ToLower(signerID), artifact, o.visibility, o.status, tx.Hash())
	return
}

func waitForTx(tx common.Hash, maxRounds uint64, pollInterval time.Duration) (timeout bool, err error) {
	client, err := ethclient.Dial(meta.MainNet())
	if err != nil {
		return false, err
	}
	for i := uint64(0); i < maxRounds; i++ {
		_, pending, err := client.TransactionByHash(context.Background(), tx)
		if err != nil {
			return false, err
		}
		if !pending {
			return false, nil
		}
		time.Sleep(pollInterval)
	}
	return true, nil
}
