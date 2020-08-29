package hardware

import "fmt"

type instruction struct { //string struct
	name     string
	operate  func(*Olc5602) uint8
	addrMode func(*Olc5602) uint8
	imp      bool
	cycles   uint8
}

// lookup instruction table
var lookup [256]instruction

type Olc5602 struct { //cpu class
	bus *Bus

	a       uint8
	x       uint8
	y       uint8
	stckPtr uint8
	pc      uint16
	status  uint8

	fetched uint8
	addrAbs uint16
	addrRel uint16
	opcode  uint8
	cycles  uint8
}

func NewCPU() *Olc5602 {
	var cpu Olc5602
	bus := new(Bus)
	cpu.connectBus(bus)
	assignLookup(&lookup)
	cpu.Reset()
	return &cpu
}

const ( //flags
	C uint8 = 1 << iota //Carry Bit				1
	Z uint8 = 1 << iota //Zero					2
	I uint8 = 1 << iota //Disable interrupts	4
	D uint8 = 1 << iota //Decimal mode			8
	B uint8 = 1 << iota //Break					16
	U uint8 = 1 << iota //Unused				32
	V uint8 = 1 << iota //Overflow				64
	N uint8 = 1 << iota //Negative				128
)

func (cpu *Olc5602) connectBus(n *Bus) {
	cpu.bus = n
}

func (cpu Olc5602) Read(a uint16) uint8 {
	return cpu.bus.Read(a, false)
}

func (cpu Olc5602) Write(a uint16, d uint8) {
	cpu.bus.Write(a, d)
}

func (cpu Olc5602) getFlag(f uint8) uint8 {

	if (cpu.status & f) > 0 {
		return 1
	} else {
		return 0
	}
}

func (cpu *Olc5602) setFlag(f uint8, v bool) {

	if v {
		cpu.status |= f
	} else {
		cpu.status &= ^f
	}
}

func (cpu *Olc5602) Clock() {
	if cpu.cycles == 0 {
		cpu.setFlag(U, true)
		cpu.opcode = cpu.Read(cpu.pc)
		cpu.pc++

		cpu.cycles = lookup[cpu.opcode].cycles
		additionalCycle1 := lookup[cpu.opcode].addrMode(cpu)
		additionalCycle2 := lookup[cpu.opcode].operate(cpu)

		cpu.cycles += additionalCycle1 & additionalCycle2
		fmt.Printf("%X, %d, %q, %X\n", cpu.opcode, cpu.opcode, lookup[cpu.opcode].name, cpu.fetched)
	}
	cpu.cycles--
}

func (cpu Olc5602) GetCycles() uint8 {
	return cpu.cycles
}

func (cpu Olc5602) Completed() bool {
	return cpu.cycles == 0
}

func (cpu Olc5602) GetReg() (a uint8, x uint8, y uint8, stckPtr uint8, pc uint16, status uint8) {
	return cpu.a, cpu.x, cpu.y, cpu.stckPtr, cpu.pc, cpu.status
}

// ADDRESSING MODES

func (cpu *Olc5602) imp() uint8 {

	cpu.fetched = cpu.a
	return 0
}

func (cpu *Olc5602) imm() uint8 {
	cpu.addrAbs = cpu.pc
	cpu.pc++
	return 0
}

func (cpu *Olc5602) zp0() uint8 {
	cpu.addrAbs = uint16(cpu.Read(cpu.pc))
	cpu.pc++
	cpu.addrAbs &= 0x00FF
	return 0
}

func (cpu *Olc5602) zpx() uint8 {
	cpu.addrAbs = uint16(cpu.Read(cpu.pc + uint16(cpu.x)))
	cpu.pc++
	cpu.addrAbs &= 0x00FF
	return 0
}

func (cpu *Olc5602) zpy() uint8 {
	cpu.addrAbs = uint16(cpu.Read(cpu.pc + uint16(cpu.y)))
	cpu.pc++
	cpu.addrAbs &= 0x00FF
	return 0
}

func (cpu *Olc5602) rel() uint8 {
	cpu.addrRel = uint16(cpu.Read(cpu.pc))
	cpu.pc++

	if cpu.addrRel > 0x80 {
		cpu.addrRel -= 0x100
	}

	return 0
}

