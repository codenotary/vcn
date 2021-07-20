/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/vchain-us/vcn/pkg/bom/artifact"
)

const (
	optional = iota
	mandatory
)

type headerLine struct {
	tag      string
	presense int
	fn       func(artifact.Artifact) (string, error)
}

var headerContent = []headerLine{
	{"SPDXVersion", mandatory, func(artifact.Artifact) (string, error) {
		return "SPDX-2.2", nil
	}},
	{"DataLicense", mandatory, func(artifact.Artifact) (string, error) {
		return "CC0-1.0", nil
	}},
	{"SPDXID", mandatory, func(artifact.Artifact) (string, error) {
		return "SPDXRef-DOCUMENT", nil
	}},
	{"DocumentName", mandatory, func(p artifact.Artifact) (string, error) {
		path, err := filepath.Abs(p.Path())
		if err != nil {
			return "", err
		}
		return filepath.Base(path), nil
	}},
	{"DocumentNamespace", mandatory, func(p artifact.Artifact) (string, error) {
		path, err := filepath.Abs(p.Path())
		if err != nil {
			return "", err
		}
		return "http://spdx.org/spdxdocs/" + filepath.Base(path) + "-" + uuid.NewString(), nil
	}},
	{"Creator", mandatory, func(artifact.Artifact) (string, error) {
		return "Tool: CodeNotary vcn", nil
	}},
	{"Created", mandatory, func(artifact.Artifact) (string, error) {
		return time.Now().UTC().Format(time.RFC3339), nil
	}},
}

type componentLine struct {
	tag      string
	presense int
	fn       func(artifact.Dependency, int) (string, error)
}

var componentContent = []componentLine{
	{"PackageName", mandatory, func(d artifact.Dependency, seq int) (string, error) {
		return d.Name, nil
	}},
	{"SPDXID", mandatory, func(d artifact.Dependency, seq int) (string, error) {
		return "SPDXRef-Package-" + strconv.Itoa(seq), nil
	}},
	{"PackageVersion", mandatory, func(d artifact.Dependency, seq int) (string, error) {
		return d.Version, nil
	}},
	{"PackageDownloadLocation", mandatory, noAssertion},
	// FilesAnalysed is optional, but by default it is true, which requires presence of many other fields
	{"FilesAnalyzed", mandatory, func(d artifact.Dependency, seq int) (string, error) {
		return "false", nil
	}},
	{"PackageChecksum", mandatory, func(d artifact.Dependency, seq int) (string, error) {
		return artifact.HashTypeName(d.HashType) + ": " + d.Hash, nil
	}},
	{"PackageLicenseConcluded", mandatory, noAssertion},
	{"PackageLicenseDeclared", mandatory, noAssertion},
	{"PackageCopyrightText", mandatory, noAssertion},
	{"PackageComment", optional, func(d artifact.Dependency, seq int) (string, error) {
		return artifact.TrustLevelName(d.TrustLevel), nil
	}},
}

func noAssertion(d artifact.Dependency, seq int) (string, error) { return "NOASSERTION", nil }

// Output info about package and its components in SPDX text (tag:value) format, according to
// SPDX spec 2.2: https://spdx.dev/wp-content/uploads/sites/41/2020/08/SPDX-specification-2-2.pdf
func OutputSpdxText(a artifact.Artifact) (string, error) {
	deps, err := a.Dependencies()
	if err != nil {
		return "", err
	}
	builder := strings.Builder{}

	// SPDX header
	for _, line := range headerContent {
		value, err := line.fn(a)
		if err != nil {
			if line.presense == mandatory {
				return "", fmt.Errorf("cannot get value for tag %s: %w", line.tag, err)
			}
			continue // optional tag - ignore error
		}
		if value == "" {
			if line.presense == mandatory {
				return "", fmt.Errorf("no value found for mandatory header tag %s", line.tag)
			}
			continue // optional
		}
		builder.WriteString(line.tag)
		builder.WriteString(": ")
		builder.WriteString(value)
		builder.WriteByte('\n')
	}

	builder.WriteString("\n##### Software components\n\n")
	if deps == nil {
		deps, err = a.Dependencies()
		if err != nil {
			return "", fmt.Errorf("cannot get package components: %w", err)
		}
	}
	for i, dep := range deps {
		for _, line := range componentContent {
			value, err := line.fn(dep, i+1)
			if err != nil {
				if line.presense == mandatory {
					return "", fmt.Errorf("cannot get value for tag %s for component %s: %w", line.tag, dep.Name, err)
				}
				continue // optional tag - ignore error
			}
			if value == "" {
				if line.presense == mandatory {
					return "", fmt.Errorf("no value found for mandatory component tag %s for component %s", line.tag, dep.Name)
				}
				continue // optional
			}
			builder.WriteString(line.tag)
			builder.WriteString(": ")
			builder.WriteString(value)
			builder.WriteByte('\n')
		}
		builder.WriteByte('\n')
	}

	return builder.String(), nil
}
