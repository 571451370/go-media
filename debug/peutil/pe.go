package peutil

import (
	"bufio"
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
)

type DOSHeader struct {
	Magic      uint16 // MZ
	LastSize   uint16 // image size mod 512, number of bytes on last page
	NumBlocks  uint16 // number of 512-byte pages in images
	NumRelocs  uint16 // count of relocation entries
	HeaderSize uint16 // size of header in paragraphs
	MinAlloc   uint16 // min required memory
	MaxAlloc   uint16 // max required memory
	SS         uint16 // stack seg offset in load module
	SP         uint16 // initial sp value
	Checksum   uint16 // one complement sum of all word in exe file
	IP         uint16 // initial ip value
	CS         uint16 // cs offset in load module
	RelocPos   uint16 // offset of first reloc item
	NoOverlay  uint16 // overlay number
	_          [4]uint16
	OEMID      uint16
	OEMInfo    uint16
	_          [10]uint16
	LFANew     uint32 // offset to pe header in windows
}

type ExportDirectory struct {
	Characteristics       uint32
	TimeDateStamp         uint32
	MajorVersion          uint16
	MinorVersion          uint16
	Name                  uint32
	Base                  uint32
	NumberOfFunctions     uint32
	NumberOfNames         uint32
	AddressOfFunctions    uint32
	AddressOfNames        uint32
	AddressOfNameOrdinals uint32
}

type Symbol struct {
	pe.Symbol
	DllName          string
	ForwardedAddress uint32
}

type File struct {
	*pe.File
	Strings []string
}

var DOSStub = [...]byte{
	// push cs
	0x0E,
	// pop ds
	0x1F,
	// mov dx, 0xe
	0xBA, 0x0E, 0x00,
	// mov ah, 0x9
	0xB4, 0x09,
	// int 0x21
	0xCD, 0x21,
	// mov ax, 0x4c01
	0xB8, 0x01, 0x4C,
	// int 0x21
	0xCD, 0x21,
	// "This program cannot be run in DOS Mode"
	0x68, 0x69, 0x73, 0x20, 0x70, 0x72, 0x6F,
	0x67, 0x72, 0x61, 0x6D, 0x20, 0x63, 0x61,
	0x6E, 0x6E, 0x6F, 0x74, 0x20, 0x62, 0x65,
	0x20, 0x72, 0x75, 0x6E, 0x20, 0x69, 0x6E,
	0x20, 0x44, 0x4F, 0x53, 0x20, 0x6D, 0x6F,
	0x64, 0x65, 0x2E, 0x0D, 0x0D, 0x0A, 0x24,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func Open(name string) (*File, error) {
	p, err := pe.Open(name)
	if err != nil {
		return nil, err
	}

	f := &File{File: p}
	for i := 4; i < len(f.StringTable); {
		str, err := f.StringTable.String(uint32(i))
		if err != nil {
			break
		}
		f.Strings = append(f.Strings, str)
		i += len(str) + 1
	}
	return f, nil
}

func (f *File) ExportedSymbols() ([]Symbol, error) {
	var d ExportDirectory
	idd := f.readDataDirectory(pe.IMAGE_DIRECTORY_ENTRY_EXPORT, &d)
	if idd == nil {
		return nil, nil
	}

	_, fp := f.sectionVA(d.AddressOfFunctions)
	_, od := f.sectionVA(d.AddressOfNameOrdinals)
	_, na := f.sectionVA(d.AddressOfNames)
	if fp == nil || od == nil {
		return nil, nil
	}

	dllName := f.readStrzVA(d.Name)
	var s []Symbol
	for i := uint32(0); i < d.NumberOfFunctions && len(od) >= 4; i, od = i+1, od[2:] {
		fn := binary.LittleEndian.Uint16(od) * 2
		if fn >= uint16(len(fp)) {
			continue
		}

		var name string
		var fwd uint32

		va := binary.LittleEndian.Uint32(fp[fn:])
		if idd.VirtualAddress <= va && va < idd.VirtualAddress+idd.Size {
			fwd = va
		}

		if len(na) < 4 {
			name = fmt.Sprintf("%s+%#x", dllName, va)
		} else {
			name = f.readStrzVA(binary.LittleEndian.Uint32(na))
			na = na[4:]
		}
		p := Symbol{
			Symbol: pe.Symbol{
				Name: name,
			},
			DllName:          dllName,
			ForwardedAddress: fwd,
		}
		s = append(s, p)
	}
	return s, nil
}

func (f *File) readStrzVA(va uint32) string {
	_, b := f.sectionVA(va)
	if b == nil {
		return ""
	}
	return readStrz(b)
}

func (f *File) sectionVA(va uint32) (*pe.Section, []byte) {
	for _, s := range f.Sections {
		if s.VirtualAddress <= va && va < s.VirtualAddress+s.VirtualSize {
			d, err := s.Data()
			if err != nil {
				return s, nil
			}
			return s, d[va-s.VirtualAddress:]
		}
	}
	return nil, nil
}

func (f *File) readDataDirectory(index int, v interface{}) *pe.DataDirectory {
	var dirlen uint32
	var idd pe.DataDirectory
	switch h := f.OptionalHeader.(type) {
	case *pe.OptionalHeader64:
		dirlen = h.NumberOfRvaAndSizes
		idd = h.DataDirectory[index]
	case *pe.OptionalHeader32:
		dirlen = h.NumberOfRvaAndSizes
		idd = h.DataDirectory[index]
	}
	if dirlen < uint32(index)+1 {
		return nil
	}

	ds, data := f.sectionVA(idd.VirtualAddress)
	if ds == nil {
		return nil
	}

	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, v)
	if err != nil {
		return nil
	}

	return &idd
}

func readStrz(b []byte) string {
	for i := range b {
		if b[i] == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

func Format(f *File, w io.Writer) error {
	b := bufio.NewWriter(w)

	// these values are common across many exe files
	dh := DOSHeader{
		Magic:      0x5a4d,
		LastSize:   0x90,
		NumBlocks:  0x03,
		HeaderSize: 0x04,
		MaxAlloc:   0xffff,
	}
	binary.Write(b, binary.LittleEndian, &dh)
	binary.Write(b, binary.LittleEndian, &DOSStub)

	peSig := uint16(0x4550)
	binary.Write(b, binary.LittleEndian, &peSig)
	binary.Write(b, binary.LittleEndian, &f.FileHeader)

	return b.Flush()
}
