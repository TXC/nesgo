// CNROM Mapper (INES Mapper 3)
package nes

import (
	"github.com/txc/nesgo/pkg/rp2ago3"
)

type CNROMRegisters struct {
	BankSelect uint8
}

type CNROM struct {
	*ROMFile
	Registers CNROMRegisters
}

func (reg *CNROMRegisters) Reset() {
	reg.BankSelect = 0x00
}

func NewCNROM(romf *ROMFile) *CNROM {
	cnrom := &CNROM{
		ROMFile: romf,
	}

	cnrom.Registers.Reset()

	return cnrom
}

func (cnrom *CNROM) String() string {
	return cnrom.ROMFile.String() + "Mapper: 3 (CNROM)"
}

func (cnrom *CNROM) Mappings(which rp2ago3.Mapping) (fetch, store []uint16) {
	fetch = []uint16{}
	store = []uint16{}

	switch which {
	case rp2ago3.PPU:
		if cnrom.CHRRomSize > 0 {
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
		if cnrom.PRGRomSize > 0 {
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

func (cnrom *CNROM) Reset() {
	cnrom.Registers.BankSelect = 0x00
}

func (cnrom *CNROM) Fetch(address uint16) (value uint8) {
	index := cnrom.getPrgIndex(address)
	switch {
	// PPU only
	//case address >= 0x0000 && address <= 0x1fff:
	case address <= 0x1fff:
		if cnrom.CHRRomSize > 0 {
			//value = cnrom.VROMBanks[cnrom.Registers.BankSelect][address]
			value = cnrom.CHRRomData[index]
		}
	// CPU only
	//case address >= 0x8000 && address <= 0xffff:
	case address >= 0x8000:
		if cnrom.PRGRomSize == 0 {
			return
		}

		switch {
		// PRG bank 1
		case address >= 0x8000 && address <= 0xbfff:
			//value = cnrom.ROMBanks[0][index]
			value = cnrom.PRGRomData[index]
		// PRG bank 2
		//case address >= 0xc000 && address <= 0xffff:
		case address >= 0xc000:
			//value = cnrom.ROMBanks[cnrom.PRGBanks-1][index]
			value = cnrom.PRGRomData[index]
		}
	}

	return
}

func (cnrom *CNROM) Store(address uint16, value uint8) (oldValue uint8) {
	index := cnrom.getPrgIndex(address)
	// PPU only
	switch {
	// CHR banks 1 & 2
	//case address >= 0x0000 && address <= 0x1fff:
	case address <= 0x1fff:
		if cnrom.CHRRomSize > 0 {
			//cnrom.VROMBanks[cnrom.Registers.BankSelect][address] = value
			cnrom.CHRRomData[index] = value
		}
	// CPU only
	// PRG banks 1 & 2
	//case address >= 0x8000 && address <= 0xffff:
	case address >= 0x8000:
		cnrom.Registers.BankSelect = value & 0x03
	}

	return
}

func (cnrom *CNROM) getPrgIndex(address uint16) (index uint16) {
	if address <= 0x2000 {
		index = (uint16(cnrom.Registers.BankSelect) * 0x2000) + address
		return
	}
	if cnrom.PRGRomSize <= 0x4000 {
		index = address & 0x3fff
	} else if cnrom.PRGRomSize <= 0x8000 {
		index = address & 0x7fff
	}
	return index
}
