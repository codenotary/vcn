/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom_package

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vchain-us/vcn/pkg/bom_component"
)

const (
	optional = iota
	mandatory
)

type headerLine struct {
	tag      string
	presense int
	fn       func(Package) (string, error)
}

var headerContent = []headerLine{
	{"SPDXVersion", mandatory, func(Package) (string, error) { return "SPDX-2.2", nil }},
	{"DataLicense", mandatory, func(Package) (string, error) { return "CC0-1.0", nil }},
	{"SPDXID", mandatory, func(Package) (string, error) { return "SPDXRef-DOCUMENT", nil }},
	{"DocumentName", mandatory, func(p Package) (string, error) {
		path, err := filepath.Abs(p.Path())
		if err != nil {
			return "", err
		}
		return filepath.Base(path), nil
	}},
	{"DocumentNamespace", mandatory, func(p Package) (string, error) {
		path, err := filepath.Abs(p.Path())
		if err != nil {
			return "", err
		}
		return "http://spdx.org/spdxdocs/" + filepath.Base(path) + "-" + uuid.NewString(), nil
	}},
	{"Creator", mandatory, func(Package) (string, error) { return "Tool: CodeNotary vcn", nil }},
	{"Created", mandatory, func(Package) (string, error) { return time.Now().UTC().Format(time.RFC3339), nil }},
}

type componentLine struct {
	tag      string
	presense int
	fn       func(bom_component.Component, int) (string, error)
}

var componentContent = []componentLine{
	{"PackageName", mandatory, func(c bom_component.Component, seq int) (string, error) { return c.Name, nil }},
	{"SPDXID", mandatory, func(c bom_component.Component, seq int) (string, error) {
		return "SPDXRef-Package-" + strconv.Itoa(seq), nil
	}},
	{"PackageVersion", mandatory, func(c bom_component.Component, seq int) (string, error) { return c.Version, nil }},
	{"PackageDownloadLocation", mandatory, noAssertion},
	// FilesAnalysed is optional, but by default it is true, which requires presence of many other fields
	{"FilesAnalyzed", mandatory, func(c bom_component.Component, seq int) (string, error) { return "false", nil }},
	{"PackageChecksum", mandatory, func(c bom_component.Component, seq int) (string, error) {
		hashType, err := c.HashTypeName()
		if err != nil {
			return "", err
		}
		return hashType + ": " + c.Hash, nil
	}},
	{"PackageLicenseConcluded", mandatory, noAssertion},
	{"PackageLicenseDeclared", mandatory, noAssertion},
	{"PackageCopyrightText", mandatory, noAssertion},
}

func noAssertion(c bom_component.Component, seq int) (string, error) { return "NOASSERTION", nil }

// Output info about package and its components in SPDX text (tag:value) format, according to
// SPDX spec 2.2: https://spdx.dev/wp-content/uploads/sites/41/2020/08/SPDX-specification-2-2.pdf
func OutputSpdxText(p Package) (string, error) {
	builder := strings.Builder{}

	// SPDX header
	for _, line := range headerContent {
		value, err := line.fn(p)
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
	components, err := p.Components()
	if err != nil {
		return "", fmt.Errorf("cannot get package components: %w", err)
	}
	for i, comp := range components {
		for _, line := range componentContent {
			value, err := line.fn(comp, i+1)
			if err != nil {
				if line.presense == mandatory {
					return "", fmt.Errorf("cannot get value for tag %s for component %s: %w", line.tag, comp.Name, err)
				}
				continue // optional tag - ignore error
			}
			if value == "" {
				if line.presense == mandatory {
					return "", fmt.Errorf("no value found for mandatory component tag %s for component %s", line.tag, comp.Name)
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
