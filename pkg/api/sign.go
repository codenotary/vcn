/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"bytes"
	"context"
	goErr "errors"
	"fmt"
	"math/big"
	"math/rand"
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

//
var WrongPassphraseErr = goErr.New("incorrect notarization password")

// Sign is invoked by the User to notarize an artifact using the given functional options,
// if successful a BlockchainVerification is returned.
// By default, the artifact is notarized using status = meta.StatusTrusted, visibility meta.VisibilityPrivate.
// At least the key (secret) must be provided using SignWithKey().
func (u User) Sign(artifact Artifact, options ...SignOption) (*BlockchainVerification, error) {
	if artifact.Hash == "" {
		return nil, makeError("hash is missing", nil)
	}
	if artifact.Size < 0 {
		return nil, makeError("invalid size", nil)
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

	// In order to handle parallel calls here there is a retry mechanism if another transaction is already in place.
	// This is a workaround. Need a proper solution to handle parallel signing
	var verification *BlockchainVerification
	for i := uint64(0); i < meta.TxVerificationRounds(); i++ {
		verification, err = u.commitTransaction(artifact, options...)
		if err != nil {
			if err.Error() == errors.BlockchainPermission {
				rand.Seed(time.Now().UnixNano())
				sleepTime := time.Second * time.Duration(int64(rand.Intn(6)))
				time.Sleep(sleepTime)
				continue
			}
			break
		}
		return verification, err
	}
	return verification, err
}

func (u User) commitTransaction(
	artifact Artifact,
	opts ...SignOption,
) (verification *BlockchainVerification, err error) {

	o, err := makeSignOpts(opts...)
	if err != nil {
		return
	}
	transactor, err := bind.NewTransactor(bytes.NewReader([]byte(o.keyin)), o.passphrase)
	if err != nil {
		if err.Error() == "could not decrypt key with given passphrase" {
			err = WrongPassphraseErr
		}
		return nil, err
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
		return nil, err
	}
	timeout, err := waitForTx(client, tx.Hash(), meta.TxVerificationRounds(), meta.PollInterval())
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

func waitForTx(client *ethclient.Client, tx common.Hash, maxRounds uint64, pollInterval time.Duration) (timeout bool, err error) {
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
