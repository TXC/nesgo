// MMC1 Mapper (INES Mapper 1)
package nes

import (
	"github.com/txc/nesgo/pkg/rp2ago3"
)

type MMC1Mirroring uint8

const (
	OneScreenLowerBank MMC1Mirroring = iota
	OneScreenUpperBank
	Vertical
	Horizontal
)

type ControlFlag uint8

const (
	Mirroring ControlFlag = 1 << iota
	_
	PRGRomBankMode
	_
	CHRRomBankMode
)

type PRGBankFlag uint8

const (
	PRGBankSelect PRGBankFlag = 1 << iota
	_
	_
	_
	PRGRAMChipEnable
)

type MMC1BoardType uint8

const (
	BoardSAROM MMC1BoardType = iota
	BoardSBROM
	BoardSCROM
	BoardSC1ROM
	BoardSEROM
	BoardSFROM
	BoardSF1ROM
	BoardSFEXPROM
	BoardSGROM
	BoardSHROM
	BoardSH1ROM
	BoardSIROM
	BoardSIEPROM
	BoardSJROM
	BoardSKROM
	BoardSKEPROM
	BoardSLROM
	BoardSL1ROM
	BoardSL2ROM
	BoardSL3ROM
	BoardSLRROM
	BoardSMROM
	BoardSNROM
	BoardSNWEPROM
	BoardSOROM
	BoardSTROM
	BoardSUROM
	BoardSXROM
	BoardSZROM
)

func (bt MMC1BoardType) String() string {
	switch bt {
	case BoardSAROM:
		return "SAROM"
	case BoardSBROM:
		return "SBROM"
	case BoardSCROM:
		return "SCROM"
	case BoardSC1ROM:
		return "SC1ROM"
	case BoardSEROM:
		return "SEROM"
	case BoardSFROM:
		return "SFROM"
	case BoardSF1ROM:
		return "SF1ROM"
	case BoardSFEXPROM:
		return "SFEXPROM"
	case BoardSGROM:
		return "SGROM"
	case BoardSHROM:
		return "SHROM"
	case BoardSH1ROM:
		return "SH1ROM"
	case BoardSIROM:
		return "SIROM"
	case BoardSIEPROM:
		return "SIEPROM"
	case BoardSJROM:
		return "SJROM"
	case BoardSKROM:
		return "SKROM"
	case BoardSKEPROM:
		return "SKEPROM"
	case BoardSLROM:
		return "SLROM"
	case BoardSL1ROM:
		return "SL1ROM"
	case BoardSL2ROM:
		return "SL2ROM"
	case BoardSL3ROM:
		return "SL3ROM"
	case BoardSLRROM:
		return "SLRROM"
	case BoardSMROM:
		return "SMROM"
	case BoardSNROM:
		return "SNROM"
	case BoardSNWEPROM:
		return "SNWEPROM"
	case BoardSOROM:
		return "SOROM"
	case BoardSTROM:
		return "STROM"
	case BoardSUROM:
		return "SUROM"
	case BoardSXROM:
		return "SXROM"
	case BoardSZROM:
		return "SZROM"
	}

	return "Unknown"
}

type MMC1 struct {
	*ROMFile
	Registers       MMC1Registers
	Board           MMC1BoardType
	_lastWriteCycle uint64
	_lastChrReg     uint16
}

type MMC1Registers struct {
	Load       uint8
	Control    uint8
	CHRBank0   uint8
	CHRBank1   uint8
	PRGBank    uint8
	Shift      uint8
	ShiftCount uint8
}

func (reg *MMC1Registers) Reset() {
	reg.Load = 0x00
	reg.Control = 0x0C // PRG mode 3: switch lower bank at $8000, fix last bank at $C000
	reg.CHRBank0 = 0x00
	reg.CHRBank1 = 0x00
	reg.PRGBank = 0x00
	reg.Shift = 0x00
	reg.ShiftCount = 0x00
}

