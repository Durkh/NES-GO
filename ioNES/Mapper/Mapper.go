package Mapper

type Mapper interface {
	CpuMapRead(uint16, uint32) bool
	CpuMapWrite(uint16, uint32) bool
	PpuMapRead(uint16, uint32) bool
	PpuMapWrite(uint16, uint32) bool
}
