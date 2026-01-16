package nes

import (
	//"encoding/hex"

	"github.com/txc/nesgo/pkg/rp2cgo2"
	"github.com/txc/nesgo/pkg/utils"
)

type romHeader struct {
	signature [4]byte
	format    Format
	input     []byte
}

func NewRomHeader() *romHeader {
	return &romHeader{
		signature: [4]byte{0x4e, 0x45, 0x53, 0x1a},
		input:     make([]byte, 16),
	}
}

func (rh *romHeader) Format(format Format) *romHeader {
	rh.format = format
	return rh
}

func (rh *romHeader) Add(position int, buf []byte) *romHeader {
	for i := range buf {
		rh.input[position+i] = buf[i]
	}
	return rh
}

func (rh *romHeader) Set(position int, buf byte) *romHeader {
	rh.input[position] = buf
	return rh
}

func (rh *romHeader) Build() []byte {
	buf := make([]byte, 16)
	copy(buf[0:4], rh.signature[0:4])
	copy(buf[4:], rh.input[4:])
	if rh.format == Format_NES20 {
		buf[7] |= 0x08
	}
	return buf
}

func (rh *romHeader) Bytes() []byte { return rh.Build() }

func (rh *romHeader) romData() []byte {
	rf := &ROMFile{}
	_, prgSize := rf.calculateSize(rh.input[4], rh.input[9]&0x0F, 0x4000)
	_, chrSize := rf.calculateSize(rh.input[5], rh.input[9]>>4, 0x2000)

	romSize := prgSize + chrSize
	offset := 16
	size := offset + int(romSize)

	buf := make([]byte, size)
	for i := offset; i < size; i++ {
		switch i {
		case offset, size - 1:
			buf[i] = 0xee
		default:
			buf[i] = 0xff
		}
	}
	return buf
}

// region iNESRomTestSuite
type iNESRomTestSuite struct {
	utils.CommonTest
	romf *ROMFile
}

func (s *iNESRomTestSuite) BeforeTest(suiteName, testName string) {
}

func (s *iNESRomTestSuite) AfterTest(_, _ string) {
	s.romf = nil
}

func (s *iNESRomTestSuite) TestConstant() {
	tests := []struct {
		name  string
		bytes []byte
		want  bool
	}{
		{
			"Valid NES Signature",
			[]byte{0x4e, 0x45, 0x53, 0x1a},
			true,
		},
		{
			"Invalid NES Signature 1",
			[]byte{0x4f, 0x45, 0x53, 0x1a},
			false,
		},
		{
			"Invalid NES Signature 2",
			[]byte{0x4e, 0x46, 0x53, 0x1a},
			false,
		},
		{
			"Invalid NES Signature 3",
			[]byte{0x4e, 0x45, 0x54, 0x1a},
			false,
		},
		{
			"Invalid NES Signature 4",
			[]byte{0x4e, 0x45, 0x53, 0x1b},
			false,
		},
	}
	for _, tt := range tests {
		buf := make([]byte, 0x10)
		copy(buf[0:4], tt.bytes)
		s.Run(tt.name, func() {
			_, err := NewROMFile(buf)
			if tt.want {
				s.NoError(err)
			} else {
				s.Error(err)
			}
		})
	}
}

func (s *iNESRomTestSuite) TestRunt() {
	buf := []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0xff, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00,
	}

	_, err := NewROMFile(buf)
	s.Error(err)
}

func (s *iNESRomTestSuite) TestPrgBanks() {
	r := s.Require()
	buf := make([]byte, 16+(16*1024)-1)

	copy(buf, []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})

	_, err := NewROMFile(buf)
	r.Error(err)

	buf = make([]byte, 16+(16*1024))

	copy(buf, []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})

	n := 16 + (16 * 1024)

	for i := 16; i < n; i++ {
		switch i {
		case 16, n - 1:
			buf[i] = 0xee
		default:
			buf[i] = 0xff
		}
	}

	rom, err := NewROMFile(buf)

	r.NoError(err)

	s.Equal(uint16(0x01), rom.PRGBanks)

	n = (16 * 1024)

	for i := 0; i < n; i++ {
		var val uint8

		switch i {
		case 0, n - 1:
			val = 0xee
		default:
			val = 0xff
		}

		s.Equal(val, rom.ROMBanks[0][i])
	}
}

