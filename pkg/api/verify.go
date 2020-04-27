/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"crypto/sha256"
	"encoding/json"
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

// BlockchainVerification represents the notarized data onto the blockchain.
type BlockchainVerification struct {
	Owner     common.Address `json:"owner" yaml:"owner"`
	Level     meta.Level     `json:"level" yaml:"level"`
	Status    meta.Status    `json:"status" yaml:"status"`
	Timestamp time.Time      `json:"timestamp" yaml:"timestamp"`
}

// Trusted returns true if v.Status is meta.StatusTrusted
func (v *BlockchainVerification) Trusted() bool {
	return v != nil && v.Status == meta.StatusTrusted
}

// Unknown returns true if v is nil or v.Status is meta.StatusUnknown
func (v *BlockchainVerification) Unknown() bool {
	return v == nil || v.Status == meta.StatusUnknown
}

func (v *BlockchainVerification) toMap() map[string]interface{} {
	if v == nil {
		return nil
	}
	return map[string]interface{}{
		"owner":     v.SignerID(),
		"level":     v.Level,
		"status":    v.Status,
		"timestamp": v.Date(),
	}
}

func (v *BlockchainVerification) fromUnmarshaler(unmarshal func(interface{}) error) error {
	if v == nil {
		v = &BlockchainVerification{}
	}
	data := struct {
		Owner     string
		Level     int64
		Status    int64
		Timestamp string
	}{}

	if err := unmarshal(&data); err != nil {
		return err
	}

	if data.Owner != "" {
		v.Owner = common.HexToAddress(data.Owner)
	}
	v.Level = meta.Level(data.Level)
	v.Status = meta.Status(data.Status)
	if data.Timestamp != "" {
		err := v.Timestamp.UnmarshalText([]byte(data.Timestamp))
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (v *BlockchainVerification) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.toMap())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *BlockchainVerification) UnmarshalJSON(b []byte) error {
	return v.fromUnmarshaler(func(value interface{}) error {
		return json.Unmarshal(b, value)
	})
}

// MarshalYAML implements the yaml.Marshaler interface.
func (v *BlockchainVerification) MarshalYAML() (interface{}, error) {
	return v.toMap(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (v *BlockchainVerification) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return v.fromUnmarshaler(unmarshal)
}

// MetaHash returns the SHA256 digest of BlockchainVerification's data.
// The returned value uniquely identify a single notarization.
func (v *BlockchainVerification) MetaHash() string {
	if v == nil {
		return ""
	}
	metadata := fmt.Sprintf("%s-%d-%d-%d",
		v.Owner.Hex(),
		int64(v.Level),
		int64(v.Status),
		v.Timestamp.Unix())
	metadataHashAsBytes := sha256.Sum256([]byte(metadata))
	metahash := fmt.Sprintf("%x", metadataHashAsBytes)
	logger().WithFields(logrus.Fields{
		"metadata": metadata,
		"metahash": metahash,
	}).Trace("Generated metahash")
	return metahash
}

// SignerID returns the public address derived from owner's public key (v.Owner), if any, otherwise an empty string.
func (v *BlockchainVerification) SignerID() string {
	if v != nil && v.Owner != common.BigToAddress(big.NewInt(0)) {
		return strings.ToLower(v.Owner.Hex())
	}
	return ""
}

// Date returns a RFC3339 formatted string of verification time (v.Timestamp), if any, otherwise an empty string.
func (v *BlockchainVerification) Date() string {
	if v != nil {
		ut := v.Timestamp.UTC()
		if ut.Unix() > 0 {
			return ut.Format(time.RFC3339)
		}
	}
	return ""
}

func callVerifyFunc(f func(*blockchain.AssetsRelay) (common.Address, *big.Int, *big.Int, *big.Int, error)) (*BlockchainVerification, error) {
	client, err := ethclient.Dial(meta.MainNet())
	if err != nil {
		return nil, err
	}
	contractAddress := common.HexToAddress(meta.AssetsRelayContractAddress())
	instance, err := blockchain.NewAssetsRelay(contractAddress, client)
	if err != nil {
		return nil, err
	}
	address, level, status, timestamp, err := f(instance)
	if err != nil {
		return nil, err
	}
	if meta.Status(status.Int64()) != meta.StatusUnknown && address != common.BigToAddress(big.NewInt(0)) {
		verification := &BlockchainVerification{
			Owner:     address,
			Level:     meta.Level(level.Int64()),
			Status:    meta.Status(status.Int64()),
			Timestamp: time.Unix(timestamp.Int64(), 0),
		}
		logger().
			WithField("verification", verification).
			Trace("Blockchain verification found")
		return verification, nil
	}

	logger().Trace("No blockchain verification found")
	return &BlockchainVerification{
		Status: meta.StatusUnknown,
	}, nil
}

// Verify returns the most recent *BlockchainVerification with highest level available for the given hash.
func Verify(hash string) (*BlockchainVerification, error) {
	logger().WithFields(logrus.Fields{
		"hash": hash,
	}).Trace("Verify")

	return callVerifyFunc(func(instance *blockchain.AssetsRelay) (common.Address, *big.Int, *big.Int, *big.Int, error) {
		return instance.Verify(nil, hash)
	})
}

// VerifyMatchingSignerIDWithFallback returns *BlockchainVerification for the hash matching a given SignerID,
// if any, otherwise it returns the same result of Verify().
func VerifyMatchingSignerIDWithFallback(hash, signerID string) (*BlockchainVerification, error) {
	logger().WithFields(logrus.Fields{
		"hash":     hash,
		"signerID": signerID,
	}).Trace("VerifyMatchingSignerIDWithFallback")

	address := common.HexToAddress(signerID)

	return callVerifyFunc(func(instance *blockchain.AssetsRelay) (common.Address, *big.Int, *big.Int, *big.Int, error) {
		return instance.VerifyAgainstPublisherWithFallback(nil, hash, address)
	})
}

// VerifyMatchingSignerID returns *BlockchainVerification for hash matching a given SignerID.
func VerifyMatchingSignerID(hash, signerID string) (*BlockchainVerification, error) {
	return VerifyMatchingSignerIDs(hash, []string{signerID})
}

// VerifyMatchingSignerIDs returns *BlockchainVerification for hash
// matching at least one of signerIDs.
func VerifyMatchingSignerIDs(hash string, signerIDs []string) (*BlockchainVerification, error) {
	logger().WithFields(logrus.Fields{
		"hash":      hash,
		"signerIDs": signerIDs,
	}).Trace("VerifyMatchingSignerIDs")

	addresses := make([]common.Address, len(signerIDs))
	for i, s := range signerIDs {
		addresses[i] = common.HexToAddress(s)
	}

	return callVerifyFunc(func(instance *blockchain.AssetsRelay) (common.Address, *big.Int, *big.Int, *big.Int, error) {
		return instance.VerifyAgainstPublishers(nil, hash, addresses)
	})
}
