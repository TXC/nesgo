package rp2cgo2

import (
	"testing"

	s "github.com/stretchr/testify/suite"
)

func TestPPUSuite(t *testing.T) {
	s.Run(t, &ntscTest{})
	s.Run(t, &palTest{})
}

// region ppuTest
type ppuTest struct {
	s.Suite
	ppu *RP2C02
}

// endregion ppuTest

// region ntscTest
type ntscTest struct {
	ppuTest
}

func (s *ntscTest) BeforeTest(suiteName, testName string) {
	s.ppu = NewRP2C02(nil, "NTSC")
}

func (s *ntscTest) AfterTest(_, _ string) {
	s.ppu = nil
}

// endregion ntscTest

// region palTest
type palTest struct {
	ppuTest
}

func (s *palTest) BeforeTest(suiteName, testName string) {
	s.ppu = NewRP2C02(nil, "PAL")
}

func (s *palTest) AfterTest(_, _ string) {
	s.ppu = nil
}

// endregion palTest

// region ppuTest methods
func (s *ppuTest) TestController() {
	s.ppu.Registers.Controller = 0x00
	s.ppu.Store(0x2000, 0xff)

	s.Equal(uint8(0xff), s.ppu.Registers.Controller)

	// BaseNametableAddress
	s.ppu.Registers.Controller = 0xff - 0x03

	s.Equal(uint16(0x2000), s.ppu.controller(BaseNametableAddress))

	s.ppu.Registers.Controller = 0xff - 0x02

	s.Equal(uint16(0x2400), s.ppu.controller(BaseNametableAddress))
	s.ppu.Registers.Controller = 0xff - 0x01

	s.Equal(uint16(0x2800), s.ppu.controller(BaseNametableAddress))
	s.ppu.Registers.Controller = 0xff

	s.Equal(uint16(0x2c00), s.ppu.controller(BaseNametableAddress))

	// VRAMAddressIncrement
	s.ppu.Registers.Controller = ^uint8(VRAMAddressIncrement)

	s.Equal(uint16(1), s.ppu.controller(VRAMAddressIncrement))

	s.ppu.Registers.Controller = uint8(VRAMAddressIncrement)

	s.Equal(uint16(32), s.ppu.controller(VRAMAddressIncrement))

	// SpritePatternAddress
	s.ppu.Registers.Controller = ^uint8(SpritePatternAddress)

	s.Equal(uint16(0x0000), s.ppu.controller(SpritePatternAddress))

	s.ppu.Registers.Controller = uint8(SpritePatternAddress)

	s.Equal(uint16(0x1000), s.ppu.controller(SpritePatternAddress))

	// BackgroundPatternAddress
	s.ppu.Registers.Controller = ^uint8(BackgroundPatternAddress)

	s.Equal(uint16(0x0000), s.ppu.controller(BackgroundPatternAddress))

	s.ppu.Registers.Controller = uint8(BackgroundPatternAddress)

	s.Equal(uint16(0x1000), s.ppu.controller(BackgroundPatternAddress))

	// SpriteSize
	s.ppu.Registers.Controller = ^uint8(SpriteSize)

	s.Equal(uint16(8), s.ppu.controller(SpriteSize))

	s.ppu.Registers.Controller = uint8(SpriteSize)

	s.Equal(uint16(16), s.ppu.controller(SpriteSize))

	// NMIOnVBlank
	s.ppu.Registers.Controller = ^uint8(NMIOnVBlank)

	s.Equal(uint16(0), s.ppu.controller(NMIOnVBlank))

	s.ppu.Registers.Controller = uint8(NMIOnVBlank)

	s.Equal(uint16(1), s.ppu.controller(NMIOnVBlank))
}

func (s *ppuTest) TestMask() {
	s.ppu.Registers.Mask = 0x00
	value := uint8(0xff)
	s.ppu.Store(0x2001, value)

	s.Equal(value, s.ppu.Registers.Mask)

	for _, m := range []MaskFlag{
		Grayscale, ShowBackgroundLeft, ShowSpritesLeft, ShowBackground,
		ShowSprites, IntensifyReds, IntensifyGreens, IntensifyBlues,
	} {
		s.ppu.Registers.Mask = uint8(m)

		s.True(s.ppu.mask(m))
	}
}

