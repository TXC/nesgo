// MMC3 Mapper (INES Mapper 4)
package nes

import (
	"fmt"

	"github.com/txc/nesgo/pkg/rp2ago3"
	"github.com/txc/nesgo/pkg/rp2cgo2"
)

type MMC3BankSelectFlag uint8

const (
	BankRegister MMC3BankSelectFlag = 1 << iota
	PRGROMBankMode
	CHRA12Inversion
)

type MMC3Registers struct {
	BankSelect    uint8
	BankData      uint8
	Mirroring     uint8
	PRGRAMProtect uint8

	IRQLatch   uint8
	IRQReload  bool
	IRQEnable  bool
	IRQCounter uint8

	CHRBank1 uint8
	CHRBank2 uint8
	CHRBank3 uint8
	CHRBank4 uint8
	CHRBank5 uint8
	CHRBank6 uint8

	PRGBankLow  uint8
	PRGBankHigh uint8
}

func (mmc3 *MMC3Registers) String() string {
	return fmt.Sprintf("Registers: "+
		"BankSelect: 0x%02X\n"+
		"BankData: 0x%02X\n"+
		"Mirroring: 0x%02X\n"+
		"PRGRAMProtect: 0x%02X\n"+

		"IRQLatch   0x%02X\n"+ //uint8
		"IRQReload  %t\n"+ //bool
		"IRQEnable  %t\n"+ //bool
		"IRQCounter 0x%02X\n"+ //uint8

		"CHRBank1: 0x%02X\n"+
		"CHRBank2: 0x%02X\n"+
		"CHRBank3: 0x%02X\n"+
		"CHRBank4: 0x%02X\n"+
		"CHRBank5: 0x%02X\n"+
		"CHRBank6: 0x%02X\n"+

		"PRGBankLow: 0x%02X\n"+
		"PRGBankHigh: 0x%02X\n",
		mmc3.BankSelect,
		mmc3.BankData,
		mmc3.Mirroring,
		mmc3.PRGRAMProtect,
		mmc3.IRQLatch,
		mmc3.IRQReload,
		mmc3.IRQEnable,
		mmc3.IRQCounter,
		mmc3.CHRBank1,
		mmc3.CHRBank2,
		mmc3.CHRBank3,
		mmc3.CHRBank4,
		mmc3.CHRBank5,
		mmc3.CHRBank6,
		mmc3.PRGBankLow,
		mmc3.PRGBankHigh,
	)
}

type MMC3 struct {
	MapperA12RisingEdge
	*ROMFile
	Registers   MMC3Registers
	A12LowClock uint64
}

func (reg *MMC3Registers) Reset() {
	reg.BankSelect = 0x00
	reg.BankData = 0x00
	reg.Mirroring = 0x00
	reg.PRGRAMProtect = 0x80 // Enable WRAM chip by default (bit 7=1), writes enabled (bit 6=0)

	reg.IRQLatch = 0x00
	reg.IRQReload = false
	reg.IRQEnable = true
	reg.IRQCounter = 0x00

	reg.CHRBank1 = 0x00
	reg.CHRBank2 = 0x02
	reg.CHRBank3 = 0x04
	reg.CHRBank4 = 0x05
	reg.CHRBank5 = 0x06
	reg.CHRBank6 = 0x07
	reg.PRGBankLow = 0x00
	reg.PRGBankHigh = 0x01
}

