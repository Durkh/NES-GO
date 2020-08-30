package hardware

import "NES/ioNES"

type PPU struct {
	cartridge *ioNES.Cartridge

	nameTable    [2][1024]uint8
	paletteTable [32]uint8
}

func (ppu *PPU) CpuRead(addr uint16, rOnly bool) uint8 {

	var data uint8

	switch addr {
	case 0x0000:
	case 0x0001:
	case 0x0002:
	case 0x0003:
	case 0x0004:
	case 0x0005:
	case 0x0006:
	case 0x0007:
	}

	return data
}

func (ppu *PPU) CpuWrite(addr uint16, data uint8) {

	switch addr {
	case 0x0000:
	case 0x0001:
	case 0x0002:
	case 0x0003:
	case 0x0004:
	case 0x0005:
	case 0x0006:
	case 0x0007:
	}

}

func (ppu *PPU) ppuRead(addr uint16, rOnly bool) uint8 {

	var data uint8 = 0x0
	addr &= 0x3FFF

	if ppu.cartridge.PpuRead(addr, data) {

	}

	return data
}

func (ppu *PPU) ppuWrite(addr uint16, data uint8) {

	if ppu.cartridge.PpuWrite(addr, data) {

	}

	addr &= 0x3FFF
}

func (ppu *PPU) ConnectCartridge(cartridge *ioNES.Cartridge) {
	ppu.cartridge = cartridge
}

func (ppu *PPU) Clock() {

}
