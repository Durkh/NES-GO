package hardware

type Cpu interface {
	connectBus(*Bus)
}

type Bus struct {
	cpu Cpu
	ram [64 * 1024]uint8
}

func (bus *Bus) Bus(cpu Cpu) { //relembrar de chamar o construtor passando a CPU
	bus.cpu = cpu
	bus.cpu.connectBus(bus)
}

func (bus *Bus) Write(addr uint16, data uint8) {
	if addr >= 0x0000 && addr <= 0xFFFF {
		bus.ram[addr] = data
	}
}

func (bus *Bus) Read(addr uint16, bReadOnly bool) uint8 {
	if addr >= 0x0000 && addr <= 0xFFFF {
		return bus.ram[addr]
	}
	return 0x0000
}
