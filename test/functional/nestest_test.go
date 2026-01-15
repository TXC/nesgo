package functional

import (
	"bufio"
	"fmt"
	"os"

	"github.com/txc/nesgo/pkg/m65go2"
	"github.com/txc/nesgo/pkg/rp2ago3"
	"github.com/txc/nesgo/pkg/utils"
)

type rp2ago3Test struct {
	utils.CommonTest
	cpu *rp2ago3.RP2A03
}

func (s *rp2ago3Test) BeforeTest(suiteName, testName string) {
	r := s.Require()
	// Load the ROM file
	romPath := romPath + "/nestest/nestest.nes"
	romData, err := os.ReadFile(romPath)
	r.NoError(err)

	//cpu := rp2ago3.NewRP2A03(44100)
	cpu := rp2ago3.NewRP2A03(1789773)
	// NES 2A03 CPU does not support decimal mode
	cpu.DisableDecimalMode()

	prgRomSize := s.ValidateRom(romData)
	// PRG ROM starts at byte 16
	prgRom := romData[16 : 16+prgRomSize]

	// Load ROM into memory
	for i := 0; i < len(prgRom) && i < 16384; i++ {
		cpu.Memory.Store(0x8000+uint16(i), prgRom[i])
		cpu.Memory.Store(0xC000+uint16(i), prgRom[i])
	}

	// Set PC to $C000 for automated test mode (bypasses PPU warmup)
	cpu.M6502.Registers.PC = 0xC000
	cpu.M6502.Registers.P = 0x24 // Set flags to match nestest initial state
	cpu.M6502.Registers.SP = 0xFD

	s.cpu = cpu
}

func (s *rp2ago3Test) AfterTest(_, _ string) {
	s.cpu = nil
}

func (s *rp2ago3Test) checkStatusCode() {
	// If we got here without hitting the end marker, test may have passed
	// Check the result byte at $0002
	s.Equal(uint8(0x00),
		s.cpu.Memory.Fetch(0x02),
		"nestest failed with status code %02X at PC:%04X",
		s.cpu.Memory.Fetch(0x02),
		s.cpu.M6502.Registers.PC)

	// Check the result byte at $0003
	s.Equal(uint8(0x00),
		s.cpu.Memory.Fetch(0x03),
		"nestest failed with status code %02X at PC:%04X",
		s.cpu.Memory.Fetch(0x03),
		s.cpu.M6502.Registers.PC)
}

// TestNestest runs the nestest ROM and compares output against the golden log
// nestest.nes is a comprehensive CPU test ROM that tests all official 6502 instructions
// See: https://wiki.nesdev.com/w/index.php/Emulator_tests
func (s *rp2ago3Test) TestNestest() {
	// Load expected log
	logPath := romPath + "nestest/nestest.log"
	logFile, err := os.Open(logPath)
	s.NoError(err)
	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)
	lineNum := 0
	failCount := 0
	maxFails := 10 // Stop after this many failures to avoid spam

	for scanner.Scan() && failCount < maxFails {
		lineNum++
		expectedLine := scanner.Text()

		// Parse expected state from log line
		// Format: "C000  4C F5 C5  JMP $C5F5                       A:00 X:00 Y:00 P:24 SP:FD CYC:  0 SL:241"
		if len(expectedLine) < 48 {
			continue // Skip short lines
		}

		// Extract expected values
		expectedPC := expectedLine[0:4]
		expectedA := expectedLine[50:52]
		expectedX := expectedLine[55:57]
		expectedY := expectedLine[60:62]
		expectedP := expectedLine[65:67]
		expectedSP := expectedLine[71:73]

		// Get actual state before execution
		actualPC := fmt.Sprintf("%04X", s.cpu.M6502.Registers.PC)
		actualA := fmt.Sprintf("%02X", s.cpu.M6502.Registers.A)
		actualX := fmt.Sprintf("%02X", s.cpu.M6502.Registers.X)
		actualY := fmt.Sprintf("%02X", s.cpu.M6502.Registers.Y)
		actualP := fmt.Sprintf("%02X", s.cpu.M6502.Registers.P)
		actualSP := fmt.Sprintf("%02X", s.cpu.M6502.Registers.SP)

		// Compare states
		s.Equal(expectedPC, actualPC, "PC mismatch at line %d", lineNum)
		s.Equal(expectedA, actualA, "A mismatch at line %d", lineNum)
		s.Equal(expectedX, actualX, "X mismatch at line %d", lineNum)
		s.Equal(expectedY, actualY, "Y mismatch at line %d", lineNum)
		s.Equal(expectedP, actualP, "P mismatch at line %d", lineNum)
		s.Equal(expectedSP, actualSP, "SP mismatch at line %d", lineNum)

		// Execute one instruction
		_, err := s.cpu.Execute()
		s.NoError(err, "Execution error at line %d", lineNum)
		// Check if we've reached the end of test (infinite loop at $C66E)
		if s.cpu.M6502.Registers.PC == 0xC66E {
			// Test completed - check status codes
			s.checkStatusCode()
			s.Logf("nestest PASSED! All official instructions working correctly")
			return
		}
	}

	s.NoError(scanner.Err(), "Error reading log file")

	// If we got here without hitting the end marker, test may have passed
	s.checkStatusCode()
}

