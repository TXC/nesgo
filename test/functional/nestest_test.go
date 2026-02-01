package functional

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/txc/nesgo/pkg/m65go2"
	//"github.com/txc/nesgo/pkg/rp2ago3"
	"github.com/txc/nesgo/pkg/nes"
	"github.com/txc/nesgo/pkg/utils"
)

type NesTestValues struct {
	PC      uint16
	A       uint8
	X       uint8
	Y       uint8
	P       m65go2.Status
	SP      uint8
	lineNum int
}

func NewNesTestValues(ns *nes.NES, lineNum int) *NesTestValues {
	return &NesTestValues{
		PC:      ns.CPU.M6502.Registers.PC,
		A:       ns.CPU.M6502.Registers.A,
		X:       ns.CPU.M6502.Registers.X,
		Y:       ns.CPU.M6502.Registers.Y,
		P:       ns.CPU.M6502.Registers.P,
		SP:      ns.CPU.M6502.Registers.SP,
		lineNum: lineNum,
	}
}

func NewNesTestValuesFromString(value string, lineNum int) *NesTestValues {
	expectedPC, err := strconv.Atoi(value[0:4])
	if err != nil {
		panic(err)
	}
	expectedA, err := strconv.Atoi(value[50:52])
	if err != nil {
		panic(err)
	}
	expectedX, err := strconv.Atoi(value[55:57])
	if err != nil {
		panic(err)
	}
	expectedY, err := strconv.Atoi(value[60:62])
	if err != nil {
		panic(err)
	}
	expectedP, err := strconv.Atoi(value[65:67])
	if err != nil {
		panic(err)
	}
	expectedSP, err := strconv.Atoi(value[71:73])
	if err != nil {
		panic(err)
	}

	return &NesTestValues{
		PC:      uint16(expectedPC),
		A:       uint8(expectedA),
		X:       uint8(expectedX),
		Y:       uint8(expectedY),
		P:       m65go2.Status(expectedP),
		SP:      uint8(expectedSP),
		lineNum: lineNum,
	}
}

func (hm *NesTestValues) String() string {
	return fmt.Sprintf("\n"+
		"PC: 0x%04X,\n"+
		"A: 0x%02X,\n"+
		"X: 0x%02X,\n"+
		"Y: 0x%02X,\n"+
		"P: 0x%02X,\n"+
		"SP: 0x%02X,\n",
		hm.PC, hm.A, hm.X, hm.Y, hm.P, hm.SP)
}

func (hm *NesTestValues) Equals(s *NesTest, actual *NesTestValues) {
	// Compare states
	s.Equal(hm.PC, actual.PC, "PC mismatch at line %d", actual.lineNum)
	s.Equal(hm.A, actual.A, "A mismatch at line %d", actual.lineNum)
	s.Equal(hm.X, actual.X, "X mismatch at line %d", actual.lineNum)
	s.Equal(hm.Y, actual.Y, "Y mismatch at line %d", actual.lineNum)
	s.Equal(hm.P, actual.P, "P mismatch at line %d", actual.lineNum)
	s.Equal(hm.SP, actual.SP, "SP mismatch at line %d", actual.lineNum)
}

type NesTest struct {
	utils.CommonTest
	//cpu *rp2ago3.RP2A03
	nes *nes.NES
}

func (s *NesTest) BeforeTest(suiteName, testName string) {
	r := s.Require()
	// Load the ROM file
	romPath := romPath + "/nestest/nestest.nes"
	options := &nes.Options{
		Region: "NTSC",
		//Recorder: "jpeg",
	}

	ns, err := nes.NewNES(romPath, options)
	r.NoError(err)

	romData, err := os.ReadFile(romPath)
	r.NoError(err)

	//cpu := rp2ago3.NewRP2A03(1789773)
	// NES 2A03 CPU does not support decimal mode
	//cpu.DisableDecimalMode()

	prgRomSize := s.ValidateRom(romData)
	// PRG ROM starts at byte 16
	prgRom := romData[16 : 16+prgRomSize]

	// Load ROM into memory
	for i := 0; i < len(prgRom) && i < 0x4000; i++ {
		ns.CPU.Memory.Store(0x8000+uint16(i), prgRom[i])
		ns.CPU.Memory.Store(0xC000+uint16(i), prgRom[i])
	}

	// Set PC to $C000 for automated test mode (bypasses PPU warmup)
	ns.CPU.M6502.Registers.PC = 0xC000
	ns.CPU.M6502.Registers.P = 0x24 // Set flags to match nestest initial state
	ns.CPU.M6502.Registers.SP = 0xFD

	s.nes = ns
}

