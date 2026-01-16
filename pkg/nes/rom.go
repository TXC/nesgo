package nes

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"archive/zip"

	"strings"

	//"github.com/txc/nesgo/nes/rom/mapper"
	"github.com/txc/nesgo/pkg/rp2ago3"
	"github.com/txc/nesgo/pkg/rp2cgo2"
)

func init() {
	gob.Register(&ROMFile{})

	gob.Register(&ANROM{})
	gob.Register(&CNROM{})
	gob.Register(&MMC1{})
	gob.Register(&MMC2{})
	gob.Register(&MMC3{})
	gob.Register(&NROM{})
	gob.Register(&UNROM{})
}

//go:generate stringer -type=Region
type Region uint8

const (
	NTSC Region = iota
	PAL
	Multi
	Dendy
	Unknown
)

//go:generate stringer -type=CartType
type CartType uint8

const (
	CartType_NES          CartType = iota // Nintendo Entertainment System
	CartType_VsSystem                     // Nintendo Vs. System
	CartType_Playchoice10                 // Nintendo Playchoice-10
	CartType_Extended                     // Extended Console Type
)

//go:generate stringer -type=Format
type Format uint8

const (
	Format_iNES       Format = iota // iNES
	Format_ArchaicNES               // Archaic iNES
	Format_NES20                    // NES2.0
	Format_Unknown    = 0xff        // Unknown
)

func RegionFromString(s string) Region {
	switch s {
	case "NTSC":
		return NTSC
	case "PAL":
		return PAL
	case "Multi":
		return Multi
	case "Dendy":
		return Dendy
	default:
		return Unknown
	}
}

const (
	PRGBankSize uint16 = 0x4000 // 16KB
	CHRBankSize uint16 = 0x2000 // 8KB
	RAMBankSize uint16 = 0x2000 // 8KB
)

type ROMFile struct {
	Format          Format
	Gamename        string
	PRGBanks        uint16 // Set by mapper
	PRGRomSize      int
	CHRBanks        uint16 // Set by mapper
	CHRRomSize      int
	Mirroring       rp2cgo2.Mirroring
	Battery         bool
	Trainer         bool
	FourScreen      bool
	CartType        CartType
	Mapper          uint16
	SubMapper       uint8
	RegionFlag      Region
	systemType      uint8
	MiscRoms        uint8
	ExpansionDevice uint8
	TrainerData     []uint8
	ROMBanks        [][]uint8 // Banked PRG-ROM Data (Set by mapper)
	VROMBanks       [][]uint8 // Banked CHR-ROM Data (Set by mapper)
	PRGRomData      []uint8
	CHRRomData      []uint8
	PRGRamSize      int
	PRGNvRamSize    int
	CHRRamSize      int
	CHRNvRamSize    int

	PRGRamData   []uint8
	PRGNvRamData []uint8
	CHRRamData   []uint8
	CHRNvRamData []uint8
	//WRAMData        []uint8
	//VRAMData        []uint8
	irq       func(state bool)
	setTables func(t0, t1, t2, t3 int)
}

type ROM interface {
	rp2ago3.MappableMemory
	Region() Region
	String() string
	GameName() string
	LoadBattery()
	SaveBattery() (err error)
	GetROMFile() *ROMFile
}

func getBuf(filename string) (buf []byte, suffix string, err error) {
	var r *zip.ReadCloser
	var rc io.ReadCloser

	switch {
	case strings.HasSuffix(filename, ".nes") || strings.HasSuffix(filename, ".NES"):
		suffix = filename[len(filename)-len(".nes"):]
		buf, err = os.ReadFile(filename)
	case strings.HasSuffix(filename, ".zip") || strings.HasSuffix(filename, ".ZIP"):
		suffix = filename[len(filename)-len(".zip"):]

		// Open a zip archive for reading.
		r, err = zip.OpenReader(filename)

		if err != nil {
			return
		}

		defer r.Close()

		// Iterate through the files in the archive,
		// printing some of their contents.
		for _, f := range r.File {
			if !strings.HasSuffix(f.Name, ".nes") && !strings.HasSuffix(f.Name, ".NES") {
				continue
			}

			rc, err = f.Open()

			if err != nil {
				return
			}

			buf, err = io.ReadAll(rc)

			if err != nil {
				return
			}

			rc.Close()
			break
		}
	default:
		err = errors.New("Unknown filetype, must be .nes, .NES, .zip or .ZIP")
	}

	return
}

