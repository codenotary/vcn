/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package sniff

import (
	"errors"
	"os"
)

type Data struct {
	Format   string `json:"format"`
	Type     string `json:"type"`
	Platform string `json:"platform"`
	Arch     string `json:"arch"`
	X64      bool   `json:"x64"`
}

func (d Data) ContentType() string {
	switch {
	case d.Platform == PlatformMachO:
		return "application/x-mach-binary"
	case d.Platform == Platform_PE:
		return "application/x-dosexec"
	case d.Format == "ELF":
		return "application/x-executable"
	}
	return "application/octet-stream"
}

var sniffers = []func(*os.File) (*Data, error){
	ELF,
	PE,
	MachO,
}

func File(file *os.File) (*Data, error) {

	for _, sniffer := range sniffers {
		if d, e := sniffer(file); e == nil {
			return d, nil
		}
	}

	return nil, errors.New("Nothing found")
}