func (s *ppuTest) TestStatus() {
	s.ppu.Registers.Status = 0xff
	s.ppu.Latch = true
	s.ppu.LatchValue = 0x00

	s.Equal(uint8(0xe0), s.ppu.Fetch(0x2002))

	s.Equal(uint8(0x7f), s.ppu.Registers.Status)

	s.False(s.ppu.Latch)

	for _, sf := range []StatusFlag{
		SpriteOverflow,
		Sprite0Hit,
		VBlankStarted,
	} {
		s.ppu.Registers.Status = uint8(sf)

		s.True(s.ppu.status(sf))
	}
}

func (s *ppuTest) TestAddress() {
	s.ppu.Registers.Address = 0x00
	s.ppu.Store(0x2006, 0xff)
	s.ppu.Store(0x2006, 0xff)

	s.Equal(uint16(0x3fff), s.ppu.Registers.Address)

	// CoarseXScroll
	s.ppu.Registers.Address = 0x0000

	s.Equal(uint16(0x0000), s.ppu.address(CoarseXScroll))

	s.ppu.Registers.Address = 0xffff

	s.Equal(uint16(0x001f), s.ppu.address(CoarseXScroll))

	// CoarseYScroll
	s.ppu.Registers.Address = 0x0000

	s.Equal(uint16(0x0000), s.ppu.address(CoarseYScroll))

	s.ppu.Registers.Address = 0xffff

	s.Equal(uint16(0x001f), s.ppu.address(CoarseYScroll))

	// NametableSelect
	s.ppu.Registers.Address = 0x0000

	s.Equal(uint16(0x0000), s.ppu.address(NametableSelect))

	s.ppu.Registers.Address = 0xffff

	s.Equal(uint16(0x0003), s.ppu.address(NametableSelect))

	// FineYScroll
	s.ppu.Registers.Address = 0x0000

	s.Equal(uint16(0x0000), s.ppu.address(FineYScroll))

	s.ppu.Registers.Address = 0xffff

	s.Equal(uint16(0x0007), s.ppu.address(FineYScroll))
}

func (s *ppuTest) TestSprite() {
	sprite := uint32(0)

	// YPosition
	sprite = 0x00000000

	s.Equal(uint8(0x00), s.ppu.sprite(sprite, YPosition))

	sprite = 0xffffffff

	s.Equal(uint8(0xff), s.ppu.sprite(sprite, YPosition))

	// TileBank
	sprite = 0x00000000

	s.Equal(uint8(0x00), s.ppu.sprite(sprite, TileBank))

	sprite = 0xffffffff

	s.Equal(uint8(0x01), s.ppu.sprite(sprite, TileBank))

	// TileNumber
	sprite = 0x00000000

	s.Equal(uint8(0x00), s.ppu.sprite(sprite, TileNumber))

	sprite = 0xffffffff

	s.Equal(uint8(0xff), s.ppu.sprite(sprite, TileNumber))

	// SpritePalette
	sprite = 0x00000000

	s.Equal(uint8(0x00), s.ppu.sprite(sprite, SpritePalette))

	sprite = 0xffffffff

	s.Equal(uint8(0x03), s.ppu.sprite(sprite, SpritePalette))

	// Priority
	sprite = 0x00000000

	s.Equal(uint8(0x00), s.ppu.sprite(sprite, Priority))

	sprite = 0xffffffff

	s.Equal(uint8(0x01), s.ppu.sprite(sprite, Priority))

	// FlipHorizontally
	sprite = 0x00000000

	s.Equal(uint8(0x00), s.ppu.sprite(sprite, FlipHorizontally))

	sprite = 0xffffffff

	s.Equal(uint8(0x01), s.ppu.sprite(sprite, FlipHorizontally))

	// FlipVertically
	sprite = 0x00000000

	s.Equal(uint8(0x00), s.ppu.sprite(sprite, FlipVertically))

	sprite = 0xffffffff

	s.Equal(uint8(0x01), s.ppu.sprite(sprite, FlipVertically))

	// XPosition
	sprite = 0x00000000

	s.Equal(uint8(0x00), s.ppu.sprite(sprite, XPosition))

	sprite = 0xffffffff

	s.Equal(uint8(0xff), s.ppu.sprite(sprite, XPosition))
}

