package ioNES

import (
	"encoding/binary"
	"errors"
	"os"
)

type Cartridge struct {
	prg, chr, sRam       []byte
	mapper, mirror       byte
	nPrgBanks, nCHRBanks byte

	bImageValid bool
}

type iNesHeader struct {
	name                 uint32
	prgChunk, chrChunk   byte
	flags1, flags2       byte
	ramSize              byte
	tvSystem1, tvSystem2 byte
	_                    [5]byte
}

func NewCartridge(path string) (*Cartridge, error) {
	var cart Cartridge

	//opening the .nes file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//reading the header
	header := iNesHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	//checking if the file is valid
	if header.name != 0x1a53454e {
		return nil, errors.New("invalid .nes file")
	}

	//mapper setting
	mapper1 := header.flags1 >> 4
	mapper2 := header.flags2 >> 4
	cart.mapper = mapper1 | (mapper2 << 4)

	//mirror getting
	mirror1 := header.flags1 & 1
	mirror2 := (header.flags1 >> 3) & 1
	cart.mirror = mirror1 | mirror2<<1

	if header.flags1&0x4 == 0x4 {
		trainer := make([]byte, 512)
		if _, err := file.Read(trainer); err != nil {
			return nil, err
		}
	}

	switch nFileType := uint8(1); nFileType {
	case 0:
	case 1:
		cart.nPrgBanks = header.prgChunk
		prgMem := make([]byte, int(cart.nPrgBanks)*16384)
		if _, err := file.Read(prgMem); err != nil {
			return nil, err
		}
		cart.prg = prgMem

		cart.nCHRBanks = header.chrChunk
		chrMem := make([]byte, int(cart.nCHRBanks)*8192)
		if _, err := file.Read(chrMem); err != nil {
			return nil, err
		}
		cart.chr = chrMem
	case 2:

	}

	cart.bImageValid = true

	return &cart, nil
}

func (cart Cartridge) CpuRead(addr uint16, data byte) bool {
	return false
}

func (cart Cartridge) CpuWrite(addr uint16, data byte) bool {
	return false
}

func (cart Cartridge) PpuRead(addr uint16, data byte) bool {
	return false
}

func (cart Cartridge) PpuWrite(addr uint16, data byte) bool {
	return false
}
