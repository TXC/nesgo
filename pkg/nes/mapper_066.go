// MHROM Mapper (INES Mapper 66)
package nes

import (
	"github.com/txc/nesgo/pkg/rp2ago3"
)

type MHROMRegisters struct {
	BankSelect uint8
}

type MHROM struct {
	*ROMFile
	Registers MHROMRegisters
}

func (reg *MHROMRegisters) Reset() {
	// Bank select will be initialized by the mapper to last PRG bank
}

func NewMHROM(romf *ROMFile) *MHROM {
	mhrom := &MHROM{
		ROMFile: romf,
	}

	mhrom.Reset()

	return mhrom
}

func (mhrom *MHROM) String() string {
	return mhrom.ROMFile.String() + "Mapper: 66 (MHROM)"
}

func (mhrom *MHROM) Mappings(which rp2ago3.Mapping) (fetch, store []uint16) {
	fetch = []uint16{}
	store = []uint16{}

	switch which {
	case rp2ago3.PPU:
		if mhrom.CHRRomSize > 0 {
			// CHR bank
			for i := uint32(0x0000); i <= 0x1fff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}
		}
	case rp2ago3.CPU:
		if mhrom.PRGRomSize > 0 {
			// WRAM region $6000-$7FFF - explicitly map to prevent default WRAM
			for i := uint32(0x6000); i < 0x8000; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}
			// PRG bank
			for i := uint32(0x8000); i <= 0xffff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}
		}
	}

	return
}

func (mhrom *MHROM) Reset() {
	mhrom.Registers.Reset()

	// Initialize to last PRG bank (bits 4-5 select PRG bank)
	// For 64KB ROM: 2 banks of 32KB, so last bank = 1
	if mhrom.PRGRomSize > 0 {
		lastBank := (mhrom.PRGRomSize / 0x8000) - 1
		mhrom.Registers.BankSelect = uint8((lastBank & 0x03) << 4)
	} else {
		mhrom.Registers.BankSelect = 0x00
	}
}

func (mhrom *MHROM) Fetch(address uint16) (value uint8) {
	switch {
	// PPU only - CHR ROM 8KB bank (bits 0-1)
	case address <= 0x1fff:
		if mhrom.CHRRomSize > 0 {
			// CHR-ROM: bank selected by lower 2 bits
			bank := mhrom.Registers.BankSelect & 0x03
			// Wrap to available CHR banks
			numBanks := mhrom.CHRRomSize / 0x2000
			if numBanks > 0 {
				bank = bank % uint8(numBanks)
			}
			offset := uint32(bank) * 0x2000
			index := uint32(address)
			value = mhrom.CHRRomData[offset+index]
		} else if mhrom.CHRRamSize > 0 {
			// CHR-RAM: bank selected by lower 2 bits
			bank := mhrom.Registers.BankSelect & 0x03
			offset := int(bank) * 0x2000
			index := int(address)
			if mhrom.CHRRamSize > 0 {
				value = mhrom.CHRRamData[(offset+index)%mhrom.CHRRamSize]
			}
		}
	// CPU WRAM region - MHROM has no WRAM
	case address >= 0x6000 && address <= 0x7fff:
		value = mhrom.FetchPrgRam(address)
	// CPU only - PRG ROM 32KB bank (bits 4-5)
	case address >= 0x8000:
		if mhrom.PRGRomSize > 0 {
			// PRG bank selected by bits 4-5
			bank := (mhrom.Registers.BankSelect >> 4) & 0x03
			// Wrap to available banks
			numBanks := mhrom.PRGRomSize / 0x8000
			if numBanks > 0 {
				bank = bank % uint8(numBanks)
			}
			offset := uint32(bank) * 0x8000
			index := uint32(address - 0x8000)
			value = mhrom.PRGRomData[offset+index]
		}
	}

	return
}

func (mhrom *MHROM) Store(address uint16, value uint8) (oldValue uint8) {
	switch {
	// PPU only - CHR ROM/RAM writes
	case address <= 0x1fff:
		if mhrom.CHRRomSize > 0 {
			// CHR-ROM: writes are typically ignored
			bank := mhrom.Registers.BankSelect & 0x03
			offset := uint32(bank) * 0x2000
			index := uint32(address)
			if int(offset+index) < mhrom.CHRRomSize {
				mhrom.CHRRomData[offset+index] = value
			}
		} else if mhrom.CHRRamSize > 0 {
			// CHR-RAM: bank selected by lower 2 bits
			bank := mhrom.Registers.BankSelect & 0x03
			offset := int(bank) * 0x2000
			index := int(address)
			if mhrom.CHRRamSize > 0 {
				mhrom.CHRRamData[(offset+index)%mhrom.CHRRamSize] = value
			}
		}
	// CPU WRAM region - MHROM has no WRAM, ignore writes
	case address >= 0x6000 && address <= 0x7fff:
		// Ignore writes to indicate no WRAM
		mhrom.StorePrgRam(address, value)
	// CPU only - Bank select register
	case address >= 0x8000:
		// Write to any address in $8000-$FFFF sets bank
		// Bits 0-1: CHR bank (8KB)
		// Bits 4-5: PRG bank (32KB)
		mhrom.Registers.BankSelect = value
	}

	return
}
