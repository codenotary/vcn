/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package javacom

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/uri"
	"github.com/vchain-us/vcn/pkg/bom/java"
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

	m := api.Metadata{}

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	checksum := h.Sum(nil)

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	ct := "text/xml; charset=utf-8"

	pb, err := ioutil.ReadFile(pomPath)
	if err != nil {
		return nil, err
	}

	var project Pom
	if err := xml.Unmarshal(pb, &project); err != nil {
		return nil, err
	}

	m["groupId"] = project.Parent.GroupId
	m["artifactId"] = project.ArtifactId
	m["version"] = project.Parent.Version
	m["packaging"] = project.Packaging

	return []*api.Artifact{{
		Kind:        java.AssetType,
		Name:        project.Name,
		Hash:        hex.EncodeToString(checksum),
		Size:        uint64(stat.Size()),
		ContentType: ct,
		Metadata:    m,
	}}, nil
}

type Pom struct {
	XMLName      xml.Name     `xml:"project"`
	ModelVersion string       `xml:"modelVersion"`
	Parent       Parent       `xml:"parent"`
	GroupId      string       `xml:"groupId"`
	ArtifactId   string       `xml:"artifactId"`
	Version      string       `xml:"version"`
	Packaging    string       `xml:"packaging"`
	Name         string       `xml:"name"`
	Dependencies []Dependency `xml:"dependencies>dependency"`
	Modules      []string     `xml:"modules>module"`
}

type Parent struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type Dependency struct {
	XMLName    xml.Name `xml:"dependency"`
	GroupId    string   `xml:"groupId"`
	ArtifactId string   `xml:"artifactId"`
	Version    string   `xml:"version"`
	Classifier string   `xml:"classifier"`
	Type       string   `xml:"type"`
	Scope      string   `xml:"scope"`
}
