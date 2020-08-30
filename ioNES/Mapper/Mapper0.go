package Mapper

type Mapper0 struct {
	nPRGBanks, nCHRBanks byte
}

func NewMapper(prgBanks, chrBanks byte) *Mapper0 {
	var mapper Mapper0

	mapper.nPRGBanks, mapper.nCHRBanks = prgBanks, chrBanks

	return &mapper
}

func (mapper *Mapper0) CpuMapRead(addr uint16, mappedAddr *uint32) bool {

	if addr > 0x8000 {
		if mapper.nPRGBanks > 1 {
			*mappedAddr = uint32(addr & 0x7FFF)
		} else {
			*mappedAddr = uint32(addr & 0x3FFF)
		}
		return true
	}
	return false
}
func (mapper *Mapper0) CpuMapWrite(addr uint16, mappedAddr *uint32) bool {

	if addr > 0x8000 {
		if mapper.nPRGBanks > 1 {
			*mappedAddr = uint32(addr & 0x7FFF)
		} else {
			*mappedAddr = uint32(addr & 0x3FFF)
		}
		return true
	}
	return false
}
func (mapper *Mapper0) PpuMapRead(addr uint16, mappedAddr *uint32) bool {

	if addr >= 0x8000 && addr <= 0x1FFF {
		*mappedAddr = uint32(addr)
		return true
	}

	return false
}

func (mapper *Mapper0) PpuMapWrite(addr uint16, mappedAddr *uint32) bool {

	if addr >= 0x8000 && addr <= 0x1FFF {

		if mapper.nCHRBanks == 0 {
			*mappedAddr = uint32(addr)
		}
		return true
	}

	return false

}
