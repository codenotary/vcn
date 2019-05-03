/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/vchain-us/vcn/internal/blockchain"
	"github.com/vchain-us/vcn/pkg/meta"
)

type BlockchainVerification struct {
	Owner     common.Address
	Level     meta.Level
	Status    meta.Status
	Timestamp time.Time
}

func (verification *BlockchainVerification) MetaHash() string {
	metadata := fmt.Sprintf("%s-%d-%d-%d",
		verification.Owner.Hex(),
		int64(verification.Level),
		int64(verification.Status),
		int64(verification.Timestamp.Unix()))
	metadataHashAsBytes := sha256.Sum256([]byte(metadata))
	logger().WithFields(logrus.Fields{
		"metahash": metadata,
	}).Trace("Generated metahash")
	return fmt.Sprintf("%x", metadataHashAsBytes)
}

func BlockChainVerify(hash string) (verification *BlockchainVerification, err error) {
	logger().WithFields(logrus.Fields{
		"hash": hash,
	}).Trace("BlockChainVerify")
	client, err := ethclient.Dial(meta.MainNetEndpoint())
	if err != nil {
		return nil, err
	}
	contractAddress := common.HexToAddress(meta.AssetsRelayContractAddress())
	instance, err := blockchain.NewAssetsRelay(contractAddress, client)
	if err != nil {
		return nil, err
	}
	address, level, status, timestamp, err := instance.Verify(nil, hash)
	if err != nil {
		return nil, err
	}
	verification = new(BlockchainVerification)
	verification.Owner = address
	verification.Level = meta.Level(level.Int64())
	verification.Status = meta.Status(status.Int64())
	verification.Timestamp = time.Unix(timestamp.Int64(), 0)
	logger().
		WithField("verification", verification).
		Trace("Received blockchain verification")
	return verification, nil
}

func BlockChainVerifyMatchingPublicKey(hash string, pubKey string) (verification *BlockchainVerification, err error) {
	return BlockChainVerifyMatchingPublicKeys(hash, []string{pubKey})
}

func BlockChainVerifyMatchingPublicKeys(hash string, pubKeys []string) (verification *BlockchainVerification, err error) {
	logger().WithFields(logrus.Fields{
		"hash": hash,
		"keys": pubKeys,
	}).Trace("BlockChainVerifyMatchingPublicKeys")

	// Connect and get verification count
	client, err := ethclient.Dial(meta.MainNetEndpoint())
	if err != nil {
		return nil, err
	}
	contractAddress := common.HexToAddress(meta.AssetsRelayContractAddress())
	instance, err := blockchain.NewAssetsRelay(contractAddress, client)
	if err != nil {
		return nil, err
	}
	count, err := instance.GetAssetCountForHash(nil, hash)
	if err != nil {
		return nil, err
	}

	// Make a map to lookup pubKey quickly
	keysMap := map[string]bool{}
	for i, pubKey := range pubKeys {
		pubKey = strings.ToLower(pubKey)
		pubKeys[i] = pubKey
		keysMap[pubKey] = true
	}

	// Iterate over verifications
	for i := count.Int64() - 1; i >= 0; i-- {
		address, level, status, timestamp, err := instance.VerifyByIndex(nil, hash, big.NewInt(i))
		if err != nil {
			return nil, err
		}
		if keysMap[strings.ToLower(address.Hex())] {
			return &BlockchainVerification{
				Owner:     address,
				Level:     meta.Level(level.Int64()),
				Status:    meta.Status(status.Int64()),
				Timestamp: time.Unix(timestamp.Int64(), 0),
			}, nil
		}
	}

	return nil, fmt.Errorf("no matching asset for hash %s and publicKeys %s", hash, pubKeys)
}
