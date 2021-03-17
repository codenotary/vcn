/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package sign

import (
	"encoding/json"
	"strings"
)

type mapOpts map[string]string

// Set adds the input value to the map, by splitting on '='.
func (m mapOpts) Set(value string) error {
	vals := strings.SplitN(value, "=", 2)
	if len(vals) == 1 {
		m[vals[0]] = ""
	} else {
		m[vals[0]] = vals[1]
	}
	return nil
}

func (m mapOpts) String() string {
	if len(m) < 1 {
		return ""
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func (m mapOpts) Type() string {
	return "key=value"
}

func (m mapOpts) StringToInterface() map[string]interface{} {
	r := map[string]interface{}{}
	for k, v := range m {
		r[k] = v
	}
	return r
}

func (m mapOpts) KeysToValues() []string {
	var as []string
	for k, _ := range m {
		as = append(as, k)
	}
	return as
}
