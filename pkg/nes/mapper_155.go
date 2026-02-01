// MMC1A Mapper (INES Mapper 155)
package nes

type MMC1A struct {
	*MMC1
}

func NewMMC1A(romf *ROMFile) *MMC1A {
	mmc1 := NewMMC1(romf)
	return &MMC1A{mmc1}
}

func (mmc1 *MMC1A) String() string {
	str := mmc1.ROMFile.String() + "Mapper: 155 (MMC1A)"
	str += "Board: " + mmc1.Board.String()
	return str
}

func (mmc1 *MMC1A) Fetch(address uint16) (value uint8) {
	return mmc1.processFetch(address)
}

func (mmc1 *MMC1A) processFetch(address uint16) (value uint8) {
	switch {
	// PPU only
	// CHR banks 1 & 2
	//case address >= 0x0000 && address <= 0x1fff:
	case address <= 0x1fff:
		index := address & 0x0fff
		lower, upper := mmc1.chrBanks()

		switch {
		// CHR bank 1
		//case address >= 0x0000 && address <= 0x0fff:
		case address <= 0x0fff:
			if mmc1.CHRBanks > 0 {
				// Wrap bank to available banks
				bank := uint16(lower) % mmc1.CHRBanks
				value = mmc1.VROMBanks[bank][index]
			}
			//if mmc1.CHRRomSize > 0 {
			//	value = mmc1.CHRROMData[index]
			//}
		// CHR bank 2
		case address >= 0x1000 && address <= 0x1fff:
			if mmc1.CHRBanks > 0 {
				// Wrap bank to available banks
				bank := uint16(upper) % mmc1.CHRBanks
				value = mmc1.VROMBanks[bank][index]
			}
			//if mmc1.CHRRomSize > 0 {
			//	value = mmc1.CHRROMData[index]
			//}
		}
	// CPU only
	// PRG RAM bank
	case address >= 0x6000 && address <= 0x7fff:
		if mmc1.isPRGRAMEnabled() {
			value = mmc1.FetchPrgRam(address)
		}
		// If no PRG RAM, return 0xFF (open bus)
		return 0 // Open bus when disabled
	case address >= 0x8000:
		lower, upper := mmc1.prgBanks()
		switch {
		// PRG bank 1
		case address >= 0x8000 && address <= 0xbfff:
			index := address & 0x3fff

			if mmc1.PRGBanks > 0 {
				// Wrap bank to available banks
				bank := lower % uint16(mmc1.PRGBanks)
				value = mmc1.ROMBanks[bank][index]
			}
		// PRG bank 2
		//case address >= 0xc000 && address <= 0xffff:
		case address >= 0xc000:
			index := address & 0x3fff

			if mmc1.PRGBanks > 0 {
				// Wrap bank to available banks
				bank := upper % uint16(mmc1.PRGBanks)
				value = mmc1.ROMBanks[bank][index]
			}
		}
	}

	return
}

func (mmc1 *MMC1A) Store(address uint16, value uint8) (oldValue uint8) {
	currentCycle := mmc1.nes.Tick
	if (value&0x80) != 0 || currentCycle-mmc1._lastWriteCycle >= 2 {
		oldValue = mmc1.processStore(address, value)
	}
	mmc1._lastWriteCycle = currentCycle
	return
}

