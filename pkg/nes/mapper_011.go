// ColorDreams Mapper (INES Mapper 11)
package nes

import (
	"github.com/txc/nesgo/pkg/rp2ago3"
)

type ColorDreamsRegisters struct {
	BankSelect uint8
}

type ColorDreams struct {
	*ROMFile
	Registers ColorDreamsRegisters
}

func (reg *ColorDreamsRegisters) Reset() {
	// Bank select will be initialized by the mapper to last PRG bank
}

func NewColorDreams(romf *ROMFile) *ColorDreams {
	cd := &ColorDreams{
		ROMFile: romf,
	}

	cd.Reset()

	return cd
}

func (cd *ColorDreams) String() string {
	return cd.ROMFile.String() + "Mapper: 11 (ColorDreams)"
}

func (cd *ColorDreams) Mappings(which rp2ago3.Mapping) (fetch, store []uint16) {
	fetch = []uint16{}
	store = []uint16{}

	switch which {
	case rp2ago3.PPU:
		if cd.CHRRomSize > 0 {
			// CHR bank
			for i := uint32(0x0000); i <= 0x1fff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}
		}
	case rp2ago3.CPU:
		if cd.PRGRomSize > 0 {
			// WRAM region $6000-$7FFF - explicitly map to prevent default WRAM
			for i := uint32(0x6000); i < 0x8000; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}
			// PRG bank 1
			for i := uint32(0x8000); i <= 0xffff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}
		}
	}

	return
}

func (cd *ColorDreams) Reset() {
	cd.Registers.Reset()
	// Initialize to last PRG bank (bits 0-1 select PRG bank)
	// For 64KB ROM: 2 banks of 32KB, so last bank = 1
	if cd.PRGRomSize > 0 {
		lastBank := (cd.PRGRomSize / 0x8000) - 1
		cd.Registers.BankSelect = uint8(lastBank & 0x03)
	} else {
		cd.Registers.BankSelect = 0x00
	}
}

func (cd *ColorDreams) Fetch(address uint16) (value uint8) {
	switch {
	// PPU only - CHR ROM/RAM 8KB bank (bits 4-7)
	case address <= 0x1fff:
		if cd.CHRRomSize > 0 {
			// CHR-ROM: bank selected by upper 4 bits
			bank := (cd.Registers.BankSelect >> 4) & 0x0F
			// Wrap to available CHR banks
			numBanks := cd.CHRRomSize / 0x2000
			if numBanks > 0 {
				bank = bank % uint8(numBanks)
			}
			offset := uint32(bank) * 0x2000
			index := uint32(address)
			value = cd.CHRRomData[offset+index]
		} else if cd.CHRRamSize > 0 {
			// CHR-RAM: bank selected by upper 4 bits
			bank := (cd.Registers.BankSelect >> 4) & 0x0F
			offset := int(bank) * 0x2000
			index := int(address)
			// Wrap around if bank exceeds CHR-RAM size
			if cd.CHRRamSize > 0 {
				value = cd.CHRRamData[(offset+index)%cd.CHRRamSize]
			}
		}
	// CPU WRAM region - ColorDreams has no WRAM
	case address >= 0x6000 && address <= 0x7fff:
		value = cd.FetchPrgRam(address)
	// CPU only - PRG ROM 32KB bank (bits 0-1)
	case address >= 0x8000:
		if cd.PRGRomSize > 0 {
			// PRG bank selected by lower 2 bits
			bank := cd.Registers.BankSelect & 0x03
			// Wrap to available banks (e.g., if ROM has 2 banks, bank 3 wraps to bank 1)
			numBanks := cd.PRGRomSize / 0x8000
			if numBanks > 0 {
				bank = bank % uint8(numBanks)
			}
			offset := uint32(bank) * 0x8000
			index := uint32(address - 0x8000)
			value = cd.PRGRomData[offset+index]
		}
	}

	return
}

func (cd *ColorDreams) Store(address uint16, value uint8) (oldValue uint8) {
	switch {
	// PPU only - CHR RAM/ROM writes
	case address <= 0x1fff:
		if cd.CHRRomSize > 0 {
			// CHR-ROM: writes are ignored but we handle them anyway
			bank := (cd.Registers.BankSelect >> 4) & 0x0F
			offset := uint32(bank) * 0x2000
			index := uint32(address)
			if int(offset+index) < cd.CHRRomSize {
				cd.CHRRomData[offset+index] = value
			}
		} else if cd.CHRRamSize > 0 {
			// CHR-RAM: bank selected by upper 4 bits
			bank := (cd.Registers.BankSelect >> 4) & 0x0F
			offset := int(bank) * 0x2000
			index := int(address)
			// Wrap around if bank exceeds CHR-RAM size
			if cd.CHRRamSize > 0 {
				cd.CHRRamData[(offset+index)%cd.CHRRamSize] = value
			}
		}
	// CPU WRAM region - ColorDreams has no WRAM, ignore writes
	case address >= 0x6000 && address <= 0x7fff:
		cd.StorePrgRam(address, value)
		return
	// CPU only - Bank select register
	case address >= 0x8000:
		// Write to any address in $8000-$FFFF sets bank
		// Bits 0-1: PRG bank (32KB)
		// Bits 4-7: CHR bank (8KB)
		cd.Registers.BankSelect = value
	}

	return
}
