/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package uri

import (
	"encoding/json"
	"fmt"
	"strings"
)

// URI represents the canonical identification of an artifact
type URI struct {
	Scheme string // Scheme identifies a kind of artifacts
	Opaque string // Rest of encoded data
}

// String implements the stringer interface
func (u *URI) String() string {
	if u.Scheme != "" {
		return fmt.Sprintf("%s:%s", u.Scheme, u.Opaque)
	}
	return u.Opaque
}

// Parse converts a rawURI string into an URI structure
func Parse(rawURI string) (*URI, error) {
	parts := strings.Split(rawURI, "://")
	l := len(parts)
	if l == 1 {
		return &URI{
			Scheme: "",
			Opaque: rawURI,
		}, nil
	}
	if l == 2 {
		return &URI{
			Scheme: parts[0],
			Opaque: "//" + parts[1],
		}, nil
	}
	return nil, fmt.Errorf("invalid URI: %s", rawURI)
}

// MarshalJSON implements the json.Marshaller interface
func (u URI) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON parses URI
func (u *URI) UnmarshalJSON(input []byte) error {
	var rawURI string
	if err := json.Unmarshal(input, &rawURI); err != nil {
		return err
	}
	pu, err := Parse(rawURI)
	if err != nil {
		return err
	}
	u.Scheme = pu.Scheme
	u.Opaque = pu.Opaque
	return nil
}