func NewMMC1(romf *ROMFile) *MMC1 {
	// Detect board type first
	board := getMMC1BoardType(romf)

	// Override WRAM for boards that shouldn't have it
	// SGROM/SLROM boards with large CHR-ROM (>32KB) don't have WRAM
	// This overrides both default and explicitly specified RAM
	if romf.CHRRomSize > 0 && romf.CHRRomSize > 0x8000 {
		// Remove WRAM for large CHR-ROM boards
		romf.PRGRamSize = 0
		romf.PRGRamData = nil
		romf.PRGNvRamSize = 0
		romf.PRGNvRamData = nil
	}

	mmc1 := &MMC1{
		ROMFile: romf,
		Board:   board,
	}

	// divide 8KB CHR banks into 4KB banks since we may be
	// swapping 4KB banks
	if romf.CHRRomSize > 0 {
		// Number of 4KB banks = (CHR ROM size / 4KB)
		banks := romf.CHRRomSize / 0x1000
		vromBanks := make([][]uint8, banks)

		for n := range banks {
			vromBanks[n] = romf.CHRRomData[n*0x1000 : (n+1)*0x1000]
		}

		romf.VROMBanks = vromBanks
		romf.CHRBanks = uint16(banks)

		if mmc1._debugLog {
			mmc1.log.Logf("[MMC1] CHR ROM: %d banks of 4KB (total $%X bytes)\n", banks, romf.CHRRomSize)
		}
	} else if romf.CHRRamSize > 0 {
		// Initialize CHR-RAM banks (4KB each for MMC1's 4KB banking)
		// MMC1 uses 4KB banking, so we need CHR RAM size / 4KB banks
		// For proper MMC1 operation, we need at least 2 banks (8KB) even if header says less
		numBanks := romf.CHRRamSize / 0x1000
		if numBanks < 2 {
			numBanks = 2             // Minimum 2 banks (8KB total) for 4K CHR mode
			romf.CHRRamSize = 0x2000 // Adjust size to match
		}

		if mmc1._debugLog {
			mmc1.log.Logf("[MMC1] CHR RAM: %d banks of 4KB (total $%X bytes)\n", numBanks, romf.CHRRamSize)
		}

		vromBanks := make([][]uint8, numBanks)
		for i := 0; i < numBanks; i++ {
			vromBanks[i] = make([]uint8, 0x1000) // 4KB per bank
		}
		romf.VROMBanks = vromBanks
		romf.CHRBanks = uint16(numBanks)
	}

	mmc1.Registers.Reset()
	mmc1.setTables(mmc1.Tables())

	return mmc1
}

func (mmc1 *MMC1) String() string {
	str := mmc1.ROMFile.String() + "Mapper: 1 (MMC1B) "
	str += "Board: " + mmc1.Board.String()
	return str
}