func (cpu *Olc5602) abs() uint8 {

	var low, high uint16

	low = uint16(cpu.Read(cpu.pc))
	cpu.pc++
	high = uint16(cpu.Read(cpu.pc))
	cpu.pc++

	cpu.addrAbs = (high << 8) | low

	return 0
}

func (cpu *Olc5602) abx() uint8 {

	var low, high uint16

	low = uint16(cpu.Read(cpu.pc))
	cpu.pc++
	high = uint16(cpu.Read(cpu.pc))
	cpu.pc++

	cpu.addrAbs = (high << 8) | low
	cpu.addrAbs += uint16(cpu.x)

	if (cpu.addrAbs & 0xFF00) != (high << 8) {
		return 1
	} else {
		return 0
	}
}

func (cpu *Olc5602) aby() uint8 {

	var low, high uint16

	low = uint16(cpu.Read(cpu.pc))
	cpu.pc++
	high = uint16(cpu.Read(cpu.pc))
	cpu.pc++

	cpu.addrAbs = (high << 8) | low
	cpu.addrAbs += uint16(cpu.y)

	if (cpu.addrAbs & 0xFF00) != (high << 8) {
		return 1
	} else {
		return 0
	}
}

func (cpu *Olc5602) ind() uint8 {
	var ptrLow, ptrHigh, ptr uint16

	ptrLow = uint16(cpu.Read(cpu.pc))
	cpu.pc++
	ptrHigh = uint16(cpu.Read(cpu.pc))
	cpu.pc++

	ptr = (ptrHigh << 8) | ptrLow

	if ptrLow == 0x00FF {
		cpu.addrAbs = (uint16(cpu.Read(ptr&0xFF00)) << 8) | uint16(cpu.Read(ptr+0))
	} else {
		cpu.addrAbs = (uint16(cpu.Read(ptr+1)) << 8) | uint16(cpu.Read(ptr+0))
	}
	return 0
}

func (cpu *Olc5602) izx() uint8 {
	var low, high, t uint16

	t = uint16(cpu.Read(cpu.pc))
	cpu.pc++

	low = uint16(cpu.Read((t + uint16(cpu.x)) & 0x00FF))
	high = uint16(cpu.Read((t + uint16(cpu.x) + 1) & 0x00FF))

	cpu.addrAbs = (high << 8) | low

	return 0
}

func (cpu *Olc5602) izy() uint8 {
	var low, high, t uint16

	t = uint16(cpu.Read(cpu.pc))
	cpu.pc++

	low = uint16(cpu.Read(t & 0x00FF))
	high = uint16(cpu.Read((t + 1) & 0x00FF))

	cpu.addrAbs = (high << 8) | low
	cpu.addrAbs += uint16(cpu.y)

	if (cpu.addrAbs & 0xFF00) != (high << 8) {
		return 1
	} else {
		return 0
	}
}

func (cpu *Olc5602) fetch() uint8 {

	if lookup[cpu.opcode].imp != true {
		cpu.fetched = cpu.Read(cpu.addrAbs)
	}
	return cpu.fetched
}

// OPCODE FUNCTIONS

func (cpu *Olc5602) adc() uint8 {

	cpu.fetch()

	temp := uint16(cpu.a) + uint16(cpu.fetched) + uint16(cpu.getFlag(C))
	cpu.setFlag(C, temp > 255)
	cpu.setFlag(Z, (temp&0x00FF) == 0)
	cpu.setFlag(V, (^(uint16(cpu.a)^uint16(cpu.fetched))&(uint16(cpu.a)^temp)&0x0080) != 0)
	cpu.setFlag(N, (temp&0x80) != 0)
	cpu.a = uint8(temp & 0x00FF)

	return 1
}

func (cpu *Olc5602) and() uint8 {
	cpu.fetch()
	cpu.a = cpu.a & cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 1
}

func (cpu *Olc5602) asl() uint8 {

	cpu.fetch()
	temp := uint16(cpu.fetched) << 1
	cpu.setFlag(C, (temp&0xFF00) > 0)
	cpu.setFlag(Z, (temp&0x00FF) == 0)
	cpu.setFlag(N, (temp&0x80) != 0)

	if lookup[cpu.opcode].imp != true {
		cpu.a = uint8(temp & 0x00FF)
	} else {
		cpu.Write(cpu.addrAbs, uint8(temp&0x00FF))
	}

	return 0
}

