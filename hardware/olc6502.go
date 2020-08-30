package hardware

import (
	"NES/ioNES"
	"fmt"
)

type instruction struct { //struct with the model of an instruction of 6502
	name     string               //the name of the instruction
	operate  func(*Nes6502) uint8 //the operation of the instruction
	addrMode func(*Nes6502) uint8 //the addressing mode of the instruction
	imp      bool                 //boolean to help catching functions with implied addressing mode
	cycles   uint8                //number of cycles required to perform the instruction
}

// lookup instruction table
var lookup [256]instruction

type Nes6502 struct { //cpu class
	bus *ioNES.Bus

	//cpu registers
	a, x, y uint8
	stckPtr uint8
	pc      uint16
	status  uint8

	//variables to help the emulation
	fetched          uint8
	addrAbs, addrRel uint16
	opcode           uint8
	cycles           uint8
}

func NewCPU() *Nes6502 { //cpu constructor
	var cpu Nes6502
	bus := new(ioNES.Bus)
	cpu.ConnectBus(bus)
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

func (cpu *Nes6502) ConnectBus(n *ioNES.Bus) {
	cpu.bus = n
}

func (cpu Nes6502) Read(a uint16) uint8 {
	return cpu.bus.CpuRead(a, false)
}

func (cpu Nes6502) Write(a uint16, d uint8) {
	cpu.bus.CpuWrite(a, d)
}

func (cpu Nes6502) getFlag(f uint8) uint8 {

	if (cpu.status & f) > 0 {
		return 1
	} else {
		return 0
	}
}

func (cpu *Nes6502) setFlag(f uint8, v bool) {

	if v {
		cpu.status |= f
	} else {
		cpu.status &= ^f
	}
}

func (cpu *Nes6502) Clock() {
	if cpu.cycles == 0 {
		cpu.setFlag(U, true)
		cpu.opcode = cpu.Read(cpu.pc)
		cpu.pc++

		cpu.cycles = lookup[cpu.opcode].cycles
		additionalCycle1 := lookup[cpu.opcode].addrMode(cpu)
		additionalCycle2 := lookup[cpu.opcode].operate(cpu)

		cpu.cycles += additionalCycle1 & additionalCycle2
		//debugging function TODO: move later to function in main
		fmt.Printf("%X, %d, %q, %X\n", cpu.opcode, cpu.opcode, lookup[cpu.opcode].name, cpu.fetched)
	}
	cpu.cycles--
}

//cycles getter
func (cpu Nes6502) GetCycles() uint8 {
	return cpu.cycles
}

//function to return true when a instruction is completed, useful to step-to-step tests
func (cpu Nes6502) Completed() bool {
	return cpu.cycles == 0
}

//registers getter
func (cpu Nes6502) GetReg() (a uint8, x uint8, y uint8, stckPtr uint8, pc uint16, status uint8) {
	return cpu.a, cpu.x, cpu.y, cpu.stckPtr, cpu.pc, cpu.status
}

// ADDRESSING MODES

/*
Implied addressing mode
Used when the instruction has no aditional data to fetch
*/
func (cpu *Nes6502) imp() uint8 {

	cpu.fetched = cpu.a
	return 0
}

/*
Immediate addressing mode
Used when the data is immediately after the instruction
*/
func (cpu *Nes6502) imm() uint8 {
	cpu.addrAbs = cpu.pc
	cpu.pc++
	return 0
}

/*
Zero page addressing mode
Get only the first byte of address, saving time reading the other byte
*/
func (cpu *Nes6502) zp0() uint8 {
	cpu.addrAbs = uint16(cpu.Read(cpu.pc))
	cpu.pc++
	cpu.addrAbs &= 0x00FF
	return 0
}

func (cpu *Nes6502) zpx() uint8 {
	cpu.addrAbs = uint16(cpu.Read(cpu.pc + uint16(cpu.x)))
	cpu.pc++
	cpu.addrAbs &= 0x00FF
	return 0
}

func (cpu *Nes6502) zpy() uint8 {
	cpu.addrAbs = uint16(cpu.Read(cpu.pc + uint16(cpu.y)))
	cpu.pc++
	cpu.addrAbs &= 0x00FF
	return 0
}

/*
relative addressing mode
can refer to -128 ~ +127 of the current address
used to branch flow of the code
*/
func (cpu *Nes6502) rel() uint8 {
	cpu.addrRel = uint16(cpu.Read(cpu.pc))
	cpu.pc++

	// catch if the address if referring to a negative position referencing the current address
	if cpu.addrRel > 0x80 {
		cpu.addrRel -= 0x100 //if true, make the hexadecimal negative by subtracting 0x100, useful when adding to PC
	}

	return 0
}

func (cpu *Nes6502) abs() uint8 {

	var low, high uint16

	low = uint16(cpu.Read(cpu.pc))
	cpu.pc++
	high = uint16(cpu.Read(cpu.pc))
	cpu.pc++

	cpu.addrAbs = (high << 8) | low

	return 0
}

func (cpu *Nes6502) abx() uint8 {

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

func (cpu *Nes6502) aby() uint8 {

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

func (cpu *Nes6502) ind() uint8 {
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

func (cpu *Nes6502) izx() uint8 {
	var low, high, t uint16

	t = uint16(cpu.Read(cpu.pc))
	cpu.pc++

	low = uint16(cpu.Read((t + uint16(cpu.x)) & 0x00FF))
	high = uint16(cpu.Read((t + uint16(cpu.x) + 1) & 0x00FF))

	cpu.addrAbs = (high << 8) | low

	return 0
}

func (cpu *Nes6502) izy() uint8 {
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

func (cpu *Nes6502) fetch() uint8 {

	if lookup[cpu.opcode].imp != true {
		cpu.fetched = cpu.Read(cpu.addrAbs)
	}
	return cpu.fetched
}

// OPCODE FUNCTIONS

func (cpu *Nes6502) adc() uint8 {

	cpu.fetch()

	temp := uint16(cpu.a) + uint16(cpu.fetched) + uint16(cpu.getFlag(C))
	cpu.setFlag(C, temp > 255)
	cpu.setFlag(Z, (temp&0x00FF) == 0)
	cpu.setFlag(V, (^(uint16(cpu.a)^uint16(cpu.fetched))&(uint16(cpu.a)^temp)&0x0080) != 0)
	cpu.setFlag(N, (temp&0x80) != 0)
	cpu.a = uint8(temp & 0x00FF)

	return 1
}

func (cpu *Nes6502) and() uint8 {
	cpu.fetch()
	cpu.a = cpu.a & cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 1
}

func (cpu *Nes6502) asl() uint8 {

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

func (cpu *Nes6502) bcc() uint8 {
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

func (cpu *Nes6502) bcs() uint8 {

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

func (cpu *Nes6502) beq() uint8 {

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

func (cpu *Nes6502) bit() uint8 {

	cpu.fetch()
	temp := cpu.a & cpu.fetched
	cpu.setFlag(Z, (temp&0x00FF) == 0)
	cpu.setFlag(N, (cpu.fetched&(1<<7)) != 0)
	cpu.setFlag(V, (cpu.fetched&(1<<6)) != 0)

	return 0
}

func (cpu *Nes6502) bmi() uint8 {

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

func (cpu *Nes6502) bne() uint8 {

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

func (cpu *Nes6502) bpl() uint8 {

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

func (cpu *Nes6502) brk() uint8 {

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

func (cpu *Nes6502) bvc() uint8 {

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

func (cpu *Nes6502) bvs() uint8 {

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

func (cpu *Nes6502) clc() uint8 {

	cpu.setFlag(C, false)
	return 0
}

func (cpu *Nes6502) cld() uint8 {
	cpu.setFlag(D, false)
	return 0
}

func (cpu *Nes6502) cli() uint8 {
	cpu.setFlag(I, false)
	return 0
}

func (cpu *Nes6502) clv() uint8 {
	cpu.setFlag(V, false)
	return 0
}

func (cpu *Nes6502) cmp() uint8 {
	cpu.fetch()

	temp := uint16(cpu.a) - uint16(cpu.fetched)
	cpu.setFlag(C, cpu.a >= cpu.fetched)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	return 1
}

func (cpu *Nes6502) cpx() uint8 {
	cpu.fetch()

	temp := uint16(cpu.x) - uint16(cpu.fetched)
	cpu.setFlag(C, cpu.x >= cpu.fetched)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	return 0
}

func (cpu *Nes6502) cpy() uint8 {
	cpu.fetch()

	temp := uint16(cpu.y) - uint16(cpu.fetched)
	cpu.setFlag(C, cpu.y >= cpu.fetched)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	return 0
}

func (cpu *Nes6502) dec() uint8 {

	cpu.fetch()
	temp := cpu.fetched - 1
	cpu.Write(cpu.addrAbs, temp&0x00FF)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	return 0
}

func (cpu *Nes6502) dex() uint8 {

	cpu.x--
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)

	return 0
}

func (cpu *Nes6502) dey() uint8 {

	cpu.y--
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)

	return 0
}

func (cpu *Nes6502) eor() uint8 {

	cpu.fetch()
	cpu.a = cpu.a ^ cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 1
}

func (cpu *Nes6502) inc() uint8 {
	cpu.fetch()

	temp := cpu.fetched + 1
	cpu.Write(cpu.addrAbs, temp&0x00FF)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)

	return 0
}

func (cpu *Nes6502) inx() uint8 {

	cpu.x++
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)

	return 0
}

func (cpu *Nes6502) iny() uint8 {

	cpu.y++
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)

	return 0
}

func (cpu *Nes6502) jmp() uint8 {

	cpu.pc = cpu.addrAbs

	return 0
}

func (cpu *Nes6502) jsr() uint8 {

	cpu.pc--

	cpu.Write(0x0100+uint16(cpu.stckPtr), uint8((cpu.pc>>8)&0x00FF))
	cpu.stckPtr--
	cpu.Write(0x0100+uint16(cpu.stckPtr), uint8((cpu.pc)&0x00FF))
	cpu.stckPtr--

	cpu.pc = cpu.addrAbs

	return 0
}

func (cpu *Nes6502) lda() uint8 {

	cpu.fetch()
	cpu.a = cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 1
}

func (cpu *Nes6502) ldx() uint8 {

	cpu.fetch()
	cpu.x = cpu.fetched
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)

	return 1
}

func (cpu *Nes6502) ldy() uint8 {

	cpu.fetch()
	cpu.y = cpu.fetched
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)

	return 1
}

func (cpu *Nes6502) lsr() uint8 {

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

func (cpu *Nes6502) nop() uint8 {

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

func (cpu *Nes6502) ora() uint8 {
	cpu.fetch()
	cpu.a = cpu.a | cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 1
}

func (cpu *Nes6502) pha() uint8 {

	cpu.Write(0x0100+uint16(cpu.stckPtr), cpu.a)
	cpu.stckPtr--
	return 0
}

func (cpu *Nes6502) php() uint8 {
	cpu.Write(0x0100+uint16(cpu.stckPtr), cpu.status|B|U)
	cpu.setFlag(B, false)
	cpu.setFlag(U, false)
	cpu.stckPtr--

	return 0
}

func (cpu *Nes6502) pla() uint8 {

	cpu.stckPtr++
	cpu.a = cpu.Read(0x0100 + uint16(cpu.stckPtr))
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 0
}

func (cpu *Nes6502) plp() uint8 {

	cpu.stckPtr++
	cpu.status = cpu.Read(0x0100 + uint16(cpu.stckPtr))
	cpu.setFlag(U, true)

	return 0
}

func (cpu *Nes6502) rol() uint8 {

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

func (cpu *Nes6502) ror() uint8 {

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

func (cpu *Nes6502) rti() uint8 {

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

func (cpu *Nes6502) rts() uint8 {

	cpu.stckPtr++
	cpu.pc = uint16(cpu.Read(0x0100 + uint16(cpu.stckPtr)))
	cpu.stckPtr++
	cpu.pc |= uint16(cpu.Read(0x0100+uint16(cpu.stckPtr))) << 8

	cpu.pc++
	return 0
}

func (cpu *Nes6502) sbc() uint8 {

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

func (cpu *Nes6502) sec() uint8 {

	cpu.setFlag(C, true)
	return 0
}

func (cpu *Nes6502) sed() uint8 {

	cpu.setFlag(D, true)
	return 0
}

func (cpu *Nes6502) sei() uint8 {

	cpu.setFlag(I, true)
	return 0
}

func (cpu *Nes6502) sta() uint8 {

	cpu.Write(cpu.addrAbs, cpu.a)
	return 0
}

func (cpu *Nes6502) stx() uint8 {

	cpu.Write(cpu.addrAbs, cpu.x)
	return 0
}

func (cpu *Nes6502) sty() uint8 {

	cpu.Write(cpu.addrAbs, cpu.y)
	return 0

}

func (cpu *Nes6502) tax() uint8 {

	cpu.x = cpu.a
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)

	return 0
}

func (cpu *Nes6502) tay() uint8 {

	cpu.y = cpu.a
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)

	return 0
}

func (cpu *Nes6502) tsx() uint8 {

	cpu.x = cpu.stckPtr
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)

	return 0
}

func (cpu *Nes6502) txa() uint8 {

	cpu.a = cpu.x
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 0

}

func (cpu *Nes6502) txs() uint8 {

	cpu.stckPtr = cpu.x

	return 0
}

func (cpu *Nes6502) tya() uint8 {

	cpu.a = cpu.y
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)

	return 0

}

func (cpu Nes6502) xxx() uint8 {
	return 0
}

func (cpu *Nes6502) Reset() {
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

func (cpu *Nes6502) irq() {

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

func (cpu *Nes6502) nmi() {

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
