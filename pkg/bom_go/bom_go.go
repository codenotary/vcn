package bom_go

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	component "github.com/codenotary/vcn/pkg/bom_component"
)

// GoPackage implements Package interface
type GoPackage struct {
	file exe
}

// New returns new GoPackage object, or nil if filename doesn't referer to ELF, built from Go source
func New(filename string) *GoPackage {
	file, err := openExe(filename)
	if err != nil {
		return nil // not a ELF binary
	}
	if file.DataStart() == 0 {
		file.Close()
		return nil // cannot find build info
	}
	return &GoPackage{file: file}
}

func (p *GoPackage) Type() string {
	return "Go"
}

func (p *GoPackage) Close() {
	if p.file != nil {
		p.file.Close()
	}
}

// The logic is copied from 'go version' utility source: https://golang.org/src/cmd/go/internal/version/version.go

// The build info blob left by the linker is identified by
// a 16-byte header, consisting of buildInfoMagic (14 bytes),
// the binary's pointer size (1 byte),
// and whether the binary is big endian (1 byte).
var buildInfoMagic = []byte("\xff Go buildinf:")

// Components returns list of go packages used during the build
func (p *GoPackage) Components() ([]component.Component, error) {
	// Read the first 64kB of text to find the build info blob.
	text := p.file.DataStart()
	data, err := p.file.ReadData(text, 64*1024)
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
			return nil, fmt.Errorf("no build info found")
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

	mod := readString(p.file, ptrSize, readPtr, readPtr(data[16+ptrSize:]))
	if len(mod) >= 33 && mod[len(mod)-17] == '\n' {
		// Strip module framing.
		mod = mod[16 : len(mod)-16]
	} else {
		return nil, fmt.Errorf("no build info found")
	}

	lines := strings.Split(mod, "\n")
	res := make([]component.Component, 0, len(lines))
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if fields[0] == "dep" {
			var comp component.Component
			switch len(fields) {
			default:
				comp.Hash = fields[3]
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