func (cpu *Olc5602) bcc() uint8 {
	if cpu.getFlag(C) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *Olc5602) bcs() uint8 {

	if cpu.getFlag(C) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *Olc5602) beq() uint8 {

	if cpu.getFlag(Z) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *Olc5602) bit() uint8 {

	cpu.fetch()
	temp := cpu.a & cpu.fetched
	cpu.setFlag(Z, (temp&0x00FF) == 0)
	cpu.setFlag(N, (cpu.fetched&(1<<7)) != 0)
	cpu.setFlag(V, (cpu.fetched&(1<<6)) != 0)

	return 0
}

func (cpu *Olc5602) bmi() uint8 {

	if cpu.getFlag(N) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *Olc5602) bne() uint8 {

	if cpu.getFlag(Z) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *Olc5602) bpl() uint8 {

	if cpu.getFlag(N) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *Olc5602) brk() uint8 {

	cpu.pc++

	cpu.setFlag(I, true)
	cpu.Write(0x0100+uint16(cpu.stckPtr), uint8((cpu.pc>>8)&0x00FF))
	cpu.stckPtr--
	cpu.Write(0x0100+uint16(cpu.stckPtr), uint8((cpu.pc)&0x00FF))
	cpu.stckPtr--

	cpu.setFlag(B, true)
	cpu.Write(0x0100+uint16(cpu.stckPtr), cpu.status)
	cpu.stckPtr--
	cpu.setFlag(B, false)

	cpu.pc = uint16(cpu.Read(0xFFFE)) | (uint16(cpu.Read(0xFFFF)) << 8)

	return 0
}

func (cpu *Olc5602) bvc() uint8 {

	if cpu.getFlag(V) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *Olc5602) bvs() uint8 {

	if cpu.getFlag(V) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel

		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *Olc5602) clc() uint8 {

	cpu.setFlag(C, false)
	return 0
}

func (cpu *Olc5602) cld() uint8 {
	cpu.setFlag(D, false)
	return 0
}

func (cpu *Olc5602) cli() uint8 {
	cpu.setFlag(I, false)
	return 0
}

func (cpu *Olc5602) clv() uint8 {
	cpu.setFlag(V, false)
	return 0
}

func (cpu *Olc5602) cmp() uint8 {
	cpu.fetch()

	temp := uint16(cpu.a) - uint16(cpu.fetched)
	cpu.setFlag(C, cpu.a >= cpu.fetched)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	return 1
}

func (cpu *Olc5602) cpx() uint8 {
	cpu.fetch()

	temp := uint16(cpu.x) - uint16(cpu.fetched)
	cpu.setFlag(C, cpu.x >= cpu.fetched)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	return 0
}

func (cpu *Olc5602) cpy() uint8 {
	cpu.fetch()

	temp := uint16(cpu.y) - uint16(cpu.fetched)
	cpu.setFlag(C, cpu.y >= cpu.fetched)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	return 0
}

func (cpu *Olc5602) dec() uint8 {

	cpu.fetch()
	temp := cpu.fetched - 1
	cpu.Write(cpu.addrAbs, temp&0x00FF)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	return 0
}

func (cpu *Olc5602) dex() uint8 {

	cpu.x--
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)

	return 0
}

func (cpu *Olc5602) dey() uint8 {

	cpu.y--
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)

	return 0
}

func (cpu *Olc5602) eor() uint8 {

	cpu.fetch()
	cpu.a = cpu.a ^ cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 1
}

func (cpu *Olc5602) inc() uint8 {
	cpu.fetch()

	temp := cpu.fetched + 1
	cpu.Write(cpu.addrAbs, temp&0x00FF)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	return 0
}

func (cpu *Olc5602) inx() uint8 {

	cpu.x++
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)

	return 0
}

func (cpu *Olc5602) iny() uint8 {

	cpu.y++
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)

	return 0
}

func (cpu *Olc5602) jmp() uint8 {

	cpu.pc = cpu.addrAbs

	return 0
}

func (cpu *Olc5602) jsr() uint8 {

	cpu.pc--

	cpu.Write(0x0100+uint16(cpu.stckPtr), uint8((cpu.pc>>8)&0x00FF))
	cpu.stckPtr--
	cpu.Write(0x0100+uint16(cpu.stckPtr), uint8((cpu.pc)&0x00FF))
	cpu.stckPtr--

	cpu.pc = cpu.addrAbs

	return 0
}

