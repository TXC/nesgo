package rp2ago3

import (
	s "github.com/stretchr/testify/suite"
)

// region cpuTest
type cpuTest struct {
	s.Suite
	cpu *RP2A03
}

func (s *cpuTest) BeforeTest(suiteName, testName string) {
	s.cpu = NewRP2A03(1789773)
	s.cpu.Reset()
}

func (s *cpuTest) AfterTest(_, _ string) {
	s.cpu = nil
}

// endregion cpuTest

func (s *cpuTest) TestStore() {
	s.cpu.APU.Pulse1.Registers[0] = 0xde
	s.cpu.Memory.Store(0x4000, 0xff)

	s.Equal(uint8(0xff), s.cpu.APU.Pulse1.Registers[0])

	s.cpu.Memory.Store(0x0800, 0xff)

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0000))

	s.cpu.Memory.Store(0x0800, 0x00)

	s.Equal(uint8(0x00), s.cpu.Memory.Fetch(0x0000))
}

type FakePPU struct {
	i      uint16
	memory [65536]uint8
}

func (ppu *FakePPU) Reset() {

}

func (ppu *FakePPU) Mappings(which Mapping) (fetch, store []uint16) {
	store = []uint16{0x2004}
	return
}

func (ppu *FakePPU) Fetch(address uint16) (value uint8) {
	return
}

func (ppu *FakePPU) Store(address uint16, value uint8) (oldValue uint8) {
	switch address {
	case 0x2004:
		ppu.memory[ppu.i] = value
		ppu.i++
	}

	return
}

func (s *cpuTest) TestDMA() {
	ppu := &FakePPU{}
	s.cpu.Memory.AddMappings(ppu, CPU)
	for address := uint32(0xff00); address <= 0xffff; address++ {
		s.cpu.Memory.Store(uint16(address), uint8(address&0x00ff))
	}

	s.cpu.Memory.Store(0x4014, 0xff)
	s.cpu.DMA.PerformDMA()

	for address := uint32(0x0000); address <= 0x00ff; address++ {
		s.Equal(uint8(address), ppu.memory[uint16(address)])
	}
}
