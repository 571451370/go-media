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
	Magic    uint16
	Cblp     uint16
	Cp       uint16
	Crlc     uint16
	Cparhdr  uint16
	MinAlloc uint16
	MaxAlloc uint16
	SS       uint16
	SP       uint16
	Checksum uint16
	IP       uint16
	CS       uint16
	Lfarlc   uint16
	Ovno     uint16
	_        [4]uint8
	OEMID    uint16
	OEMInfo  uint16
	_        [10]uint16
	Lfanew   uint32
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

	dh := DOSHeader{Magic: 0x5a4d}
	binary.Write(b, binary.LittleEndian, &dh)

	return b.Flush()
}
