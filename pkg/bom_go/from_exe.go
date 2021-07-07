/*
 * Copyright (c) 2021 CodeNotary, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package bom_go

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/vchain-us/vcn/pkg/bom_component"
)

// The logic is copied from 'go version' utility source: https://golang.org/src/cmd/go/internal/version/version.go

// The build info blob left by the linker is identified by
// a 16-byte header, consisting of buildInfoMagic (14 bytes),
// the binary's pointer size (1 byte),
// and whether the binary is big endian (1 byte).
var buildInfoMagic = []byte("\xff Go buildinf:")

// Components returns list of go packages used during the build
func exeComponents(x exe) ([]bom_component.Component, error) {
	// Read the first 64kB of text to find the build info blob.
	text := x.DataStart()
	data, err := x.ReadData(text, 64*1024)
	if err != nil {
		return nil, err
	}
	for ; !bytes.HasPrefix(data, buildInfoMagic); data = data[32:] {
		if len(data) < 32 {
			return nil, err
		}
	}
	// find where build info actually starts
	for ; !bytes.HasPrefix(data, buildInfoMagic); data = data[32:] {
		if len(data) < 32 {
			return nil, errors.New("no build info found")
		}
	}

	// Decode the blob.
	ptrSize := int(data[14])
	bigEndian := data[15] != 0
	var bo binary.ByteOrder
	if bigEndian {
		bo = binary.BigEndian
	} else {
		bo = binary.LittleEndian
	}
	var readPtr func([]byte) uint64
	if ptrSize == 4 {
		readPtr = func(b []byte) uint64 { return uint64(bo.Uint32(b)) }
	} else {
		readPtr = bo.Uint64
	}

	mod := readString(x, ptrSize, readPtr, readPtr(data[16+ptrSize:]))
	if len(mod) >= 33 && mod[len(mod)-17] == '\n' {
		// Strip module framing.
		mod = mod[16 : len(mod)-16]
	} else {
		return nil, errors.New("no build info found")
	}

	lines := strings.Split(mod, "\n")
	res := make([]bom_component.Component, 0, len(lines))
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if fields[0] == "dep" {
			var comp bom_component.Component
			switch len(fields) {
			default:
				comp.Hash, comp.HashType, err = goModHash(fields[3])
				if err != nil {
					return nil, fmt.Errorf("cannot decode hash: %w", err)
				}
				fallthrough
			case 3:
				comp.Version = fields[2]
				fallthrough
			case 2:
				comp.Name = fields[1]
			case 1:
				continue
			}
			res = append(res, comp)
		}
	}

	return res, nil
}

// readString returns the string at address addr in the executable x.
func readString(x exe, ptrSize int, readPtr func([]byte) uint64, addr uint64) string {
	hdr, err := x.ReadData(addr, uint64(2*ptrSize))
	if err != nil || len(hdr) < 2*ptrSize {
		return ""
	}
	dataAddr := readPtr(hdr)
	dataLen := readPtr(hdr[ptrSize:])
	data, err := x.ReadData(dataAddr, dataLen)
	if err != nil || uint64(len(data)) < dataLen {
		return ""
	}
	return string(data)
}