func (s *NesTest) AfterTest(_, _ string) {
	//s.cpu = nil
	s.nes = nil
}

func (s *NesTest) checkStatusCode() {
	// If we got here without hitting the end marker, test may have passed
	// Check the result byte at $0002
	s.Equal(uint8(0x00),
		s.nes.CPU.Memory.Fetch(0x02),
		"nestest failed with status code %02X at PC:%04X",
		s.nes.CPU.Memory.Fetch(0x02),
		s.nes.CPU.M6502.Registers.PC)

	// Check the result byte at $0003
	s.Equal(uint8(0x00),
		s.nes.CPU.Memory.Fetch(0x03),
		"nestest failed with status code %02X at PC:%04X",
		s.nes.CPU.Memory.Fetch(0x03),
		s.nes.CPU.M6502.Registers.PC)
}

// TestNestest runs the nestest ROM and compares output against the golden log
// nestest.nes is a comprehensive CPU test ROM that tests all official 6502 instructions
// See: https://wiki.nesdev.com/w/index.php/Emulator_tests
func (s *NesTest) TestNestest() {
	// Load expected log
	logPath := romPath + "/nestest/nestest.log"
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

		expected := NewNesTestValuesFromString(expectedLine, lineNum)
		NewNesTestValues(s.nes, lineNum).Equals(s, expected)

		// Execute one instruction
		_, err := s.nes.CPU.Execute()
		s.NoError(err, "Execution error at line %d", lineNum)
		// Check if we've reached the end of test (infinite loop at $C66E)
		if s.nes.CPU.M6502.Registers.PC == 0xC66E {
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
		pc := s.nes.CPU.M6502.Registers.PC

		// Check for infinite loop (test end marker)
		if pc == 0xC66E {
			s.checkStatusCode()
			return
		}

		_, err := s.nes.CPU.Execute()
		s.NoError(err, "Execution error at PC=%04X", pc)
	}

	s.Fail("Test did not complete after %d instructions (PC=%04X)",
		maxInstructions,
		s.nes.CPU.M6502.Registers.PC)
}

// TestWithTrace runs nestest and outputs trace for debugging
func (s *NesTest) TestWithTrace() {
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
		pc := s.nes.CPU.M6502.Registers.PC

		// Format: PC  BYTES  INSTRUCTION  A:XX X:XX Y:XX P:XX SP:XX
		opcode := s.nes.CPU.Memory.Fetch(pc)
		inst := s.nes.CPU.Instructions.GetOp(m65go2.OpCode(opcode))
		s.NotNil(inst)

		bytes := fmt.Sprintf("%02X", opcode)

		// Add operand bytes
		switch inst.AddressingMode {
		case m65go2.Immediate, m65go2.ZeroPage, m65go2.ZeroPageX,
			m65go2.ZeroPageY, m65go2.IndirectX, m65go2.IndirectY,
			m65go2.Relative:
			if pc+1 < 0xFFFF {
				bytes += fmt.Sprintf(" %02X", s.nes.CPU.Memory.Fetch(pc+1))
			}
		case m65go2.Absolute, m65go2.AbsoluteX, m65go2.AbsoluteY, m65go2.Indirect:
			if pc+2 < 0xFFFF {
				bytes += fmt.Sprintf(" %02X %02X", s.nes.CPU.Memory.Fetch(pc+1), s.nes.CPU.Memory.Fetch(pc+2))
			}
		}

		trace := fmt.Sprintf("%04X  %-8s  A:%02X X:%02X Y:%02X P:%02X SP:%02X\n",
			pc, bytes,
			s.nes.CPU.M6502.Registers.A, s.nes.CPU.M6502.Registers.X, s.nes.CPU.M6502.Registers.Y,
			s.nes.CPU.M6502.Registers.P, s.nes.CPU.M6502.Registers.SP)
		writer.WriteString(trace)

		if pc == 0xC66E {
			statusCode := s.nes.CPU.Memory.Fetch(0x02)
			fmt.Fprintf(writer, "\nTest completed with status: %02X\n", statusCode)
			break
		}

		_, err := s.nes.CPU.Execute()
		s.NoError(err)
	}

	s.Logf("Trace written to %s", logPath)
}
