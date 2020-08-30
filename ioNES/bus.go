package ioNES

import (
	"NES/hardware"
	"fmt"
)

type Bus struct {
	cpu       *hardware.Nes6502
	ppu       *hardware.PPU
	cartridge *Cartridge

	cpuRam [2048]uint8

	systemClockCounter uint32
}

func (bus *Bus) NewBus(ROMPath string) {
	cpu := hardware.NewCPU()
	bus.cpu = cpu
	bus.cpu.ConnectBus(bus)
	defer bus.cpu.Reset()

	cart, err := NewCartridge(ROMPath)
	if err != nil {
		fmt.Print(err)
		return
	}
	bus.ConnectCartridge(cart)

}

func (bus *Bus) CpuWrite(addr uint16, data uint8) {

	if bus.cartridge.CpuWrite(addr, data) {

	} else if addr <= 0x1FFF {
		bus.cpuRam[addr&0x7FF] = data
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		bus.ppu.CpuWrite(addr&0x0007, data)
	}
}

func (bus *Bus) CpuRead(addr uint16, bReadOnly bool) uint8 {
	var data uint8 = 0x0

	if bus.cartridge.CpuWrite(addr, data) {

	} else if addr <= 0x1FFF {
		return bus.cpuRam[addr&0x7FF]
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		data = bus.ppu.CpuRead(addr&0x0007, bReadOnly)
	}

	return data
}

func (bus *Bus) Reset() {
	bus.cpu.Reset()
	bus.systemClockCounter = 0
}

func (bus *Bus) clock() {

	bus.ppu.Clock()

	if bus.systemClockCounter%3 == 0 {
		bus.cpu.Clock()
	}

	bus.systemClockCounter++
}

func (bus *Bus) ConnectCartridge(cartridge *Cartridge) {
	bus.cartridge = cartridge
	bus.ppu.ConnectCartridge(cartridge)
}
