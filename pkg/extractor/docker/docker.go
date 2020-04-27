/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package docker

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/uri"
)

// Scheme for docker (default)
const Scheme = "docker"

// SchemePodman is the scheme for podman (Docker-compatible CLI interface)
const SchemePodman = "podman"

var schemes = map[string]bool{Scheme: true, SchemePodman: true}

// Artifact returns a file *api.Artifact from a given u
func Artifact(u *uri.URI, options ...extractor.Option) (*api.Artifact, error) {

	if !schemes[u.Scheme] {
		return nil, nil
	}

	id := strings.TrimPrefix(u.Opaque, "//")
	images, err := inspect(u.Scheme, id)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect %s image: %s", u.Scheme, err)
	}
	if len(images) < 1 {
		return nil, fmt.Errorf("no %s image found for: %s", u.Scheme, id)
	}

	i := images[0]

	m := api.Metadata{
		"architecture": i.Architecture,
		"platform":     i.Os,
	}

	if version := i.inferVer(); version != "" {
		m["version"] = version
	}

	m[u.Scheme] = i
	return &api.Artifact{
		Kind:     u.Scheme,
		Name:     u.Scheme + "://" + i.name(),
		Hash:     i.hash(),
		Size:     i.Size,
		Metadata: m,
	}, nil
}

type image struct {
	ID            string      `json:"Id"`
	RepoTags      []string    `json:"RepoTags"`
	RepoDigests   []string    `json:"RepoDigests"`
	Comment       string      `json:"Comment"`
	Created       string      `json:"Created"`
	DockerVersion string      `json:"DockerVersion"`
	Author        string      `json:"Author"`
	Architecture  string      `json:"Architecture"`
	Os            string      `json:"Os"`
	VirtualSize   uint64      `json:"VirtualSize"`
	Size          uint64      `json:"Size"`
	Metadata      interface{} `json:"Metadata"`
}

func (i *image) hash() string {
	return strings.TrimSpace(strings.Replace(fmt.Sprint(i.ID), "sha256:", "", 1))
}

func (i *image) name() string {
	if len(i.RepoTags) > 0 {
		return i.RepoTags[0]
	}
	return i.hash()
}

func (i *image) inferVer() string {
	if len(i.RepoTags) > 0 {
		parts := strings.SplitN(i.RepoTags[0], ":", 2)
		if len(parts) > 1 && parts[1] != "latest" {
			return parts[1]
		}
	}

	return ""
}

func inspect(executable, arg string) ([]image, error) {
	cmd := exec.Command(executable, "inspect", arg, "--type", "image")
	cmdOutput, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
			return nil, fmt.Errorf(string(ee.Stderr))
		}
		return nil, err
	}
	data := []image{}
	if err = json.Unmarshal(cmdOutput, &data); err != nil {
		return nil, err
	}
	return data, nil
}