func NewMMC3(romf *ROMFile) *MMC3 {
	mmc3 := &MMC3{
		ROMFile: romf,
	}

	// MMC3 typically has 8KB of WRAM even if not specified in ROM header
	// This is needed for compatibility with older iNES ROMs
	if mmc3.PRGRamSize == 0 && mmc3.PRGNvRamSize == 0 {
		mmc3.PRGRamSize = 8192 // 8KB default WRAM
		mmc3.PRGRamData = make([]uint8, 8192)
	}

	// divide 8KB CHR banks into 1KB banks
	if mmc3.CHRRomSize > 0 {
		offset := 0x0400
		vromBanks := make([][]uint8, uint16(romf.CHRBanks)*8)

		for n := 0; n < int(romf.CHRBanks); n++ {
			for i := range 8 {
				vromBanks[(8*n)+i] = romf.VROMBanks[n][(offset * i):((offset * i) + offset)]
			}
		}

		romf.VROMBanks = vromBanks
		romf.CHRBanks *= 8
	} else if romf.CHRRamSize > 0 {
		// Initialize CHR-RAM banks (1KB each)
		// Use actual CHR-RAM size from header, but ensure minimum of 8 banks (8KB)
		numBanks := max(romf.CHRRamSize/0x0400, 8)
		vromBanks := make([][]uint8, numBanks)
		for i := range numBanks {
			vromBanks[i] = make([]uint8, 0x0400) // 1KB per bank
		}
		romf.VROMBanks = vromBanks
		romf.CHRBanks = uint16(numBanks)
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

	mmc3.Registers.Reset()
	mmc3.setTables(mmc3.Tables())

	return mmc3
}

func (mmc3 *MMC3) String() string {
	return mmc3.ROMFile.String() + "Mapper: 4 (MMC3)"
}

func (mmc3 *MMC3) Mappings(which rp2ago3.Mapping) (fetch, store []uint16) {
	fetch = []uint16{}
	store = []uint16{}

	switch which {
	case rp2ago3.PPU:
		if mmc3.CHRRomSize > 0 || mmc3.CHRRamSize > 0 {
			// CHR bank 1
			for i := uint32(0x0000); i <= 0x03ff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}

			// CHR bank 2
			for i := uint32(0x0400); i <= 0x07ff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}

			// CHR bank 3
			for i := uint32(0x0800); i <= 0x0bff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}

			// CHR bank 4
			for i := uint32(0x0c00); i <= 0x0fff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}

			// CHR bank 5
			for i := uint32(0x1000); i <= 0x13ff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}

			// CHR bank 6
			for i := uint32(0x1400); i <= 0x17ff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}

			// CHR bank 7
			for i := uint32(0x1800); i <= 0x1bff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}

			// CHR bank 8
			for i := uint32(0x1c00); i <= 0x1fff; i++ {
				fetch = append(fetch, uint16(i))
				store = append(store, uint16(i))
			}
		}
	case rp2ago3.CPU:
		if mmc3.PRGRamSize > 0 || mmc3.PRGNvRamSize > 0 {
			// PRG RAM bank
			for i := uint32(0x6000); i <= 0x7fff; i++ {
				store = append(store, uint16(i))
				fetch = append(fetch, uint16(i))
			}
		}

		//if mmc3.PRGBanks > 0 {
		if mmc3.PRGRomSize > 0 {
			// PRG bank 1
			for i := uint32(0x8000); i <= 0x9fff; i++ {
				store = append(store, uint16(i))
				fetch = append(fetch, uint16(i))
			}

			// PRG bank 2
			for i := uint32(0xa000); i <= 0xbfff; i++ {
				store = append(store, uint16(i))
				fetch = append(fetch, uint16(i))
			}

			// PRG bank 3
			for i := uint32(0xc000); i <= 0xdfff; i++ {
				store = append(store, uint16(i))
				fetch = append(fetch, uint16(i))
			}

			// PRG bank 4
			for i := uint32(0xe000); i <= 0xffff; i++ {
				store = append(store, uint16(i))
				fetch = append(fetch, uint16(i))
			}
		}
	}

	return
}

func (mmc3 *MMC3) Reset() {
	mmc3.Registers.Reset()
}

