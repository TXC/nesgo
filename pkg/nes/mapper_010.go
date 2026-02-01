// MMC4 Mapper (INES Mapper 10)
package nes

import (
	"github.com/txc/nesgo/pkg/rp2ago3"
	"github.com/txc/nesgo/pkg/rp2cgo2"
)

type MMC4Registers struct {
	PRGBank   uint8
	CHRBank0  uint8
	CHRBank1  uint8
	CHRBank2  uint8
	CHRBank3  uint8
	Mirroring uint8
	Latch0    uint8
	Latch1    uint8
}

type MMC4 struct {
	*ROMFile
	Registers MMC4Registers
}

func (reg *MMC4Registers) Reset() {
	reg.PRGBank = 0x00
	reg.CHRBank0 = 0x00
	reg.CHRBank1 = 0x00
	reg.CHRBank2 = 0x00
	reg.CHRBank3 = 0x00
	reg.Mirroring = 0x00
	reg.Latch0 = 0xfd
	reg.Latch1 = 0xfe
}

func NewMMC4(romf *ROMFile) *MMC4 {
	mmc4 := &MMC4{
		ROMFile: romf,
	}

	// divide 8KB CHR banks into 4KB banks since we may be
	// swapping 4KB banks
	if romf.CHRRomSize > 0 {
		// Calculate number of 4KB banks (CHR comes in 8KB units typically)
		banks := uint16(romf.CHRBanks) * 2
		vromBanks := make([][]uint8, banks)

		for n := 0; n < int(romf.CHRBanks); n++ {
			vromBanks[2*n] = romf.CHRRomData[n*0x2000 : n*0x2000+0x1000]
			vromBanks[(2*n)+1] = romf.CHRRomData[n*0x2000+0x1000 : (n+1)*0x2000]
		}

		romf.VROMBanks = vromBanks
		romf.CHRBanks = banks
	}

	// divide 16KB PRG banks into 8KB banks since we may be
	// swapping 8KB banks
	if romf.PRGBanks > 0 {
		romBanks := make([][]uint8, uint16(romf.PRGBanks)*2)

		for n := 0; n < int(romf.PRGBanks); n++ {
			romBanks[2*n] = romf.ROMBanks[n][0x0000:0x2000]
			romBanks[(2*n)+1] = romf.ROMBanks[n][0x2000:0x4000]
		}

		romf.ROMBanks = romBanks
		romf.PRGBanks *= 2
	}

	mmc4.Registers.Reset()
	mmc4.setTables(mmc4.Tables())

	return mmc4
}

func (mmc4 *MMC4) String() string {
	return mmc4.ROMFile.String() + "Mapper: 10 (MMC4)"
}