func (mmc1 *MMC1A) processStore(address uint16, value uint8) (oldValue uint8) {
	switch {
	// PPU only
	// CHR banks 1 & 2
	//case address >= 0x0000 && address <= 0x1fff:
	case address <= 0x1fff:
		index := address & 0x0fff
		lower, upper := mmc1.chrBanks()

		switch {
		// CHR bank 1
		//case address >= 0x0000 && address <= 0x0fff:
		case address <= 0x0fff:
			if mmc1.CHRBanks > 0 {
				// Wrap bank to available banks
				bank := uint16(lower) % mmc1.CHRBanks
				mmc1.VROMBanks[bank][index] = value
			}
		// CHR bank 2
		case address >= 0x1000 && address <= 0x1fff:
			if mmc1.CHRBanks > 0 {
				// Wrap bank to available banks
				bank := uint16(upper) % mmc1.CHRBanks
				mmc1.VROMBanks[bank][index] = value
			}
		}
	// CPU only
	// PRG RAM bank
	case address >= 0x6000 && address <= 0x7fff:
		if mmc1.isPRGRAMEnabled() {
			mmc1.StorePrgRam(address, value)
		}
	// PRG banks 1 & 2
	//case address >= 0x8000 && address <= 0xffff:
	case address >= 0x8000:
		oldMirrors := mmc1.control(Mirroring)

		oldValue = mmc1.Registers.Load

		if (value & 0x80) != 0 {
			mmc1.Registers.Load = 0x00

			mmc1.Registers.Shift = 0x00
			mmc1.Registers.ShiftCount = 0x00

			mmc1.Registers.Control |= 0x0c
		} else {
			mmc1.Registers.Load = value

			mmc1.Registers.Shift |= (value & 0x01) << mmc1.Registers.ShiftCount
			mmc1.Registers.ShiftCount++

			if mmc1.Registers.ShiftCount == 0x05 {
				switch (address >> 13) & 0x0003 {
				case 0x0000:
					mmc1.Registers.Control = mmc1.Registers.Shift
				case 0x0001:
					mmc1.Registers.CHRBank0 = mmc1.Registers.Shift
				case 0x0002:
					mmc1.Registers.CHRBank1 = mmc1.Registers.Shift
				case 0x0003:
					mmc1.Registers.PRGBank = mmc1.Registers.Shift
				}

				mmc1.Registers.Shift = 0x00
				mmc1.Registers.ShiftCount = 0x00
			}
		}

		if mmc1.control(Mirroring) != oldMirrors {
			mmc1.setTables(mmc1.Tables())
		}
	}

	return
}

// isPRGRAMEnabled returns whether PRG-RAM is accessible
// MMC1A: Always enabled (ignores bit 4)
// MMC1B: Controlled by bit 4 (0 = enabled, 1 = disabled)
func (mmc1 *MMC1A) isPRGRAMEnabled() bool {
	return true // PRG-RAM always enabled
}

// |++--- PRG ROM bank mode (0, 1: switch 32 KB at $8000, ignoring low bit of bank number;
// |                         2: fix first bank at $8000 and switch 16 KB bank at $C000;
// |                         3: fix last bank at $C000 and switch 16 KB bank at $8000)
func (mmc1 *MMC1A) prgBanks() (lower, upper uint16) {
	prgBank := mmc1.prgBank(PRGBankSelect)

	// MMC1A: When bit 4 is set, bit 3 directly controls A17 across entire $8000-$FFFF
	if mmc1.Registers.PRGBank&0x10 != 0 {
		// Bit 3 controls A17, bits 2-0 select the 16KB bank within the 128KB half
		base := uint16((mmc1.Registers.PRGBank & 0x08) << 1) // bit 3 -> A17 (bit 4 of address)
		bank := uint16(prgBank & 0x07)

		switch mmc1.control(PRGRomBankMode) {
		case 0, 1: // 32 KB mode
			lower = base | (bank & 0x0e)
			upper = lower | 0x01
		case 2: // Fix first, switch upper
			lower = base
			upper = base | bank
		case 3: // Switch lower, fix last
			lower = base | bank
			upper = base | 0x07
		}
		return
	}

	// MMC1A with bit 4 clear, fallback to standard MMC1B behavior
	switch mmc1.control(PRGRomBankMode) {
	// 32 KB
	case 0, 1:
		lower = uint16(prgBank)
		upper = lower | 0x01
	// 16 KB
	case 2:
		lower = 0
		upper = uint16(prgBank)
	case 3:
		lower = uint16(prgBank)
		upper = mmc1.PRGBanks - 1
	}

	return
}
