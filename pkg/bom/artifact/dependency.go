/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package artifact

import (
	"fmt"
	"errors"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/meta"
)

type HashType int
const (
	HashInvalid HashType = iota
	HashSHA1
	HashSHA224
	HashSHA256
	HashSHA384
	HashSHA512
	HashMD2
	HashMD4
	HashMD5
	HashMD6
	minHash = HashSHA1
	maxHash = HashMD6
)
var hashText = [maxHash+1]string{"Invalid", "SHA1", "SHA224", "SHA256", "SHA384", "SHA512", "MD2", "MD4", "MD5", "MD6"}

type TrustLevel int
const (
	Untrusted TrustLevel = iota
	Unsupported
	Trusted
	MinTrustLevel = Untrusted
	MaxTrustLevel = Trusted
)
var levelText = [MaxTrustLevel+1]string{"Untrusted", "Unsupported", "Trusted"}

const MaxGoroutines = 8 // used by other packages that query components from external sources

// Dependency is a single building block, used for building the Artifact
type Dependency struct {
	Name     string
	Version  string
	Hash     string
	HashType HashType
	TrustLevel TrustLevel	// set by Notorize/Authenticate
	SignerID string			// set by Notorize/Authenticate
}

func HashTypeName(hashType HashType) string {
	if hashType < minHash || hashType > maxHash {
		return hashText[HashInvalid]
	}
	return hashText[hashType]
}

func TrustLevelName(level TrustLevel) string {
	if level < MinTrustLevel || level >= MaxTrustLevel {
		return levelText[Untrusted]
	}
	return levelText[level]
}

// Notarize() notarizes the dependency
func (dep *Dependency) Notarize(lcUser *api.LcUser, kind string) error {
	_, _, err := lcUser.Sign(api.Artifact{
			Kind: kind,
			Name: dep.Name,
			Hash: dep.Hash,
			Size: 10,	// FIXME
			ContentType: "text/plain; charset=utf-8",
			Metadata: map[string]interface{}{
				"version": dep.Version,
				"hashType": HashTypeName(dep.HashType)},	// FIXME - add more metadata
		},
		api.LcSignWithStatus(meta.StatusTrusted))
	if err != nil {
		return fmt.Errorf("notarization of %s failed: %w", dep.Name, err)
	}
	dep.TrustLevel = Trusted
	dep.SignerID = api.GetSignerIDByApiKey(lcUser.Client.ApiKey)

	return nil
}

// Authenticate() checks dependency trust level, sets TrustLevel and Signer fields
func (dep *Dependency) Authenticate(lcUser *api.LcUser, signerID string) error {
	level := Unsupported

	ar, verified, err := lcUser.LoadArtifact(dep.Hash, signerID, "", 0)
	if err == nil {
		switch {
		case !verified:
			return errors.New("the ledger is compromised")
		case ar.Revoked != nil && !ar.Revoked.IsZero():
			level = Untrusted
		case ar.Status == meta.StatusUnsupported:
			level = Unsupported
		default:
			level = Trusted
		}
	} else if err != api.ErrNotFound {
		return err
	}

	dep.TrustLevel = level
	if level != Unsupported {
		dep.SignerID = signerID
	}

	return nil
}
