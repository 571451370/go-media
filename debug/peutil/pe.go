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

type ImportDescriptor struct {
	OriginalFirstThunk uint32
	TimeDateStamp      uint32
	ForwarderChain     uint32
	Name               uint32
	FirstThunk         uint32
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
	DllNameOff       uint64
	NameOff          uint64
	OriginalThunkOff uint64
	ThunkOff         uint64
	Auxillary        interface{}
}

type Section struct {
	*pe.Section
	Data []byte
}

type File struct {
	*pe.File
	RawSizeOfHeaders int
	SizeOfHeaders    int
	ArchSize         int
	DOSHeader        DOSHeader
	DOSStub          []byte
	Sections         []*Section
	Strings          []string
	r                io.ReaderAt
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
	sr := io.NewSectionReader(r, 0, math.MaxInt64)
	err := binary.Read(sr, binary.LittleEndian, &dh)
	if err == nil {
		f.DOSHeader = dh
	}
	stub := make([]byte, f.DOSHeader.LFANew-0x40)
	_, err = io.ReadAtLeast(sr, stub, len(stub))
	if err == nil {
		f.DOSStub = stub
	}

	for _, s := range f.File.Sections {
		p := &Section{}
		p.Section = s
		p.Data, err = s.Data()
		if err != nil {
			return nil, err
		}
		f.Sections = append(f.Sections, p)
	}
	f.RawSizeOfHeaders, f.SizeOfHeaders = f.calcSizeOfHeaders()

	switch f.Machine {
	case pe.IMAGE_FILE_MACHINE_AMD64:
		f.ArchSize = 64
	default:
		f.ArchSize = 32
	}

	return f, nil
}

func (f *File) RemoveDLLImport(dllName string) error {
	idd, dt, _, err := f.ReadImportTable()
	if err != nil {
		return err
	}

	var nt []ImportDescriptor
	for i := range dt {
		if f.readStrzVA(uint64(dt[i].Name)) != dllName {
			nt = append(nt, dt[i])
		}
	}
	if len(nt) == len(dt) {
		return fmt.Errorf("no dll import %q", dllName)
	}

	off := uint64(idd.VirtualAddress)
	for i := range nt {
		f.putVA(off, &nt[i])
		off += uint64(reflect.TypeOf(nt[i]).Size())
	}
	f.putVA(off, &ImportDescriptor{})

	return nil
}

func (f *File) RemoveSection(name string) error {
	var (
		found        bool
		size, offset uint32
	)
	for i, s := range f.Sections {
		if s.Name == name {
			size = s.Size
			offset = s.Offset
			copy(f.Sections[i:], f.Sections[i+1:])
			f.Sections = f.Sections[:len(f.Sections)-1]
			f.RawSizeOfHeaders, f.SizeOfHeaders = f.calcSizeOfHeaders()
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("section %q does not exist", name)
	}

	for _, s := range f.Sections {
		if s.Offset >= offset {
			s.Offset -= size
		}
	}
	return nil
}

func (f *File) calcSizeOfHeaders() (rawSizeOfHeaders, sizeOfHeaders int) {
	rawSizeOfHeaders += int(reflect.TypeOf(f.DOSHeader).Size())
	rawSizeOfHeaders += len(f.DOSStub)
	// PE signature
	rawSizeOfHeaders += 4
	rawSizeOfHeaders += int(reflect.TypeOf(f.FileHeader).Size())

	switch oh := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		sizeOfHeaders = int(oh.SizeOfHeaders)
		rawSizeOfHeaders += int(reflect.TypeOf(*oh).Size())
	case *pe.OptionalHeader64:
		sizeOfHeaders = int(oh.SizeOfHeaders)
		rawSizeOfHeaders += int(reflect.TypeOf(*oh).Size())
	}

	for _, s := range f.Sections {
		sh := f.sectionHeader32(&s.SectionHeader)
		rawSizeOfHeaders += int(reflect.TypeOf(sh).Size())
	}
	if rawSizeOfHeaders > sizeOfHeaders {
		sizeOfHeaders = rawSizeOfHeaders
	}
	return
}

func (f *File) ReadImportTable() (*pe.DataDirectory, []ImportDescriptor, []Symbol, error) {
	idd, err := f.readDataDirectory(pe.IMAGE_DIRECTORY_ENTRY_IMPORT, nil)
	if err != nil {
		return idd, nil, nil, err
	}
	if idd == nil {
		return idd, nil, nil, fmt.Errorf("no import table")
	}

	_, data := f.sectionVA(uint64(idd.VirtualAddress))
	if len(data) == 0 {
		return idd, nil, nil, fmt.Errorf("no import table")
	}

	var (
		dt []ImportDescriptor
		d  ImportDescriptor
	)
	r := bytes.NewReader(data)
	for {
		err := binary.Read(r, binary.LittleEndian, &d)
		if err != nil {
			break
		}
		if d.OriginalFirstThunk == 0 {
			break
		}
		dt = append(dt, d)
	}

	var s []Symbol
	for _, d := range dt {
		dll := f.readStrzVA(uint64(d.Name))
		_, ft := f.sectionVA(uint64(d.OriginalFirstThunk))
		ftoff := uint64(d.OriginalFirstThunk)
		thoff := uint64(d.FirstThunk)
		for len(ft) > 0 {
			var (
				ftsz uint64
				na   uint64
				mask uint64
			)
			switch f.Machine {
			case pe.IMAGE_FILE_MACHINE_AMD64:
				ftsz = 8
				na = binary.LittleEndian.Uint64(ft)
				mask = 1 << 63
			default:
				ftsz = 4
				na = uint64(binary.LittleEndian.Uint32(ft))
				mask = 1 << 31
			}
			if uint64(na)&mask > 0 {
				panic("dynimport ordinals unimplemented")
			}
			if na == 0 {
				break
			}

			p := Symbol{
				Symbol: pe.Symbol{
					Name: f.readStrzVA(na + 2),
				},
				DllName:          dll,
				DllNameOff:       uint64(d.Name),
				NameOff:          na,
				OriginalThunkOff: ftoff,
				ThunkOff:         thoff,
				Auxillary:        d,
			}
			s = append(s, p)
			ft = ft[ftsz:]
			ftoff += ftsz
			thoff += ftsz
		}
	}

	return idd, dt, s, nil
}

