package m65go2

import (
	"fmt"
	"os"
	//s "github.com/stretchr/testify/suite"
)

// region FunctionalTest
type FunctionalTest struct {
	baseTest
}

func (s *FunctionalTest) BeforeTest(suiteName, testName string) {
	r := s.Require()
	// Load the ROM file
	romPath := "test-roms/klaus_dormann/6502_functional_test.bin"
	binData, err := os.ReadFile(romPath)
	r.NoError(err)

	memory := NewBasicMemory(65536)
	cpu := NewM6502(memory)

	// Load ROM into memory at $0400 (code_segment address from listing)
	// The binary is the assembled code starting from code_segment ($0400)
	for i := 0; i < len(binData); i++ {
		memory.Store(0x0400+uint16(i), binData[i])
	}

	// Set PC to $0400 where the test code starts
	cpu.Registers.PC = 0x0400
	cpu.Registers.P = 0x24 // Set flags to match initial state
	cpu.Registers.SP = 0xFD
	//cpu.EnableDecode()

	s.cpu = cpu
}

// endregion FunctionalTest

func (s *FunctionalTest) TestFunctional() {
	s.SkipShort()
	r := s.Require()

	const successPC = 0x3463       // Address indicating successful test completion
	const maxCycles = 100_000_000  // Maximum cycles to run before giving up
	const trapDetectionCycles = 10 // Number of cycles with same PC to detect trap

	var cycleCount uint64
	lastPC := s.cpu.Registers.PC
	stuckCount := 0

	for cycleCount < maxCycles {
		// Execute one instruction
		cycles, err := s.cpu.Execute()
		r.NoError(err, "CPU execution error at PC=$%04X, test_case=$%02X",
			lastPC, s.cpu.Memory.Fetch(0x0200))

		cycleCount += uint64(cycles)
		currentPC := s.cpu.Registers.PC

		if currentPC == successPC {
			s.Logf("Test completed successfully after %d cycles", cycleCount)
			return
		}

		// Check if we're stuck at the same address (trap detection)
		if currentPC == lastPC {
			stuckCount++
			if stuckCount >= trapDetectionCycles {
				// We've hit a trap - read the test case that failed
				testCase := s.cpu.Memory.Fetch(0x0200)
				r.Fail(fmt.Sprintf("Test trapped at PC=$%04X after %d cycles, test_case=$%02X",
					currentPC, cycleCount, testCase))
				return
			}
		} else {
			stuckCount = 0
			lastPC = currentPC
		}

		// Check if we've reached the success loop
		// The functional test ends with a success loop at a specific address
		// You may need to adjust this based on the actual test completion address
		// For now, we'll rely on the trap detection
	}

	// If we get here, we ran out of cycles without completing
	testCase := s.cpu.Memory.Fetch(0x0200)
	r.Fail(fmt.Sprintf("Test did not complete after %d cycles, test_case=$%02X", cycleCount, testCase))
}

// region InterruptTest
type InterruptTest struct {
	baseTest
	portValue uint8
}

func (s *InterruptTest) handleInterrupt() {
	value := s.cpu.Memory.Fetch(0xBFFC)

	// Check if writing to the interrupt feedback port ($BFFC)
	if s.portValue != value {
		oldValue := s.portValue
		s.portValue = value

		// Bit 0 controls IRQ (rising edge triggers)
		if (value&0x01) != 0 && (oldValue&0x01) == 0 {
			s.cpu.Interrupt(Irq, true)
		} else if (value&0x01) == 0 && (oldValue&0x01) != 0 {
			s.cpu.Interrupt(Irq, false)
		}

		// Bit 1 controls NMI (rising edge triggers)
		if (value&0x02) != 0 && (oldValue&0x02) == 0 {
			s.cpu.Interrupt(Nmi, true)
		} else if (value&0x02) == 0 && (oldValue&0x02) != 0 {
			s.cpu.Interrupt(Nmi, false)
		}
	}
}

func (s *InterruptTest) BeforeTest(suiteName, testName string) {
	r := s.Require()
	// Load the ROM file
	romPath := "test-roms/klaus_dormann/6502_interrupt_test.bin"
	binData, err := os.ReadFile(romPath)
	r.NoError(err)

	memory := NewBasicMemory(65536)
	cpu := NewM6502(memory)

	// Load ROM into memory at $0400 (code_segment address from listing)
	// The binary is the assembled code starting from code_segment ($0400)
	for i := range binData {
		memory.Store(0x0400+uint16(i), binData[i])
	}

	// Set PC to $0400 where the test code starts
	cpu.Registers.PC = 0x0400
	cpu.Registers.P = 0x24 // Set flags to match initial state
	cpu.Registers.SP = 0xFD
	//cpu.EnableDecode()
	s.cpu = cpu
}