// TestAutomated is a simpler version that just checks the final status
func (s *rp2ago3Test) TestAutomated() {
	// Run until we hit the infinite loop at $C66E or max instructions
	maxInstructions := 10000
	for range maxInstructions {
		pc := s.cpu.M6502.Registers.PC

		// Check for infinite loop (test end marker)
		if pc == 0xC66E {
			s.checkStatusCode()
			return
		}

		_, err := s.cpu.Execute()
		s.NoError(err, "Execution error at PC=%04X", pc)
	}

	s.Fail("Test did not complete after %d instructions (PC=%04X)",
		maxInstructions,
		s.cpu.M6502.Registers.PC)
}

// TestWithTrace runs nestest and outputs trace for debugging
func (s *rp2ago3Test) TestWithTrace() {
	s.SkipShort()

	logPath := "../../nestest_trace.log"

	// Create trace output file
	traceFile, err := os.Create(logPath)
	s.NoError(err)
	defer traceFile.Close()

	writer := bufio.NewWriter(traceFile)
	defer writer.Flush()

	maxInstructions := 10000
	for range make([]struct{}, maxInstructions) {
		pc := s.cpu.M6502.Registers.PC

		// Format: PC  BYTES  INSTRUCTION  A:XX X:XX Y:XX P:XX SP:XX
		opcode := s.cpu.Memory.Fetch(pc)
		inst := s.cpu.Instructions.Get(m65go2.OpCode(opcode))
		s.NotNil(inst)

		bytes := fmt.Sprintf("%02X", opcode)

		// Add operand bytes
		switch inst.AddressingMode {
		case m65go2.Immediate, m65go2.ZeroPage, m65go2.ZeroPageX,
			m65go2.ZeroPageY, m65go2.IndirectX, m65go2.IndirectY,
			m65go2.Relative:
			if pc+1 < 0xFFFF {
				bytes += fmt.Sprintf(" %02X", s.cpu.Memory.Fetch(pc+1))
			}
		case m65go2.Absolute, m65go2.AbsoluteX, m65go2.AbsoluteY, m65go2.Indirect:
			if pc+2 < 0xFFFF {
				bytes += fmt.Sprintf(" %02X %02X", s.cpu.Memory.Fetch(pc+1), s.cpu.Memory.Fetch(pc+2))
			}
		}

		trace := fmt.Sprintf("%04X  %-8s  A:%02X X:%02X Y:%02X P:%02X SP:%02X\n",
			pc, bytes,
			s.cpu.M6502.Registers.A, s.cpu.M6502.Registers.X, s.cpu.M6502.Registers.Y,
			s.cpu.M6502.Registers.P, s.cpu.M6502.Registers.SP)
		writer.WriteString(trace)

		if pc == 0xC66E {
			statusCode := s.cpu.Memory.Fetch(0x02)
			fmt.Fprintf(writer, "\nTest completed with status: %02X\n", statusCode)
			break
		}

		_, err := s.cpu.Execute()
		s.NoError(err)
	}

	s.Logf("Trace written to %s", logPath)
}

