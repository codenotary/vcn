/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package cli

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/meta"
	"gopkg.in/yaml.v2"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
	"time"
)

func PrintLc(output string, r *types.LcResult) error {
	switch output {
	case "":
		WriteLcResultTo(r, colorable.NewColorableStdout())
	case "attachments":
		WriteLcResultTo(r, colorable.NewColorableStdout())
	case "yaml":
		b, err := yaml.Marshal(r)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "json":
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		return outputNotSupportedErr(output)
	}
	return nil
}

func WriteLcResultTo(r *types.LcResult, out io.Writer) (n int64, err error) {
	if r == nil {
		return 0, nil
	}

	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 0, '\t', 0)

	printf := func(format string, a ...interface{}) error {
		m, err := fmt.Fprintf(w, format, a...)
		n += int64(m)
		return err
	}

	s := reflect.ValueOf(r).Elem()
	s = s.FieldByName("LcArtifact")
	typeOfT := s.Type()

	for i, l := 0, s.NumField(); i < l; i++ {
		f := s.Field(i)
		if key, ok := typeOfT.Field(i).Tag.Lookup("vcn"); ok {
			var value string
			switch key {
			case "Size":
				if size, ok := f.Interface().(uint64); ok && size > 0 {
					value = humanize.Bytes(size)
				}
			case "Metadata":
				if metadata, ok := f.Interface().(api.Metadata); ok {
					for k, v := range metadata {
						if v == "" {
							continue
						}
						if vv, err := json.MarshalIndent(v, "\t", "    "); err == nil {
							value += fmt.Sprintf("\n\t\t%s=%s", k, string(vv))
						}
					}
					value = strings.TrimPrefix(value, "\n")
				}
			case "Apikey revoked":
				if f.IsZero() {
					value = color.New(meta.StyleWarning()).Sprintf("not available")
				} else {
					if revoked, ok := f.Interface().(*time.Time); ok {
						if revoked.IsZero() {
							value = color.New(meta.StyleAffordance()).Sprintf("no")
						} else {
							value = color.New(meta.StyleError()).Sprintf(revoked.Format(time.UnixDate))
						}
					}
				}
			case "Attachments":
				if attachments, ok := f.Interface().([]api.Attachment); ok {
					for _, attach := range attachments {
						if attach.Filename == "" {
							continue
						}
						v := reflect.ValueOf(attach)

						for i := 0; i < v.NumField(); i++ {
							sep := "\n\t  %s:\t%s"
							if i == 0 {
								sep = "\n\t- %s:\t%s"
							}
							value += fmt.Sprintf(sep, v.Type().Field(i).Name, v.Field(i).String())
						}
					}
				}
			case "Status":
				err = printf("Status:\t%s\n", meta.StatusNameStyled(r.Status))
				if err != nil {
					return
				}
			default:
				value = fmt.Sprintf("%s", f.Interface())
			}
			if value != "" {
				err = printf("%s:\t%s\n", key, value)
				if err != nil {
					return
				}
			}
		}
	}

	// here extra data when --verbose flag is provided
	if r.Verbose != nil {
		err = printf("\nAdditional details:\n")
		if err != nil {
			return
		}
		s = reflect.ValueOf(r.Verbose).Elem()
		typeOfT = s.Type()
		for i, l := 0, s.NumField(); i < l; i++ {
			if key, ok := typeOfT.Field(i).Tag.Lookup("vcn"); ok {
				switch key {
				case "LedgerName":
					err = printf("Ledger Name:\t%s\n", r.Verbose.LedgerName)
					if err != nil {
						return
					}
				case "LocalSID":
					err = printf("Local SignerID:\t%s\n", r.Verbose.LocalSID)
					if err != nil {
						return
					}
				case "ApiKey":
					err = printf("Api-key:\t%s\n", r.Verbose.ApiKey)
					if err != nil {
						return
					}
				}
			}
		}
	}

	for _, e := range r.Errors {
		err = printf("Error:\t%s\n", color.New(meta.StyleError()).Sprintf(e.Error()))
		if err != nil {
			return
		}
	}

	return n, w.Flush()
}

func PrintLcSlice(output string, rs []*types.LcResult) error {
	switch output {
	case "":
		for _, r := range rs {
			WriteLcResultTo(r, colorable.NewColorableStdout())
			fmt.Println()
		}
	case "yaml":
		b, err := yaml.Marshal(rs)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "json":
		b, err := json.MarshalIndent(rs, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		return outputNotSupportedErr(output)
	}
	return nil
}
