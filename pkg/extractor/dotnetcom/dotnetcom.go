package dotnetcom

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/bom/artifact"
	"github.com/vchain-us/vcn/pkg/bom/dotnet"
	"github.com/vchain-us/vcn/pkg/extractor"
	"github.com/vchain-us/vcn/pkg/uri"
)

// Scheme for Go component
const Scheme = "dotnetcom"

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

	pkgCacheDirs, err := dotnet.GetPackageHCacheDirs()
	if err != nil {
		return nil, err
	}

	hash, err := dotnet.GetPackageHash(fields[0], fields[1], pkgCacheDirs)
	if err != nil {
		return nil, fmt.Errorf("cannot get checksum for module %s: %w", path, err)
	}

	return []*api.Artifact{artifact.ToApiArtifact(dotnet.AssetType, fields[0], fields[1], hash, artifact.HashSHA512)}, nil
}
