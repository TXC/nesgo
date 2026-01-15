package nes

import (
	s "github.com/stretchr/testify/suite"
	"github.com/txc/nesgo/pkg/rp2cgo2"
)

// region mirroringTest
type mirroringTest struct {
	s.Suite
	ppu *rp2cgo2.RP2C02
}

func (s *mirroringTest) BeforeTest(suiteName, testName string) {
	s.ppu = rp2cgo2.NewRP2C02(nil, "NTSC")
}

func (s *mirroringTest) AfterTest(_, _ string) {
	s.ppu = nil
}

// endregion mirroringTest

func (s *mirroringTest) TestVerticalMirroring() {
	s.ppu.Nametable.SetTables(0, 1, 0, 1)

	// Mirror nametable #2 to #0
	for i := uint16(0x2800); i <= 0x2bff; i++ {
		s.ppu.Memory.Store(i-0x0800, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i-0x0800, 0x00)
		s.Equal(uint8(0x00), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i-0x0800))

		s.ppu.Memory.Store(i, 0x00)
		s.Equal(uint8(0x00), s.ppu.Memory.Fetch(i))
	}

	// Mirror nametable #3 to #1
	for i := uint16(0x2c00); i <= 0x2fff; i++ {
		s.ppu.Memory.Store(i-0x0800, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i-0x0800))

		s.ppu.Memory.Store(i-0x0800, 0x00)
		s.Equal(uint8(0x00), s.ppu.Memory.Fetch(i-0x0800))

		s.ppu.Memory.Store(i, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i, 0x00)
		s.Equal(uint8(0x00), s.ppu.Memory.Fetch(i))
	}

	// Mirror nametable #2 to #0
	for i := uint16(0x3000); i <= 0x33ff; i++ {
		s.ppu.Memory.Store(i-0x1000, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i-0x1000, 0x00)
		s.Equal(uint8(0x00), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i-0x1000))

		s.ppu.Memory.Store(i, 0x00)
		s.Equal(uint8(0x00), s.ppu.Memory.Fetch(i))
	}

	// Mirror nametable #3 to #1
	for i := uint16(0x3400); i <= 0x37ff; i++ {
		s.ppu.Memory.Store(i-0x1000, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i-0x1000, 0x00)
		s.Equal(uint8(0x00), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i-0x1000))

		s.ppu.Memory.Store(i, 0x00)
		s.Equal(uint8(0x00), s.ppu.Memory.Fetch(i))
	}
}

func (s *mirroringTest) TestHorizontalMirroring() {
	s.ppu.Nametable.SetTables(0, 0, 1, 1)

	// Mirror nametable #1 to #0
	for i := uint16(0x2400); i <= 0x27ff; i++ {
		s.ppu.Memory.Store(i-0x0400, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i-0x0400, 0x00)
		s.ppu.Memory.Store(i, 0xff)

		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i-0x0400))

		s.ppu.Memory.Store(i, 0x00)
	}

	// Mirror nametable #3 to #2
	for i := uint16(0x2c00); i <= 0x2fff; i++ {
		s.ppu.Memory.Store(i-0x0400, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i-0x0400, 0x00)
		s.ppu.Memory.Store(i, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i-0x0400))

		s.ppu.Memory.Store(i, 0x00)
	}

	// Mirror nametable #1 to #0
	for i := uint16(0x3400); i <= 0x37ff; i++ {
		s.ppu.Memory.Store(i-0x1400, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i-0x1400, 0x00)
		s.ppu.Memory.Store(i, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i-0x1400))

		s.ppu.Memory.Store(i, 0x00)
	}

	// Mirror nametable #3 to #2
	for i := uint16(0x3c00); i <= 0x3eff; i++ {
		s.ppu.Memory.Store(i-0x1400, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i-0x1400, 0x00)
		s.ppu.Memory.Store(i, 0xff)
		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i-0x1400))

		s.ppu.Memory.Store(i, 0x00)
	}
}
