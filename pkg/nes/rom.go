package nes

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
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
	Unknown
)

func RegionFromString(s string) Region {
	switch s {
	case "NTSC":
		return NTSC
	case "PAL":
		return PAL
	default:
		return Unknown
	}
}

type ROMFile struct {
	Gamename    string
	PRGBanks    uint16
	CHRBanks    uint16
	Mirroring   rp2cgo2.Mirroring
	Battery     bool
	Trainer     bool
	FourScreen  bool
	VSCart      bool
	Mapper      uint8
	RAMBanks    uint8
	RegionFlag  Region
	TrainerData []uint8
	WRAMBanks   [][]uint8
	ROMBanks    [][]uint8
	VROMBanks   [][]uint8
	irq         func(state bool)
	setTables   func(t0, t1, t2, t3 int)
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

	if buf[7]&0x0C == 0x08 {
		romf.parseNES20(buf)
	} else {
		romf.parseiNES(buf)
	}
	//i += 6
	i := 15

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

	offset = 1024 * 16

	if len(buf) < (i + (offset * int(romf.PRGBanks))) {
		romf = nil
		err = errors.New("Invalid ROM: EOF in ROM bank data")
		return
	}

	romf.ROMBanks = make([][]uint8, romf.PRGBanks)

	for n := 0; n < int(romf.PRGBanks); n++ {
		romf.ROMBanks[n] = buf[i : i+offset]
		i += offset
	}

	offset = 1024 * 8

	if len(buf) < (i + (offset * int(romf.CHRBanks))) {
		romf = nil
		err = errors.New("Invalid ROM: EOF in VROM bank data")
		return
	}

	romf.VROMBanks = make([][]uint8, romf.CHRBanks)

	for n := 0; n < int(romf.CHRBanks); n++ {
		romf.VROMBanks[n] = buf[i : i+offset]
		i += offset
	}

	offset = 1024 * 8

	romf.WRAMBanks = make([][]uint8, romf.RAMBanks)

	for n := 0; n < int(romf.RAMBanks); n++ {
		romf.WRAMBanks[n] = make([]uint8, offset)
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
	i := 4

	for ; i < 10; i++ {
		b := buf[i]

		switch i {
		case 4:
			romf.PRGBanks = uint16(b)
		case 5:
			romf.CHRBanks = uint16(b)
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

			romf.Mapper = (b >> 4) & 0x0f
		case 7:
			if b&0x01 != 0 {
				romf.VSCart = true
			}

			romf.Mapper |= b & 0xf0

		case 8:
			romf.RAMBanks = b

			if romf.RAMBanks == 0 {
				romf.RAMBanks = 1
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
}

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
		"VS Cart: %v\n"+
		"RAM Banks: %v\n"+
		"Region: %v\n",
		romf.PRGBanks, romf.CHRBanks, romf.Mirroring,
		romf.Battery, romf.Trainer, romf.FourScreen,
		romf.VSCart, romf.RAMBanks, romf.RegionFlag)
}

func (romf *ROMFile) GameName() string {
	return romf.Gamename
}

func (romf *ROMFile) LoadBattery() {
	var ram []byte

	if !romf.Battery || romf.RAMBanks == 0 {
		return
	}

	savename := romf.Gamename + ".sav"
	ram, err := os.ReadFile(savename)

	if err != nil {
		return
	}

	fmt.Println("*** Loading battery from " + savename)

	for b := range romf.WRAMBanks {
		for i := range uint16(0x2000) {
			romf.WRAMBanks[b][i] = ram[i]
		}
	}
	//return
}

func (romf *ROMFile) SaveBattery() (err error) {
	if !romf.Battery || romf.RAMBanks == 0 {
		return
	}

	savename := romf.Gamename + ".sav"

	fmt.Println("*** Saving battery to " + savename)

	buf := bytes.Buffer{}

	for b := range romf.WRAMBanks {
		for i := range uint16(0x2000) {
			buf.WriteByte(romf.WRAMBanks[b][i])
		}
	}

	err = os.WriteFile(savename, buf.Bytes(), 0644)

	return
}

func (romf *ROMFile) GetROMFile() *ROMFile {
	return romf
}