// endregion InterruptTest

func (s *InterruptTest) TestInterrupt() {
	s.SkipShort()
	r := s.Require()

	const successPC = 0x06EC       // NMOS 6502 success address
	const successWAI = 0x0706      // W65C02 WAI successful address
	const successWAI_IRQ = 0x0723  // W65C02 WAI interrupts enabled successful address
	const maxCycles = 100_000_000  // Maximum cycles to run before giving up
	const trapDetectionCycles = 10 // Number of cycles with same PC to detect trap

	var cycleCount uint64
	lastPC := s.cpu.Registers.PC
	stuckCount := 0

	for cycleCount < maxCycles {
		// Execute one instruction
		cycles, err := s.cpu.Execute()
		r.NoError(err, "CPU execution error at PC=$%04X", lastPC)
		s.handleInterrupt()

		cycleCount += uint64(cycles)
		currentPC := s.cpu.Registers.PC

		// Check if we've reached a success address (jmp * instructions)
		// Success addresses from listing: $06EC, $0706, $0723
		if currentPC == successPC || currentPC == successWAI || currentPC == successWAI_IRQ {
			// Verify it's actually a "jmp *" instruction (stuck loop)
			if currentPC == lastPC {
				stuckCount++
				if stuckCount >= 3 {
					s.Logf("Test completed successfully at PC=$%04X after %d cycles", currentPC, cycleCount)
					return
				}
			}
		}

		// Check if we're stuck at the same address (trap detection for failures)
		if currentPC == lastPC {
			stuckCount++
			if stuckCount >= trapDetectionCycles {
				// Check if it's a success address
				if currentPC == successPC || currentPC == successWAI || currentPC == successWAI_IRQ {
					s.Logf("Test completed successfully at PC=$%04X after %d cycles", currentPC, cycleCount)
					return
				}
				// We've hit a trap - read interrupt test state
				// $0200=nmi_count, $0201=irq_count, $0202=brk_count, $0203=I_src
				nmiCount := s.cpu.Memory.Fetch(0x0200)
				irqCount := s.cpu.Memory.Fetch(0x0201)
				brkCount := s.cpu.Memory.Fetch(0x0202)
				iSrc := s.cpu.Memory.Fetch(0x0203)
				r.Fail(fmt.Sprintf("Test trapped at PC=$%04X after %d cycles (nmi=$%02X irq=$%02X brk=$%02X I_src=$%02X)",
					currentPC, cycleCount, nmiCount, irqCount, brkCount, iSrc))
				return
			}
		} else {
			stuckCount = 0
			lastPC = currentPC
		}
	}

	// If we get here, we ran out of cycles without completing
	iSrc := s.cpu.Memory.Fetch(0x0203)
	r.Fail(fmt.Sprintf("Test did not complete after %d cycles, I_src=$%02X", cycleCount, iSrc))
}

// region DecimalTest
type DecimalTest struct {
	baseTest
}

func (s *DecimalTest) BeforeTest(suiteName, testName string) {
	r := s.Require()
	// Load the ROM file
	romPath := "test-roms/klaus_dormann/6502_decimal_test.bin"
	binData, err := os.ReadFile(romPath)
	r.NoError(err)

	memory := NewBasicMemory(65536)
	cpu := NewM6502(memory)

	// Load ROM into memory at $0400 (code_segment address from listing)
	// The binary is the assembled code starting from code_segment ($0400)
	for i := range binData {
		memory.Store(0x0200+uint16(i), binData[i])
	}

	// Set PC to $0400 where the test code starts
	cpu.Registers.PC = 0x0200
	cpu.Registers.P = 0x24 // Set flags to match initial state
	cpu.Registers.SP = 0xFD
	//cpu.EnableDecode()
	s.cpu = cpu
}

// endregion DecimalTest

type decimal struct {
	N1    uint8
	N2    uint8
	HA    uint8
	HNVZC uint8
	DA    uint8
	DNVZC uint8
	AR    uint8
	NF    uint8
	VF    uint8
	ZF    uint8
	CF    uint8
	ERROR uint8
	N1L   uint8
	N1H   uint8
	N2L   uint8
	N2H   uint16
}

