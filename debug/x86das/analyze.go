package x86das

import (
	"debug/pe"
	"fmt"
	"io"
	"os"

	"github.com/qeedquan/go-media/debug/peutil"
	"github.com/qeedquan/go-media/debug/xed"
)

const (
	S_R = 1 << iota
	S_W
	S_X
)

type Prog struct {
	Format   interface{}
	Mode     xed.MachineMode
	Width    xed.AddressWidth
	Entry    uint64
	Sections []*Section
}

type Section struct {
	Name  string
	Addr  uint64
	Size  uint64
	Align uint64
	Flag  uint64
	Data  []byte
}

type BasicBlock struct {
}

func AnalyzeFile(name string) (*Prog, error) {
	pe, err := peutil.Open(name)
	if err == nil {
		return readPE(pe)
	}
	if os.IsNotExist(err) || os.IsPermission(err) {
		return nil, err
	}
	return nil, fmt.Errorf("unknown file format")
}

func AnalyzeExec(r io.ReaderAt) (*Prog, error) {
	pe, err := peutil.NewFile(r)
	if err == nil {
		return readPE(pe)
	}
	return nil, fmt.Errorf("unknown file format")
}

func readPE(f *peutil.File) (*Prog, error) {
	var mode xed.MachineMode
	var width xed.AddressWidth
	switch f.Machine {
	case pe.IMAGE_FILE_MACHINE_AMD64:
		mode = xed.MACHINE_MODE_LONG_64
		width = xed.ADDRESS_WIDTH_64b
	case pe.IMAGE_FILE_MACHINE_I386:
		mode = xed.MACHINE_MODE_LONG_COMPAT_32
		width = xed.ADDRESS_WIDTH_32b
	default:
		return nil, fmt.Errorf("unsupported machine type %d", f.Machine)
	}

	var align uint64
	var base uint64
	var entry uint64
	switch h := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		align = uint64(h.SectionAlignment)
		base = uint64(h.ImageBase)
		entry = uint64(h.AddressOfEntryPoint) + base
	case *pe.OptionalHeader64:
		align = uint64(h.SectionAlignment)
		base = uint64(h.ImageBase)
		entry = uint64(h.AddressOfEntryPoint) + base
	}

	var sections []*Section
	for _, s := range f.Sections {
		var flag uint64
		if s.Characteristics&peutil.IMAGE_SCN_MEM_READ != 0 {
			flag |= S_R
		}
		if s.Characteristics&peutil.IMAGE_SCN_MEM_WRITE != 0 {
			flag |= S_W
		}
		if s.Characteristics&peutil.IMAGE_SCN_MEM_EXECUTE != 0 {
			flag |= S_X
		}

		data, err := s.Data()
		if err != nil {
			return nil, err
		}
		section := &Section{
			Name:  s.Name,
			Addr:  uint64(s.VirtualAddress) + base,
			Size:  uint64(s.VirtualSize),
			Align: align,
			Flag:  flag,
			Data:  data,
		}
		sections = append(sections, section)
	}

	return &Prog{
		Format:   f,
		Mode:     mode,
		Width:    width,
		Entry:    entry,
		Sections: sections,
	}, nil
}