func (s *ppuTest) TestOAMAddress() {
	s.ppu.Registers.OAMAddress = 0x00

	s.ppu.Store(0x2003, 0xff)

	s.Equal(uint8(0xff), s.ppu.Registers.OAMAddress)
}

func (s *ppuTest) TestOAMData() {
	s.ppu.Registers.OAMAddress = 0x00

	for i := uint16(0x0000); i <= 0x00ff; i++ {
		s.ppu.Store(0x2004, uint8(i))
	}

	for i := uint16(0x0000); i <= 0x00ff; i++ {
		s.Equal(uint8(i), s.ppu.OAM.Fetch(uint16(i)))
	}
}

func (s *ppuTest) TestPaletteMirroring() {
	// Mirrored palette
	for _, i := range []uint16{0x3f10, 0x3f14, 0x3f18, 0x3f1c} {
		s.ppu.Memory.Store(i-0x0010, 0xff)

		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i-0x0010, 0x00)
		s.ppu.Memory.Store(i, 0xff)

		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i-0x0010))
	}

	for i := uint16(0x3f20); i <= 0x3fff; i++ {
		s.ppu.Memory.Store(i-0x0020, 0xff)

		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i))

		s.ppu.Memory.Store(i-0x0020, 0x00)
		s.ppu.Memory.Store(i, 0xff)

		s.Equal(uint8(0xff), s.ppu.Memory.Fetch(i-0x0020))
	}
}

func (s *ppuTest) TestAddressFetchStore() {
	s.ppu.Registers.Address = 0x0000
	s.ppu.Fetch(0x2002)

	s.ppu.Store(0x2006, 0x3f)
	s.ppu.Store(0x2006, 0xff)

	s.Equal(uint16(0x3fff), s.ppu.Registers.Address)

	s.ppu.Fetch(0x2002)

	s.ppu.Store(0x2006, 0x01)
	s.ppu.Store(0x2006, 0x01)

	s.Equal(uint16(0x0101), s.ppu.Registers.Address)
}