func (f *File) ExportedSymbols() ([]Symbol, error) {
	var d ExportDirectory
	idd, err := f.readDataDirectory(pe.IMAGE_DIRECTORY_ENTRY_EXPORT, &d)
	if err != nil {
		return nil, err
	}
	if idd == nil {
		return nil, nil
	}

	_, fp := f.sectionVA(uint64(d.AddressOfFunctions))
	_, od := f.sectionVA(uint64(d.AddressOfNameOrdinals))
	_, na := f.sectionVA(uint64(d.AddressOfNames))
	no := uint32(0)
	if fp == nil || od == nil {
		return nil, nil
	}

	dllName := f.readStrzVA(uint64(d.Name))
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

		nameoff := va
		if len(na) < 4 {
			name = fmt.Sprintf("%s+%#x", dllName, va)
		} else {
			name = f.readStrzVA(uint64(binary.LittleEndian.Uint32(na)))
			nameoff = no
			na = na[4:]
			no += 4
		}
		p := Symbol{
			Symbol: pe.Symbol{
				Name: name,
			},
			DllName:          dllName,
			ForwardedAddress: fwd,
			NameOff:          uint64(nameoff),
			Auxillary:        idd,
		}
		s = append(s, p)
	}
	return s, nil
}

func (f *File) putByteVA(va uint64, val byte) error {
	s, d := f.sectionVA(va)
	if s == nil || len(d) == 0 {
		return fmt.Errorf("virtual address %#x does not exist", va)
	}
	d[0] = val
	return nil
}

func (f *File) putWordVA(va uint64, val uint64) error {
	switch f.ArchSize {
	case 64:
		return f.putVA(va, uint64(val))
	case 32:
		return f.putVA(va, uint32(val))
	default:
		panic("unsupported arch size")
	}
}

func (f *File) putVA(va uint64, val interface{}) error {
	var b []byte
	switch v := val.(type) {
	case uint8:
		b = make([]byte, 1)
		b[0] = v
	case uint16:
		b = make([]byte, 2)
		binary.LittleEndian.PutUint16(b, v)
	case uint32:
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(b, v)
	case uint64:
		b = make([]byte, 8)
		binary.LittleEndian.PutUint64(b, v)
	case *ImportDescriptor:
		p := new(bytes.Buffer)
		binary.Write(p, binary.LittleEndian, v)
		b = p.Bytes()
	default:
		panic(fmt.Errorf("unsupported type %T", v))
	}

	for i := range b {
		err := f.putByteVA(va+uint64(i), b[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *File) readStrzVA(va uint64) string {
	_, b := f.sectionVA(va)
	if b == nil {
		return ""
	}
	return readStrz(b)
}

func (f *File) sectionVA(va uint64) (*Section, []byte) {
	for _, s := range f.Sections {
		if uint64(s.VirtualAddress) <= va && va < uint64(s.VirtualAddress+s.VirtualSize) {
			return s, s.Data[va-uint64(s.VirtualAddress):]
		}
	}
	return nil, nil
}

func (f *File) readDataDirectory(index int, v interface{}) (*pe.DataDirectory, error) {
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
		return nil, nil
	}

	ds, data := f.sectionVA(uint64(idd.VirtualAddress))
	if ds == nil {
		return nil, nil
	}

	if v != nil {
		r := bytes.NewReader(data)
		err := binary.Read(r, binary.LittleEndian, v)
		if err != nil {
			return nil, err
		}
	}

	return &idd, nil
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

func Format(f *File, w io.Writer) error {
	b := bufio.NewWriter(w)

	binary.Write(b, binary.LittleEndian, &f.DOSHeader)
	b.Write(f.DOSStub)

	peSig := [...]byte{'P', 'E', 0x00, 0x00}
	binary.Write(b, binary.LittleEndian, peSig)
	binary.Write(b, binary.LittleEndian, &f.FileHeader)

	switch oh := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		binary.Write(b, binary.LittleEndian, oh)
	case *pe.OptionalHeader64:
		binary.Write(b, binary.LittleEndian, oh)
	}

	for _, s := range f.Sections {
		sh := f.sectionHeader32(&s.SectionHeader)
		binary.Write(b, binary.LittleEndian, &sh)
	}
	pad := make([]byte, f.SizeOfHeaders-f.RawSizeOfHeaders)
	b.Write(pad)

	for _, s := range f.Sections {
		b.Write(s.Data)
	}

	return b.Flush()
}