func (s *iNESRomTestSuite) TestChrBanks() {
	r := s.Require()
	buf := make([]byte, 16+0x2000-1)

	copy(buf, []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})

	_, err := NewROMFile(buf)
	r.Error(err)

	buf = make([]byte, 16+0x2000)

	copy(buf, []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})

	n := 16 + 0x2000

	for i := 16; i < n; i++ {
		switch i {
		case 16:
		case n - 1:
			buf[i] = 0xee
		default:
			buf[i] = 0xff
		}
	}

	for i := 16; i < n; i++ {
		switch i {
		case 16, n - 1:
			buf[i] = 0xee
		default:
			buf[i] = 0xff
		}
	}

	rom, err := NewROMFile(buf)
	s.NoError(err)

	s.Equal(uint16(0x01), rom.CHRBanks)

	n = (8 * 1024)

	for i := 0; i < n; i++ {
		var val uint8

		switch i {
		case 0, n - 1:
			val = 0xee
		default:
			val = 0xff
		}

		s.Equal(val, rom.VROMBanks[0][i])
	}
}

func (s *iNESRomTestSuite) TestMirroring() {
	buf := []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err := NewROMFile(buf)
	s.NoError(err)
	s.Equal(rp2cgo2.Horizontal, rom.Mirroring)

	buf = []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err = NewROMFile(buf)
	s.NoError(err)
	s.Equal(rp2cgo2.Vertical, rom.Mirroring)
}

func (s *iNESRomTestSuite) TestBattery() {
	buf := []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err := NewROMFile(buf)
	s.NoError(err)
	s.False(rom.Battery)

	buf = []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x02, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err = NewROMFile(buf)
	s.NoError(err)
	s.True(rom.Battery)
}

func (s *iNESRomTestSuite) TestTrainer() {
	buf := []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err := NewROMFile(buf)
	s.NoError(err)
	s.False(rom.Trainer)

	buf = make([]byte, 16+512-1)

	copy(buf, []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x04, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})

	rom, err = NewROMFile(buf)
	s.Error(err)
	//s.True(rom.Battery)

	buf = make([]byte, 16+512)

	copy(buf, []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x04, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})

	for i := 16; i < 16+512; i++ {
		buf[i] = 0xff
	}

	rom, err = NewROMFile(buf)
	s.NoError(err)
	s.True(rom.Trainer)

	for i := 0; i < 512; i++ {
		s.Equal(uint8(0xff), rom.TrainerData[i])
	}
}

func (s *iNESRomTestSuite) TestFourScreen() {
	buf := []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err := NewROMFile(buf)
	s.NoError(err)
	s.False(rom.FourScreen)

	buf = []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x08, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err = NewROMFile(buf)
	s.NoError(err)
	s.True(rom.FourScreen)
}

func (s *iNESRomTestSuite) TestMapperLow() {
	buf := []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err := NewROMFile(buf)
	s.NoError(err)
	s.Equal(uint16(0x00), rom.Mapper)

	buf = []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0xf0, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err = NewROMFile(buf)
	s.NoError(err)
	s.Equal(uint16(0x0f), rom.Mapper)
}

func (s *iNESRomTestSuite) TestVsCart() {
	buf := []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err := NewROMFile(buf)
	s.NoError(err)
	//s.False(rom.VSCart)
	s.NotEqual(rom.CartType, CartType_VsSystem)

	buf = []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x08, 0x01,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err = NewROMFile(buf)
	s.NoError(err)
	//s.True(rom.VSCart)
	s.Equal(rom.CartType, CartType_VsSystem)
}

func (s *iNESRomTestSuite) TestMapperHigh() {
	buf := []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err := NewROMFile(buf)
	s.NoError(err)
	s.Equal(uint16(0x00), rom.Mapper)

	buf = []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x00, 0xf0,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err = NewROMFile(buf)
	s.NoError(err)
	s.Equal(uint16(0xf0), rom.Mapper)
}

func (s *iNESRomTestSuite) TestMapper() {
	buf := []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err := NewROMFile(buf)
	s.NoError(err)
	s.Equal(uint16(0x00), rom.Mapper)

	buf = []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0xf0, 0xf0,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err = NewROMFile(buf)
	s.NoError(err)
	s.Equal(uint16(0xff), rom.Mapper)
}