func NewROM(filename string, irq func(state bool), setTables func(t0, t1, t2, t3 int)) (rom ROM, err error) {
	buf, suffix, err := getBuf(filename)

	if err != nil {
		return nil, err
	}

	return NewROMFromBuf(buf, filename, suffix, irq, setTables)
}

func NewROMFromBuf(buf []byte, filename, suffix string, irq func(state bool), setTables func(t0, t1, t2, t3 int)) (rom ROM, err error) {
	romf, err := NewROMFile(buf)

	if err != nil {
		return
	}

	romf.irq = irq
	romf.setTables = setTables

	romf.setTables(romf.Tables())
	romf.Gamename = strings.TrimSuffix(filename, suffix)

	switch romf.Mapper {
	case 0x00, 0x40, 0x41:
		rom = NewNROM(romf)
	case 0x01:
		rom = NewMMC1(romf)
	case 0x02, 0x42:
		rom = NewUNROM(romf)
	case 0x03, 0x43:
		rom = NewCNROM(romf)
	case 0x04:
		rom = NewMMC3(romf)
	case 0x07:
		rom = NewANROM(romf)
	case 0x09:
		rom = NewMMC2(romf)
	default:
		err = fmt.Errorf("Unsupported mapper type %v", romf.Mapper)
	}

	return
}

func NewROMFile(buf []byte) (romf *ROMFile, err error) {
	var offset int

	if len(buf) < 16 {
		err = errors.New("Invalid ROM: Missing 16-byte header")
		return
	}

	if buf[0] != 0x4E || buf[1] != 0x45 || buf[2] != 0x53 || buf[3] != 0x1a {
		err = errors.New("Invalid ROM: Missing 'NES' constant in header")
		return
	}

	romf = &ROMFile{}

	switch buf[7] & 0x0C {
	case 0x08:
		romf.Format = Format_NES20
		romf.parseNES20(buf)
	case 0x04:
		romf.Format = Format_ArchaicNES
		romf.parseiNES(buf)
	default:
		if buf[7]&0x0C == 0x00 && buf[12] == 0 && buf[13] == 0 && buf[14] == 0 && buf[15] == 0 {
			romf.Format = Format_iNES
		} else {
			romf.Format = Format_Unknown
		}
		romf.parseiNES(buf)
	}
	//i += 6
	i := 16

	if romf.Trainer {
		offset = 512

		if len(buf) < (i + offset) {
			romf = nil
			err = errors.New("Invalid ROM: EOF in trainer data")
			return
		}

		romf.TrainerData = buf[i : i+offset]
		i += offset
	}
	bufLen := len(buf)

	if bufLen < (i + int(romf.PRGRomSize)) {
		errmsg := fmt.Sprintf(
			"Invalid ROM: EOF in PRG-ROM bank data. "+
				"Expected %v bytes, got %d bytes",
			romf.PRGRomSize, bufLen-i)
		romf = nil
		err = errors.New(errmsg)
		return
	}

	romf.ROMBanks = make([][]uint8, romf.PRGBanks)

	for n := 0; n < int(romf.PRGBanks); n++ {
		romf.ROMBanks[n] = buf[i : i+int(PRGBankSize)]
		i += int(PRGBankSize)
	}

	if bufLen < (i + int(romf.CHRRomSize)) {
		errmsg := fmt.Sprintf("Invalid ROM: EOF in CHR-ROM bank data. "+
			"Expected %v bytes, got %d bytes",
			romf.CHRRomSize, bufLen-i)
		romf = nil
		err = errors.New(errmsg)
		return
	}

	romf.VROMBanks = make([][]uint8, romf.CHRBanks)

	for n := 0; n < int(romf.CHRBanks); n++ {
		romf.VROMBanks[n] = buf[i : i+int(CHRBankSize)]
		i += int(CHRBankSize)
	}

	// PRGRamSize
	if romf.PRGRamSize > 0 {
		romf.PRGRamData = make([]uint8, int(romf.PRGRamSize))
	}

	// PRGNvRamSize
	if romf.PRGNvRamSize > 0 {
		romf.PRGNvRamData = make([]uint8, int(romf.PRGNvRamSize))
	}

	// CHRRamSize
	if romf.CHRRamSize > 0 {
		romf.CHRRamData = make([]uint8, int(romf.CHRRamSize))
	}

	// CHRNvRamSize
	if romf.CHRNvRamSize > 0 {
		romf.CHRNvRamData = make([]uint8, int(romf.CHRNvRamSize))
	}

	return
}

