package functional

import (
	"os"

	"github.com/txc/nesgo/pkg/m65go2"
	//"github.com/txc/nesgo/pkg/rp2cgo2"
	//"github.com/txc/nesgo/pkg/utils"
)

// region CycleTest
type CycleTest struct {
	baseTest
}

func (s *CycleTest) BeforeTest(suiteName, testName string) {
	r := s.Require()
	// Load the ROM file
	romPath := romPath + "/ed_spittles/6502_cycle_test.bin"
	binData, err := os.ReadFile(romPath)
	r.NoError(err)

	memory := m65go2.NewBasicMemory(65536)
	cpu := m65go2.NewM6502(memory)

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

	// cycle_test_data.csv

	s.cpu = cpu
}

// endregion CycleTest

func (s *CycleTest) TestCycles() {
	s.SkipShort()
	r := s.Require()

	const successPC = 0xFF3F   // Address indicating successful test completion
	const maxCycles = 2000     // Maximum cycles to run before giving up
	const successCycles = 1141 // Expected cycles for successful completion
	var cycleCount uint64

	for cycleCount < maxCycles {
		// Execute one instruction
		cycles, err := s.cpu.Execute()
		r.NoError(err)

		cycleCount += uint64(cycles)
		currentPC := s.cpu.Registers.PC

		if currentPC == successPC {
			s.Logf("Test completed successfully after %d cycles", cycleCount)
			return
		}
	}

	r.Equal(successCycles, cycleCount)
}