func (s *iNESRomTestSuite) TestRamBanks() {
	tests := []struct {
		name        string
		data        *romHeader
		expectRam   int
		expectNvRam int
	}{
		{
			name:        "Zero RAM Banks",
			data:        NewRomHeader().Add(8, []byte{0x00}),
			expectRam:   int(0x01 * RAMBankSize),
			expectNvRam: 0x00,
		},
		{
			name:        "Two RAM Banks",
			data:        NewRomHeader().Add(8, []byte{0x02}),
			expectRam:   int(0x02 * RAMBankSize),
			expectNvRam: 0x00,
		},
		{
			name:        "Four RAM Banks",
			data:        NewRomHeader().Add(8, []byte{0x04}),
			expectRam:   int(0x04 * RAMBankSize),
			expectNvRam: 0x00,
		},
		{
			name:        "Eight RAM Banks",
			data:        NewRomHeader().Add(8, []byte{0x08}),
			expectRam:   int(0x08) * int(RAMBankSize),
			expectNvRam: 0x00,
		},
		{
			name:        "256 RAM Banks",
			data:        NewRomHeader().Add(8, []byte{0xff}),
			expectRam:   int(0xff) * int(RAMBankSize),
			expectNvRam: 0x00,
		},
		{
			name:        "Check NvRAM Banks",
			data:        NewRomHeader().Add(6, []byte{0x02}).Add(8, []byte{0x08}),
			expectRam:   0x00,
			expectNvRam: 0x08 * int(RAMBankSize),
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			romf, err := NewROMFile(tt.data.Build())
			s.NoError(err)
			if tt.expectRam > 0 && tt.expectNvRam == 0 {
				s.Equal(tt.expectRam, int(romf.PRGRamSize))
				s.Equal(int(romf.PRGRamSize), len(romf.PRGRamData))
			} else if tt.expectRam == 0 && tt.expectNvRam > 0 {
				s.Equal(tt.expectNvRam, int(romf.PRGNvRamSize))
				s.Equal(int(romf.PRGNvRamSize), len(romf.PRGNvRamData))
			}
		})
	}
}

func (s *iNESRomTestSuite) TestRegion() {
	buf := []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err := NewROMFile(buf)
	s.NoError(err)
	s.Equal(NTSC, rom.RegionFlag)

	buf = []byte{
		0x4e, 0x45, 0x53, 0x1a,
		0x00, 0x00, 0x01, 0x00,
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	rom, err = NewROMFile(buf)
	s.NoError(err)
	s.Equal(PAL, rom.RegionFlag)
}

// endregion iNESRomTestSuite

// region NES20RomTestSuite
type NES20RomTestSuite struct {
	utils.CommonTest
	romf *ROMFile
}

func (s *NES20RomTestSuite) BeforeTest(suiteName, testName string) {
}

func (s *NES20RomTestSuite) AfterTest(_, _ string) {
	s.romf = nil
}

func (s *NES20RomTestSuite) TestNES20Signature() {
	tests := []struct {
		name        string
		header      [16]byte
		want        Format
		shouldError bool
	}{
		{"Valid iNES 2.0", [16]byte{'N', 'E', 'S', 0x1A, 0, 0, 0, 0x08}, Format_NES20, false},
		{"Legacy iNES", [16]byte{'N', 'E', 'S', 0x1A, 0, 0, 0, 0x00}, Format_iNES, false},
		{"Invalid Signature", [16]byte{'B', 'A', 'D', 0x1A, 0, 0, 0, 0x08}, Format_Unknown, true},
		{"Corrupt iNES 2.0 (Bits 2-3 are 11)", [16]byte{'N', 'E', 'S', 0x1A, 0, 0, 0, 0x0C}, Format_Unknown, false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			buf := make([]byte, 16)
			copy(buf, tt.header[:])

			romf, err := NewROMFile(buf)
			if tt.shouldError {
				s.Error(err)
				s.Nil(romf)
			} else {
				s.NotNil(romf)
				s.Equal(tt.want, romf.Format)
			}
		})
	}
}

func (s *NES20RomTestSuite) TestNES20PrgRomSize() {
	tests := []struct {
		name string
		data *romHeader
		want int // Expected PRG banks
	}{
		{
			name: "16 KiB (Small)",
			data: NewRomHeader().Format(Format_NES20).Add(4, []byte{0x01}).Add(9, []byte{0x00}),
			want: 1,
		},
		{
			name: "32 KiB (Standard)",
			data: NewRomHeader().Format(Format_NES20).Add(4, []byte{0x02}).Add(9, []byte{0x00}),
			want: 2,
		},
		{
			name: "Large ROM (MSB used)",
			data: NewRomHeader().Format(Format_NES20).Add(4, []byte{0x00}).Add(9, []byte{0x01}),
			want: 256,
		},
		{
			name: "Exponent Case (2^7 * 1)",
			data: NewRomHeader().Format(Format_NES20).Add(4, []byte{0x1C}).Add(9, []byte{0x0F}),
			want: 1,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			buf := tt.data.Build()

			rf := &ROMFile{}
			banks, bytes := rf.calculateSize(buf[4], buf[9]&0x0F, 0x4000)

			s.Equal(tt.want, int(banks))
			// For exponent-multiplier, bytes can be less than one bank size
			// For standard format, bytes should equal banks * unitSize
			if (buf[9] & 0x0F) != 0x0F {
				s.Equal(tt.want*0x4000, bytes)
			}
			rf.parseNES20(buf)
			s.Equal(uint16(tt.want), rf.PRGBanks)
		})
	}
}

