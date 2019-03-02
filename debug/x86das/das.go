package x86das

import "github.com/qeedquan/go-media/debug/xed"

const (
	RIP = iota
)

type Mach struct {
	MachineMode    xed.MachineMode
	StackAddrWidth xed.AddressWidth
	Reg            [256]uint64
	Mem            []byte
}

func (m *Mach) Interpret(inst *xed.DecodedInst) {
}

func (m *Mach) Step() {
}