type NesTest struct {
	utils.CommonTest
	cpu *rp2ago3.RP2A03
}

func (s *NesTest) checkStatusCode() {
	// If we got here without hitting the end marker, test may have passed
	// Check the result byte at $0002
	s.Equal(uint8(0x00),
		s.cpu.Memory.Fetch(0x02),
		"nestest failed with status code %02X at PC:%04X",
		s.cpu.Memory.Fetch(0x02),
		s.cpu.M6502.Registers.PC)

	// Check the result byte at $0003
	s.Equal(uint8(0x00),
		s.cpu.Memory.Fetch(0x03),
		"nestest failed with status code %02X at PC:%04X",
		s.cpu.Memory.Fetch(0x03),
		s.cpu.M6502.Registers.PC)
}

func (s *NesTest) TestNestest() {
	// Load expected log
	logPath := romPath + "nestest/nestest.log"
	logFile, err := os.Open(logPath)
	s.NoError(err)
	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)
	lineNum := 0
	failCount := 0
	maxFails := 10 // Stop after this many failures to avoid spam

	for scanner.Scan() && failCount < maxFails {
		lineNum++
		expectedLine := scanner.Text()

		// Parse expected state from log line
		// Format: "C000  4C F5 C5  JMP $C5F5                       A:00 X:00 Y:00 P:24 SP:FD CYC:  0 SL:241"
		if len(expectedLine) < 48 {
			continue // Skip short lines
		}

		// Extract expected values
		expectedPC := expectedLine[0:4]
		expectedA := expectedLine[50:52]
		expectedX := expectedLine[55:57]
		expectedY := expectedLine[60:62]
		expectedP := expectedLine[65:67]
		expectedSP := expectedLine[71:73]

		// Get actual state before execution
		actualPC := fmt.Sprintf("%04X", s.cpu.M6502.Registers.PC)
		actualA := fmt.Sprintf("%02X", s.cpu.M6502.Registers.A)
		actualX := fmt.Sprintf("%02X", s.cpu.M6502.Registers.X)
		actualY := fmt.Sprintf("%02X", s.cpu.M6502.Registers.Y)
		actualP := fmt.Sprintf("%02X", s.cpu.M6502.Registers.P)
		actualSP := fmt.Sprintf("%02X", s.cpu.M6502.Registers.SP)

		// Compare states
		s.Equal(expectedPC, actualPC, "PC mismatch at line %d", lineNum)
		s.Equal(expectedA, actualA, "A mismatch at line %d", lineNum)
		s.Equal(expectedX, actualX, "X mismatch at line %d", lineNum)
		s.Equal(expectedY, actualY, "Y mismatch at line %d", lineNum)
		s.Equal(expectedP, actualP, "P mismatch at line %d", lineNum)
		s.Equal(expectedSP, actualSP, "SP mismatch at line %d", lineNum)

		// Execute one instruction
		_, err := s.cpu.Execute()
		s.NoError(err, "Execution error at line %d", lineNum)
		// Check if we've reached the end of test (infinite loop at $C66E)
		if s.cpu.M6502.Registers.PC == 0xC66E {
			// Test completed - check status codes
			s.checkStatusCode()
			s.Logf("nestest PASSED! All official instructions working correctly")
			return
		}
	}

	s.NoError(scanner.Err(), "Error reading log file")

	// If we got here without hitting the end marker, test may have passed
	s.checkStatusCode()
}
