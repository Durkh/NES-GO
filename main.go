package main

import (
	"NES/hardware"
	"NES/ioNES"
	"encoding/hex"
	"fmt"
)

func main() {

	fmt.Println("Debug mode:")
	fmt.Println("1 with ROM, 2 Manual program")
	var choice uint8

	switch fmt.Scan(choice); choice {
	case 1:
		var ROMPath string
		var bus ioNES.Bus
		bus.NewBus(ROMPath)
	case 2:
		cpu := hardware.NewCPU()
		fmt.Println("hex string:")
		var str string
		fmt.Scan(str)
		LoadRAM([]byte(str), cpu)
		Program(cpu)
	default:
		return
	}

}

// first function used to debug, keep here just to remember how the project is getting bigger
func ProgramRAM(cpu *hardware.Nes6502) bool {

	var counter uint16 = 0x8000

	// Load Program (assembled at https://www.masswerk.at/6502/assembler.html)
	//shout-out to javidx for this piece of code, truly useful to debugging
	/*
		*=$8000
		LDX #10
		STX $0000
		LDX #3
		STX $0001
		LDY $0000
		LDA #0
		CLC
		loop
		ADC $0001
		DEY
		BNE loop
		STA $0002
		NOP
		NOP
		NOP
	*/

	str := []byte("A20A8E0000A2038E0100AC0000A900186D010088D0FA8D0200EAEAEA")
	hexCode := make([]byte, hex.DecodedLen(len(str)))

	_, err := hex.Decode(hexCode, str)
	if err != nil {
		return false
	}

	cpu.Write(0xFFFC, 0x00)
	cpu.Write(0xFFFD, 0x80)

	for _, i := range hexCode {
		cpu.Write(counter, i)
		counter++
	}

	//fmt.Printf("% X", ReadRAM(100, cpu))

	return true
}

func LoadRAM(code []byte, cpu *hardware.Nes6502) bool {

	//the codes will be stored on memory location 0x800
	var counter uint16 = 0x8000
	hexCode := make([]byte, hex.DecodedLen(len(code)))

	//transform the string into hex digits
	_, err := hex.Decode(hexCode, code)
	if err != nil {
		return false
	}

	//tell the CPU were to look for the code
	cpu.Write(0xFFFC, 0x00)
	cpu.Write(0xFFFD, 0x80)

	//write the code into RAM
	for _, i := range hexCode {
		cpu.Write(counter, i)
		counter++
	}

	return true
}

func ReadRAM(size uint32, cpu *hardware.Nes6502) (slice []byte) {

	slice = make([]byte, size)

	for i := uint32(0); i < size; i++ {
		slice[i] = cpu.Read(uint16(i))
	}

	return
}

func DrawScreen(cpu *hardware.Nes6502) {

	var values [5]byte
	var addr uint16

	values[0], values[1], values[2], values[3], addr, values[4] = cpu.GetReg()

	fmt.Printf("\n\n")
	fmt.Printf("Status:\nC Z I D B - V N\n")
	for i := 0; i < 8; i++ {
		if values[4]&(1<<i) != 0 {
			fmt.Printf("%b ", 1)
		} else {
			fmt.Printf("%b ", 0)
		}
	}
	fmt.Printf("\t%8b", values[4])
	fmt.Printf("\na:%#X\tx:%#X\ty:%#X\nstack Ptr:%#X\nPC:%#X\n\n", values[0], values[1], values[2], values[3], addr)

	arr := ReadRAM(10, cpu)
	for i := 0; i < len(arr); i++ {
		fmt.Printf("%X ", arr[i])
	}

}

func Program(cpu *hardware.Nes6502) {

	var tecla string

loop:
	for {
		DrawScreen(cpu)
		_, err := fmt.Scan(&tecla)

		if err != nil {
			fmt.Println("Erro lendo")
			return
		}

		switch tecla {
		case "r":
			cpu.Reset()
		case "q":
			break loop
		case "w":
			for {
				cpu.Clock()
				if cpu.Completed() {
					continue loop
				}
			}
		}
	}

}