func (cpu *Olc5602) lda() uint8 {

	cpu.fetch()
	cpu.a = cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 1
}

func (cpu *Olc5602) ldx() uint8 {

	cpu.fetch()
	cpu.x = cpu.fetched
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)

	return 1
}

func (cpu *Olc5602) ldy() uint8 {

	cpu.fetch()
	cpu.y = cpu.fetched
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)

	return 1
}

func (cpu *Olc5602) lsr() uint8 {

	cpu.fetch()
	cpu.setFlag(C, (cpu.fetched&0x0001) != 0)
	temp := uint16(cpu.fetched >> 1)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	if lookup[cpu.opcode].imp != true {
		cpu.a = uint8(temp & 0x00FF)
	} else {
		cpu.Write(cpu.addrAbs, uint8(temp&0x00FF))
	}

	return 0
}

func (cpu *Olc5602) nop() uint8 {

	switch cpu.opcode {
	case 0x1C:
		fallthrough
	case 0x3C:
		fallthrough
	case 0x5C:
		fallthrough
	case 0x7C:
		fallthrough
	case 0xDC:
		fallthrough
	case 0xFC:
		return 1
	}
	return 0
}

func (cpu *Olc5602) ora() uint8 {
	cpu.fetch()
	cpu.a = cpu.a | cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 1
}

func (cpu *Olc5602) pha() uint8 {

	cpu.Write(0x0100+uint16(cpu.stckPtr), cpu.a)
	cpu.stckPtr--
	return 0
}

func (cpu *Olc5602) php() uint8 {
	cpu.Write(0x0100+uint16(cpu.stckPtr), cpu.status|B|U)
	cpu.setFlag(B, false)
	cpu.setFlag(U, false)
	cpu.stckPtr--

	return 0
}

func (cpu *Olc5602) pla() uint8 {

	cpu.stckPtr++
	cpu.a = cpu.Read(0x0100 + uint16(cpu.stckPtr))
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 0
}

func (cpu *Olc5602) plp() uint8 {

	cpu.stckPtr++
	cpu.status = cpu.Read(0x0100 + uint16(cpu.stckPtr))
	cpu.setFlag(U, true)

	return 0
}

func (cpu *Olc5602) rol() uint8 {

	cpu.fetch()
	temp := uint16(cpu.fetched<<1) | uint16(cpu.getFlag(C))
	cpu.setFlag(C, (temp&0xFF00) != 0)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	if lookup[cpu.opcode].imp != true {
		cpu.a = uint8(temp & 0x00FF)
	} else {
		cpu.Write(cpu.addrAbs, uint8(temp&0x00FF))
	}

	return 0
}

func (cpu *Olc5602) ror() uint8 {

	cpu.fetch()
	temp := (uint16(cpu.getFlag(C)) << 7) | uint16(cpu.fetched>>1)
	cpu.setFlag(C, (cpu.fetched&0x01) != 0)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	if lookup[cpu.opcode].imp != true {
		cpu.a = uint8(temp & 0x00FF)
	} else {
		cpu.Write(cpu.addrAbs, uint8(temp&0x00FF))
	}

	return 0
}

func (cpu *Olc5602) rti() uint8 {

	cpu.stckPtr++
	cpu.status = cpu.Read(0x100 + uint16(cpu.stckPtr))
	cpu.status &= ^B
	cpu.status &= ^U

	cpu.stckPtr++
	cpu.pc = uint16(cpu.Read(0x0100 + uint16(cpu.stckPtr)))
	cpu.stckPtr++
	cpu.pc |= uint16(cpu.Read(0x0100+uint16(cpu.stckPtr))) << 8
	return 0
}

func (cpu *Olc5602) rts() uint8 {

	cpu.stckPtr++
	cpu.pc = uint16(cpu.Read(0x0100 + uint16(cpu.stckPtr)))
	cpu.stckPtr++
	cpu.pc |= uint16(cpu.Read(0x0100+uint16(cpu.stckPtr))) << 8

	cpu.pc++
	return 0
}

