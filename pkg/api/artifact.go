/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vchain-us/vcn/pkg/meta"
)

type Artifact struct {
	Kind        string
	Name        string
	Hash        string
	Size        uint64
	ContentType string
	Metadata
}

func (a Artifact) toRequest() *ArtifactRequest {
	aR := &ArtifactRequest{
		// root fields
		Kind:        a.Kind,
		Name:        a.Name,
		Hash:        a.Hash,
		Size:        a.Size,
		ContentType: a.ContentType,

		// custom metadata
		Metadata: a.Metadata,
	}

	// promote url from custom metadata to root
	aR.Url = a.Metadata.swipeString("url")

	return aR
}

type ArtifactRequest struct {
	// root fields
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Hash        string `json:"hash"`
	Size        uint64 `json:"size,omitempty"`
	ContentType string `json:"contentType"`
	Url         string `json:"url"`

	// custom metadata
	Metadata Metadata `json:"metadata"`

	// ArtifactRequest specific
	Visibility string `json:"visibility"`
	Status     string `json:"status"`
	MetaHash   string `json:"metaHash"`
}

type PagedArtifactResponse struct {
	Content []ArtifactResponse `json:"content"`
}

type ArtifactResponse struct {
	// root fields
	Kind        string `json:"kind" vcn:"Kind"`
	Name        string `json:"name" vcn:"Name"`
	Hash        string `json:"hash" vcn:"Hash"`
	Size        uint64 `json:"size" vcn:"Size"`
	ContentType string `json:"contentType" vcn:"ContentType"`
	Url         string `json:"url" vcn:"URL"`

	// custom metadata
	Metadata Metadata `json:"metadata" vcn:"Metadata"`

	// ArtifactResponse specific
	Level             int64  `json:"level"`
	Visibility        string `json:"visibility" vcn:"Visibility"`
	Status            string `json:"status"`
	CreatedAt         string `json:"createdAt"`
	VerificationCount uint64 `json:"verificationCount"`
	PublisherCount    uint64 `json:"publisherCount"`

	// Publisher info
	Publisher           string `json:"publisher" vcn:"Signer"`
	PublisherCompany    string `json:"publisherCompany" vcn:"Company"`
	PublisherWebsiteUrl string `json:"publisherWebsiteUrl" vcn:"Website"`
}

func (a ArtifactResponse) String() string {
	return fmt.Sprintf("Name:\t%s\nHash:\t%s\nStatus:\t%s\n\n",
		a.Name, a.Hash, a.Status)
}

// Artifact returns an new *Artifact from a
func (a ArtifactResponse) Artifact() *Artifact {
	return &Artifact{
		// root fields
		Kind:        a.Kind,
		Name:        a.Name,
		Hash:        a.Hash,
		Size:        a.Size,
		ContentType: a.ContentType,

		// custom metadata
		Metadata: a.Metadata,
	}
}

func (u User) createArtifact(verification *BlockchainVerification, walletAddress string,
	artifact Artifact, visibility meta.Visibility, status meta.Status) error {

	hasAuth, err := u.IsAuthenticated()
	if err != nil {
		return err
	}
	if !hasAuth {
		return makeAuthRequiredError()
	}

	aR := artifact.toRequest()
	aR.Visibility = meta.VisibilityName(visibility)
	aR.Status = meta.StatusName(status)
	aR.MetaHash = verification.MetaHash()

	restError := new(Error)
	r, err := newSling(u.token()).
		Post(meta.ArtifactEndpointForWallet(walletAddress)).
		BodyJSON(aR).Receive(nil, restError)
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("request failed: %s (%d)", restError.Message,
			restError.Status)
	}
	return nil
}

func (u *User) LoadAllArtifacts() ([]ArtifactResponse, error) {
	ret := []ArtifactResponse{}
	for _, pubKey := range u.cfg.PubKeys() {
		chunk, err := u.LoadArtifacts(pubKey)
		if err != nil {
			return nil, err
		}
		ret = append(ret, chunk...)
	}
	return ret, nil
}

func (u *User) LoadArtifacts(walletAddress string) ([]ArtifactResponse, error) {
	response := new(PagedArtifactResponse)
	restError := new(Error)
	r, err := newSling(u.token()).
		Get(meta.ArtifactEndpointForWallet(walletAddress)).
		Receive(&response, restError)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, fmt.Errorf("request failed: %s (%d)",
			restError.Message, restError.Status)
	}
	return response.Content, nil
}

// LoadArtifact returns an *ArtifactResponse for the given hash and current u, if any
func (u *User) LoadArtifact(hash string) (*ArtifactResponse, error) {
	notFound := func() (*ArtifactResponse, error) {
		return nil, fmt.Errorf("no asset matching hash %s signed by %s found", hash, u.Email())
	}
	response := new(PagedArtifactResponse)
	restError := new(Error)
	r, err := newSling(u.token()).
		Get(meta.ArtifactEndpoint()+"/"+hash+"?scope=CURRENT_USER&size=1&sort=createdAt,desc").
		Receive(&response, restError)
	if err != nil {
		return nil, err
	}

	switch r.StatusCode {
	case 200:
		if len(response.Content) < 1 {
			return notFound()
		}
	case 404:
		return notFound()
	default:
		return nil, fmt.Errorf("request failed: %s (%d)",
			restError.Message, restError.Status)
	}

	return &response.Content[0], nil
}

func LoadArtifactForHash(user *User, hash string, metahash string) (*ArtifactResponse, error) {
	response := new(ArtifactResponse)
	restError := new(Error)
	r, err := newSling(user.token()).
		Get(meta.ArtifactEndpoint()+"/"+hash+"/"+metahash).
		Receive(&response, restError)
	logger().WithFields(logrus.Fields{
		"response":  response,
		"err":       err,
		"restError": restError,
	}).Trace("LoadArtifactForHash")
	if err != nil {
		return nil, err
	}
	switch r.StatusCode {
	case 200:
		return response, nil
	case 404:
		return nil, fmt.Errorf("no asset matching hash %s/%s found", hash, metahash)
	}
	return nil, fmt.Errorf("loading artifact for hash failed: %+v", restError)
}