func (mmc1 *MMC1) Mappings(which rp2ago3.Mapping) (fetch, store []uint16) {
	fetch = []uint16{}
	store = []uint16{}

	switch which {
	case rp2ago3.PPU:
		// Map CHR for both ROM and RAM
		if mmc1.CHRRomSize > 0 || mmc1.CHRRamSize > 0 {
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
		// Map WRAM
		for i := uint32(0x6000); i <= 0x7fff; i++ {
			store = append(store, uint16(i))
			fetch = append(fetch, uint16(i))
		}

		if mmc1.PRGRomSize > 0 {
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

func (mmc1 *MMC1) Reset() {
	mmc1.Registers.Reset()
}

func (mmc1 *MMC1) SetDebug(enable bool) {
	mmc1._debugLog = enable
	if enable {
		mmc1.log.Logf("[MMC1] Debug logging enabled. Board: %s, PRG ROM: %d KB, CHR ROM: %d KB, CHR RAM: %d KB, PRG RAM: %d KB\n",
			mmc1.Board.String(), mmc1.PRGRomSize/1024, mmc1.CHRRomSize/1024, mmc1.CHRRamSize/1024, mmc1.PRGRamSize/1024)
		mmc1.log.Logf("[MMC1] Initial state: Control=$%02X, CHRBank0=$%02X, CHRBank1=$%02X, PRGBank=$%02X\n",
			mmc1.Registers.Control, mmc1.Registers.CHRBank0, mmc1.Registers.CHRBank1, mmc1.Registers.PRGBank)
		lower, upper := mmc1.prgBanks()
		mmc1.log.Logf("[MMC1] Initial PRG banks: lower=%d, upper=%d\n", lower, upper)
	}
}

func (mmc1 *MMC1) Fetch(address uint16) (value uint8) {
	return mmc1.processFetch(address)
}

func (mmc1 *MMC1) processFetch(address uint16) (value uint8) {
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
		} else {
			// Open bus when PRG RAM is disabled
			value = 0
		}
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

func (mmc1 *MMC1) Store(address uint16, value uint8) (oldValue uint8) {
	currentCycle := mmc1.nes.Tick
	if (value&0x80) != 0 || currentCycle-mmc1._lastWriteCycle >= 2 {
		oldValue = mmc1.processStore(address, value)
	}
	mmc1._lastWriteCycle = currentCycle
	return
}

func (mmc1 *MMC1) processStore(address uint16, value uint8) (oldValue uint8) {
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
				regSelect := (address >> 13) & 0x0003
				switch regSelect {
				case 0x0000:
					mmc1.Registers.Control = mmc1.Registers.Shift
				case 0x0001:
					mmc1.Registers.CHRBank0 = mmc1.Registers.Shift
					mmc1._lastChrReg = address
				case 0x0002:
					mmc1.Registers.CHRBank1 = mmc1.Registers.Shift
					mmc1._lastChrReg = address
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

// 4bit0
// -----
// CPPMM
// |||||
// |||++- Mirroring (0: one-screen, lower bank; 1: one-screen, upper bank;
// |||               2: vertical; 3: horizontal)
// |++--- PRG ROM bank mode (0, 1: switch 32 KB at $8000, ignoring low bit of bank number;
// |                         2: fix first bank at $8000 and switch 16 KB bank at $C000;
// |                         3: fix last bank at $C000 and switch 16 KB bank at $8000)
// +----- CHR ROM bank mode (0: switch 8 KB at a time; 1: switch two separate 4 KB banks)
func (mmc1 *MMC1) control(flag ControlFlag) (value uint8) {
	reg := mmc1.Registers.Control

	switch flag {
	case Mirroring:
		value = reg & 0x03
	case PRGRomBankMode:
		value = (reg >> 2) & 0x03
	case CHRRomBankMode:
		value = (reg >> 4) & 0x01
	}

	return
}

// 4bit0
// -----
// CCCCC
// |||||
// +++++- Select 4 KB or 8 KB CHR bank at PPU $0000 (low bit ignored in 8 KB mode)
func (mmc1 *MMC1) chrBank0() (bank uint8) {
	reg := mmc1.Registers.CHRBank0

	switch mmc1.control(CHRRomBankMode) {
	// 8 KB at a time
	case 0:
		bank = reg & 0x1e
	// switch two separate 4 KB banks
	case 1:
		bank = reg & 0x1f
	}

	return
}

// 4bit0
// -----
// CCCCC
// |||||
// +++++- Select 4 KB CHR bank at PPU $1000 (ignored in 8 KB mode)
func (mmc1 *MMC1) chrBank1() (bank uint8) {
	bank = mmc1.Registers.CHRBank1 & 0x1f
	return
}

// 4bit0
// -----
// RPPPP
// |||||
// |++++- Select 16 KB PRG ROM bank (low bit ignored in 32 KB mode)
// +----- PRG RAM chip enable (0: enabled; 1: disabled; ignored on MMC1A)
// +----- MMC1A:
//
//	0: fixed bank affects A17..A14
//	1: fixed bank only affects A16..A14,
//	   bit 3 directly controls A17 across the entire $8000-$FFFF address range
//	MMC1B:
//	0: PRG-RAM enabled
//	1: PRG-RAM disabled
func (mmc1 *MMC1) prgBank(flag PRGBankFlag) (value uint8) {
	reg := mmc1.Registers.PRGBank

	switch flag {
	case PRGBankSelect:
		// (low bit ignored in 32 KB mode)
		if mmc1.control(PRGRomBankMode) < 2 {
			value = reg & 0x0e
		} else {
			value = reg & 0x0f
		}
	case PRGRAMChipEnable:
		value = (reg >> 4) & 0x01
	}

	return
}

// isPRGRAMEnabled returns whether PRG-RAM is accessible
// MMC1A: Always enabled (ignores bit 4)
// MMC1B: Controlled by bit 4 of PRG bank register (0 = enabled, 1 = disabled)
// SNROM: Additionally uses CHR bank 0 register bit 4 to disable WRAM via CHR A12
func (mmc1 *MMC1) isPRGRAMEnabled() bool {
	// First check PRG bank register bit 4 (MMC1B/C protection)
	if mmc1.prgBank(PRGRAMChipEnable) != 0 {
		// bit 4 = 1 means RAM is disabled
		return false
	}

	// For SNROM specifically, CHR bank 0 register bit 4 also controls RAM access
	// When CHR A12 (bit 4 of CHR bank 0) is high, it connects to PRG RAM /CE and disables WRAM
	if mmc1.Board == BoardSNROM {
		// On SNROM: CHR bank 0 bit 4 (CHR A12) disables WRAM when set
		if (mmc1.Registers.CHRBank0 & 0x10) != 0 {
			return false
		}
	}

	// Otherwise PRG RAM is enabled
	return true
}

func (mmc1 *MMC1) chrBanks() (lower, upper uint8) {
	// Calculate mask based on actual CHR bank count
	// For 2 banks (8KB), mask is 0x01 (only bit 0 matters)
	// For 4 banks (16KB), mask is 0x03 (bits 0-1 matter)
	// etc.
	mask := uint8(mmc1.CHRBanks - 1)

	switch mmc1.control(CHRRomBankMode) {
	// 8 KB
	case 0:
		lower = mmc1.chrBank0() & mask
		upper = (lower | 0x01) & mask
	// 4 KB
	case 1:
		lower = mmc1.chrBank0() & mask
		upper = mmc1.chrBank1() & mask
	}

	return
}

// |++--- PRG ROM bank mode (0, 1: switch 32 KB at $8000, ignoring low bit of bank number;
// |                         2: fix first bank at $8000 and switch 16 KB bank at $C000;
// |                         3: fix last bank at $C000 and switch 16 KB bank at $8000)
func (mmc1 *MMC1) prgBanks() (lower, upper uint16) {
	prgBank := mmc1.prgBank(PRGBankSelect)

	if mmc1.SubMapper == 5 {
		lower = 0
		upper = 1
		mmc1.log.Logf("[MMC1] prgBanks SubMapper 5: lower=%d, upper=%d\n", lower, upper)
		return
	}

	var extraBit uint16 = 0
	// For 512KB PRG ROM (SUROM), use CHR bank bit 4 for PRG A18
	if mmc1.PRGRomSize == 0x80000 {
		_chrMode := mmc1.control(CHRRomBankMode) == 1
		var extraReg uint8
		// Determine which CHR bank register to use for extra bit
		if mmc1._lastChrReg == 0xC000 && _chrMode {
			extraReg = mmc1.Registers.CHRBank1
		} else {
			extraReg = mmc1.Registers.CHRBank0
		}
		// Use bit 4 of CHR register as PRG A18 (the 5th PRG address bit)
		// This gives us 5 bits total: 4 from PRG register + 1 from CHR register
		extraBit = uint16((extraReg>>4)&0x01) << 4
	}

	// Standard MMC1B behavior
	switch mmc1.control(PRGRomBankMode) {
	// 32 KB
	case 0, 1:
		lower = uint16(prgBank) | extraBit
		upper = (uint16(prgBank) | 0x01) | extraBit
	// 16 KB
	case 2:
		lower = 0 | extraBit
		upper = uint16(prgBank) | extraBit
	case 3:
		lower = uint16(prgBank) | extraBit
		// Last bank - for 512KB, this needs the extra bit too
		upper = (mmc1.PRGBanks - 1) | extraBit
	}

	return
}

func (mmc1 *MMC1) Tables() (t0, t1, t2, t3 int) {
	switch MMC1Mirroring(mmc1.control(Mirroring)) {
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

func getMMC1BoardType(romf *ROMFile) (board MMC1BoardType) {
	// Calculate total PRG RAM (regular + battery-backed)
	totalPrgRam := romf.PRGRamSize + romf.PRGNvRamSize

	switch romf.PRGRomSize {
	case 0x8000: // 32 KB
		if totalPrgRam > 0 {
			board = BoardSIROM
			return
		} else if romf.CHRRomSize > 0x10000 {
			board = BoardSHROM
			return
		}
		return BoardSEROM
	case 0x10000: // 64 KB
		if totalPrgRam > 0 {
			board = BoardSAROM
			return
		} else if romf.CHRRomSize > 0x10000 {
			board = BoardSCROM
			return
		}
		return BoardSBROM
	case 0x20000: // 128 KB
		// Total PRG RAM (regular + battery-backed)
		totalPrgRam := romf.PRGRamSize + romf.PRGNvRamSize

		// Check for SNROM: 128KB PRG, 8KB PRG RAM, 8KB CHR (ROM or RAM)
		if totalPrgRam == 0x2000 && (romf.CHRRomSize == 0x2000 || romf.CHRRamSize == 0x2000) {
			board = BoardSNROM
			return
		}
		// Also check with smaller CHR RAM (4KB) for test ROMs
		if totalPrgRam == 0x2000 && (romf.CHRRomSize == 0x1000 || romf.CHRRamSize == 0x1000) {
			board = BoardSNROM
			return
		}
		// Check for SNROM with 32KB CHR-ROM: 128KB PRG, 8KB PRG RAM, 32KB CHR-ROM
		if totalPrgRam == 0x2000 && romf.CHRRomSize == 0x8000 {
			board = BoardSNROM
			return
		}
		// Check for SGROM: 128KB PRG, no PRG RAM, 8KB CHR RAM
		if totalPrgRam == 0 && (romf.CHRRomSize == 0x2000 || romf.CHRRamSize == 0x2000) {
			board = BoardSGROM
			return
		}
		// Also check with smaller CHR RAM for test ROMs (no PRG RAM specified means 8KB will be added)
		if totalPrgRam == 0 && (romf.CHRRomSize == 0x1000 || romf.CHRRamSize == 0x1000) {
			board = BoardSGROM // Will get 8KB PRG RAM added later
			return
		}
		// Default 128KB configuration with RAM (SNROM-like)
		if totalPrgRam > 0 {
			board = BoardSNROM
			return
		}
		// Default 128KB configuration without RAM (SGROM-like)
		board = BoardSGROM
		return
	case 0x40000: // 256 KB
		if romf.CHRRomSize == 0x2000 && totalPrgRam == 0 {
			board = BoardSMROM
			return
		} else if romf.CHRRomSize == 0x10000 && totalPrgRam == 0 {
			board = BoardSF1ROM
			return
		} else if romf.CHRRomSize == 0x20000 && totalPrgRam == 0x2000 {
			board = BoardSKEPROM
			return
		} else if romf.CHRRomSize == 0x20000 {
			board = BoardSL3ROM
			return
		}
	case 0x80000: // 512 KB
		// SUROM: 512KB PRG with extended banking via CHR register
		// Can have various RAM configurations (8KB, 16KB, 32KB, etc.)
		board = BoardSUROM
		return
	}

	switch true {
	// Generic detection for 128 / 256 / 512 KB
	// SGROM: No PRG RAM, 8KB CHR (ROM or RAM)
	case totalPrgRam == 0 && (romf.CHRRomSize == 0x2000 || romf.CHRRamSize == 0x2000):
		board = BoardSGROM
		return
	// SFROM: No PRG RAM, CHR ROM between 8KB and 128KB
	case totalPrgRam == 0 && romf.CHRRomSize > 0 && romf.CHRRomSize < 0x20000:
		board = BoardSFROM
		return
	// SLROM: No PRG RAM, 128KB CHR ROM
	case totalPrgRam == 0 && romf.CHRRomSize == 0x20000:
		board = BoardSLROM
		return
	// SNROM: 8KB PRG RAM, 8KB CHR (ROM or RAM)
	case totalPrgRam == 0x2000 && (romf.CHRRomSize == 0x2000 || romf.CHRRamSize == 0x2000):
		board = BoardSNROM
		return
	// SJROM: 8KB PRG RAM, CHR ROM between 16KB and 128KB
	case totalPrgRam == 0x2000 && romf.CHRRomSize > 0x4000 && romf.CHRRomSize < 0x20000:
		board = BoardSJROM
		return

	// SKROM: 8KB PRG RAM, 128KB CHR ROM
	case totalPrgRam == 0x2000 && romf.CHRRomSize == 0x20000:
		board = BoardSKROM
		return

	// SOROM: 16KB PRG RAM, 8KB CHR (ROM or RAM)
	case totalPrgRam == 0x4000 && (romf.CHRRomSize == 0x2000 || romf.CHRRamSize == 0x2000):
		board = BoardSOROM
		return

	// SZROM: >16KB PRG RAM, CHR ROM between 16KB and 64KB
	case totalPrgRam > 0x4000 && romf.CHRRomSize > 0x4000 && romf.CHRRomSize < 0x10000:
		board = BoardSZROM
		return

	// SXROM: 32KB PRG RAM, 8KB CHR (ROM or RAM)
	case totalPrgRam == 0x8000 && (romf.CHRRomSize == 0x2000 || romf.CHRRamSize == 0x2000):
		board = BoardSXROM
		return
	default:
		return BoardSAROM
	}
}
