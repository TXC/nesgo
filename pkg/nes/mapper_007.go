// ANROM Mapper (INES Mapper 7)
package nes

import (
	"github.com/txc/nesgo/pkg/rp2ago3"
)

type ANROMRegisters struct {
	BankSelect uint8
}

type ANROM struct {
	*ROMFile
	Registers ANROMRegisters
}

func (reg *ANROMRegisters) Reset() {
	reg.BankSelect = 0x00
}

func NewANROM(romf *ROMFile) *ANROM {
	anrom := &ANROM{
		ROMFile: romf,
	}

	anrom.Registers.Reset()
	anrom.setTables(anrom.Tables())
	if anrom.CHRRamSize == 0 {
		anrom.CHRRamSize = 0x2000
	}
	if len(anrom.CHRRamData) != anrom.CHRRamSize {
		anrom.CHRRamData = make([]uint8, anrom.CHRRamSize)
	}

	return anrom
}

func (anrom *ANROM) String() string {
	return anrom.ROMFile.String() + "Mapper: 7 (ANROM)"
}

func (anrom *ANROM) Mappings(which rp2ago3.Mapping) (fetch, store []uint16) {
	fetch = []uint16{}
	store = []uint16{}

	switch which {
	case rp2ago3.PPU:
		//if anrom.CHRBanks > 0 {
		if anrom.CHRRomSize > 0 {
			// CHR bank 1
			for i := uint32(0x0000); i <= 0x0fff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}

			// CHR bank 2
			for i := uint32(0x1000); i <= 0x1fff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}
		}
	case rp2ago3.CPU:
		// WRAM at $6000-$7FFF
		if anrom.PRGNvRamSize > 0 || anrom.PRGRamSize > 0 {
			for i := uint32(0x6000); i <= 0x7fff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}
		}
		//if anrom.PRGBanks > 0 {
		if anrom.PRGRomSize > 0 {
			// PRG bank 1
			for i := uint32(0x8000); i <= 0xbfff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}

			// PRG bank 2
			for i := uint32(0xc000); i <= 0xffff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}
		}
	}

	return
}

func (anrom *ANROM) Reset() {
	anrom.Registers.BankSelect = 0x00
}

func (anrom *ANROM) Fetch(address uint16) (value uint8) {
	switch {
	// PPU only
	//case address >= 0x0000 && address <= 0x1fff:
	case address <= 0x1fff:
		if anrom.CHRRamSize > 0 {
			value = anrom.CHRRamData[address]
		} else if anrom.CHRRomSize > 0 {
			value = anrom.CHRRomData[address]
		}
	// CPU only
	case address >= 0x6000 && address <= 0x7fff:
		value = anrom.FetchPrgRam(address)
	//case address >= 0x8000 && address <= 0xffff:
	case address >= 0x8000:
		// ANROM: 32KB banks, bank select bits 0-2
		// Mask to actual number of banks to prevent overflow
		bankMask := uint32((anrom.PRGRomSize / 0x8000) - 1)
		bank := uint32(anrom.Registers.BankSelect&0x07) & bankMask
		index := bank*0x8000 + uint32(address&0x7fff)
		value = anrom.PRGRomData[index]
		/*
			index := address & 0x3fff
			lower, upper := anrom.prgBanks()

			switch {
			// PRG bank 1
			case address >= 0x8000 && address <= 0xbfff:
				//if anrom.PRGBanks > 0 {
				if anrom.PRGRomSize > 0 {
					//value = anrom.ROMBanks[lower][index]
					value = anrom.ROMBanks[lower][index]
				}
			// PRG bank 2
			//case address >= 0xc000 && address <= 0xffff:
			case address >= 0xc000:
				//if anrom.PRGBanks > 0 {
				if anrom.PRGRomSize > 0 {
					value = anrom.ROMBanks[upper][index]
				}
			}
		*/
	}

	return
}

func (anrom *ANROM) Store(address uint16, value uint8) (oldValue uint8) {
	// PPU only
	switch {
	// CHR banks 1 & 2
	//case address >= 0x0000 && address <= 0x1fff:
	case address <= 0x1fff:
		if anrom.CHRRamSize > 0 {
			anrom.CHRRamData[address] = value
		}
		// CHR ROM writes are ignored
	// CPU WRAM at $6000-$7FFF
	case address >= 0x6000 && address <= 0x7fff:
		anrom.StorePrgRam(address, value)
	// CPU Bank select register
	// PRG banks 1 & 2
	//case address >= 0x8000 && address <= 0xffff:
	case address >= 0x8000:
		oldMirrors := anrom.mirroring()
		anrom.Registers.BankSelect = value

		if anrom.mirroring() != oldMirrors {
			anrom.setTables(anrom.Tables())
		}
	}

	return
}

func (anrom *ANROM) mirroring() uint8 {
	return (anrom.Registers.BankSelect >> 4) & 0x01
}

func (anrom *ANROM) prgBanks() (lower, upper uint8) {
	bank := (anrom.Registers.BankSelect & 0x07) << 1

	lower = bank
	upper = lower + 1

	return
}

func (anrom *ANROM) Tables() (t0, t1, t2, t3 int) {
	switch MMC1Mirroring(anrom.mirroring()) {
	case OneScreenLowerBank:
		t0, t1, t2, t3 = 0, 0, 0, 0
	case OneScreenUpperBank:
		t0, t1, t2, t3 = 1, 1, 1, 1
	case Vertical:
		t0, t1, t2, t3 = 0, 1, 0, 1
	case Horizontal:
		t0, t1, t2, t3 = 0, 0, 1, 1
	}

	return
}