func (mmc3 *MMC3) Fetch(address uint16) (value uint8) {
	switch {
	// PPU only
	// CHR banks 1-8
	//case address >= 0x0000 && address <= 0x1fff:
	case address <= 0x1fff:
		if mmc3.CHRRomSize > 0 || mmc3.CHRRamSize > 0 {
			index := address & 0x03ff
			bank1, bank2, bank3, bank4, bank5, bank6, bank7, bank8 := mmc3.chrBanks()

			// For CHR-RAM, wrap banks to available banks for actual memory access
			// For CHR-ROM, banks are already masked in chrBanks()
			if mmc3.CHRRamSize > 0 {
				bankCount := uint8(len(mmc3.VROMBanks))
				bank1 %= bankCount
				bank2 %= bankCount
				bank3 %= bankCount
				bank4 %= bankCount
				bank5 %= bankCount
				bank6 %= bankCount
				bank7 %= bankCount
				bank8 %= bankCount
			}

			switch {
			// CHR bank 1
			//case address >= 0x0000 && address <= 0x03ff:
			case address <= 0x03ff:
				value = mmc3.VROMBanks[bank1][index]
			// CHR bank 2
			case address >= 0x0400 && address <= 0x07ff:
				value = mmc3.VROMBanks[bank2][index]
			// CHR bank 3
			case address >= 0x0800 && address <= 0x0bff:
				value = mmc3.VROMBanks[bank3][index]
			// CHR bank 4
			case address >= 0x0c00 && address <= 0x0fff:
				value = mmc3.VROMBanks[bank4][index]
			// CHR bank 5
			case address >= 0x1000 && address <= 0x13ff:
				value = mmc3.VROMBanks[bank5][index]
			// CHR bank 6
			case address >= 0x1400 && address <= 0x17ff:
				value = mmc3.VROMBanks[bank6][index]
			// CHR bank 7
			case address >= 0x1800 && address <= 0x1bff:
				value = mmc3.VROMBanks[bank7][index]
			// CHR bank 8
			case address >= 0x1c00 && address <= 0x1fff:
				value = mmc3.VROMBanks[bank8][index]
			}
		}

	// CPU only
	case address >= 0x6000 && address <= 0x7fff:
		if chipEnable, _ := mmc3.prgRAMProtect(); chipEnable {
			value = mmc3.FetchPrgRam(address)
			//index := address & 0x1fff
			//if mmc3.PRGRamSize > 0 {
			//	value = mmc3.PRGRamData[index]
			//} else if mmc3.PRGRomSize > 0 {
			//	value = mmc3.PRGRomData[index]
			//}

			////value = mmc3.WRAMBanks[0][index]
			//value = mmc3.PRGRamData[index]
		}
	//case address >= 0x8000 && address <= 0xffff:
	case address >= 0x8000:
		index := address & 0x1fff
		bank1, bank2, bank3, bank4 := mmc3.prgBanks()

		switch {
		// PRG bank 1
		case address >= 0x8000 && address <= 0x9fff:
			value = mmc3.ROMBanks[bank1][index]
		// PRG bank 2
		case address >= 0xa000 && address <= 0xbfff:
			value = mmc3.ROMBanks[bank2][index]
		// PRG bank 3
		case address >= 0xc000 && address <= 0xdfff:
			value = mmc3.ROMBanks[bank3][index]
		// PRG bank 4
		//case address >= 0xe000 && address <= 0xffff:
		case address >= 0xe000:
			value = mmc3.ROMBanks[bank4][index]
		}
	}

	return
}

func (mmc3 *MMC3) Store(address uint16, value uint8) (oldValue uint8) {
	switch {
	// PPU only
	// CHR banks 1-8
	//case address >= 0x0000 && address <= 0x1fff:
	case address <= 0x1fff:
		if mmc3.CHRRamSize > 0 {
			index := address & 0x03ff
			bank1, bank2, bank3, bank4, bank5, bank6, bank7, bank8 := mmc3.chrBanks()

			// Wrap banks to available banks for actual memory access
			bankCount := uint8(len(mmc3.VROMBanks))
			bank1 %= bankCount
			bank2 %= bankCount
			bank3 %= bankCount
			bank4 %= bankCount
			bank5 %= bankCount
			bank6 %= bankCount
			bank7 %= bankCount
			bank8 %= bankCount

			switch {
			// CHR bank 1
			//case address >= 0x0000 && address <= 0x03ff:
			case address <= 0x03ff:
				mmc3.VROMBanks[bank1][index] = value
			// CHR bank 2
			case address >= 0x0400 && address <= 0x07ff:
				mmc3.VROMBanks[bank2][index] = value
			// CHR bank 3
			case address >= 0x0800 && address <= 0x0bff:
				mmc3.VROMBanks[bank3][index] = value
			// CHR bank 4
			case address >= 0x0c00 && address <= 0x0fff:
				mmc3.VROMBanks[bank4][index] = value
			// CHR bank 5
			case address >= 0x1000 && address <= 0x13ff:
				mmc3.VROMBanks[bank5][index] = value
			// CHR bank 6
			case address >= 0x1400 && address <= 0x17ff:
				mmc3.VROMBanks[bank6][index] = value
			// CHR bank 7
			case address >= 0x1800 && address <= 0x1bff:
				mmc3.VROMBanks[bank7][index] = value
			// CHR bank 8
			case address >= 0x1c00 && address <= 0x1fff:
				mmc3.VROMBanks[bank8][index] = value
			}
		}
	// CPU only
	// PRG RAM bank
	case address >= 0x6000 && address <= 0x7fff:
		if chipEnable, allowWrites := mmc3.prgRAMProtect(); chipEnable && allowWrites {
			mmc3.StorePrgRam(address, value)
			//index := address & 0x1fff
			//if mmc3.PRGRamSize > 0 {
			//	mmc3.PRGRamData[index] = value
			//} else if mmc3.PRGRomSize > 0 {
			//	mmc3.PRGRomData[index] = value
			//}

			////mmc3.WRAMBanks[0][index] = value
			//mmc3.PRGRamData[index] = value
		}
	// Bank select (even) / Bank data (odd)
	case address >= 0x8000 && address <= 0x9fff:
		if (address & 0x0001) == 0x0000 { // even
			mmc3.Registers.BankSelect = value
		} else { // odd
			switch mmc3.bankSelect(BankRegister) {
			case 0:
				// Registers 0 and 1 ignore bit 0 (2KB banks)
				mmc3.Registers.CHRBank1 = value & 0xFE
			case 1:
				mmc3.Registers.CHRBank2 = value & 0xFE
			case 2:
				mmc3.Registers.CHRBank3 = value
			case 3:
				mmc3.Registers.CHRBank4 = value
			case 4:
				mmc3.Registers.CHRBank5 = value
			case 5:
				mmc3.Registers.CHRBank6 = value
			case 6:
				mmc3.Registers.PRGBankLow = value & 0x3F
			case 7:
				mmc3.Registers.PRGBankHigh = value & 0x3F
			}

			mmc3.Registers.BankData = value
		}
	// PRG RAM protect (odd) / Mirroring (even)
	case address >= 0xa000 && address <= 0xbfff:
		if (address & 0x0001) == 0x0001 { // odd
			mmc3.Registers.PRGRAMProtect = value
		} else { // even
			oldMirroring := mmc3.mirroring()
			mmc3.Registers.Mirroring = value

			if mmc3.mirroring() != oldMirroring {
				mmc3.setTables(mmc3.Tables())
			}
		}
	// IRQ latch (even) / IRQ reload (odd)
	case address >= 0xc000 && address <= 0xdfff:
		if (address & 0x0001) == 0x0000 { // even
			mmc3.Registers.IRQLatch = value
		} else { // odd
			mmc3.Registers.IRQReload = true
		}
	// IRQ enable (odd) / IRQ disable (even)
	//case address >= 0xe000 && address <= 0xffff:
	case address >= 0xe000:
		if (address & 0x0001) == 0x0001 { // odd
			mmc3.Registers.IRQEnable = true
		} else { // even
			mmc3.Registers.IRQEnable = false
			mmc3.irq(false)
		}
	}

	return
}

