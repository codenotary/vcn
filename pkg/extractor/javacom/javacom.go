/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package javacom

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/bom/artifact"
	"github.com/vchain-us/vcn/pkg/bom/java"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/uri"
)

// Scheme for java component
const Scheme = "javacom"

// Artifact returns a file *api.Artifact from a given u
func Artifact(u *uri.URI, options ...extractor.Option) ([]*api.Artifact, error) {

	if u.Scheme != Scheme {
		return nil, nil
	}

	path := strings.TrimPrefix(u.Opaque, "//")

	pomPath, err := java.GetPOM(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(pomPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	checksum := h.Sum(nil)

	pb, err := ioutil.ReadFile(pomPath)
	if err != nil {
		return nil, err
	}

	var project Pom
	if err := xml.Unmarshal(pb, &project); err != nil {
		return nil, err
	}

	return []*api.Artifact{artifact.ToApiArtifact(java.AssetType, project.Name, project.Parent.Version, hex.EncodeToString(checksum), artifact.HashSHA256)}, nil
}

type Pom struct {
	XMLName      xml.Name `xml:"project"`
	ModelVersion string   `xml:"modelVersion"`
	Parent       Parent   `xml:"parent"`
	Version      string   `xml:"version"`
	Name         string   `xml:"name"`
	Modules      []string `xml:"modules>module"`
}

type Parent struct {
	Version string `xml:"version"`
}
