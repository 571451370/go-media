package peutil

import (
	"bufio"
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"

	"github.com/qeedquan/go-media/xio"
)

const (
	IMAGE_SCN_TYPE_NO_PAD            = 0x00000008
	IMAGE_SCN_CNT_CODE               = 0x00000020
	IMAGE_SCN_CNT_INITIALIZED_DATA   = 0x00000040
	IMAGE_SCN_CNT_UNINITIALIZED_DATA = 0x00000080
	IMAGE_SCN_LNK_OTHER              = 0x00000100
	IMAGE_SCN_LNK_INFO               = 0x00000200
	IMAGE_SCN_LNK_REMOVE             = 0x00000800
	IMAGE_SCN_LNK_COMDAT             = 0x00001000
	IMAGE_SCN_GPREL                  = 0x00008000
	IMAGE_SCN_MEM_PURGEABLE          = 0x00020000
	IMAGE_SCN_MEM_16BIT              = 0x00020000
	IMAGE_SCN_MEM_LOCKED             = 0x00040000
	IMAGE_SCN_MEM_PRELOAD            = 0x00080000
	IMAGE_SCN_ALIGN_1BYTES           = 0x00100000
	IMAGE_SCN_ALIGN_2BYTES           = 0x00200000
	IMAGE_SCN_ALIGN_8BYTES           = 0x00400000
	IMAGE_SCN_ALIGN_16BYTES          = 0x00500000
	IMAGE_SCN_ALIGN_32BYTES          = 0x00600000
	IMAGE_SCN_ALIGN_64BYTES          = 0x00700000
	IMAGE_SCN_ALIGN_128BYTES         = 0x00800000
	IMAGE_SCN_ALIGN_256BYTES         = 0x00900000
	IMAGE_SCN_ALIGN_512BYTES         = 0x00A00000
	IMAGE_SCN_ALIGN_1024BYTES        = 0x00B00000
	IMAGE_SCN_ALIGN_2048BYTES        = 0x00C00000
	IMAGE_SCN_ALIGN_4096BYTES        = 0x00D00000
	IMAGE_SCN_ALIGN_8192BYTES        = 0x00E00000
	IMAGE_SCN_LNK_NRELOC_OVFL        = 0x01000000
	IMAGE_SCN_MEM_DISCARDABLE        = 0x02000000
	IMAGE_SCN_MEM_NOT_CACHED         = 0x04000000
	IMAGE_SCN_MEM_NOT_PAGED          = 0x08000000
	IMAGE_SCN_MEM_SHARED             = 0x10000000
	IMAGE_SCN_MEM_EXECUTE            = 0x20000000
	IMAGE_SCN_MEM_READ               = 0x40000000
	IMAGE_SCN_MEM_WRITE              = 0x80000000
)

const (
	IMAGE_FILE_RELOCS_STRIPPED         = 0x0001
	IMAGE_FILE_EXECUTABLE_IMAGE        = 0x0002
	IMAGE_FILE_LINE_NUMS_STRIPPED      = 0x0004
	IMAGE_FILE_LOCAL_SYMS_STRIPPED     = 0x0008
	IMAGE_FILE_AGGRESSIVE_WS_TRIM      = 0x0010
	IMAGE_FILE_LARGE_ADDRESS_AWARE     = 0x0020
	IMAGE_FILE_BYTES_REVERSED_LO       = 0x0080
	IMAGE_FILE_32BIT_MACHINE           = 0x0100
	IMAGE_FILE_DEBUG_STRIPPED          = 0x0200
	IMAGE_FILE_REMOVABLE_RUN_FROM_SWAP = 0x0400
	IMAGE_FILE_NET_RUN_FROM_SWAP       = 0x0800
	IMAGE_FILE_SYSTEM                  = 0x1000
	IMAGE_FILE_DLL                     = 0x2000
	IMAGE_FILE_UP_SYSTEM_ONLY          = 0x4000
	IMAGE_FILE_BYTES_REVERSED_HI       = 0x8000
)

const (
	IMAGE_SUBSYSTEM_UNKNOWN                  = 0
	IMAGE_SUBSYSTEM_NATIVE                   = 1
	IMAGE_SUBSYSTEM_WINDOWS_GUI              = 2
	IMAGE_SUBSYSTEM_WINDOWS_CUI              = 3
	IMAGE_SUBSYSTEM_OS2_CUI                  = 5
	IMAGE_SUBSYSTEM_POSIX_CUI                = 7
	IMAGE_SUBSYSTEM_NATIVE_WINDOWS           = 8
	IMAGE_SUBSYSTEM_WINDOWS_CE_GUI           = 9
	IMAGE_SUBSYSTEM_EFI_APPLICATION          = 10
	IMAGE_SUBSYSTEM_EFI_BOOT_SERVICE_DRIVER  = 11
	IMAGE_SUBSYSTEM_EFI_RUNTIME_DRIVER       = 12
	IMAGE_SUBSYSTEM_EFI_ROM                  = 13
	IMAGE_SUBSYSTEM_XBOX                     = 14
	IMAGE_SUBSYSTEM_WINDOWS_BOOT_APPLICATION = 16
)