func (mmc3 *MMC3) scanlineCounter() {
	// Clock the counter based on MMC3 Rev B behavior
	// From Mesen: counter is reloaded if it's 0 or reload flag is set
	if mmc3.Registers.IRQCounter == 0x00 || mmc3.Registers.IRQReload {
		mmc3.Registers.IRQCounter = mmc3.Registers.IRQLatch
	} else {
		mmc3.Registers.IRQCounter--
	}

	// Trigger IRQ if counter is now 0 and IRQs are enabled
	if mmc3.Registers.IRQCounter == 0x00 && mmc3.Registers.IRQEnable {
		mmc3.irq(true)
	}

	// Clear reload flag after counter operation
	mmc3.Registers.IRQReload = false
}

func (mmc3 *MMC3) bankSelect(flag MMC3BankSelectFlag) (value uint8) {
	reg := mmc3.Registers.BankSelect

	switch flag {
	case BankRegister:
		value = reg & 0x07
	case PRGROMBankMode:
		value = (reg >> 6) & 0x01
	case CHRA12Inversion:
		value = (reg >> 7) & 0x01
	}

	return
}

func (mmc3 *MMC3) mirroring() rp2cgo2.Mirroring {
	switch mmc3.Registers.Mirroring & 0x01 {
	case 0:
		return rp2cgo2.Vertical
	default:
		return rp2cgo2.Horizontal
	}
}

func (mmc3 *MMC3) prgRAMProtect() (chipEnable, allowWrites bool) {
	if (mmc3.Registers.PRGRAMProtect & 0x80) != 0 {
		chipEnable = true
	}

	if (mmc3.Registers.PRGRAMProtect & 0x40) == 0 {
		allowWrites = true
	}

	return
}