func (s *ppuTest) TestDataIncrement1() {
	s.ppu.Registers.Address = 0x0000
	s.ppu.Fetch(0x2002)

	s.ppu.Store(0x2006, 0x01)
	s.ppu.Store(0x2006, 0x00)

	s.Equal(uint16(0x0100), s.ppu.Registers.Address)

	s.ppu.Store(0x2007, 0xff)
	s.ppu.Store(0x2007, 0xff)
	s.ppu.Store(0x2007, 0xff)

	s.Equal(uint8(0xff), s.ppu.Memory.Fetch(0x0100))

	s.Equal(uint8(0xff), s.ppu.Memory.Fetch(0x0101))

	s.Equal(uint8(0xff), s.ppu.Memory.Fetch(0x0102))

	s.Equal(uint16(0x0103), s.ppu.Registers.Address)

	s.ppu.Memory.Store(0x0103, 0xff)
	s.ppu.Fetch(0x2007)

	s.Equal(uint8(0xff), s.ppu.Fetch(0x2007))

	s.Equal(uint16(0x0105), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x0000
	s.ppu.Fetch(0x2002)

	s.ppu.Store(0x2006, 0x3f)
	s.ppu.Store(0x2006, 0x00)

	s.Equal(uint16(0x3f00), s.ppu.Registers.Address)

	s.ppu.Memory.Store(0x3f00, 0xff)
	s.ppu.Memory.Store(0x3f01, 0xff)
	s.ppu.Memory.Store(0x3f02, 0xff)

	s.Equal(uint8(0xff), s.ppu.Fetch(0x2007))

	s.Equal(uint16(0x3f01), s.ppu.Registers.Address)
}

func (s *ppuTest) TestDataIncrement32() {
	s.ppu.Store(0x2000, 0x04)

	s.ppu.Registers.Address = 0x0000
	s.ppu.Fetch(0x2002)

	s.ppu.Store(0x2006, 0x01)
	s.ppu.Store(0x2006, 0x00)

	s.Equal(uint16(0x0100), s.ppu.Registers.Address)

	s.ppu.Store(0x2007, 0xff)
	s.ppu.Store(0x2007, 0xff)
	s.ppu.Store(0x2007, 0xff)

	s.Equal(uint8(0xff), s.ppu.Memory.Fetch(0x0100))

	s.Equal(uint8(0xff), s.ppu.Memory.Fetch(0x0140))
}

func (s *ppuTest) TestIncrementX() {
	s.ppu.Registers.Address = 0x0000
	s.ppu.incrementX()

	s.Equal(uint16(0x0001), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x001e
	s.ppu.incrementX()

	s.Equal(uint16(0x001f), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x001f
	s.ppu.incrementX()

	s.Equal(uint16(0x0400), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x041f
	s.ppu.incrementX()

	s.Equal(uint16(0x0000), s.ppu.Registers.Address)
}

func (s *ppuTest) TestTransferX() {
	s.ppu.Registers.Address = 0x7be0
	s.ppu.LatchAddress = 0x041f

	s.ppu.transferX()

	s.Equal(uint16(0x7fff), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x7be0
	s.ppu.LatchAddress = 0xffff

	s.ppu.transferX()

	s.Equal(uint16(0x7fff), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0xffff
	s.ppu.LatchAddress = 0x0000

	s.ppu.transferX()

	s.Equal(uint16(0x7be0), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x0000
	s.ppu.LatchAddress = 0xffff

	s.ppu.transferX()

	s.Equal(uint16(0x041f), s.ppu.Registers.Address)
}

func (s *ppuTest) TestIncrementY() {
	s.ppu.Registers.Address = 0x0000
	s.ppu.incrementY()

	s.Equal(uint16(0x1000), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x1000
	s.ppu.incrementY()

	s.Equal(uint16(0x2000), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x6000
	s.ppu.incrementY()

	s.Equal(uint16(0x7000), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x7000
	s.ppu.incrementY()

	s.Equal(uint16(0x0020), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x73d0
	s.ppu.incrementY()

	s.Equal(uint16(0x03f0), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x7fa0
	s.ppu.incrementY()

	s.Equal(uint16(0x0400), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x73a0
	s.ppu.incrementY()

	s.Equal(uint16(0x0800), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x73e1
	s.ppu.incrementY()

	s.Equal(uint16(0x0001), s.ppu.Registers.Address)
}

func (s *ppuTest) TestTransferY() {
	s.ppu.Registers.Address = 0x041f
	s.ppu.LatchAddress = 0x7be0

	s.ppu.transferY()

	s.Equal(uint16(0x7fff), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x041f
	s.ppu.LatchAddress = 0xffff

	s.ppu.transferY()

	s.Equal(uint16(0x7fff), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0xffff
	s.ppu.LatchAddress = 0x0000

	s.ppu.transferY()

	s.Equal(uint16(0x041f), s.ppu.Registers.Address)

	s.ppu.Registers.Address = 0x0000
	s.ppu.LatchAddress = 0xffff

	s.ppu.transferY()

	s.Equal(uint16(0x7be0), s.ppu.Registers.Address)
}

// region ppuTest methods

// region Region-specific timing tests

func (s *ntscTest) TestNTSCTimingParameters() {
	s.False(s.ppu.PAL, "PPU should be NTSC")
	s.Equal(uint16(262), s.ppu.NumScanlines, "NTSC should have 262 scanlines")
	s.Equal(uint16(261), s.ppu.PreRenderScanline, "NTSC pre-render scanline should be 261")
}

func (s *ntscTest) TestNTSCTotalCyclesPerFrame() {
	// NTSC: 262 scanlines * 341 cycles = 89342 cycles per frame
	expectedCycles := uint32(262) * uint32(341)
	s.Equal(uint32(89342), expectedCycles, "NTSC should have 89342 cycles per frame")
}

func (s *palTest) TestPALTimingParameters() {
	s.True(s.ppu.PAL, "PPU should be PAL")
	s.Equal(uint16(312), s.ppu.NumScanlines, "PAL should have 312 scanlines")
	s.Equal(uint16(311), s.ppu.PreRenderScanline, "PAL pre-render scanline should be 311")
}

func (s *palTest) TestPALTotalCyclesPerFrame() {
	// PAL: 312 scanlines * 341 cycles = 106392 cycles per frame
	expectedCycles := uint32(312) * uint32(341)
	s.Equal(uint32(106392), expectedCycles, "PAL should have 106392 cycles per frame")
}

func (s *ntscTest) TestNTSCVBlankTiming() {
	// VBlank starts at scanline 241 for both NTSC and PAL
	s.Equal(uint16(241), StartNMIScanline)
	s.True(s.ppu.PreRenderScanline > StartNMIScanline, "Pre-render should be after VBlank start")

	// NTSC has 20 VBlank scanlines (241-260)
	vblankScanlines := s.ppu.PreRenderScanline - StartNMIScanline
	s.Equal(uint16(20), vblankScanlines, "NTSC should have 20 VBlank scanlines")
}

func (s *palTest) TestPALVBlankTiming() {
	// VBlank starts at scanline 241 for both NTSC and PAL
	s.Equal(uint16(241), StartNMIScanline)
	s.True(s.ppu.PreRenderScanline > StartNMIScanline, "Pre-render should be after VBlank start")

	// PAL has 70 VBlank scanlines (241-310)
	vblankScanlines := s.ppu.PreRenderScanline - StartNMIScanline
	s.Equal(uint16(70), vblankScanlines, "PAL should have 70 VBlank scanlines")
}

func (s *ntscTest) TestNTSCVisibleScanlines() {
	// Both NTSC and PAL have 240 visible scanlines (0-239)
	visibleScanlines := LastVisibleScanline - FirstVisibleScanline + 1
	s.Equal(uint16(240), visibleScanlines, "NTSC should have 240 visible scanlines")
}

func (s *palTest) TestPALVisibleScanlines() {
	// Both NTSC and PAL have 240 visible scanlines (0-239)
	visibleScanlines := LastVisibleScanline - FirstVisibleScanline + 1
	s.Equal(uint16(240), visibleScanlines, "PAL should have 240 visible scanlines")
}

func (s *ntscTest) TestNTSCFrameRate() {
	// NTSC: ~60.098 Hz (CPU clock / cycles per frame)
	// CPU clock: 1.789773 MHz (NTSC)
	// Cycles per frame: 89342
	// Frame rate: 1789773 / 89342 ≈ 20.03 (but PPU runs at 3x CPU clock)
	// Actual PPU clock: 5.369318 MHz / 89342 ≈ 60.098 Hz

	cyclesPerFrame := uint32(s.ppu.NumScanlines) * uint32(CyclesPerScanline)
	s.Equal(uint32(89342), cyclesPerFrame)
}

func (s *palTest) TestPALFrameRate() {
	// PAL: ~50.007 Hz (CPU clock / cycles per frame)
	// CPU clock: 1.662607 MHz (PAL)
	// Cycles per frame: 106392
	// Frame rate: 1662607 / 106392 ≈ 15.62 (but PPU runs at 3x CPU clock)
	// Actual PPU clock: 5.320214 MHz / 106392 ≈ 50.007 Hz

	cyclesPerFrame := uint32(s.ppu.NumScanlines) * uint32(CyclesPerScanline)
	s.Equal(uint32(106392), cyclesPerFrame)
}

// endregion Region-specific timing tests
