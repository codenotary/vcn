/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package nodecom

import (
	"encoding/json"
	"errors"
	"github.com/vchain-us/vcn/pkg/bom_node"
	"io/ioutil"
	"strings"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/uri"
)

// Scheme for java component
const Scheme = "nodecom"

// Artifact returns a file *api.Artifact from a given u
func Artifact(u *uri.URI, options ...extractor.Option) ([]*api.Artifact, error) {

	if u.Scheme != Scheme {
		return nil, nil
	}

	componentFolder := strings.TrimPrefix(u.Opaque, "//")

	pJSONP, err := bom_node.GetPackageJsonPath(componentFolder)
	if err != nil {
		return nil, err
	}
	o1, err := ioutil.ReadFile(pJSONP)
	if err != nil {
		return nil, err
	}

	var plJson map[string]interface{}
	err = json.Unmarshal(o1, &plJson)
	if err != nil {
		return nil, err
	}

	name, ok := plJson["name"].(string)
	if !ok {
		return nil, errors.New("malformed package json")
	}
	version, ok := plJson["version"].(string)
	if !ok {
		return nil, errors.New("malformed package json")
	}

	digest, err := bom_node.GetNodeComDigest(componentFolder)
	if err != nil {
		return nil, err
	}

	m := api.Metadata{}
	m["version"] = version
	m["name"] = name

	return []*api.Artifact{{
		Kind:     "nodecom",
		Name:     name + "-" + version,
		Hash:     digest.Encoded(),
		Metadata: m,
	}}, nil
}
