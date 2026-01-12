package m65go2

import (
	"bufio"
	"fmt"
	"os"
)

type NesTest struct {
	baseTest
}

func (s *NesTest) BeforeTest(suiteName, testName string) {
	r := s.Require()
	// Load the ROM file
	romPath := "test-roms/nestest/nestest.nes"
	romData, err := os.ReadFile(romPath)
	r.NoError(err)

	// NES ROM format: 16-byte header + PRG ROM + CHR ROM
	// Header format: NES\x1a + 16KB PRG banks + 8KB CHR banks + flags...
	r.Greater(len(romData), 16)
	r.Equal([]uint8{0x4E, 0x45, 0x53, 0x1A}, romData[0:4])

	prgRomSize := int(romData[4]) * 16384 // Number of 16KB banks
	r.Greater(len(romData), 16+prgRomSize)
	if err != nil {
		s.Fail("Skipping nestest: ROM file not found")
		return
	}

	if len(romData) < 16 || string(romData[0:4]) != "NES\x1a" {
		s.FailNow("Invalid NES ROM")
	}

	// PRG ROM starts at byte 16
	prgRom := romData[16 : 16+prgRomSize]

	memory := NewBasicMemory(65536)
	cpu := NewM6502(memory)
	// NES 2A03 CPU does not support decimal mode
	cpu.DisableDecimalMode()

	// Load ROM into memory
	for i := 0; i < len(prgRom) && i < 16384; i++ {
		memory.Store(0x8000+uint16(i), prgRom[i])
		memory.Store(0xC000+uint16(i), prgRom[i])
	}

	// Set PC to $C000 for automated test mode (bypasses PPU warmup)
	cpu.Registers.PC = 0xC000
	cpu.Registers.P = 0x24 // Set flags to match nestest initial state
	cpu.Registers.SP = 0xFD

	s.cpu = cpu
}

func (s *NesTest) checkStatusCode() {
	// If we got here without hitting the end marker, test may have passed
	// Check the result byte at $0002
	s.Equal(uint8(0x00), s.cpu.Memory.Fetch(0x02), "nestest failed with status code %02X at PC:%04X", s.cpu.Memory.Fetch(0x02), s.cpu.Registers.PC)

	// Check the result byte at $0003
	s.Equal(uint8(0x00), s.cpu.Memory.Fetch(0x03), "nestest failed with status code %02X at PC:%04X", s.cpu.Memory.Fetch(0x03), s.cpu.Registers.PC)
}

// TestNestest runs the nestest ROM and compares output against the golden log
// nestest.nes is a comprehensive CPU test ROM that tests all official 6502 instructions
// See: https://wiki.nesdev.com/w/index.php/Emulator_tests
func (s *NesTest) TestNestest() {
	// Load expected log
	logPath := "test-roms/nestest/nestest.log"
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
		actualPC := fmt.Sprintf("%04X", s.cpu.Registers.PC)
		actualA := fmt.Sprintf("%02X", s.cpu.Registers.A)
		actualX := fmt.Sprintf("%02X", s.cpu.Registers.X)
		actualY := fmt.Sprintf("%02X", s.cpu.Registers.Y)
		actualP := fmt.Sprintf("%02X", s.cpu.Registers.P)
		actualSP := fmt.Sprintf("%02X", s.cpu.Registers.SP)

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
		if s.cpu.Registers.PC == 0xC66E {
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
func (s *NesTest) TestAutomated() {
	// Run until we hit the infinite loop at $C66E or max instructions
	maxInstructions := 10000
	for range maxInstructions {
		pc := s.cpu.Registers.PC

		// Check for infinite loop (test end marker)
		if pc == 0xC66E {
			s.checkStatusCode()
			return
		}

		_, err := s.cpu.Execute()
		s.NoError(err, "Execution error at PC=%04X", pc)
	}

	s.Fail("Test did not complete after %d instructions (PC=%04X)", maxInstructions, s.cpu.Registers.PC)
}

// TestWithTrace runs nestest and outputs trace for debugging
func (s *NesTest) TestWithTrace() {
	s.SkipShort()

	// Create trace output file
	traceFile, err := os.Create("test-roms/nestest/nestest_trace.log")
	s.NoError(err)
	defer traceFile.Close()

	writer := bufio.NewWriter(traceFile)
	defer writer.Flush()

	maxInstructions := 10000
	for range make([]struct{}, maxInstructions) {
		pc := s.cpu.Registers.PC

		// Format: PC  BYTES  INSTRUCTION  A:XX X:XX Y:XX P:XX SP:XX
		opcode := s.cpu.Memory.Fetch(pc)
		inst := s.cpu.Instructions.Get(OpCode(opcode))
		s.NotNil(inst)

		bytes := fmt.Sprintf("%02X", opcode)

		// Add operand bytes
		switch inst.AddressingMode {
		case Immediate, ZeroPage, ZeroPageX, ZeroPageY, IndirectX, IndirectY, Relative:
			if pc+1 < 0xFFFF {
				bytes += fmt.Sprintf(" %02X", s.cpu.Memory.Fetch(pc+1))
			}
		case Absolute, AbsoluteX, AbsoluteY, Indirect:
			if pc+2 < 0xFFFF {
				bytes += fmt.Sprintf(" %02X %02X", s.cpu.Memory.Fetch(pc+1), s.cpu.Memory.Fetch(pc+2))
			}
		}

		trace := fmt.Sprintf("%04X  %-8s  A:%02X X:%02X Y:%02X P:%02X SP:%02X\n",
			pc, bytes,
			s.cpu.Registers.A, s.cpu.Registers.X, s.cpu.Registers.Y,
			s.cpu.Registers.P, s.cpu.Registers.SP)
		writer.WriteString(trace)

		if pc == 0xC66E {
			statusCode := s.cpu.Memory.Fetch(0x02)
			fmt.Fprintf(writer, "\nTest completed with status: %02X\n", statusCode)
			break
		}

		_, err := s.cpu.Execute()
		s.NoError(err)
	}

	s.Log("Trace written to test-roms/nestest/nestest_trace.log")
}
