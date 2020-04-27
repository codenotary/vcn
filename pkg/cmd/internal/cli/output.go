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
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"

	"gopkg.in/yaml.v2"

	"github.com/dustin/go-humanize"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
	"github.com/vchain-us/vcn/pkg/meta"
)

func WriteResultTo(r *types.Result, out io.Writer) (n int64, err error) {
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
	s = s.FieldByName("ArtifactResponse")
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
			case "Signer":
				// todo(leogr): this will not happen anymore with the new platform APIs.
				// Still retained to accommodate future improvements.
				if f.Interface() != r.Verification.SignerID() {
					value = fmt.Sprintf("%s", f.Interface())
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

	if bv := r.Verification; bv != nil {
		if key := bv.SignerID(); key != "" {
			err = printf("SignerID:\t%s\n", bv.SignerID())
			if err != nil {
				return
			}
		}
		if bv.Level != 0 {
			err = printf("Level:\t%s\n", bv.Level.String())
			if err != nil {
				return
			}
		}
		if date := bv.Date(); date != "" {
			err = printf("Date:\t%s\n", date)
			if err != nil {
				return
			}
		}
		err = printf("Status:\t%s\n", meta.StatusNameStyled(r.Verification.Status))
		if err != nil {
			return
		}
	}

	for _, e := range r.Errors {
		c, s := meta.StyleError()
		err = printf("Error:\t%s\n", color.New(c, s).Sprintf(e.Error()))
		if err != nil {
			return
		}
	}

	return n, w.Flush()
}

func Print(output string, r *types.Result) error {
	switch output {
	case "":
		WriteResultTo(r, colorable.NewColorableStdout())
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

func PrintSlice(output string, rs []types.Result) error {
	switch output {
	case "":
		for i := range rs {
			_, err := WriteResultTo(&rs[i], colorable.NewColorableStdout())
			if err != nil {
				return err
			}
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

func PrintList(output string, artifacts []api.ArtifactResponse) error {
	switch output {
	case "":
		for i := range artifacts {
			fmt.Print(artifacts[i])
		}
	case "yaml":
		b, err := yaml.Marshal(artifacts)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "json":
		b, err := json.MarshalIndent(artifacts, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		return outputNotSupportedErr(output)
	}
	return nil
}

func PrintError(output string, err *types.Error) error {
	if err == nil {
		return nil
	}
	switch output {
	case "":
		fmt.Fprintf(os.Stderr, "\nError: %s\n", err)
	case "yaml":
		b, err := yaml.Marshal(err)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "json":
		b, err := json.MarshalIndent(err, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		return outputNotSupportedErr(output)
	}
	return nil
}

func PrintWarning(output, message string) error {
	switch output {
	case "":
		fallthrough
	case "yaml":
		fallthrough
	case "json":
		fmt.Fprintf(os.Stderr, "\nWarning: %s\n", message)
	default:
		return outputNotSupportedErr(output)
	}
	return nil
}

func outputNotSupportedErr(output string) error {
	return fmt.Errorf("output format not supported: %s", output)
}

func PrintObjects(output string, out interface{}) error {
	switch output {
	case "":
		fallthrough
	case "yaml":
		b, err := yaml.Marshal(out)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "json":
		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		return outputNotSupportedErr(output)
	}
	return nil
}