// Parse iNES format ROM
//
//	Bytes Description
//	0-3   Constant $4E $45 $53 $1A (ASCII "NES" followed by MS-DOS end-of-file)
//	4     PRG ROM size in 16 KB units
//	5     CHR ROM size in 8 KB units (value 0 means the board uses CHR RAM)
//	6     Flags 6 – Mapper, mirroring, battery, trainer
//	7     Flags 7 – Mapper, VS/Playchoice, NES 2.0
//	8     Flags 8 – PRG-RAM size (rarely used extension)
//	9     Flags 9 – TV system (rarely used extension)
//	10    Flags 10 – TV system, PRG-RAM presence (unofficial, rarely used extension)
//	11-15 Unused padding (should be filled with zero, but some rippers put their name across bytes 7-15)
func (romf *ROMFile) parseiNES(buf []byte) {
	// Use calculateSize for consistent handling between iNES and NES 2.0
	// For iNES, MSB is 0, so this just converts bytes to banks
	romf.PRGBanks, romf.PRGRomSize = romf.calculateSize(buf[4], 0, int(PRGBankSize))
	romf.CHRBanks, romf.CHRRomSize = romf.calculateSize(buf[5], 0, int(CHRBankSize))

	//i := 4
	i := 6
	for ; i < 10; i++ {
		b := buf[i]

		switch i {
		case 4, 5:
			// PRG/CHR banks already handled above
			continue
		case 6:
			for j := range 4 {
				if b&(0x01<<uint8(j)) != 0 {
					switch j {
					case 0:
						romf.Mirroring = rp2cgo2.Vertical
					case 1:
						romf.Battery = true
					case 2:
						romf.Trainer = true
					case 3:
						romf.FourScreen = true
						romf.Mirroring = rp2cgo2.FourScreen
					}
				}
			}

			romf.Mapper = uint16((b >> 4) & 0x0f)
		case 7:
			//if b&0x01 != 0 {
			//	romf.VSCart = true
			//}
			if b&0x01 != 0 && b&0x02 == 0 {
				romf.CartType = CartType_VsSystem
			}
			// The PlayChoice-10 bit is not part of the official specification,
			// and most emulators simply ignore the extra 8 KB of data.
			// PlayChoice games are designed to look good with the 2C03 RGB PPU,
			// which handles color emphasis differently from a standard NES PPU.
			if b&0x01 == 0 && b&0x02 != 0 {
				romf.CartType = CartType_Playchoice10
			}

			romf.Mapper |= uint16(b & 0xf0)

		case 8:
			if romf.Battery {
				if b > 0 {
					romf.PRGNvRamSize = int(b) * int(RAMBankSize)
				} else {
					romf.PRGNvRamSize = int(RAMBankSize)
				}
			} else {
				if b > 0 {
					romf.PRGRamSize = int(b) * int(RAMBankSize)
				} else {
					romf.PRGRamSize = int(RAMBankSize)
				}
			}

		case 9:
			if b&0x01 != 0 {
				romf.RegionFlag = PAL
			}
		}
	}
}