const (
	IMAGE_DLLCHARACTERISTICS_HIGH_ENTROPY_VA       = 0x0020
	IMAGE_DLLCHARACTERISTICS_DYNAMIC_BASE          = 0x0040
	IMAGE_DLLCHARACTERISTICS_FORCE_INTEGRITY       = 0x0080
	IMAGE_DLLCHARACTERISTICS_NX_COMPAT             = 0x0100
	IMAGE_DLLCHARACTERISTICS_NO_ISOLATION          = 0x0200
	IMAGE_DLLCHARACTERISTICS_NO_SEH                = 0x0400
	IMAGE_DLLCHARACTERISTICS_NO_BIND               = 0x0800
	IMAGE_DLLCHARACTERISTICS_APPCONTAINER          = 0x1000
	IMAGE_DLLCHARACTERISTICS_WDM_DRIVER            = 0x2000
	IMAGE_DLLCHARACTERISTICS_GUARD_CF              = 0x4000
	IMAGE_DLLCHARACTERISTICS_TERMINAL_SERVER_AWARE = 0x8000
)

const (
	IMAGE_DEBUG_TYPE_UNKNOWN       = 0
	IMAGE_DEBUG_TYPE_COFF          = 1
	IMAGE_DEBUG_TYPE_CODEVIEW      = 2
	IMAGE_DEBUG_TYPE_FPO           = 3
	IMAGE_DEBUG_TYPE_MISC          = 4
	IMAGE_DEBUG_TYPE_EXCEPTION     = 5
	IMAGE_DEBUG_TYPE_FIXUP         = 6
	IMAGE_DEBUG_TYPE_OMAP_TO_SRC   = 7
	IMAGE_DEBUG_TYPE_OMAP_FROM_SRC = 8
	IMAGE_DEBUG_TYPE_BORLAND       = 9
	IMAGE_DEBUG_TYPE_RESERVED10    = 10
	IMAGE_DEBUG_TYPE_CLSID         = 11
	IMAGE_DEBUG_TYPE_REPRO         = 16
)

const (
	IMAGE_REL_AMD64_ABSOLUTE = 0x0000
	IMAGE_REL_AMD64_ADDR64   = 0x0001
	IMAGE_REL_AMD64_ADDR32   = 0x0002
	IMAGE_REL_AMD64_ADDR32NB = 0x0003
	IMAGE_REL_AMD64_REL32    = 0x0004
	IMAGE_REL_AMD64_REL32_1  = 0x0005
	IMAGE_REL_AMD64_REL32_2  = 0x0006
	IMAGE_REL_AMD64_REL32_3  = 0x0007
	IMAGE_REL_AMD64_REL32_4  = 0x0008
	IMAGE_REL_AMD64_REL32_5  = 0x0009
	IMAGE_REL_AMD64_SECTION  = 0x000A
	IMAGE_REL_AMD64_SECREL   = 0x000B
	IMAGE_REL_AMD64_SECREL7  = 0x000C
	IMAGE_REL_AMD64_TOKEN    = 0x000D
	IMAGE_REL_AMD64_SREL32   = 0x000E
	IMAGE_REL_AMD64_PAIR     = 0x000F
	IMAGE_REL_AMD64_SSPAN32  = 0x0010
)

const (
	IMAGE_REL_I386_ABSOLUTE = 0x0000
	IMAGE_REL_I386_DIR16    = 0x0001
	IMAGE_REL_I386_REL16    = 0x0002
	IMAGE_REL_I386_DIR32    = 0x0006
	IMAGE_REL_I386_DIR32NB  = 0x0007
	IMAGE_REL_I386_SEG12    = 0x0009
	IMAGE_REL_I386_SECTION  = 0x000A
	IMAGE_REL_I386_SECREL   = 0x000B
	IMAGE_REL_I386_TOKEN    = 0x000C
	IMAGE_REL_I386_SECREL7  = 0x000D
	IMAGE_REL_I386_REL32    = 0x0014
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
	DOSHeader DOSHeader
	DOSStub   []byte
	Strings   []string
	r         io.ReaderAt
}

// these values are common across many exe files
var DOSHdr = DOSHeader{
	Magic:      0x5a4d,
	LastSize:   0x90,
	NumBlocks:  0x03,
	HeaderSize: 0x04,
	MaxAlloc:   0xffff,
	SP:         0xb8,
	RelocPos:   0x40,
	LFANew:     0x40 + uint32(len(DOSStub)),
}

