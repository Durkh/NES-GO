package main

import (
	"NES/hardware"
	"encoding/hex"
	"fmt"
)

func main() {

	cpu := hardware.NewCPU()

	if ProgramRAM(cpu) != true {
		println("error populating memory")
	}

	cpu.Reset()

	Program(cpu)

}

func ProgramRAM(cpu *hardware.Olc5602) bool {

	var counter uint16 = 0x8000

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

	/*dbg := make([]byte, hex.DecodedLen(len(str)))

	for i,_ := range dbg {
		dbg[i] = cpu.Read(uint16(0x8000 + i))
	}

	for i,_ := range dbg {
		fmt.Printf("%#X ", dbg[i])
	}*/

	return true
}

func DrawScreen(cpu *hardware.Olc5602) {

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

	for i := 0; i < 10; i++ {
		fmt.Printf("%X ", cpu.Read(uint16(i)))
	}

}

func Program(cpu *hardware.Olc5602) {

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
