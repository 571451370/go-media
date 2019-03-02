package x86das

import "github.com/qeedquan/go-media/debug/xed"

const (
	REG_NIP = xed.REG_ZMM_LAST + iota
	REG_NSP
	REG_NAX
	REG_NEX
	REG_NDX
	REG_NBX
	REG_NDI
	REG_NSI
)

type Mach struct {
	Mode  xed.MachineMode
	Width xed.AddressWidth
	Reg   [256]uint64
}

func (m *Mach) ReadReg(reg xed.Reg) uint64 {
	return 0
}

func (m *Mach) ReadMem(addr uint64, size int) uint64 {
	return 0
}

func (m *Mach) ReadBuffer(addr uint64, buf []byte) {
	for i := range buf {
		buf[i] = byte(m.ReadMem(addr+uint64(i), 1))
	}
}

func (m *Mach) WriteReg(reg xed.Reg, val uint64) {
}

func (m *Mach) WriteMem(addr, val uint64, size int) {
}

func (m *Mach) WriteBuffer(addr uint64, buf []byte) {
	for i := range buf {
		m.WriteMem(addr+uint64(i), uint64(buf[i]), 1)
	}
}

func (m *Mach) Step() error {
	var inst xed.DecodedInst
	m.fetch(&inst)
	return m.Op(&inst)
}

func (m *Mach) Op(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) fetch(inst *xed.DecodedInst) {
	var buf [16]byte
	m.ReadBuffer(m.ReadReg(REG_NIP), buf[:])

	inst.Zero()
	inst.SetMode(m.Mode, m.Width)
	inst.Decode(buf[:])
}