func (d *decimal) update(cpu *M6502) {
	if cpu == nil {
		return
	}

	// operands - register Y = carry in
	// 0x0000 1 byte
	d.N1 = cpu.Memory.Fetch(0x0000)
	// 0x0001 1 byte
	d.N2 = cpu.Memory.Fetch(0x0001)
	// binary result
	// 0x0002 1 byte
	d.HA = cpu.Memory.Fetch(0x0002)
	// 0x0003 1 byte
	d.HNVZC = cpu.Memory.Fetch(0x0003)
	// decimal result
	// 0x0004 1 byte
	d.DA = cpu.Memory.Fetch(0x0004)
	// 0x0005 1 byte
	d.DNVZC = cpu.Memory.Fetch(0x0005)
	// predicted results
	// 0x0006 1 byte
	d.AR = cpu.Memory.Fetch(0x0006)
	// 0x0007 1 byte
	d.NF = cpu.Memory.Fetch(0x0007)

	// 0x0008 1 byte
	d.VF = cpu.Memory.Fetch(0x0008)
	// 0x0009 1 byte
	d.ZF = cpu.Memory.Fetch(0x0009)
	// 0x000a 1 byte
	d.CF = cpu.Memory.Fetch(0x000a)
	// 0x000b 1 byte
	d.ERROR = cpu.Memory.Fetch(0x000b)
	// workspace
	// 0x000c 1 byte
	d.N1L = cpu.Memory.Fetch(0x000c)
	// 0x000d 1 byte
	d.N1H = cpu.Memory.Fetch(0x000d)
	// 0x000e 1 byte
	d.N2L = cpu.Memory.Fetch(0x000e)
	// 0x000f 2 byte
	d.N2H = uint16(cpu.Memory.Fetch(0x000f)) |
		(uint16(cpu.Memory.Fetch(0x0010)) << 8)
}

func (d *decimal) String() string {
	return fmt.Sprintf(
		"Decimal Test State:\n"+
			"ERROR=$%02X\n"+
			"N1=$%02X            N2=$%02X\n"+
			"N1L=$%02X N1H=$%02X   N2L=$%02X N2H=$%04X\n"+
			"Actual Accumulator & Flag Results (Decimal Mode):\n"+
			"DA=$%02X  DNVZC=$%02X\n"+
			"Actual Accumulator & Flag Results (Binary Mode):\n"+
			"HA=$%02X  HNVZC=$%02X\n"+

			"Predicted Results:\n"+
			"AR=$%02X  NF=$%02X "+
			"VF=$%02X  ZF=$%02X "+
			"CF=$%02X\n",
		d.ERROR,
		d.N1, d.N2,
		d.N1L, d.N1H, d.N2L, d.N2H,
		d.DA, d.DNVZC,
		d.HA, d.HNVZC,
		d.AR, d.NF, d.VF, d.ZF, d.CF)
}

func (s *DecimalTest) TestDecimal() {
	s.SkipShort()
	r := s.Require()

	const successPC = 0x025b       // Address indicating successful test completion
	const maxCycles = 0x2500000    // Maximum cycles to run before giving up
	const trapDetectionCycles = 10 // Number of cycles with same PC to detect trap
	var cycleCount uint64
	lastPC := s.cpu.Registers.PC
	stuckCount := 0
	d := &decimal{}

	for cycleCount < maxCycles {
		// Execute one instruction
		cycles, err := s.cpu.Execute()
		currentPC := s.cpu.Registers.PC

		cycleCount += uint64(cycles)
		d.update(s.cpu)
		if currentPC == successPC {
			if d.ERROR == 0 {
				s.Logf("Test completed successfully after %d cycles", cycleCount)
				return
			} else {
				r.Fail(fmt.Sprintf("Test failed at PC=$%04X after %d cycles,\n%s\n%s",
					currentPC, cycleCount, s.cpu.decode.String(), d.String()))
				return
			}
		}
		r.NoError(err, "CPU execution error at PC=$%04X", lastPC)

		// Check if we're stuck at the same address (trap detection)
		if currentPC == lastPC {
			stuckCount++
			if stuckCount >= trapDetectionCycles {
				// We've hit a trap - read the test case that failed
				testCase := s.cpu.Memory.Fetch(0x0200)
				r.Fail(fmt.Sprintf("Test trapped at PC=$%04X after %d cycles, test_case=$%02X",
					currentPC, cycleCount, testCase))
				return
			}
		} else {
			stuckCount = 0
			lastPC = currentPC
		}
	}

	// If we get here, we ran out of cycles without completing
	testCase := s.cpu.Memory.Fetch(0x0200)
	r.Fail(fmt.Sprintf("Test did not complete after %d cycles, test_case=$%02X", cycleCount, testCase))
}
