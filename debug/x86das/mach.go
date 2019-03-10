package x86das

import (
	"fmt"

	"github.com/qeedquan/go-media/debug/xed"
	"github.com/qeedquan/go-media/math/mathutil"
)

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

type Seg struct {
	Name string
	Virt uint64
	Phys uint64
	Size uint64
	Flag uint64
}

type Mach struct {
	Prog *Prog
	Page int
	Reg  [256]uint64
	Seg  []*Seg
	Mem  []byte
}

func (m *Mach) LoadProg(prog *Prog) error {
	var stackAddr uint64

	m.Prog = prog
	m.Page = 4096
	m.Seg = m.Seg[:0]
	for _, s := range prog.Sections {
		size := uint64(mathutil.Multiple(int(s.Size), m.Page))
		sg, err := m.Map(s.Addr, size, s.Flag)
		sg.Name = s.Name
		if err != nil {
			return err
		}
		copy(m.Mem[sg.Phys:], s.Data)
		if stackAddr < s.Addr+size {
			stackAddr = s.Addr + size
		}
	}
	sg, err := m.Map(stackAddr, 32*1024*1024, S_R|S_W|S_X)
	if err != nil {
		return err
	}
	m.WriteReg(REG_NSP, sg.Virt+sg.Size)
	m.WriteReg(REG_NIP, prog.Entry)
	return nil
}

func (m *Mach) Map(addr, size, flag uint64) (*Seg, error) {
	m.Seg = append(m.Seg, &Seg{
		Virt: addr,
		Phys: uint64(len(m.Mem)),
		Size: size,
		Flag: flag,
	})
	m.Mem = append(m.Mem, make([]byte, size)...)
	return m.Seg[len(m.Seg)-1], nil
}

func (m *Mach) PushN(val uint64, size int) {
	for i := 0; i < size; i++ {
		sp := m.ReadReg(REG_NSP) - 1
		m.WriteMem(sp, val&0xff, 1)
		m.WriteReg(REG_NSP, sp)
		val >>= 8
	}
}

func (m *Mach) PopN(size int) uint64 {
	var val uint64
	for i := 0; i < size; i++ {
		sp := m.ReadReg(REG_NSP)
		val |= uint64(m.ReadMem(sp, 1)) << (8 * uint(i))
		m.WriteReg(REG_NSP, sp+1)
	}
	return val
}

func (m *Mach) Push(val uint64) {
	switch m.Prog.Width {
	case xed.ADDRESS_WIDTH_16b:
		m.PushN(val, 2)
	case xed.ADDRESS_WIDTH_32b:
		m.PushN(val, 4)
	case xed.ADDRESS_WIDTH_64b:
		m.PushN(val, 8)
	default:
		panic("invalid address width")
	}
}

func (m *Mach) Pop() uint64 {
	switch m.Prog.Width {
	case xed.ADDRESS_WIDTH_16b:
		return m.PopN(2)
	case xed.ADDRESS_WIDTH_32b:
		return m.PopN(4)
	case xed.ADDRESS_WIDTH_64b:
		return m.PopN(8)
	default:
		panic("invalid address width")
	}
}

func (m *Mach) decodeReg(reg xed.Reg) (idx int, shift, mask uint64) {
	const (
		mask8  = 0xff
		mask16 = 0xffff
		mask32 = 0xffffffff
		mask64 = 0xffffffffffffffff
	)

	switch m.Prog.Width {
	case xed.ADDRESS_WIDTH_16b:
		mask = mask16
	case xed.ADDRESS_WIDTH_32b:
		mask = mask32
	case xed.ADDRESS_WIDTH_64b:
		mask = mask64
	default:
		panic("invalid address width")
	}

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

	case xed.REG_EIP:
		idx = 0
		mask = mask32
	case xed.REG_RIP:
		idx = 0
		mask = mask64

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
	return
}

func (m *Mach) PhysAddr(vaddr uint64) uint64 {
	for _, sg := range m.Seg {
		if sg.Virt <= vaddr && vaddr < sg.Virt+sg.Size {
			return sg.Phys + vaddr - sg.Virt
		}
	}
	return 0
}

func (m *Mach) ReadReg(reg xed.Reg) uint64 {
	idx, shift, mask := m.decodeReg(reg)
	return (m.Reg[idx] >> shift) & mask
}

func (m *Mach) ReadMem(addr uint64, size int) uint64 {
	var v uint64
	for i := 0; i < size; i++ {
		n := m.PhysAddr(addr + uint64(i))
		v |= uint64(m.Mem[n]) << (8 * uint(i))
	}
	return v
}

func (m *Mach) ReadBuffer(addr uint64, buf []byte) {
	for i := range buf {
		buf[i] = byte(m.ReadMem(addr+uint64(i), 1))
	}
}

