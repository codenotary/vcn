package pythoncom

import (
	"strings"
	"errors"
	"fmt"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/bom_component"
	"github.com/vchain-us/vcn/pkg/bom_python"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/uri"
)

// Scheme for Go component
const Scheme = "pythoncom"

// Artifact returns a file *api.Artifact from a given uri
func Artifact(u *uri.URI, options ...extractor.Option) ([]*api.Artifact, error) {
	if u.Scheme != Scheme {
		return nil, nil
	}

	path := strings.TrimPrefix(u.Opaque, "//")

	fields := strings.Split(path, "@")
	if len(fields) != 2 {
		return nil, errors.New("component path format is <name>@<version>")
	}

	hash, hashType, err := bom_python.QueryHash(fields[0], fields[1])
	if err != nil {
		return nil, fmt.Errorf("cannot get checksum for module %s: %w", path, err)
	}

	m := api.Metadata{}

	m["path"] = fields[0]
	m["version"] = fields[1]
	m["hashType"], _ = bom_component.HashTypeName(hashType)

	return []*api.Artifact{{
		Kind:        Scheme,
		Name:        fields[0],
		Hash:        hash,
		Size:        0,
		ContentType: "text/json; charset=utf-8",
		Metadata:    m,
	}}, nil
}