func (mmc3 *MMC3) prgBanks() (bank1, bank2, bank3, bank4 uint16) {
	// Mask PRG banks to available ROM banks
	prgBankMask := mmc3.PRGBanks - 1

	mmc3.log.Logf("PRG Bank Mode: %d\nPRG Banks: %d\nPRG Bank Mask: 0x%02X", mmc3.bankSelect(PRGROMBankMode), mmc3.PRGBanks, prgBankMask)
	mmc3.log.Logf("%s", mmc3.Registers.String())

	switch mmc3.bankSelect(PRGROMBankMode) {
	// $8000-$9fff swappable,
	// $c000-$dfff fixed to second-last bank
	case 0:
		bank1 = uint16(mmc3.Registers.PRGBankLow) & prgBankMask
		bank2 = uint16(mmc3.Registers.PRGBankHigh) & prgBankMask
		bank3 = mmc3.PRGBanks - 2
		bank4 = mmc3.PRGBanks - 1
	// $c000-$dfff swappable,
	// $8000-$9fff fixed to second-last bank
	case 1:
		bank1 = mmc3.PRGBanks - 2
		bank2 = uint16(mmc3.Registers.PRGBankHigh) & prgBankMask
		bank3 = uint16(mmc3.Registers.PRGBankLow) & prgBankMask
		bank4 = mmc3.PRGBanks - 1
	}

	return
}

func (mmc3 *MMC3) chrBanks() (bank1, bank2, bank3, bank4, bank5, bank6, bank7, bank8 uint8) {
	// For CHR-ROM, mask to available banks; for CHR-RAM, keep full bank numbers
	// and let Fetch/Store handle the wrapping
	var chrBankMask uint8
	if mmc3.CHRRomSize > 0 {
		chrBankMask = uint8(len(mmc3.VROMBanks) - 1)
	} else {
		chrBankMask = 0xFF // No masking for CHR-RAM
	}

	switch mmc3.bankSelect(CHRA12Inversion) {
	// two 2 KB banks at $0000-$0FFF,
	// four 1 KB banks at $1000-$1FFF
	case 0:
		bank1 = (mmc3.Registers.CHRBank1 & 0xfe) & chrBankMask
		bank2 = (mmc3.Registers.CHRBank1 | 0x01) & chrBankMask
		bank3 = (mmc3.Registers.CHRBank2 & 0xfe) & chrBankMask
		bank4 = (mmc3.Registers.CHRBank2 | 0x01) & chrBankMask
		bank5 = mmc3.Registers.CHRBank3 & chrBankMask
		bank6 = mmc3.Registers.CHRBank4 & chrBankMask
		bank7 = mmc3.Registers.CHRBank5 & chrBankMask
		bank8 = mmc3.Registers.CHRBank6 & chrBankMask
	// two 2 KB banks at $1000-$1FFF,
	// four 1 KB banks at $0000-$0FFF
	case 1:
		bank1 = mmc3.Registers.CHRBank3 & chrBankMask
		bank2 = mmc3.Registers.CHRBank4 & chrBankMask
		bank3 = mmc3.Registers.CHRBank5 & chrBankMask
		bank4 = mmc3.Registers.CHRBank6 & chrBankMask
		bank5 = (mmc3.Registers.CHRBank1 & 0xfe) & chrBankMask
		bank6 = (mmc3.Registers.CHRBank1 | 0x01) & chrBankMask
		bank7 = (mmc3.Registers.CHRBank2 & 0xfe) & chrBankMask
		bank8 = (mmc3.Registers.CHRBank2 | 0x01) & chrBankMask
	}

	return
}

func (mmc3 *MMC3) Tables() (t0, t1, t2, t3 int) {
	switch mmc3.mirroring() {
	case rp2cgo2.Vertical:
		t0, t1, t2, t3 = 0, 1, 0, 1
	case rp2cgo2.Horizontal:
		t0, t1, t2, t3 = 0, 0, 1, 1
	}

	return
}

// IsA12RisingEdge checks if address line A12 has risen (0->1)
// Returns true if A12 went from low to high
func (mmc3 *MMC3) IsA12RisingEdge(addr uint16, cpuClock uint64) bool {
	if addr&0x1000 != 0 {
		// A12 is high
		// Check if we previously tracked a low state
		// Mesen uses 3 CPU cycles minimum for low period
		if mmc3.A12LowClock > 0 && (cpuClock-mmc3.A12LowClock) >= 3 {
			mmc3.A12LowClock = 0
			return true
		}
		mmc3.A12LowClock = 0
	} else if mmc3.A12LowClock == 0 {
		// A12 is low, record the clock
		mmc3.A12LowClock = cpuClock
	}
	return false
}

// NotifyVRAMAddressChange is called when the PPU address bus changes
// This is used to detect A12 rising edges for IRQ counter clocking
func (mmc3 *MMC3) NotifyVRAMAddressChange(addr uint16, cpuClock uint64) {
	if mmc3.IsA12RisingEdge(addr, cpuClock) {
		mmc3.scanlineCounter()
	}
}