func (m *Mach) WriteReg(reg xed.Reg, val uint64) {
	idx, shift, mask := m.decodeReg(reg)
	m.Reg[idx] = (((m.Reg[idx] >> shift) &^ mask) | (val & mask)) << shift
}

func (m *Mach) WriteMem(addr, val uint64, size int) {
	var v uint64
	for i := 0; i < size; i++ {
		n := m.PhysAddr(addr + uint64(i))
		m.Mem[n] = byte(v>>uint(8*i)) & 0xff
	}
}

func (m *Mach) WriteBuffer(addr uint64, buf []byte) {
	for i := range buf {
		m.WriteMem(addr+uint64(i), uint64(buf[i]), 1)
	}
}

func (m *Mach) Step() error {
	var inst xed.DecodedInst
	fmt.Println(m.Disasm(m.ReadReg(REG_NIP)))
	m.fetch(&inst)
	return m.Op(&inst)
}

func (m *Mach) Op(inst *xed.DecodedInst) error {
	switch inst.Iform() {
	case xed.IFORM_CALL_NEAR_RELBRz:
		return m.callNearRelbrz(inst)
	case xed.IFORM_PUSH_GPRv_50:
		return m.pushGprv_50(inst)
	case xed.IFORM_MOV_GPRv_MEMv:
		return m.movGprv_Memv(inst)
	case xed.IFORM_SUB_OrAX_IMMz:
		return m.subOrAxImmz(inst)
	case xed.IFORM_MOV_GPRv_GPRv_8B:
		return m.mov_Gprv_Gprv_8b(inst)
	case xed.IFORM_AND_GPRv_IMMb:
		return m.and_Gprv_Immb(inst)
	case xed.IFORM_SUB_GPRv_IMMz:
		return m.sub_Gprv_Immz(inst)
	case xed.IFORM_PUSH_IMMz:
		return m.pushImmz(inst)
	case xed.IFORM_LEA_GPRv_AGEN:
		return m.lea_Gprv_Agen(inst)
	case xed.IFORM_TEST_GPRv_GPRv:
		return m.test_Gprv_Gprv(inst)
	case xed.IFORM_JZ_RELBRb:
		return m.jz_Relbrb(inst)
	case xed.IFORM_SUB_GPRv_GPRv_2B:
		return m.sub_Gprv_Gprv_2b(inst)
	case xed.IFORM_MOV_GPR8_MEMb:
		return m.mov_Gpr8_Memb(inst)
	case xed.IFORM_MOV_MEMb_GPR8:
		return m.mov_Memb_Gpr8(inst)
	case xed.IFORM_DEC_GPRv_48:
		return m.dec_Gprv_48(inst)
	default:
		return fmt.Errorf("unsupported instruction")
	}
}

func (m *Mach) dec_Gprv_48(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) mov_Memb_Gpr8(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) mov_Gpr8_Memb(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) sub_Gprv_Gprv_2b(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) jz_Relbrb(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) test_Gprv_Gprv(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) lea_Gprv_Agen(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) pushImmz(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) sub_Gprv_Immz(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) and_Gprv_Immb(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) mov_Gprv_Gprv_8b(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) subOrAxImmz(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) movGprv_Memv(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) pushGprv_50(inst *xed.DecodedInst) error {
	return nil
}

func (m *Mach) callNearRelbrz(inst *xed.DecodedInst) error {
	disp := inst.BranchDisplacement()
	m.Push(m.ReadReg(REG_NIP))
	m.WriteReg(REG_NIP, m.ReadReg(REG_NIP)+uint64(disp))
	return nil
}

func (m *Mach) ReadInst(addr uint64, inst *xed.DecodedInst) {
	var buf [16]byte
	m.ReadBuffer(addr, buf[:])
	inst.Zero()
	inst.SetMode(m.Prog.Mode, m.Prog.Width)
	inst.Decode(buf[:])
}

func (m *Mach) Disasm(addr uint64) string {
	var inst xed.DecodedInst
	m.ReadInst(addr, &inst)
	istr, _ := xed.FormatContext(xed.SYNTAX_INTEL, &inst, 0, nil)
	ostr := ""
	for i := uint(0); i < inst.Length(); i++ {
		ostr += fmt.Sprintf("%02x ", inst.Byte(i))
	}
	iform := inst.Iform()
	return fmt.Sprintf("%x %20s %32s %20s", addr, ostr, istr, iform)
}

func (m *Mach) fetch(inst *xed.DecodedInst) {
	m.ReadInst(m.ReadReg(REG_NIP), inst)
	m.WriteReg(REG_NIP, m.ReadReg(REG_NIP)+uint64(inst.Length()))
}