// Parse NES 2.0 format ROM
//
//	Bytes Description
//	0-3   Constant $4E $45 $53 $1A (ASCII "NES" followed by MS-DOS end-of-file)
//	4     PRG ROM size LSB
//	5     CHR ROM size LSB
//	6     Flags 6 – Mapper LSB (first 4 bits), mirroring, battery, trainer
//	7     Flags 7 – Mapper LSB (second 4 bits), VS/Playchoice, NES 2.0
//	8     Flags 8 – Submapper, Mapper MSB
//	9     Flags 9 – PRG ROM MSB, CHR ROM MSB
//	10    Flags 10 – PRG-RAM shift count, PRG-NVRAM shift count
//	11    Flags 11 - CHR-RAM size
//	12    Flags 12 - CPU/PPU timing
//	13    Flags 13 - Vs. System Type / Extended Console Type
//	14    Flags 14 - Misc. ROMs
//	15    Flags 15 - Default Expansion Device
func (romf *ROMFile) parseNES20(buf []byte) {
	romf.PRGBanks, romf.PRGRomSize = romf.calculateSize(buf[4], buf[9]&0x0F, int(PRGBankSize))
	romf.CHRBanks, romf.CHRRomSize = romf.calculateSize(buf[5], buf[9]>>4, int(CHRBankSize))

	//i := 4
	i := 6
	for ; i < 16; i++ {
		b := buf[i]

		switch i {
		case 4, 5:
			// PRG/CHR banks already handled above
			continue
		case 6:
			for j := range 4 {
				if b&(0x01<<uint8(j)) != 0 {
					switch j {
					case 0:
						romf.Mirroring = rp2cgo2.Vertical
					case 1:
						romf.Battery = true
					case 2:
						romf.Trainer = true
					case 3:
						romf.FourScreen = true
						romf.Mirroring = rp2cgo2.FourScreen
					}
				}
			}

			romf.Mapper = uint16((b >> 4) & 0x0f)
		case 7:
			switch b & 0x03 {
			case 0x00:
				romf.CartType = CartType_NES
			case 0x01:
				romf.CartType = CartType_VsSystem
			case 0x02:
				romf.CartType = CartType_Playchoice10
			case 0x03:
				romf.CartType = CartType_Extended
			}

			romf.Mapper |= uint16(b & 0xf0)
		case 8:
			//romf.RAMBanks = b
			//
			//if romf.RAMBanks == 0 {
			//	romf.RAMBanks = 1
			//}
			romf.Mapper |= uint16(b&0x0f) << 8
			romf.SubMapper = uint8(b&0xf0) >> 4
		case 9:
			continue
		case 10:
			romf.PRGRamSize = romf.calculateRAMSize(b & 0x0F)
			romf.PRGNvRamSize = romf.calculateRAMSize(b & 0xF0 >> 4)
		case 11:
			romf.CHRRamSize = romf.calculateRAMSize(b & 0x0F)
			romf.CHRNvRamSize = romf.calculateRAMSize(b & 0xF0 >> 4)
		case 12:
			switch b & 0x03 {
			case 0:
				romf.RegionFlag = NTSC
			case 1:
				romf.RegionFlag = PAL
			case 2:
				romf.RegionFlag = Multi
			case 3:
				romf.RegionFlag = Dendy
			default:
				romf.RegionFlag = Unknown
			}
		case 13:
			romf.systemType = b
		case 14:
			romf.MiscRoms = b
		case 15:
			romf.ExpansionDevice = b
		}
	}
}

// calculateSize handles the iNES 2.0 exponent-multiplier formula
//
// Returns size in banks & size in bytes
func (romf *ROMFile) calculateSize(lsb byte, msb byte, unitSize int) (banks uint16, bytes int) {
	// If MSB nibble is not 0xF, use standard units
	if msb != 0x0F {
		banks = (uint16(msb) << 8) | uint16(lsb)
		bytes = int(banks) * unitSize
		return
	}

	// Exponent-multiplier notation
	exponent := float64((lsb >> 2) & 0x3F)
	multiplier := uint16(lsb & 0x03)

	// Formula: 2^E * (MM * 2 + 1)
	//bytes = uint64(uint16(math.Pow(2, exponent)) * (multiplier*2 + 1))
	bytes = int(math.Pow(2, exponent)) * int(multiplier*2+1)

	// Round up: if there's any data at all, we need at least 1 bank
	// Use ceiling division: (bytes + unitSize - 1) / unitSize
	banks = uint16(math.Max(float64(bytes/unitSize), 1))
	//if bytes > 0 {
	//	banks = uint16((bytes + unitSize - 1) / unitSize)
	//}
	return
}

func (romf *ROMFile) calculateRAMSize(nibble byte) int {
	if nibble == 0 {
		return 0
	}
	// 64 * 2^(n-1) is the same as shifting 64 left by (n-1)
	return int(64) << (nibble - 1)
}

// Region returns the region flag of the ROM
//
// For licensed games, this is usually NTSC or PAL.
//   - NTSC (RP2C02) - is used when a particular PRG-ROM were released in
//     North America and Japan
//   - PAL (RP2C07) - is used when a particular PRG-ROM were released in
//     Europe, Australia, Oceania, and Hong Kong versions.
//   - Multi - is used when the exact PRG-ROM was released in both NTSC
//     and PAL regions.
//   - Dendy - is not used for licensed games at all.
func (romf *ROMFile) Region() Region {
	return romf.RegionFlag
}

