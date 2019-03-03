package x86das

import "github.com/qeedquan/go-media/debug/xed"

const (
	REG_NIP xed.Reg = xed.REG_ZMM_LAST + iota
	REG_NSP
	REG_NAX
	REG_NCX
	REG_NDX
	REG_NBX
	REG_NSI
	REG_NDI
)

type Mach struct {
	Mode  xed.MachineMode
	Width xed.AddressWidth
	Reg   [256]uint64
}

func (m *Mach) ReadReg(reg xed.Reg) uint64 {
	const (
		mask8  = 0xff
		mask16 = 0xffff
		mask32 = 0xffffffff
		mask64 = 0xffffffffffffffff
	)
	var mask uint64
	switch m.Width {
	case xed.ADDRESS_WIDTH_16b:
		mask = mask16
	case xed.ADDRESS_WIDTH_32b:
		mask = mask32
	case xed.ADDRESS_WIDTH_64b:
		mask = mask64
	default:
		panic("invalid address width")
	}

	var idx int
	var shift uint
	switch reg {
	case REG_NIP:
		idx = 0
	case REG_NSP:
		idx = 1
	case REG_NAX:
		idx = 2
	case REG_NCX:
		idx = 3
	case REG_NDX:
		idx = 4
	case REG_NBX:
		idx = 5
	case REG_NSI:
		idx = 6
	case REG_NDI:
		idx = 7

	case xed.REG_AL:
		idx = 2
		mask = mask8
	case xed.REG_AH:
		idx = 2
		mask = mask8
		shift = 8
	case xed.REG_AX:
		idx = 2
		mask = mask16
	case xed.REG_EAX:
		idx = 2
		mask = mask32
	case xed.REG_RAX:
		idx = 2
		mask = mask64
	}
	return (m.Reg[idx] >> shift) & mask
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