var DOSStub = []byte{
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
	0x54, 0x68, 0x69, 0x73, 0x20, 0x70, 0x72, 0x6F,
	0x67, 0x72, 0x61, 0x6D, 0x20, 0x63, 0x61, 0x6E,
	0x6E, 0x6F, 0x74, 0x20, 0x62, 0x65, 0x20, 0x72,
	0x75, 0x6E, 0x20, 0x69, 0x6E, 0x20, 0x44, 0x4F,
	0x53, 0x20, 0x6D, 0x6F, 0x64, 0x65, 0x2E, 0x0D,
	0x0D, 0x0A, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00,
}

func Open(name string) (*File, error) {
	p, err := pe.Open(name)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return newFile(p, f)
}

func NewFile(r io.ReaderAt) (*File, error) {
	p, err := pe.NewFile(r)
	if err != nil {
		return nil, err
	}
	return newFile(p, r)
}

func newFile(p *pe.File, r io.ReaderAt) (*File, error) {
	f := &File{
		File:      p,
		DOSHeader: DOSHdr,
		DOSStub:   DOSStub,
	}
	for i := 4; i < len(f.StringTable); {
		str, err := f.StringTable.String(uint32(i))
		if err != nil {
			break
		}
		f.Strings = append(f.Strings, str)
		i += len(str) + 1
	}

	var dh DOSHeader
	sr := io.NewSectionReader(r, 0, math.MaxInt32)
	err := binary.Read(sr, binary.LittleEndian, &dh)
	if err == nil {
		f.DOSHeader = dh
	}
	stub := make([]byte, f.DOSHeader.LFANew-0x40)
	_, err = io.ReadAtLeast(sr, stub, len(stub))
	if err == nil {
		f.DOSStub = stub
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

func (f *File) sectionHeader32(s *pe.SectionHeader) pe.SectionHeader32 {
	h := pe.SectionHeader32{
		VirtualSize:          s.VirtualSize,
		VirtualAddress:       s.VirtualAddress,
		SizeOfRawData:        s.Size,
		PointerToRawData:     s.Offset,
		PointerToRelocations: s.PointerToRelocations,
		PointerToLineNumbers: s.PointerToLineNumbers,
		NumberOfRelocations:  s.NumberOfRelocations,
		NumberOfLineNumbers:  s.NumberOfLineNumbers,
		Characteristics:      s.Characteristics,
	}

	name := s.Name
	if len(s.Name) > len(h.Name) {
		n := bytes.Index(f.StringTable, []byte(s.Name))
		if n >= 0 {
			name = fmt.Sprintf("/%d", n+4)
		}
	}
	copy(h.Name[:], name[:])
	return h
}

func readStrz(b []byte) string {
	for i := range b {
		if b[i] == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

func (f *File) DuplicateSection(name string) (*pe.Section, error) {
	return f.DuplicateRawSection(f.Section(name))
}

func (f *File) DuplicateRawSection(s *pe.Section) (*pe.Section, error) {
	if s == nil {
		return nil, nil
	}
	p := &pe.Section{
		SectionHeader: s.SectionHeader,
		Relocs:        make([]pe.Reloc, len(s.Relocs)),
	}
	copy(p.Relocs, s.Relocs)
	data, err := s.Data()
	if err != nil {
		return p, err
	}
	p.ReaderAt = xio.NewBuffer(data)
	return p, nil
}

func Format(f *File, w io.Writer) error {
	b := bufio.NewWriter(w)

	binary.Write(b, binary.LittleEndian, &f.DOSHeader)
	b.Write(f.DOSStub)

	size := int(reflect.TypeOf(f.DOSHeader).Size()) + len(f.DOSStub)
	sizeOfHeaders := size

	peSig := [...]byte{'P', 'E', 0x00, 0x00}
	binary.Write(b, binary.LittleEndian, peSig)
	binary.Write(b, binary.LittleEndian, &f.FileHeader)
	size += len(peSig) + int(reflect.TypeOf(f.FileHeader).Size())

	switch oh := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		binary.Write(b, binary.LittleEndian, oh)
		sizeOfHeaders = int(oh.SizeOfHeaders)
		size += int(reflect.TypeOf(*oh).Size())
	case *pe.OptionalHeader64:
		binary.Write(b, binary.LittleEndian, oh)
		sizeOfHeaders = int(oh.SizeOfHeaders)
		size += int(reflect.TypeOf(*oh).Size())
	}

	for _, s := range f.Sections {
		sh := f.sectionHeader32(&s.SectionHeader)
		binary.Write(b, binary.LittleEndian, &sh)
		size += int(reflect.TypeOf(sh).Size())
	}
	if size < sizeOfHeaders {
		pad := make([]byte, sizeOfHeaders-size)
		b.Write(pad)
	}

	for _, s := range f.Sections {
		data, err := s.Data()
		if err != nil {
			return err
		}
		b.Write(data)
	}

	return b.Flush()
}
