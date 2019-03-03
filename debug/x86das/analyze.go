package x86das

import (
	"fmt"
	"io"

	"github.com/qeedquan/go-media/debug/peutil"
)

const (
	S_R = 1 << iota
	S_W
	S_X
)

type Prog struct {
	Format   interface{}
	Sections []*Section
}

type Section struct {
	Name string
	Addr uint64
	Flag uint64
	Data []byte
}

type BasicBlock struct {
}

func AnalyzeFile(name string) (*Prog, error) {
	pe, err := peutil.Open(name)
	if err == nil {
		return readPE(pe)
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
		section := &Section{
			Name: s.Name,
			Flag: flag,
		}
		sections = append(sections, section)
	}

	return &Prog{
		Format:   f,
		Sections: sections,
	}, nil
}