func (cpu *Olc5602) sbc() uint8 {

	cpu.fetch()
	value := uint16(cpu.fetched) ^ 0x00FF
	temp := uint16(cpu.a) + value + uint16(cpu.getFlag(C))
	cpu.setFlag(C, (temp&0xFF00) != 0)
	cpu.setFlag(Z, (temp&0x00FF) == 0)
	cpu.setFlag(V, ((temp^uint16(cpu.a))&(temp&value)&0x0080) != 0)
	cpu.setFlag(N, (temp&0x0080) != 0)
	cpu.a = uint8(temp & 0x00FF)

	return 1
}

func (cpu *Olc5602) sec() uint8 {

	cpu.setFlag(C, true)
	return 0
}

func (cpu *Olc5602) sed() uint8 {

	cpu.setFlag(D, true)
	return 0
}

func (cpu *Olc5602) sei() uint8 {

	cpu.setFlag(I, true)
	return 0
}

func (cpu *Olc5602) sta() uint8 {

	cpu.Write(cpu.addrAbs, cpu.a)
	return 0
}

func (cpu *Olc5602) stx() uint8 {

	cpu.Write(cpu.addrAbs, cpu.x)
	return 0
}

func (cpu *Olc5602) sty() uint8 {

	cpu.Write(cpu.addrAbs, cpu.y)
	return 0

}

func (cpu *Olc5602) tax() uint8 {

	cpu.x = cpu.a
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)

	return 0
}

func (cpu *Olc5602) tay() uint8 {

	cpu.y = cpu.a
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)

	return 0
}

func (cpu *Olc5602) tsx() uint8 {

	cpu.x = cpu.stckPtr
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)

	return 0
}

func (cpu *Olc5602) txa() uint8 {

	cpu.a = cpu.x
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 0

}

func (cpu *Olc5602) txs() uint8 {

	cpu.stckPtr = cpu.x

	return 0
}

func (cpu *Olc5602) tya() uint8 {

	cpu.a = cpu.y
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 0

}

func (cpu Olc5602) xxx() uint8 {
	return 0
}

func (cpu *Olc5602) Reset() {
	cpu.a = 0
	cpu.x = 0
	cpu.y = 0
	cpu.stckPtr = 0xFD
	cpu.status = 0x00 | U

	cpu.addrAbs = 0xFFFC
	low := uint16(cpu.Read(cpu.addrAbs + 0))
	high := uint16(cpu.Read(cpu.addrAbs + 1))

	cpu.pc = (high << 8) | low

	cpu.addrRel = 0x0000
	cpu.addrAbs = 0x0000
	cpu.fetched = 0x00

	cpu.cycles = 8
}

func (cpu *Olc5602) irq() {

	if cpu.getFlag(I) == 0 {

		cpu.Write(0x0100+uint16(cpu.stckPtr), uint8((cpu.pc>>8)&0x000FF))
		cpu.stckPtr--
		cpu.Write(0x0100+uint16(cpu.stckPtr), uint8(cpu.pc&0x00FF))
		cpu.stckPtr--

		cpu.setFlag(B, false)
		cpu.setFlag(U, true)
		cpu.setFlag(I, true)
		cpu.Write(0x0100+uint16(cpu.stckPtr), cpu.status)
		cpu.stckPtr--

		cpu.addrAbs = 0xFFFE
		low := uint16(cpu.Read(cpu.addrAbs + 0))
		high := uint16(cpu.Read(cpu.addrAbs + 1))
		cpu.pc = (high << 8) | low

		cpu.cycles = 7
	}
}

func (cpu *Olc5602) nmi() {

	cpu.Write(0x0100+uint16(cpu.stckPtr), uint8((cpu.pc>>8)&0x000FF))
	cpu.stckPtr--
	cpu.Write(0x0100+uint16(cpu.stckPtr), uint8(cpu.pc&0x00FF))
	cpu.stckPtr--

	cpu.setFlag(B, false)
	cpu.setFlag(U, true)
	cpu.setFlag(I, true)
	cpu.Write(0x0100+uint16(cpu.stckPtr), cpu.status)
	cpu.stckPtr--

	cpu.addrAbs = 0xFFFA
	low := uint16(cpu.Read(cpu.addrAbs + 0))
	high := uint16(cpu.Read(cpu.addrAbs + 1))
	cpu.pc = (high << 8) | low

	cpu.cycles = 7
}