func (mmc4 *MMC4) Mappings(which rp2ago3.Mapping) (fetch, store []uint16) {
	fetch = []uint16{}
	store = []uint16{}

	switch which {
	case rp2ago3.PPU:
		if mmc4.CHRBanks > 0 {
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
		if mmc4.PRGBanks > 0 {
			// PRG RAM/NV RAM bank at $6000-$7FFF
			for i := uint32(0x6000); i <= 0x7fff; i++ {
				store = append(store, uint16(i))
				fetch = append(fetch, uint16(i))
			}

			// PRG bank 1
			for i := uint32(0x8000); i <= 0xbfff; i++ {
				store = append(store, uint16(i))
				fetch = append(fetch, uint16(i))
			}

			// PRG bank 2
			for i := uint32(0xc000); i <= 0xffff; i++ {
				store = append(store, uint16(i))
				fetch = append(fetch, uint16(i))
			}
		}
	}

	return
}

func (mmc4 *MMC4) Reset() {
	mmc4.Registers.Reset()
}

func (mmc4 *MMC4) Fetch(address uint16) (value uint8) {
	switch {
	// PPU only
	// CHR banks 1 & 2
	//case address >= 0x0000 && address <= 0x1fff:
	case address <= 0x1fff:
		index := address & 0x0fff
		lower, upper := mmc4.chrBanks()

		switch {
		// CHR bank 1
		//case address >= 0x0000 && address <= 0x0fff:
		case address <= 0x0fff:
			if mmc4.CHRBanks > 0 {
				value = mmc4.VROMBanks[lower][index]
			}
		// CHR bank 2
		case address >= 0x1000 && address <= 0x1fff:
			if mmc4.CHRBanks > 0 {
				value = mmc4.VROMBanks[upper][index]
			}
		}

		switch {
		case address >= 0x0fd8 && address <= 0x0fdf:
			mmc4.Registers.Latch0 = 0xfd
		case address >= 0x0fe8 && address <= 0x0fef:
			mmc4.Registers.Latch0 = 0xfe
		case address >= 0x1fd8 && address <= 0x1fdf:
			mmc4.Registers.Latch1 = 0xfd
		case address >= 0x1fe8 && address <= 0x1fef:
			mmc4.Registers.Latch1 = 0xfe
		}
	// CPU only
	// PRG RAM bank
	case address >= 0x6000 && address <= 0x7fff:
		value = mmc4.FetchPrgRam(address)
		//index := address & 0x1fff
		////value = mmc4.WRAMBanks[0][index]
		//if mmc4.Battery {
		//	value = mmc4.PRGNvRamData[index]
		//} else {
		//	value = mmc4.PRGRamData[index]
		//}
	// CPU only
	//case address >= 0x8000 && address <= 0xffff:
	case address >= 0x8000:
		// MMC4: 8KB switchable at $8000-$9FFF, then 24KB fixed at $A000-$FFFF
		switch {
		// Switchable 8KB bank at $8000-$9FFF
		case address >= 0x8000 && address <= 0x9fff:
			if mmc4.PRGBanks > 0 {
				bank := mmc4.prgBank()
				if bank >= uint8(mmc4.PRGBanks) {
					//bank = bank % uint8(mmc4.PRGBanks)
					return
				}
				index := address & 0x1fff
				value = mmc4.ROMBanks[bank][index]
			}
		// Fixed: last 24KB at $A000-$FFFF (last 3 8KB banks)
		case address >= 0xa000:
			if mmc4.PRGBanks > 0 {
				// Calculate offset within the 24KB region
				offset := address - 0xa000
				// Determine which of the last 3 8KB banks
				bankOffset := offset / 0x2000 // 0, 1, or 2
				bank := mmc4.PRGBanks - 3 + uint16(bankOffset)
				index := offset % 0x2000
				value = mmc4.ROMBanks[bank][index]
			}
		}
	}

	return
}

func (mmc4 *MMC4) Store(address uint16, value uint8) (oldValue uint8) {
	switch {
	// PPU only
	// CHR banks 1 & 2
	//case address >= 0x0000 && address <= 0x1fff:
	case address <= 0x1fff:
		index := address & 0x0fff
		lower, upper := mmc4.chrBanks()

		switch {
		// CHR bank 1
		//case address >= 0x0000 && address <= 0x0fff:
		case address <= 0x0fff:
			if mmc4.CHRBanks > 0 {
				mmc4.VROMBanks[lower][index] = value
			}
		// CHR bank 2
		case address >= 0x1000 && address <= 0x1fff:
			if mmc4.CHRBanks > 0 {
				mmc4.VROMBanks[upper][index] = value
			}
		}
	// CPU only
	// PRG RAM bank
	case address >= 0x6000 && address <= 0x7fff:
		mmc4.StorePrgRam(address, value)
		//index := address & 0x1fff
		////mmc4.WRAMBanks[0][index] = value
		//if mmc4.Battery && mmc4.PRGNvRamSize > 0 {
		//	mmc4.PRGNvRamData[index] = value
		//} else if !mmc4.Battery && mmc4.PRGRamSize > 0 {
		//	mmc4.PRGRamData[index] = value
		//}
	// PRG banks 1 & 2
	case address >= 0x8000 && address <= 0x9fff:
		bank := mmc4.prgBank()
		index := address & 0x1fff

		if mmc4.PRGBanks > 0 {
			mmc4.ROMBanks[bank][index] = value
		}
	// PRG bank select
	case address >= 0xa000 && address <= 0xafff:
		mmc4.Registers.PRGBank = value
	// CHR $fd/0000 bank select
	case address >= 0xb000 && address <= 0xbfff:
		mmc4.Registers.CHRBank0 = value
	// CHR $fe/0000 bank select
	case address >= 0xc000 && address <= 0xcfff:
		mmc4.Registers.CHRBank1 = value
	// CHR $fd/1000 bank select
	case address >= 0xd000 && address <= 0xdfff:
		mmc4.Registers.CHRBank2 = value
	// CHR $fe/1000 bank select
	case address >= 0xe000 && address <= 0xefff:
		mmc4.Registers.CHRBank3 = value
	// Mirroring
	//case address >= 0xf000 && address <= 0xffff:
	case address >= 0xf000:
		oldMirroring := mmc4.mirroring()
		mmc4.Registers.Mirroring = value

		if mmc4.mirroring() != oldMirroring {
			mmc4.setTables(mmc4.Tables())
		}
	}

	return
}

func (mmc4 *MMC4) prgBank() uint8 {
	bank := mmc4.Registers.PRGBank & 0x0f
	// Wrap to available banks
	if mmc4.PRGBanks > 0 {
		bank = bank % uint8(mmc4.PRGBanks)
	}
	return bank
}

func (mmc4 *MMC4) chrBank0() (bank uint8) {
	bank = mmc4.Registers.CHRBank0 & 0x1f
	if mmc4.CHRBanks > 0 {
		bank = bank % uint8(mmc4.CHRBanks)
	}
	return
}

func (mmc4 *MMC4) chrBank1() (bank uint8) {
	bank = mmc4.Registers.CHRBank1 & 0x1f
	if mmc4.CHRBanks > 0 {
		bank = bank % uint8(mmc4.CHRBanks)
	}
	return
}

func (mmc4 *MMC4) chrBank2() (bank uint8) {
	bank = mmc4.Registers.CHRBank2 & 0x1f
	if mmc4.CHRBanks > 0 {
		bank = bank % uint8(mmc4.CHRBanks)
	}
	return
}

func (mmc4 *MMC4) chrBank3() (bank uint8) {
	bank = mmc4.Registers.CHRBank3 & 0x1f
	if mmc4.CHRBanks > 0 {
		bank = bank % uint8(mmc4.CHRBanks)
	}
	return
}

func (mmc4 *MMC4) mirroring() rp2cgo2.Mirroring {
	switch mmc4.Registers.Mirroring & 0x01 {
	case 0:
		return rp2cgo2.Vertical
	default:
		return rp2cgo2.Horizontal
	}
}

func (mmc4 *MMC4) chrBanks() (lower, upper uint8) {
	switch mmc4.Registers.Latch0 {
	case 0xfd:
		lower = mmc4.chrBank0()
	case 0xfe:
		lower = mmc4.chrBank1()
	}

	switch mmc4.Registers.Latch1 {
	case 0xfd:
		upper = mmc4.chrBank2()
	case 0xfe:
		upper = mmc4.chrBank3()
	}

	return
}

func (mmc4 *MMC4) Tables() (t0, t1, t2, t3 int) {
	switch mmc4.mirroring() {
	case rp2cgo2.Vertical:
		t0, t1, t2, t3 = 0, 1, 0, 1
	case rp2cgo2.Horizontal:
		t0, t1, t2, t3 = 0, 0, 1, 1
	}

	return
}