func (s *NES20RomTestSuite) TestNES20ChrRomSize() {
	tests := []struct {
		name string
		data *romHeader
		want int // Expected PRG banks
	}{
		{
			name: "8 KiB (Small)",
			data: NewRomHeader().
				Format(Format_NES20).
				Add(5, []byte{0x01}).
				Add(9, []byte{0x00}),
			want: 1,
		},
		{
			name: "16 KiB (Standard)",
			data: NewRomHeader().
				Format(Format_NES20).
				Add(5, []byte{0x02}).
				Add(9, []byte{0x00}),
			want: 2,
		},
		{
			name: "Large ROM (MSB used)",
			data: NewRomHeader().
				Format(Format_NES20).
				Add(5, []byte{0x00}).
				Add(9, []byte{0x10}),
			want: 256,
		},
		{
			name: "Exponent Case (2^7 * 1)",
			data: NewRomHeader().
				Format(Format_NES20).
				Add(5, []byte{0x1C}).
				Add(9, []byte{0xF0}),
			want: 1,
		},
	}

	for _, tt := range tests {
		buf := tt.data.Build()
		s.Run(tt.name, func() {
			rf := &ROMFile{}
			banks, bytes := rf.calculateSize(buf[5], buf[9]>>4, 0x2000)

			s.Equal(tt.want, int(banks))
			// For exponent-multiplier, bytes can be less than one bank size
			// For standard format, bytes should equal banks * unitSize
			if (buf[9] & 0xF0) != 0xF0 {
				s.Equal(tt.want*0x2000, bytes)
			}
			rf.parseNES20(buf)
			s.Equal(uint16(tt.want), rf.CHRBanks)
		})
	}
}

func (s *NES20RomTestSuite) TestNES20SimplePRGAndCHR() {
	header := NewRomHeader().
		Format(Format_NES20).
		Set(4, 0x02).
		Set(5, 0x01).
		Set(9, 0x00)

	buf := header.Build()
	buf = append(buf, header.romData()...)

	rf := &ROMFile{}
	prgBank, _ := rf.calculateSize(buf[4], buf[9]&0x0F, 0x4000)
	chrBank, _ := rf.calculateSize(buf[5], buf[9]>>4, 0x2000)

	s.Logf("PRG Bank: %d", prgBank)
	s.Logf("CHR Bank: %d", chrBank)

	romf, err := NewROMFile(buf)
	s.NoError(err)
	s.NotNil(romf)
	s.Logf("PRG Size: %d bytes, PRG Banks: %d, Length: %d", romf.PRGRomSize, romf.PRGBanks, len(romf.ROMBanks))
	s.Logf("CHR Size: %d bytes, CHR Banks: %d, Length: %d", romf.CHRRomSize, romf.CHRBanks, len(romf.VROMBanks))

}

func (s *NES20RomTestSuite) TestRamBanks() {
	tests := []struct {
		name        string
		data        *romHeader
		expectRam   int
		expectNvRam int
	}{
		{
			name: "Zero PRG-RAM",
			data: NewRomHeader().
				Format(Format_NES20).
				Add(10, []byte{0x00}),
			expectRam:   0x00,
			expectNvRam: 0x00,
		},
		{
			name: "Two PRG-RAM",
			data: NewRomHeader().
				Format(Format_NES20).
				Add(10, []byte{0x02}),
			expectRam:   0x80,
			expectNvRam: 0x00,
		},
		{
			name: "Two PRG-NvRAM",
			data: NewRomHeader().
				Format(Format_NES20).
				Add(10, []byte{0x20}),
			expectRam:   0x00,
			expectNvRam: 0x80,
		},
		{
			name: "Two PRG-NvRAM & PRG-RAM",
			data: NewRomHeader().
				Format(Format_NES20).
				Add(10, []byte{0x22}),
			expectRam:   0x80,
			expectNvRam: 0x80,
		},
		{
			name: "Four PRG-RAM",
			data: NewRomHeader().
				Format(Format_NES20).
				Add(10, []byte{0x04}),
			expectRam:   0x0200,
			expectNvRam: 0x00,
		},
		{
			name: "Eight PRG-RAM",
			data: NewRomHeader().
				Format(Format_NES20).
				Add(10, []byte{0x08}),
			expectRam:   0x2000,
			expectNvRam: 0x00,
		},
		{
			name: "256 PRG-(Nv)RAM",
			data: NewRomHeader().
				Format(Format_NES20).
				Add(10, []byte{0xff}),
			expectRam:   0x100000,
			expectNvRam: 0x100000,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			buf := tt.data.Build()
			//s.Log(hex.Dump(buf))
			romf, err := NewROMFile(buf)
			s.NoError(err)
			s.Equal(tt.expectRam, romf.PRGRamSize)
			s.Equal(tt.expectRam, len(romf.PRGRamData))
			s.Equal(tt.expectNvRam, romf.PRGNvRamSize)
			s.Equal(tt.expectNvRam, len(romf.PRGNvRamData))
		})
	}
}

// endregion NES20RomTestSuite