func (romf *ROMFile) Tables() (t0, t1, t2, t3 int) {
	switch romf.Mirroring {
	case rp2cgo2.Horizontal:
		t0, t1, t2, t3 = 0, 0, 1, 1
	case rp2cgo2.Vertical:
		t0, t1, t2, t3 = 0, 1, 0, 1
	}

	return
}

func (romf *ROMFile) String() string {
	return fmt.Sprintf("PRG Banks: %v\n"+
		"CHR Banks: %v\n"+
		"Mirroring: %v\n"+
		"Battery: %v\n"+
		"Trainer: %v\n"+
		"FourScreen: %v\n"+
		"Cart Type: %v\n"+
		"PRG RAM size: %v bytes\n"+
		"CHR RAM size: %v bytes\n"+
		"Region: %v\n",
		romf.PRGBanks, romf.CHRBanks, romf.Mirroring,
		romf.Battery, romf.Trainer, romf.FourScreen,
		romf.CartType, romf.PRGRamSize, romf.CHRRamSize,
		romf.RegionFlag)
}

func (romf *ROMFile) GameName() string {
	return romf.Gamename
}

func (romf *ROMFile) SystemType() {
	if romf.CartType == CartType_VsSystem {
		switch romf.systemType & 0x0F {
		case 0x00:
			// Any RP2C03/RC2C03 variant
		case 0x02:
			// RP2C04-0001
		case 0x03:
			// RP2C04-0002
		case 0x04:
			// RP2C04-0003
		case 0x05:
			// RP2C04-0004
		case 0x08:
			// RC2C05-01 (signature unknown)
		case 0x09:
			// RC2C05-02 ($2002 AND $3F =$3D)
		case 0x0A:
			// RC2C05-03 ($2002 AND $1F =$1C)
		case 0x0B:
			// RC2C05-04 ($2002 AND $1F =$1B)
		}

		switch (romf.systemType & 0xF0) >> 4 {
		case 0x00:
			// Vs. Unisystem (normal)
		case 0x01:
			// Vs. Unisystem (RBI Baseball protection)
		case 0x02:
			// Vs. Unisystem (TKO Boxing protection)
		case 0x03:
			// Vs. Unisystem (Super Xevious protection)
		case 0x04:
			// Vs. Unisystem (Vs. Ice Climber Japan protection)
		case 0x05:
			// Vs. Dual System (normal)
		case 0x06:
			// Vs. Dual System (Raid on Bungeling Bay protection)
		}
	}

	if romf.CartType == CartType_Extended {
		switch romf.systemType & 0x0F {
		case 0x00:
			// Regular NES/Famicom/Dendy
		case 0x01:
			// Nintendo Vs. System
		case 0x02:
			// Playchoice 10
		case 0x03:
			// Regular Famiclone, but with CPU that supports Decimal Mode
		case 0x04:
			// Regular NES/Famicom with EPSM module or plug-through cartridge
		case 0x05:
			// V.R. Technology VT01 with red/cyan STN palette
		case 0x06:
			// V.R. Technology VT02
		case 0x07:
			// V.R. Technology VT03
		case 0x08:
			// V.R. Technology VT09
		case 0x09:
			// V.R. Technology VT32
		case 0x0A:
			// V.R. Technology VT369
		case 0x0B:
			// UMC UM6578
		case 0x0C:
			// Famicom Network System
		}
	}
}

func (romf *ROMFile) LoadBattery() {
	var ram []byte

	if !romf.Battery || romf.PRGNvRamSize == 0 {
		return
	}

	savename := romf.Gamename + ".sav"
	ram, err := os.ReadFile(savename)

	if err != nil {
		return
	}

	fmt.Println("*** Loading battery from " + savename)

	for i := range romf.PRGNvRamSize {
		romf.PRGNvRamData[i] = ram[i]
	}
	//return
}

func (romf *ROMFile) SaveBattery() (err error) {
	if !romf.Battery || romf.PRGNvRamSize == 0 {
		return
	}

	savename := romf.Gamename + ".sav"

	fmt.Println("*** Saving battery to " + savename)

	buf := bytes.Buffer{}

	for i := range romf.PRGNvRamSize {
		buf.WriteByte(romf.PRGNvRamData[i])
	}

	err = os.WriteFile(savename, buf.Bytes(), 0644)

	return
}

func (romf *ROMFile) GetROMFile() *ROMFile {
	return romf
}
