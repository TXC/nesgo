package m65go2

import (
	"fmt"
	"log"
	"testing"

	s "github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	s.Run(t, &InitializationTest{})
	s.Run(t, &CPUTest{})
	s.Run(t, &InstructionADC{})
	s.Run(t, &InstructionAND{})
	s.Run(t, &InstructionASL{})
	s.Run(t, &InstructionBCC{})
	s.Run(t, &InstructionBCS{})
	s.Run(t, &InstructionBEQ{})
	s.Run(t, &InstructionBIT{})
	s.Run(t, &InstructionBMI{})
	s.Run(t, &InstructionBNE{})
	s.Run(t, &InstructionBPL{})
	s.Run(t, &InstructionBVC{})
	s.Run(t, &InstructionBVS{})
	s.Run(t, &InstructionCLC{})
	s.Run(t, &InstructionCLD{})
	s.Run(t, &InstructionCLI{})
	s.Run(t, &InstructionCLV{})
	s.Run(t, &InstructionCMP{})
	s.Run(t, &InstructionCPX{})
	s.Run(t, &InstructionCPY{})
	s.Run(t, &InstructionDEC{})
	s.Run(t, &InstructionDEX{})
	s.Run(t, &InstructionDEY{})
	s.Run(t, &InstructionEOR{})
	s.Run(t, &InstructionINC{})
	s.Run(t, &InstructionINX{})
	s.Run(t, &InstructionINY{})
	s.Run(t, &InstructionJMP{})
	s.Run(t, &InstructionJSR{})
	s.Run(t, &InstructionLDA{})
	s.Run(t, &InstructionLDX{})
	s.Run(t, &InstructionLDY{})
	s.Run(t, &InstructionLSR{})
	s.Run(t, &InstructionORA{})
	s.Run(t, &InstructionPHA{})
	s.Run(t, &InstructionPHP{})
	s.Run(t, &InstructionPLA{})
	s.Run(t, &InstructionPLP{})
	s.Run(t, &InstructionROL{})
	s.Run(t, &InstructionROR{})
	s.Run(t, &InstructionRTI{})
	s.Run(t, &InstructionRTS{})
	s.Run(t, &InstructionSBC{})
	s.Run(t, &InstructionSEC{})
	s.Run(t, &InstructionSED{})
	s.Run(t, &InstructionSEI{})
	s.Run(t, &InstructionSTA{})
	s.Run(t, &InstructionSTX{})
	s.Run(t, &InstructionSTY{})
	s.Run(t, &InstructionTAX{})
	s.Run(t, &InstructionTSX{})
	s.Run(t, &InstructionTXS{})
	s.Run(t, &AccumulatorAddress{})
	s.Run(t, &IndexAddress{})
	s.Run(t, &CompareAddress{})
	s.Run(t, &DecrementIncrementAddress{})
	s.Run(t, &StoreInMemoryAddress{})
	s.Run(t, &ClockCyclesTest{})
	s.Run(t, &ProgramCounterTest{})
	// Unofficial Instructions
	s.Run(t, &InstructionLAX{})
	s.Run(t, &InstructionSAX{})
	s.Run(t, &InstructionDCP{})
	s.Run(t, &InstructionISB{})
	s.Run(t, &InstructionSLO{})
	s.Run(t, &InstructionSRE{})
	s.Run(t, &InstructionRLA{})
	s.Run(t, &InstructionRRA{})
	s.Run(t, &InstructionALR{})
	s.Run(t, &InstructionANC{})
	s.Run(t, &InstructionARR{})
	s.Run(t, &InstructionSBX{})
	s.Run(t, &InstructionLAS{})
}

// region baseTest
type baseTest struct {
	s.Suite
	cpu *M6502
	// Subtest Instruction To Test
	_instruction *Instruction
	// Previous log prefix
	_oldPrefix string
}

func (s *baseTest) BeforeTest(suiteName, testName string) {
	s.cpu = NewM6502(NewBasicMemory(DefaultMemorySize))
	s.cpu.Reset()
	s.cpu.breakError = true

	//s.cpu.decode.enabled = true
	//s.cpu.debug |= DebugModeCycleDump
}

func (s *baseTest) SetupSubTest() {
	s.cpu.Reset()

	if s._instruction == nil {
		return
	}
	s._oldPrefix = s.cpu.log.Prefix()
	s.cpu.log.SetPrefix(fmt.Sprintf("%s/%s: ", s._instruction.Mnemonic, s._instruction.AddressingMode))
}

func (s *baseTest) AfterTest(_, _ string) {
	s.cpu = nil
}

func (s *baseTest) TearDownSubTest() {
	if s._instruction != nil {
		s._instruction = nil
	}
	if s._oldPrefix != "" {
		s.cpu.log.SetPrefix(s._oldPrefix)
	}
}

func (s *baseTest) loadProgram(address uint16, values []uint8) {
	for i, v := range values {
		s.cpu.Memory.Store(address+uint16(i), v)
	}
	hi, lo := s.cpu.GetHighLow(address)
	s.cpu.Memory.Store(0xFFFC, lo)
	s.cpu.Memory.Store(0xFFFD, hi)
	s.cpu.PerformRst()
	s.cpu.currentCycles = 0
}

func (s *baseTest) Log(args ...any) { s.T().Log(args...) }

func (s *baseTest) Logf(format string, args ...any) { s.T().Logf(format, args...) }

// endregion baseTest

// region InitializationTest
type InitializationTest struct {
	baseTest
}

func (s *InitializationTest) TestProcessorStatusFlagsInitializedCorrectly() {
	s.Assertions.False(s.cpu.getFlag(C))
	s.Assertions.False(s.cpu.getFlag(Z))
	s.Assertions.True(s.cpu.getFlag(I))
	s.Assertions.False(s.cpu.getFlag(D))
	s.Assertions.False(s.cpu.getFlag(B))
	s.Assertions.True(s.cpu.getFlag(U))
	s.Assertions.False(s.cpu.getFlag(V))
	s.Assertions.False(s.cpu.getFlag(N))
}

func (s *InitializationTest) TestProcessorStatusFlagsSetsAndUnsetsCorrectly() {
	// Test all flags except B and U which have special behavior
	subtests := []Status{C, Z, I, D, V, N}
	for _, flag := range subtests {
		s.Run(fmt.Sprintf("Flag %v", flag), func() {
			s.cpu.Registers.P = 0x00
			s.False(s.cpu.getFlag(flag))
			s.cpu.Registers.P |= flag
			s.True(s.cpu.getFlag(flag))
			s.cpu.Registers.P &^= flag
			s.False(s.cpu.getFlag(flag))
		})
	}
}

func (s *InitializationTest) TestProcessorRegistersInitializedCorrectly() {
	s.cpu.Reset()
	s.Equal(uint8(0x00), s.cpu.Registers.A)
	s.Equal(uint8(0x00), s.cpu.Registers.X)
	s.Equal(uint8(0x00), s.cpu.Registers.Y)
	s.Equal(uint16(0xffff), s.cpu.Registers.PC)
	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(s.cpu.Registers.PC))
}

func (s *InitializationTest) TestStackPointerInitializesCorrectly() {
	s.Equal(uint8(0xFD), s.cpu.Registers.SP)
}

func (s *InitializationTest) TestStackPointerInitializesToDefaultValueAfterReset() {
	s.cpu.Registers.SP = 0x00
	s.cpu.Reset()
	s.Equal(uint8(0xFD), s.cpu.Registers.SP)
}

// endregion InitializationTest

// region CPUTest
type CPUTest struct {
	baseTest
}

func (s *CPUTest) TestCPUFromBCD() {
	s.Equal(uint8(0), s.cpu.FromBCD(0))
	s.Equal(uint8(5), s.cpu.FromBCD(0x05))
	s.Equal(uint8(11), s.cpu.FromBCD(0x11))
	s.Equal(uint8(99), s.cpu.FromBCD(0x99))
}

func (s *CPUTest) TestCPUToBCD() {
	s.Equal(uint8(0), s.cpu.ToBCD(0))
	s.Equal(uint8(0x05), s.cpu.ToBCD(5))
	s.Equal(uint8(0x11), s.cpu.ToBCD(11))
	s.Equal(uint8(0x99), s.cpu.ToBCD(99))
}

func (s *CPUTest) TestCPUFromTwosComplement() {
	s.Equal(int8(0), s.cpu.fromTwosCom(0x00))
	s.Equal(int8(1), s.cpu.fromTwosCom(0x01))
	s.Equal(int8(127), s.cpu.fromTwosCom(0x7F))
	s.Equal(int8(-1), s.cpu.fromTwosCom(0xFF))
	s.Equal(int8(-128), s.cpu.fromTwosCom(0x80))
}

func (s *CPUTest) TestCPUToTwosComplement() {
	s.Equal(uint8(0x00), s.cpu.toTwosCom(0))
	s.Equal(uint8(0x01), s.cpu.toTwosCom(1))
	s.Equal(uint8(0x7F), s.cpu.toTwosCom(127))
	s.Equal(uint8(0xFF), s.cpu.toTwosCom(-1))
	s.Equal(uint8(0x80), s.cpu.toTwosCom(-128))
}

func (s *CPUTest) TestCPUAddressingModeCyclesCorrect() {
	subtests := []struct {
		mode   AddressingMode
		cycles uint16
	}{
		{Accumulator, 1},
		{Immediate, 0},
		{Implied, 1},
		{ZeroPage, 1},
		{ZeroPageX, 2},
		{ZeroPageY, 2},
		{Absolute, 2},
		{AbsoluteX, 2},
		{AbsoluteY, 2},
		{Indirect, 4},
		{IndirectX, 4},
		{IndirectY, 3},
		{Relative, 1},
	}

	for _, data := range subtests {
		s.Run(data.mode.String(), func() {
			startCycles := s.cpu.currentCycles
			s.cpu.getAddress(data.mode, nil)
			actual := s.cpu.currentCycles - startCycles
			s.EqualValues(data.cycles, actual)
		})
	}
}

// endregion CPUTest

// region InstructionADC
type InstructionADC struct {
	baseTest
}

func (s *InstructionADC) TestADCAccumulatorCorrectWhenNotInBCDMode() {
	subtests := []struct {
		a        uint8
		value    uint8
		carryIn  bool
		expected uint8
	}{
		{0x00, 0x00, false, 0x00},
		{0x00, 0x01, false, 0x01},
		{0x01, 0x02, false, 0x03},
		{0xFF, 0x01, false, 0x00},
		{0xFE, 0x01, false, 0xFF},
		{0xFF, 0x00, false, 0xFF},
		{0x00, 0x00, true, 0x01},
		{0x00, 0x01, true, 0x02},
		{0x01, 0x02, true, 0x04},
		{0xFE, 0x01, true, 0x00},
		{0xFD, 0x01, true, 0xFF},
		{0xFE, 0x00, true, 0xFF},
		{0xFF, 0xFF, true, 0xFF},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X CarryIn:%v", data.a, data.value, data.carryIn), func() {
			program := []uint8{0xA9, data.a, 0x69, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.setFlag(C, data.carryIn)
			s.cpu.Execute() // LDA
			s.EqualValues(data.a, s.cpu.Registers.A)
			s.cpu.Execute() // ADC
			s.EqualValues(data.expected, s.cpu.Registers.A)
		})
	}
}

func (s *InstructionADC) TestADCAccumulatorCorrectWhenInBCDMode() {
	// BCD is a bit special: 0x99 + 0x99 = 0x98, which equals: 99 + 99 = 198 = 98 (with carry)
	// it's hex in dec form: hex(99) = dec(99)
	subtests := []struct {
		a        uint8
		value    uint8
		carryIn  bool
		expected uint8
	}{
		{0x99, 0x99, false, 0x98},
		{0x99, 0x99, true, 0x99},
		{0x90, 0x99, false, 0x89},
		{0x01, 0x69, true, 0x71},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X CarryIn:%v", data.a, data.value, data.carryIn), func() {
			program := []uint8{0xF8, 0xA9, data.a, 0x69, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.setFlag(C, data.carryIn)
			s.cpu.Execute() // SED
			s.cpu.Execute() // LDA
			s.EqualValues(data.a, s.cpu.Registers.A)
			s.cpu.Execute() // ADC
			s.EqualValues(data.expected, s.cpu.Registers.A)
		})
	}
}

func (s *InstructionADC) TestADCCarryCorrectWhenNotInBCDMode() {
	subtests := []struct {
		a            uint8
		value        uint8
		carryIn      bool
		expectedFlag bool
	}{
		{0xFE, 0x01, false, false},
		{0xFE, 0x01, true, true},
		{0xFD, 0x01, true, false},
		{0xFF, 0x01, false, true},
		{0xFF, 0x01, true, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X CarryIn:%v", data.a, data.value, data.carryIn), func() {
			program := []uint8{0xA9, data.a, 0x69, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.setFlag(C, data.carryIn)
			s.cpu.Execute() // LDA
			s.EqualValues(data.a, s.cpu.Registers.A)
			s.cpu.Execute() // ADC
			s.Equal(data.expectedFlag, s.cpu.getFlag(C))
		})
	}
}

func (s *InstructionADC) TestADCCarryCorrectWhenInBCDMode() {
	subtests := []struct {
		a            uint8
		value        uint8
		carryIn      bool
		expectedFlag bool
	}{
		{0x62, 0x01, false, false},
		{0x62, 0x01, true, false},
		{0x63, 0x01, false, false},
		{0x63, 0x01, true, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X CarryIn:%v", data.a, data.value, data.carryIn), func() {
			program := []uint8{0xA9, data.a, 0x69, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.setFlag(C, data.carryIn)
			s.cpu.Execute() // LDA
			s.EqualValues(data.a, s.cpu.Registers.A)
			s.cpu.Execute() // ADC
			s.Equal(data.expectedFlag, s.cpu.getFlag(C))
		})
	}
}

func (s *InstructionADC) TestADCZeroFlagCorrectWhenNotInBCDMode() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x00, 0x00, true},
		{0xFF, 0x01, true},
		{0x00, 0x01, false},
		{0x01, 0x00, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0x69, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.EqualValues(data.a, s.cpu.Registers.A)
			s.cpu.Execute() // ADC
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionADC) TestADCNegativeFlagCorrect() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x7E, 0x01, false},
		{0x01, 0x7E, false},
		{0x01, 0x7F, true},
		{0x7F, 0x01, true},
		{0x01, 0xFE, true},
		{0xFE, 0x01, true},
		{0x01, 0xFF, false},
		{0xFF, 0x01, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0x69, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.EqualValues(data.a, s.cpu.Registers.A)
			s.cpu.Execute() // ADC
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

func (s *InstructionADC) TestADCOverflowFlagCorrect() {
	subtests := []struct {
		a            uint8
		value        uint8
		carryIn      bool
		expectedFlag bool
	}{
		{0x00, 0x7F, false, false},
		{0x00, 0x80, false, false},
		{0x01, 0x7F, false, true},
		{0x01, 0x80, false, false},
		{0x7F, 0x01, false, true},
		{0x7F, 0x7F, false, true},
		{0x80, 0x7F, false, false},
		{0x80, 0x80, false, true},
		{0x80, 0x81, false, true},
		{0x80, 0xFF, false, true},
		{0xFF, 0x00, false, false},
		{0xFF, 0x01, false, false},
		{0xFF, 0x7F, false, false},
		{0xFF, 0x80, false, true},
		{0xFF, 0xFF, false, false},
		{0x00, 0x7F, true, true},
		{0x00, 0x80, true, false},
		{0x01, 0x7F, true, true},
		{0x01, 0x80, true, false},
		{0x7F, 0x01, true, true},
		{0x7F, 0x7F, true, true},
		{0x80, 0x7F, true, false},
		{0x80, 0x80, true, true},
		{0x80, 0x81, true, true},
		{0x80, 0xFF, true, false},
		{0xFF, 0x00, true, false},
		{0xFF, 0x01, true, false},
		{0xFF, 0x7F, true, false},
		{0xFF, 0x80, true, false},
		{0xFF, 0xFF, true, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X CarryIn:%v", data.a, data.value, data.carryIn), func() {
			program := []uint8{0xA9, data.a, 0x69, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.setFlag(C, data.carryIn)
			s.cpu.Execute() // LDA
			s.EqualValues(data.a, s.cpu.Registers.A)
			s.cpu.Execute() // ADC
			s.Equal(data.expectedFlag, s.cpu.getFlag(V))
		})
	}
}

// endregion InstructionADC

// region InstructionAND
type InstructionAND struct {
	baseTest
}

func (s *InstructionAND) TestANDAccumulatorCorrect() {
	subtests := []struct {
		a        uint8
		value    uint8
		expected uint8
	}{
		{0, 0, 0},
		{255, 255, 255},
		{255, 254, 254},
		{170, 85, 0},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0x29, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // AND
			s.EqualValues(data.expected, s.cpu.Registers.A)
		})
	}
}

// endregion InstructionAND

// region InstructionASL
type InstructionASL struct {
	baseTest
}

func (s *InstructionASL) TestASLCorrectValueStored() {
	subtests := []struct {
		opcode   uint8
		value    uint8
		expected uint8
		address  uint8
	}{
		{0x0A, 0x6D, 0xDA, 0x00}, // ASL Accumulator
		{0x0A, 0x6C, 0xD8, 0x00}, // ASL Accumulator
		{0x06, 0x6D, 0xDA, 0x01}, // ASL Zero Page
		{0x16, 0x6D, 0xDA, 0x01}, // ASL Zero Page X
		{0x0E, 0x6D, 0xDA, 0x01}, // ASL Absolute
		{0x1E, 0x6D, 0xDA, 0x01}, // ASL Absolute X
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X Value:%02X", data.opcode, data.value), func() {
			switch data.opcode {
			case 0x0A:
				// Accumulator mode
				program := []uint8{0xA9, data.value, data.opcode}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // LDA
				s.EqualValues(data.value, s.cpu.Registers.A)
				s.cpu.Execute() // ASL
				s.EqualValues(data.expected, s.cpu.Registers.A)
			case 0x0E, 0x1E:
				// Absolute addressing mode - need 16-bit address
				s.cpu.Memory.Store(uint16(data.address), data.value)
				program := []uint8{data.opcode, data.address, 0x00}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // ASL
				s.EqualValues(data.expected, s.cpu.Memory.Fetch(uint16(data.address)))
			default:
				// Zero page addressing mode
				s.cpu.Memory.Store(uint16(data.address), data.value)
				program := []uint8{data.opcode, data.address}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // ASL
				s.EqualValues(data.expected, s.cpu.Memory.Fetch(uint16(data.address)))
			}
		})
	}
}

func (s *InstructionASL) TestASLCarrySetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x7F, false},
		{0x80, true},
		{0xFF, true},
		{0x00, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value, 0x0A}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.EqualValues(data.value, s.cpu.Registers.A)
			s.cpu.Execute() // ASL
			s.Equal(data.expectedFlag, s.cpu.getFlag(C))
		})
	}
}

func (s *InstructionASL) TestASLNegativeSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x3F, false},
		{0x40, true},
		{0x7F, true},
		{0x80, false},
		{0x00, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value, 0x0A}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.EqualValues(data.value, s.cpu.Registers.A)
			s.cpu.Execute() // ASL
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

func (s *InstructionASL) TestASLZeroSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x7F, false},
		{0x80, true},
		{0x00, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value, 0x0A}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.EqualValues(data.value, s.cpu.Registers.A)
			s.cpu.Execute() // ASL
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

// endregion InstructionASL

// region InstructionBCC
type InstructionBCC struct {
	baseTest
}

func (s *InstructionBCC) TestBCCProgramCounterCorrect() {
	subtests := []struct {
		pc       uint16
		offset   uint8
		expected uint16
	}{
		{0x00, 0x01, 0x03},
		{0x80, 0x80, 0x02},
		{0x00, 0x03, 0x05},
		{0x00, 0xFD, 0xFFFF},
		{0x7D, 0x80, 0xFFFF},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("PC:%04X Offset:%02X", data.pc, data.offset), func() {
			program := []uint8{0x90, data.offset}
			s.loadProgram(data.pc, program)
			s.cpu.Execute() // BCC
			s.EqualValues(data.expected, s.cpu.Registers.PC)
		})
	}
}

// endregion InstructionBCC

// region InstructionBCS
type InstructionBCS struct {
	baseTest
}

func (s *InstructionBCS) TestBCSProgramCounterCorrect() {
	subtests := []struct {
		pc       uint16
		offset   uint8
		expected uint16
	}{
		{0x00, 0x01, 0x04},
		{0x80, 0x80, 0x03},
		{0x00, 0xFC, 0xFFFF},
		{0x7C, 0x80, 0xFFFF},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("PC:%04X Offset:%02X", data.pc, data.offset), func() {
			program := []uint8{0x38, 0xB0, data.offset}
			s.loadProgram(data.pc, program)
			s.cpu.Execute() // SEC
			s.cpu.Execute() // BCS
			s.EqualValues(data.expected, s.cpu.Registers.PC)
		})
	}
}

// endregion InstructionBCS

// region InstructionBCS
type InstructionBEQ struct {
	baseTest
}

func (s *InstructionBEQ) TestBEQProgramCounterCorrect() {
	subtests := []struct {
		pc       uint16
		offset   uint8
		expected uint16
	}{
		{0x00, 0x01, 0x05},
		{0x80, 0x80, 0x04},
		{0x00, 0xFB, 0xFFFF},
		{0x7B, 0x80, 0xFFFF},
		{0x02, 0xFE, 0x04},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("PC:%04X Offset:%02X", data.pc, data.offset), func() {
			program := []uint8{0xA9, 0x00, 0xF0, data.offset}
			s.loadProgram(data.pc, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // BEQ
			s.EqualValues(data.expected, s.cpu.Registers.PC)
		})
	}
}

// endregion InstructionBCS

// region InstructionBIT
type InstructionBIT struct {
	baseTest
}

func (s *InstructionBIT) TestBITNegativeSetWhenComparisonIsNegativeNumber() {
	subtests := []struct {
		opcode       uint8
		aValue       uint8
		memValue     uint8
		expectedFlag bool
	}{
		{0x24, 0x7F, 0x7F, false}, // BIT Zero Page
		{0x24, 0x80, 0x7F, false}, // BIT Zero Page
		{0x24, 0x7F, 0x80, true},  // BIT Zero Page
		{0x24, 0x80, 0xFF, true},  // BIT Zero Page
		{0x24, 0xFF, 0x80, true},  // BIT Zero Page
		{0x2C, 0x7F, 0x7F, false}, // BIT Absolute
		{0x2C, 0x80, 0x7F, false}, // BIT Absolute
		{0x2C, 0x7F, 0x80, true},  // BIT Absolute
		{0x2C, 0x80, 0xFF, true},  // BIT Absolute
		{0x2C, 0xFF, 0x80, true},  // BIT Absolute
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X A:%02X Mem:%02X", data.opcode, data.aValue, data.memValue), func() {
			// Store the memory value at the address that BIT will test
			s.cpu.Memory.Store(0x0006, data.memValue)
			if data.opcode == 0x24 {
				// Zero page addressing
				program := []uint8{0xA9, data.aValue, data.opcode, 0x06}
				s.loadProgram(0x0600, program)
			} else {
				// Absolute addressing
				program := []uint8{0xA9, data.aValue, data.opcode, 0x06, 0x00}
				s.loadProgram(0x0600, program)
			}
			s.cpu.Execute() // LDA
			s.cpu.Execute() // BIT
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

func (s *InstructionBIT) TestBITOverflowSetByBitSix() {
	subtests := []struct {
		opcode       uint8
		aValue       uint8
		memValue     uint8
		expectedFlag bool
	}{
		{0x24, 0x3F, 0x3F, false}, // BIT Zero Page
		{0x24, 0x3F, 0x40, true},  // BIT Zero Page
		{0x24, 0x40, 0x3F, false}, // BIT Zero Page
		{0x24, 0x40, 0x7F, true},  // BIT Zero Page
		{0x24, 0x7F, 0x40, true},  // BIT Zero Page
		{0x24, 0x7F, 0x80, false}, // BIT Zero Page
		{0x24, 0x80, 0x7F, true},  // BIT Zero Page
		{0x24, 0xC0, 0xDF, true},  // BIT Zero Page
		{0x24, 0xDF, 0xC0, true},  // BIT Zero Page
		{0x24, 0x3F, 0x3F, false}, // BIT Zero Page
		{0x24, 0xC0, 0xFF, true},  // BIT Zero Page
		{0x24, 0xFF, 0xC0, true},  // BIT Zero Page
		{0x24, 0x40, 0xFF, true},  // BIT Zero Page
		{0x24, 0xFF, 0x40, true},  // BIT Zero Page
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X A:%02X Mem:%02X", data.opcode, data.aValue, data.memValue), func() {
			// Store the memory value at the address that BIT will test
			s.cpu.Memory.Store(0x0006, data.memValue)
			program := []uint8{0xA9, data.aValue, data.opcode, 0x06}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // BIT
			s.Equal(data.expectedFlag, s.cpu.getFlag(V))
		})
	}
}

func (s *InstructionBIT) TestBITZeroSetWhenComparisonIsZero() {
	subtests := []struct {
		opcode       uint8
		aValue       uint8
		memValue     uint8
		expectedFlag bool
	}{
		{0x24, 0x00, 0x00, true},  // BIT Zero Page
		{0x24, 0xFF, 0xFF, false}, // BIT Zero Page
		{0x24, 0xAA, 0x55, true},  // BIT Zero Page
		{0x24, 0x55, 0xAA, true},  // BIT Zero Page
		{0x2C, 0x00, 0x00, true},  // BIT Absolute
		{0x2C, 0xFF, 0xFF, false}, // BIT Absolute
		{0x2C, 0xAA, 0x55, true},  // BIT Absolute
		{0x2C, 0x55, 0xAA, true},  // BIT Absolute
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X A:%02X Mem:%02X", data.opcode, data.aValue, data.memValue), func() {
			// Store the memory value at the address that BIT will test
			s.cpu.Memory.Store(0x0006, data.memValue)
			if data.opcode == 0x24 {
				// Zero page addressing
				program := []uint8{0xA9, data.aValue, data.opcode, 0x06}
				s.loadProgram(0x0600, program)
			} else {
				// Absolute addressing
				program := []uint8{0xA9, data.aValue, data.opcode, 0x06, 0x00}
				s.loadProgram(0x0600, program)
			}
			s.cpu.Execute() // LDA
			s.cpu.Execute() // BIT
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

// endregion InstructionBIT

// region InstructionBMI
type InstructionBMI struct {
	baseTest
}

func (s *InstructionBMI) TestBMIProgramCounterCorrect() {
	subtests := []struct {
		pc       uint16
		offset   uint8
		expected uint16
	}{
		{0x00, 0x01, 0x05},
		{0x80, 0x80, 0x04},
		{0x00, 0xFB, 0xFFFF},
		{0x7B, 0x80, 0xFFFF},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("PC:%04X Offset:%02X", data.pc, data.offset), func() {
			program := []uint8{0xA9, 0x80, 0x30, data.offset}
			s.loadProgram(data.pc, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // BMI
			s.EqualValues(data.expected, s.cpu.Registers.PC)
		})
	}
}

// endregion InstructionBMI

// region InstructionBNE
type InstructionBNE struct {
	baseTest
}

func (s *InstructionBNE) TestBNEProgramCounterCorrect() {
	subtests := []struct {
		pc       uint16
		offset   uint8
		expected uint16
	}{
		{0x00, 0x01, 0x05},
		{0x80, 0x80, 0x04},
		{0x00, 0xFB, 0xFFFF},
		{0x7B, 0x80, 0xFFFF},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("PC:%04X Offset:%02X", data.pc, data.offset), func() {
			program := []uint8{0xA9, 0x80, 0xD0, data.offset}
			s.loadProgram(data.pc, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // BNE
			s.EqualValues(data.expected, s.cpu.Registers.PC)
		})
	}
}

// endregion InstructionBNE

// region InstructionBPL
type InstructionBPL struct {
	baseTest
}

func (s *InstructionBPL) TestBPLProgramCounterCorrect() {
	subtests := []struct {
		pc       uint16
		offset   uint8
		expected uint16
	}{
		{0x00, 0x01, 0x05},
		{0x80, 0x80, 0x04},
		{0x00, 0xFB, 0xFFFF},
		{0x7B, 0x80, 0xFFFF},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("PC:%04X Offset:%02X", data.pc, data.offset), func() {
			program := []uint8{0xA9, 0x79, 0x10, data.offset}
			s.loadProgram(data.pc, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // BPL
			s.EqualValues(data.expected, s.cpu.Registers.PC)
		})
	}
}

// endregion InstructionBPL

// region InstructionBVC
type InstructionBVC struct {
	baseTest
}

func (s *InstructionBVC) TestBVCProgramCounterCorrect() {
	subtests := []struct {
		pc       uint16
		offset   uint8
		expected uint16
	}{
		{0x00, 0x01, 0x03},
		{0x80, 0x80, 0x02},
		{0x00, 0xFD, 0xFFFF},
		{0x7D, 0x80, 0xFFFF},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("PC:%04X Offset:%02X", data.pc, data.offset), func() {
			program := []uint8{0x50, data.offset}
			s.loadProgram(data.pc, program)
			s.cpu.Execute() // BVC
			s.EqualValues(data.expected, s.cpu.Registers.PC)
		})
	}
}

// endregion InstructionBVC

// region InstructionBVS
type InstructionBVS struct {
	baseTest
}

func (s *InstructionBVS) TestBVSProgramCounterCorrect() {
	subtests := []struct {
		pc       uint16
		offset   uint8
		expected uint16
	}{
		{0x00, 0x01, 0x07},
		{0x80, 0x80, 0x06},
		{0x00, 0xF9, 0xFFFF},
		{0x79, 0x80, 0xFFFF},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("PC:%04X Offset:%02X", data.pc, data.offset), func() {
			program := []uint8{0xA9, 0x01, 0x69, 0x7F, 0x70, data.offset}
			s.loadProgram(data.pc, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // ADC
			s.cpu.Execute() // BVS
			s.EqualValues(data.expected, s.cpu.Registers.PC)
		})
	}
}

// endregion InstructionBVS

// region InstructionCLC
type InstructionCLC struct {
	baseTest
}

func (s *InstructionCLC) TestCLCCarryFlagClearedCorrectly() {
	program := []uint8{0x18}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // CLC
	s.False(s.cpu.getFlag(C))
}

// endregion InstructionCLC

// region InstructionCLD
type InstructionCLD struct {
	baseTest
}

func (s *InstructionCLD) TestCLDDecimalFlagSetAndClearedCorrectly() {
	program := []uint8{0xF8, 0xD8}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // SED
	s.cpu.Execute() // CLD
	s.False(s.cpu.getFlag(D))
}

// endregion InstructionCLD

// region InstructionCLI
type InstructionCLI struct {
	baseTest
}

func (s *InstructionCLI) TestCLIInterruptFlagClearedCorrectly() {
	program := []uint8{0x58}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // CLI
	s.False(s.cpu.getFlag(I))
}

// endregion InstructionCLI

// region InstructionCLV
type InstructionCLV struct {
	baseTest
}

func (s *InstructionCLV) TestCLVOverflowFlagClearedCorrectly() {
	program := []uint8{0xB8}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // CLV
	s.False(s.cpu.getFlag(V))
}

// endregion InstructionCLV

// region InstructionCMP
type InstructionCMP struct {
	baseTest
}

func (s *InstructionCMP) TestCMPZeroFlagSetWhenValuesMatch() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x00, 0x00, true},
		{0xFF, 0x00, false},
		{0x00, 0xFF, false},
		{0xFF, 0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0xC9, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // CMP
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionCMP) TestCMPCarryFlagSetWhenAccumulatorIsGreaterThanOrEqual() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x00, 0x00, true},
		{0xFF, 0x00, true},
		{0x00, 0xFF, false},
		{0x00, 0x01, false},
		{0xFF, 0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0xC9, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // CMP
			s.Equal(data.expectedFlag, s.cpu.getFlag(C))
		})
	}
}

func (s *InstructionCMP) TestCMPNegativeFlagSetWhenResultIsNegative() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0xFE, 0xFF, true},
		{0x81, 0x01, true},
		{0x81, 0x02, false},
		{0x79, 0x01, false},
		{0x00, 0x01, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0xC9, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // CMP
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionCMP

// region InstructionCPX
type InstructionCPX struct {
	baseTest
}

func (s *InstructionCPX) TestCPXZeroFlagSetWhenValuesMatch() {
	subtests := []struct {
		x            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x00, 0x00, true},
		{0xFF, 0x00, false},
		{0x00, 0xFF, false},
		{0xFF, 0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("X:%02X Value:%02X", data.x, data.value), func() {
			program := []uint8{0xA2, data.x, 0xE0, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // CPX
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionCPX) TestCPXCarryFlagSetWhenXIsGreaterThanOrEqual() {
	subtests := []struct {
		x            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x00, 0x00, true},
		{0xFF, 0x00, true},
		{0x00, 0xFF, false},
		{0x00, 0x01, false},
		{0xFF, 0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("X:%02X Value:%02X", data.x, data.value), func() {
			program := []uint8{0xA2, data.x, 0xE0, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // CPX
			s.Equal(data.expectedFlag, s.cpu.getFlag(C))
		})
	}
}

func (s *InstructionCPX) TestCPXNegativeFlagSetWhenResultIsNegative() {
	subtests := []struct {
		x            uint8
		value        uint8
		expectedFlag bool
	}{
		{0xFE, 0xFF, true},
		{0x81, 0x01, true},
		{0x81, 0x02, false},
		{0x79, 0x01, false},
		{0x00, 0x01, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("X:%02X Value:%02X", data.x, data.value), func() {
			program := []uint8{0xA2, data.x, 0xE0, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // CPX
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionCPX

// region InstructionCPY
type InstructionCPY struct {
	baseTest
}

func (s *InstructionCPY) TestCPYZeroFlagSetWhenValuesMatch() {
	subtests := []struct {
		y            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x00, 0x00, true},
		{0xFF, 0x00, false},
		{0x00, 0xFF, false},
		{0xFF, 0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Y:%02X Value:%02X", data.y, data.value), func() {
			program := []uint8{0xA0, data.y, 0xC0, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.cpu.Execute() // CPY
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionCPY) TestCPYCarryFlagSetWhenYIsGreaterThanOrEqual() {
	subtests := []struct {
		y            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x00, 0x00, true},
		{0xFF, 0x00, true},
		{0x00, 0xFF, false},
		{0x00, 0x01, false},
		{0xFF, 0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Y:%02X Value:%02X", data.y, data.value), func() {
			program := []uint8{0xA0, data.y, 0xC0, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.cpu.Execute() // CPY
			s.Equal(data.expectedFlag, s.cpu.getFlag(C))
		})
	}
}

func (s *InstructionCPY) TestCPYNegativeFlagSetWhenResultIsNegative() {
	subtests := []struct {
		y            uint8
		value        uint8
		expectedFlag bool
	}{
		{0xFE, 0xFF, true},
		{0x81, 0x01, true},
		{0x81, 0x02, false},
		{0x79, 0x01, false},
		{0x00, 0x01, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Y:%02X Value:%02X", data.y, data.value), func() {
			program := []uint8{0xA0, data.y, 0xC0, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.cpu.Execute() // CPY
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionCPY

// region InstructionDEC
type InstructionDEC struct {
	baseTest
}

func (s *InstructionDEC) TestDECMemoryHasCorrectValue() {
	subtests := []struct {
		initial  uint8
		expected uint8
	}{
		{0x00, 0xFF},
		{0xFF, 0xFE},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			s.cpu.Memory.Store(0x03, data.initial)
			program := []uint8{0xC6, 0x03}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // DEC
			s.EqualValues(data.expected, s.cpu.Memory.Fetch(0x03))
		})
	}
}

func (s *InstructionDEC) TestDECZeroHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x00, false},
		{0x01, true},
		{0x02, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			s.cpu.Memory.Store(0x03, data.initial)
			program := []uint8{0xC6, 0x03}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // DEC
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionDEC) TestDECNegativeHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x80, false},
		{0x81, true},
		{0x00, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			s.cpu.Memory.Store(0x03, data.initial)
			program := []uint8{0xC6, 0x03}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // DEC
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionDEC

// region InstructionDEX
type InstructionDEX struct {
	baseTest
}

func (s *InstructionDEX) TestDEXXRegisterHasCorrectValue() {
	subtests := []struct {
		initial  uint8
		expected uint8
	}{
		{0x00, 0xFF},
		{0xFF, 0xFE},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA2, data.initial, 0xCA}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // DEX
			s.EqualValues(data.expected, s.cpu.Registers.X)
		})
	}
}

func (s *InstructionDEX) TestDEXZeroHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x00, false},
		{0x01, true},
		{0x02, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA2, data.initial, 0xCA}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // DEX
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionDEX) TestDEXNegativeHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x80, false},
		{0x81, true},
		{0x00, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA2, data.initial, 0xCA}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // DEX
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionDEX

// region InstructionDEY
type InstructionDEY struct {
	baseTest
}

func (s *InstructionDEY) TestDEYYRegisterHasCorrectValue() {
	subtests := []struct {
		initial  uint8
		expected uint8
	}{
		{0x00, 0xFF},
		{0xFF, 0xFE},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA0, data.initial, 0x88}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.cpu.Execute() // DEY
			s.EqualValues(data.expected, s.cpu.Registers.Y)
		})
	}
}

func (s *InstructionDEY) TestDEYZeroHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x00, false},
		{0x01, true},
		{0x02, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA0, data.initial, 0x88}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.cpu.Execute() // DEY
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionDEY) TestDEYNegativeHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x80, false},
		{0x81, true},
		{0x00, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA0, data.initial, 0x88}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.cpu.Execute() // DEY
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionDEY

// region InstructionEOR
type InstructionEOR struct {
	baseTest
}

func (s *InstructionEOR) TestEORAccumulatorCorrect() {
	subtests := []struct {
		a        uint8
		value    uint8
		expected uint8
	}{
		{0x00, 0x00, 0x00},
		{0xFF, 0x00, 0xFF},
		{0x00, 0xFF, 0xFF},
		{0x55, 0xAA, 0xFF},
		{0xFF, 0xFF, 0x00},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0x49, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // EOR
			s.EqualValues(data.expected, s.cpu.Registers.A)
		})
	}
}

func (s *InstructionEOR) TestEORNegativeFlagCorrect() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0xFF, 0xFF, false},
		{0x80, 0x7F, true},
		{0x40, 0x3F, false},
		{0xFF, 0x7F, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0x49, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // EOR
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

func (s *InstructionEOR) TestEORZeroFlagCorrect() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0xFF, 0xFF, true},
		{0x80, 0x7F, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0x49, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // EOR
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

// endregion InstructionEOR

// region InstructionINC
type InstructionINC struct {
	baseTest
}

func (s *InstructionINC) TestINCMemoryHasCorrectValue() {
	subtests := []struct {
		initial  uint8
		expected uint8
	}{
		{0x00, 0x01},
		{0xFF, 0x00},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			s.cpu.Memory.Store(0x03, data.initial)
			program := []uint8{0xE6, 0x03}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // INC
			s.EqualValues(data.expected, s.cpu.Memory.Fetch(0x03))
		})
	}
}

func (s *InstructionINC) TestINCZeroHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x00, false},
		{0xFF, true},
		{0xFE, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			s.cpu.Memory.Store(0x03, data.initial)
			program := []uint8{0xE6, 0x03}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // INC
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionINC) TestINCNegativeHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x78, false},
		{0x80, true},
		{0x00, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			s.cpu.Memory.Store(0x03, data.initial)
			program := []uint8{0xE6, 0x03}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // INC
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionINC

// region InstructionINX
type InstructionINX struct {
	baseTest
}

func (s *InstructionINX) TestINXXRegisterHasCorrectValue() {
	subtests := []struct {
		initial  uint8
		expected uint8
	}{
		{0x00, 0x01},
		{0xFF, 0x00},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA2, data.initial, 0xE8}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // INX
			s.EqualValues(data.expected, s.cpu.Registers.X)
		})
	}
}

func (s *InstructionINX) TestINXZeroHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x00, false},
		{0xFF, true},
		{0xFE, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA2, data.initial, 0xE8}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // INX
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionINX) TestINXNegativeHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x78, false},
		{0x80, true},
		{0x00, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA2, data.initial, 0xE8}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // INX
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionINX

// region InstructionINY
type InstructionINY struct {
	baseTest
}

func (s *InstructionINY) TestINYYRegisterHasCorrectValue() {
	subtests := []struct {
		initial  uint8
		expected uint8
	}{
		{0x00, 0x01},
		{0xFF, 0x00},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA0, data.initial, 0xC8}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.cpu.Execute() // INY
			s.EqualValues(data.expected, s.cpu.Registers.Y)
		})
	}
}

func (s *InstructionINY) TestINYZeroHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x00, false},
		{0xFF, true},
		{0xFE, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA0, data.initial, 0xC8}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.cpu.Execute() // INY
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionINY) TestINYNegativeHasCorrectValue() {
	subtests := []struct {
		initial      uint8
		expectedFlag bool
	}{
		{0x78, false},
		{0x80, true},
		{0x00, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Initial:%02X", data.initial), func() {
			program := []uint8{0xA0, data.initial, 0xC8}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.cpu.Execute() // INY
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionINY

// region InstructionJMP
type InstructionJMP struct {
	baseTest
}

func (s *InstructionJMP) TestJMPProgramCounterSetCorrectlyAfterJump() {
	program := []uint8{0x4C, 0x08, 0x00}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // JMP
	s.EqualValues(0x08, s.cpu.Registers.PC)
}

func (s *InstructionJMP) TestJMPProgramCounterSetCorrectlyAfterIndirectJump() {
	program := []uint8{0x6C, 0x03, 0x00}
	s.loadProgram(0x0600, program)
	// Store target address at $0003 (little endian)
	s.cpu.Memory.Store(0x0003, 0x08)
	s.cpu.Memory.Store(0x0004, 0x00)
	s.cpu.Execute() // JMP
	s.EqualValues(0x08, s.cpu.Registers.PC)
}

func (s *InstructionJMP) TestJMPIndirectWrapsCorrectIfMSBIsFF() {
	// Test the 6502 page boundary wrap bug
	// JMP ($01FF) should read LSB from $01FF and MSB from $0100 (not $0200)
	program := []uint8{0x6C, 0xFF, 0x01}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x01FF, 0x03) // LSB of target address
	s.cpu.Memory.Store(0x0100, 0x02) // MSB of target address (wraps within page)
	s.cpu.Execute()                  // JMP ($01FF)
	s.EqualValues(0x0203, s.cpu.Registers.PC)
}

// endregion InstructionJMP

// region InstructionJSR
type InstructionJSR struct {
	baseTest
}

func (s *InstructionJSR) TestJSRStackLoadsCorrectValue() {
	program := []uint8{0x20, 0xCC, 0xCC}
	s.loadProgram(0xBBAA, program)
	stackLocation := s.cpu.Registers.SP
	s.cpu.Execute() // JSR
	s.EqualValues(0xBB, s.cpu.Memory.Fetch(0x0100+uint16(stackLocation)))
	s.EqualValues(0xAC, s.cpu.Memory.Fetch(0x0100+uint16(stackLocation-1)))
}

func (s *InstructionJSR) TestJSRProgramCounterCorrect() {
	program := []uint8{0x20, 0xCC, 0xCC}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // JSR
	s.EqualValues(0xCCCC, s.cpu.Registers.PC)
}

func (s *InstructionJSR) TestJSRStackPointerCorrect() {
	program := []uint8{0x20, 0xCC, 0xCC}
	s.loadProgram(0x0600, program)
	stackLocation := s.cpu.Registers.SP
	s.cpu.Execute() // JSR
	s.EqualValues(stackLocation-2, s.cpu.Registers.SP)
}

// endregion InstructionJSR

// region InstructionLDA
type InstructionLDA struct {
	baseTest
}

func (s *InstructionLDA) TestLDAAccumulatorHasCorrectValue() {
	subtests := []uint8{0x03, 0x91, 0x01, 0xFF}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data), func() {
			program := []uint8{0xA9, data}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.EqualValues(data, s.cpu.Registers.A)
		})
	}
}

func (s *InstructionLDA) TestLDAZeroSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x00, true},
		{0x03, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionLDA) TestLDANegativeSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x00, false},
		{0x79, false},
		{0x80, true},
		{0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionLDA

// region InstructionLDX
type InstructionLDX struct {
	baseTest
}

func (s *InstructionLDX) TestLDXXRegisterValueHasCorrectValue() {
	subtests := []uint8{0x00, 0x79, 0x80, 0xFF}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data), func() {
			program := []uint8{0xA2, data}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.EqualValues(data, s.cpu.Registers.X)
		})
	}
}

func (s *InstructionLDX) TestLDXNegativeFlagSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x00, false},
		{0x79, false},
		{0x80, true},
		{0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA2, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

func (s *InstructionLDX) TestLDXZeroSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x00, true},
		{0x03, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA2, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

// endregion InstructionLDX

// region InstructionLDY
type InstructionLDY struct {
	baseTest
}

func (s *InstructionLDY) TestLDYYRegisterValueHasCorrectValue() {
	subtests := []uint8{0x00, 0x79, 0x80, 0xFF}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data), func() {
			program := []uint8{0xA0, data}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.EqualValues(data, s.cpu.Registers.Y)
		})
	}
}

func (s *InstructionLDY) TestLDYNegativeFlagSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x00, false},
		{0x79, false},
		{0x80, true},
		{0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA0, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

func (s *InstructionLDY) TestLDYZeroSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x00, true},
		{0x03, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA0, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDY
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

// endregion InstructionLDY

// region InstructionLSR
type InstructionLSR struct {
	baseTest
}

func (s *InstructionLSR) TestLSRNegativeSetCorrectly() {
	subtests := []struct {
		value        uint8
		setCarry     bool
		expectedFlag bool
	}{
		{0xFF, false, false},
		{0xFE, false, false},
		{0xFF, true, false},
		{0x00, true, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X SetCarry:%v", data.value, data.setCarry), func() {
			program := []uint8{0xA9, data.value, 0x4A}
			s.loadProgram(0x0600, program)
			if data.setCarry {
				s.cpu.setFlag(C, true)
			} else {
				s.cpu.setFlag(C, false)
			}
			s.cpu.Execute() // LDA
			s.cpu.Execute() // LSR
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

func (s *InstructionLSR) TestLSRZeroSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x01, true},
		{0x02, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value, 0x4A}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // LSR
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionLSR) TestLSRCarryFlagSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x01, true},
		{0x02, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value, 0x4A}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // LSR
			s.Equal(data.expectedFlag, s.cpu.getFlag(C))
		})
	}
}

func (s *InstructionLSR) TestLSRCorrectValueStored() {
	subtests := []struct {
		opcode   uint8
		value    uint8
		expected uint8
		address  uint8
	}{
		{0x4A, 0xFF, 0x7F, 0x00}, // LSR Accumulator
		{0x4A, 0xFD, 0x7E, 0x00}, // LSR Accumulator
		{0x46, 0xFF, 0x7F, 0x01}, // LSR Zero Page
		{0x56, 0xFF, 0x7F, 0x01}, // LSR Zero Page X
		{0x4E, 0xFF, 0x7F, 0x01}, // LSR Absolute
		{0x5E, 0xFF, 0x7F, 0x01}, // LSR Absolute X
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X Value:%02X", data.opcode, data.value), func() {
			switch data.opcode {
			case 0x4A:
				// Accumulator mode
				program := []uint8{0xA9, data.value, data.opcode}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // LDA
				s.cpu.Execute() // LSR
				s.EqualValues(data.expected, s.cpu.Registers.A)
			case 0x4E, 0x5E:
				// Absolute addressing mode
				s.cpu.Memory.Store(uint16(data.address), data.value)
				program := []uint8{data.opcode, data.address, 0x00}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // LSR
				s.EqualValues(data.expected, s.cpu.Memory.Fetch(uint16(data.address)))
			default:
				// Zero page addressing mode
				s.cpu.Memory.Store(uint16(data.address), data.value)
				program := []uint8{data.opcode, data.address}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // LSR
				s.EqualValues(data.expected, s.cpu.Memory.Fetch(uint16(data.address)))
			}
		})
	}
}

// endregion InstructionLSR

// region InstructionORA
type InstructionORA struct {
	baseTest
}

func (s *InstructionORA) TestORAAccumulatorCorrect() {
	subtests := []struct {
		a        uint8
		value    uint8
		expected uint8
	}{
		{0x00, 0x00, 0x00},
		{0xFF, 0xFF, 0xFF},
		{0x55, 0xAA, 0xFF},
		{0xAA, 0x55, 0xFF},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0x09, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // ORA
			s.EqualValues(data.expected, s.cpu.Registers.A)
		})
	}
}

func (s *InstructionORA) TestORAZeroFlagCorrect() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x00, 0x00, true},
		{0xFF, 0xFF, false},
		{0x00, 0x01, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0x09, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // ORA
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionORA) TestORANegativeFlagCorrect() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x7F, 0x80, true},
		{0x79, 0x00, false},
		{0xFF, 0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0x09, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // ORA
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionORA

// region InstructionPHA
type InstructionPHA struct {
	baseTest
}

func (s *InstructionPHA) TestPHAStackHasCorrectValue() {
	program := []uint8{0xA9, 0x03, 0x48}
	s.loadProgram(0x0600, program)
	stackLocation := s.cpu.Registers.SP
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.EqualValues(0x03, s.cpu.Memory.Fetch(0x0100+uint16(stackLocation)))
}

func (s *InstructionPHA) TestPHAStackPointerHasCorrectValue() {
	program := []uint8{0xA9, 0x03, 0x48}
	s.loadProgram(0x0600, program)
	stackLocation := s.cpu.Registers.SP
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.EqualValues(stackLocation-1, s.cpu.Registers.SP)
}

// endregion InstructionPHA

// region InstructionPHP
type InstructionPHP struct {
	baseTest
}

func (s *InstructionPHP) TestPHPStackSetFlagOperationsCorrectly() {
	subtests := []struct {
		opcode   uint8
		expected uint8
	}{
		{0x38, 0x31}, // SEC Carry Flag Test
		{0xF8, 0x38}, // SED Decimal Flag Test
		{0x78, 0x34}, // SEI Interrupt Flag Test
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", data.opcode), func() {
			program := []uint8{0x58, data.opcode, 0x08}
			s.loadProgram(0x0600, program)
			stackLocation := s.cpu.Registers.SP
			s.cpu.Execute() // CLI
			s.cpu.Execute() // SE*
			s.cpu.Execute() // PHP
			s.EqualValues(data.expected, s.cpu.Memory.Fetch(0x0100+uint16(stackLocation)))
		})
	}
}

func (s *InstructionPHP) TestPHPStackNonSetFlagOperationsCorrectly() {
	subtests := []struct {
		a        uint8
		value    uint8
		expected uint8
	}{
		{0x01, 0x80, 0xB0}, // Negative
		{0x01, 0x7F, 0xF0}, // Overflow + Negative
		{0x00, 0x00, 0x32}, // Zero
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0x58, 0xA9, data.a, 0x69, data.value, 0x08}
			s.loadProgram(0x0600, program)
			stackLocation := s.cpu.Registers.SP
			s.cpu.Execute() // CLI
			s.cpu.Execute() // LDA
			s.cpu.Execute() // ADC
			s.cpu.Execute() // PHP
			s.EqualValues(data.expected, s.cpu.Memory.Fetch(0x0100+uint16(stackLocation)))
		})
	}
}

func (s *InstructionPHP) TestPHPStackPointerHasCorrectValue() {
	program := []uint8{0x08}
	s.loadProgram(0x0600, program)
	stackLocation := s.cpu.Registers.SP
	s.cpu.Execute() // PHP
	s.EqualValues(stackLocation-1, s.cpu.Registers.SP)
}

// endregion InstructionPHP

// region InstructionPLA
type InstructionPLA struct {
	baseTest
}

func (s *InstructionPLA) TestPLAAccumulatorHasCorrectValue() {
	program := []uint8{0xA9, 0x03, 0x48, 0xA9, 0x00, 0x68}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PLA
	s.EqualValues(0x03, s.cpu.Registers.A)
}

func (s *InstructionPLA) TestPLAZeroFlagHasCorrectValue() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x00, true},
		{0x01, false},
		{0xFF, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value, 0x48, 0x68}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // PHA
			s.cpu.Execute() // PLA
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionPLA) TestPLANegativeFlagHasCorrectValue() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x7F, false},
		{0x80, true},
		{0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value, 0x48, 0x68}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // PHA
			s.cpu.Execute() // PLA
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionPLA

// region InstructionPLP
type InstructionPLP struct {
	baseTest
}

func (s *InstructionPLP) TestPLPCarryFlagSetCorrectly() {
	program := []uint8{0xA9, 0x01, 0x48, 0x28}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // PLP
	s.True(s.cpu.getFlag(C))
}

func (s *InstructionPLP) TestPLPZeroFlagSetCorrectly() {
	program := []uint8{0xA9, 0x02, 0x48, 0x28}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // PLP
	s.True(s.cpu.getFlag(Z))
}

func (s *InstructionPLP) TestPLPDecimalFlagSetCorrectly() {
	program := []uint8{0xA9, 0x08, 0x48, 0x28}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // PLP
	s.True(s.cpu.getFlag(D))
}

func (s *InstructionPLP) TestPLPInterruptFlagSetCorrectly() {
	program := []uint8{0xA9, 0x04, 0x48, 0x28}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // PLP
	s.True(s.cpu.getFlag(I))
}

func (s *InstructionPLP) TestPLPOverflowFlagSetCorrectly() {
	program := []uint8{0xA9, 0x40, 0x48, 0x28}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // PLP
	s.True(s.cpu.getFlag(V))
}

func (s *InstructionPLP) TestPLPNegativeFlagSetCorrectly() {
	program := []uint8{0xA9, 0x80, 0x48, 0x28}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // PLP
	s.True(s.cpu.getFlag(N))
}

// endregion InstructionPLP

// region InstructionROL
type InstructionROL struct {
	baseTest
}

func (s *InstructionROL) TestROLNegativeSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x40, true},
		{0x3F, false},
		{0x80, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value, 0x2A}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // ROL
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

func (s *InstructionROL) TestROLZeroSetCorrectly() {
	subtests := []struct {
		opcode       uint8
		expectedFlag bool
	}{
		{0x38, false},
		{0x18, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", data.opcode), func() {
			program := []uint8{0xA9, 0x00, data.opcode, 0x2A}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // SEC/CLC
			s.cpu.Execute() // ROL
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionROL) TestROLCarryFlagSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x80, true},
		{0x7F, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value, 0x2A}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // ROL
			s.Equal(data.expectedFlag, s.cpu.getFlag(C))
		})
	}
}

func (s *InstructionROL) TestROLCorrectValueStored() {
	subtests := []struct {
		opcode   uint8
		value    uint8
		expected uint8
		address  uint8
	}{
		{0x2A, 0x55, 0xAA, 0x00}, // ROL Accumulator
		{0x2A, 0x55, 0xAA, 0x00}, // ROL Accumulator
		{0x26, 0x55, 0xAA, 0x01}, // ROL Zero Page
		{0x36, 0x55, 0xAA, 0x01}, // ROL Zero Page X
		{0x2E, 0x55, 0xAA, 0x01}, // ROL Absolute
		{0x3E, 0x55, 0xAA, 0x01}, // ROL Absolute X
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X Value:%02X", data.opcode, data.value), func() {
			switch data.opcode {
			case 0x2A:
				// Accumulator mode
				program := []uint8{0xA9, data.value, data.opcode}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // LDA
				s.cpu.Execute() // ROL
				s.EqualValues(data.expected, s.cpu.Registers.A)
			case 0x2E, 0x3E:
				// Absolute addressing mode
				s.cpu.Memory.Store(uint16(data.address), data.value)
				program := []uint8{data.opcode, data.address, 0x00}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // ROL
				s.EqualValues(data.expected, s.cpu.Memory.Fetch(uint16(data.address)))
			default:
				// Zero page addressing mode
				s.cpu.Memory.Store(uint16(data.address), data.value)
				program := []uint8{data.opcode, data.address}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // ROL
				s.EqualValues(data.expected, s.cpu.Memory.Fetch(uint16(data.address)))
			}
		})
	}
}

// endregion InstructionROL

// region InstructionROR
type InstructionROR struct {
	baseTest
}

func (s *InstructionROR) TestRORNegativeSetCorrectly() {
	subtests := []struct {
		value        uint8
		opcode       uint8
		expectedFlag bool
	}{
		{0xFF, 0x18, false},
		{0xFE, 0x18, false},
		{0xFF, 0x38, true},
		{0x00, 0x38, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X Opcode:%02X", data.value, data.opcode), func() {
			program := []uint8{0xA9, data.value, 0x6A}
			s.loadProgram(0x0600, program)
			if data.opcode == 0x38 {
				s.cpu.setFlag(C, true)
			} else {
				s.cpu.setFlag(C, false)
			}
			s.cpu.Execute() // LDA
			s.cpu.Execute() // ROR
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

func (s *InstructionROR) TestRORZeroSetCorrectly() {
	subtests := []struct {
		value        uint8
		opcode       uint8
		expectedFlag bool
	}{
		{0x00, 0x18, true},
		{0x00, 0x38, false},
		{0x01, 0x18, true},
		{0x01, 0x38, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X Opcode:%02X", data.value, data.opcode), func() {
			program := []uint8{0xA9, data.value, 0x6A}
			s.loadProgram(0x0600, program)
			if data.opcode == 0x38 {
				s.cpu.setFlag(C, true)
			} else {
				s.cpu.setFlag(C, false)
			}
			s.cpu.Execute() // LDA
			s.cpu.Execute() // ROR
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionROR) TestRORCarryFlagSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x01, true},
		{0x02, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA9, data.value, 0x6A}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // ROR
			s.Equal(data.expectedFlag, s.cpu.getFlag(C))
		})
	}
}

func (s *InstructionROR) TestRORCorrectValueStored() {
	subtests := []struct {
		opcode   uint8
		value    uint8
		expected uint8
		address  uint8
	}{
		{0x6A, 0xAA, 0x55, 0x00}, // ROR Accumulator
		{0x6A, 0xAA, 0x55, 0x00}, // ROR Accumulator
		{0x66, 0xAA, 0x55, 0x01}, // ROR Zero Page
		{0x76, 0xAA, 0x55, 0x01}, // ROR Zero Page X
		{0x6E, 0xAA, 0x55, 0x01}, // ROR Absolute
		{0x7E, 0xAA, 0x55, 0x01}, // ROR Absolute X
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X Value:%02X", data.opcode, data.value), func() {
			switch data.opcode {
			case 0x6A:
				// Accumulator mode
				program := []uint8{0xA9, data.value, data.opcode}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // LDA
				s.cpu.Execute() // ROR
				s.EqualValues(data.expected, s.cpu.Registers.A)
			case 0x6E, 0x7E:
				// Absolute addressing mode
				s.cpu.Memory.Store(uint16(data.address), data.value)
				program := []uint8{data.opcode, data.address, 0x00}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // ROR
				s.EqualValues(data.expected, s.cpu.Memory.Fetch(uint16(data.address)))
			default:
				// Zero page addressing mode
				s.cpu.Memory.Store(uint16(data.address), data.value)
				program := []uint8{data.opcode, data.address}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // ROR
				s.EqualValues(data.expected, s.cpu.Memory.Fetch(uint16(data.address)))
			}
		})
	}
}

// endregion InstructionROR

// region InstructionRTI
type InstructionRTI struct {
	baseTest
}

func (s *InstructionRTI) TestRTIProgramCounterCorrect() {
	// Set PC to 0xABCD and store BRK there
	s.cpu.Registers.PC = 0xABCD
	s.cpu.Memory.Store(0xABCD, 0x00) // BRK
	// Store RTI at 0x0000 (where IRQ vector points by default)
	s.cpu.Memory.Store(0x0000, 0x40) // RTI
	// Set IRQ vector to point to 0x0000
	s.cpu.Memory.Store(0xFFFE, 0x00)
	s.cpu.Memory.Store(0xFFFF, 0x00)
	s.EqualValues(0xABCD, s.cpu.Registers.PC)
	s.cpu.Execute() // BRK - pushes PC+2 (0xABCF), jumps to IRQ vector
	s.cpu.Execute() // RTI
	s.EqualValues(0xABCF, s.cpu.Registers.PC)
}

func (s *InstructionRTI) TestRTICarryFlagSetCorrectly() {
	program := []uint8{0xA9, uint8(C), 0x48, 0x40}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // RTI
	s.True(s.cpu.getFlag(C))
}

func (s *InstructionRTI) TestRTIZeroFlagSetCorrectly() {
	program := []uint8{0xA9, uint8(Z), 0x48, 0x40}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // RTI
	s.True(s.cpu.getFlag(Z))
}

func (s *InstructionRTI) TestRTIDecimalFlagSetCorrectly() {
	program := []uint8{0xA9, uint8(D), 0x48, 0x40}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // RTI
	s.True(s.cpu.getFlag(D))
}

func (s *InstructionRTI) TestRTIInterruptFlagSetCorrectly() {
	program := []uint8{0xA9, uint8(I), 0x48, 0x40}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // RTI
	s.True(s.cpu.getFlag(I))
}

func (s *InstructionRTI) TestRTIOverflowFlagSetCorrectly() {
	program := []uint8{0xA9, uint8(V), 0x48, 0x40}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // RTI
	s.True(s.cpu.getFlag(V))
}

func (s *InstructionRTI) TestRTINegativeFlagSetCorrectly() {
	program := []uint8{0xA9, uint8(N), 0x48, 0x40}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // PHA
	s.cpu.Execute() // RTI
	s.True(s.cpu.getFlag(N))
}

// endregion InstructionRTI

// region InstructionRTS
type InstructionRTS struct {
	baseTest
}

func (s *InstructionRTS) TestRTSProgramCounterHasCorrectValue() {
	// Load at 0x0000 so JSR $0004 points to the RTS at offset 4
	program := []uint8{0x20, 0x04, 0x00, 0x00, 0x60}
	s.loadProgram(0x0000, program)
	s.cpu.Execute() // JSR $0004
	s.cpu.Execute() // RTS
	s.EqualValues(0x0003, s.cpu.Registers.PC)
}

func (s *InstructionRTS) TestRTSStackPointerHasCorrectValue() {
	s.cpu.Memory.Store(0xBBAA, 0x60)
	s.cpu.Registers.PC = 0xBBAA
	stackLocation := s.cpu.Registers.SP
	s.cpu.Execute() // RTS
	s.EqualValues((stackLocation+2)&0xFF, s.cpu.Registers.SP)
}

// endregion InstructionRTS

// region InstructionSBC
type InstructionSBC struct {
	baseTest
}

func (s *InstructionSBC) TestSBCAccumulatorCorrectWhenNotInBCDMode() {
	subtests := []struct {
		a        uint8
		value    uint8
		carryIn  bool
		expected uint8
	}{
		{0x00, 0x00, false, 0xFF},
		{0x00, 0x00, true, 0x00},
		{0x50, 0xF0, false, 0x5F},
		{0x50, 0xB0, true, 0xA0},
		{0xFF, 0xFF, false, 0xFF},
		{0xFF, 0xFF, true, 0x00},
		{0xFF, 0x80, false, 0x7E},
		{0xFF, 0x80, true, 0x7F},
		{0x80, 0xFF, false, 0x80},
		{0x80, 0xFF, true, 0x81},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X CarryIn:%v", data.a, data.value, data.carryIn), func() {
			program := []uint8{0xA9, data.a, 0xE9, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.setFlag(C, data.carryIn)
			s.cpu.Execute() // SBC
			s.EqualValues(data.expected, s.cpu.Registers.A)
		})
	}
}

func (s *InstructionSBC) TestSBCAccumulatorCorrectWhenInBCDMode() {
	subtests := []struct {
		a        uint8
		value    uint8
		carryIn  bool
		expected uint8
	}{
		{0x00, 0x99, false, 0x00},
		{0x00, 0x99, true, 0x01},
		{0x45, 0x12, false, 0x32},
		{0x99, 0x99, false, 0x99},
		{0x99, 0x99, true, 0x00},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X CarryIn:%v", data.a, data.value, data.carryIn), func() {
			program := []uint8{0xF8, 0xA9, data.a, 0xE9, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // SED
			s.cpu.Execute() // LDA
			s.cpu.setFlag(C, data.carryIn)
			s.cpu.Execute() // SBC
			s.EqualValues(data.expected, s.cpu.Registers.A)
		})
	}
}

func (s *InstructionSBC) TestSBCOverflowCorrectWhenNotInBCDMode() {
	subtests := []struct {
		a            uint8
		value        uint8
		carryIn      bool
		expectedFlag bool
	}{
		{0xFF, 0x01, false, false},
		{0xFF, 0x00, false, false},
		{0x80, 0x00, false, true},
		{0x80, 0x00, true, false},
		{0x81, 0x01, false, true},
		{0x81, 0x01, true, false},
		{0x00, 0x80, false, false},
		{0x00, 0x80, true, true},
		{0x01, 0x80, true, true},
		{0x01, 0x7F, false, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X CarryIn:%v", data.a, data.value, data.carryIn), func() {
			program := []uint8{0xA9, data.a, 0xE9, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.setFlag(C, data.carryIn)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // SBC
			s.Equal(data.expectedFlag, s.cpu.getFlag(V))
		})
	}
}

func (s *InstructionSBC) TestSBCCarryCorrect() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x00, 0x00, false},
		{0x00, 0x01, false},
		{0x01, 0x00, true},
		{0x02, 0x01, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0xE9, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.Execute() // SBC
			s.Equal(data.expectedFlag, s.cpu.getFlag(C))
		})
	}
}

func (s *InstructionSBC) TestSBCZeroCorrect() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x00, 0x00, false},
		{0x00, 0x01, false},
		{0x01, 0x00, true},
		{0x01, 0x01, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0xE9, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.setFlag(C, false)
			s.cpu.Execute() // SBC
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

func (s *InstructionSBC) TestSBCNegativeCorrect() {
	subtests := []struct {
		a            uint8
		value        uint8
		expectedFlag bool
	}{
		{0x80, 0x01, false},
		{0x81, 0x01, false},
		{0x00, 0x01, true},
		{0x01, 0x01, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("A:%02X Value:%02X", data.a, data.value), func() {
			program := []uint8{0xA9, data.a, 0xE9, data.value}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDA
			s.cpu.setFlag(C, false)
			s.cpu.Execute() // SBC
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

// endregion InstructionSBC

// region InstructionSEC
type InstructionSEC struct {
	baseTest
}

func (s *InstructionSEC) TestSECCarryFlagSetCorrectly() {
	program := []uint8{0x38}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // SEC
	s.True(s.cpu.getFlag(C))
}

// endregion InstructionSEC

// region InstructionSED
type InstructionSED struct {
	baseTest
}

func (s *InstructionSED) TestSEDDecimalModeSetCorrectly() {
	program := []uint8{0xF8}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // SED
	s.True(s.cpu.getFlag(D))
}

type InstructionSEI struct {
	baseTest
}

func (s *InstructionSEI) TestSEIInterruptFlagSetCorrectly() {
	program := []uint8{0x78}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // SEI
	s.True(s.cpu.getFlag(I))
}

type InstructionSTA struct {
	baseTest
}

func (s *InstructionSTA) TestSTAMemoryHasCorrectValue() {
	program := []uint8{0xA9, 0x03, 0x85, 0x05}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDA
	s.cpu.Execute() // STA
	s.EqualValues(0x03, s.cpu.Memory.Fetch(0x05))
}

// endregion InstructionSTA

// region InstructionSTX
type InstructionSTX struct {
	baseTest
}

func (s *InstructionSTX) TestSTXMemoryHasCorrectValue() {
	program := []uint8{0xA2, 0x03, 0x86, 0x05}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDX
	s.cpu.Execute() // STX
	s.EqualValues(0x03, s.cpu.Memory.Fetch(0x05))
}

// endregion InstructionSTX

// region InstructionSTY
type InstructionSTY struct {
	baseTest
}

func (s *InstructionSTY) TestSTYMemoryHasCorrectValue() {
	program := []uint8{0xA0, 0x03, 0x84, 0x05}
	s.loadProgram(0x0600, program)
	s.cpu.Execute() // LDY
	s.cpu.Execute() // STY
	s.EqualValues(0x03, s.cpu.Memory.Fetch(0x05))
}

// endregion InstructionSTY

// region InstructionTAX
type InstructionTAX struct {
	baseTest
}

func (s *InstructionTAX) TestTAXTransferCorrectValueSet() {
	subtests := []struct {
		opcode uint8
		value  uint8
		from   string
		to     string
	}{
		{0xAA, 0x03, "A", "X"},
		{0xA8, 0x03, "A", "Y"},
		{0x8A, 0x03, "X", "A"},
		{0x98, 0x03, "Y", "A"},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("%s->%s Value:%02X", data.from, data.to, data.value), func() {
			var program []uint8
			switch data.from {
			case "A":
				program = []uint8{0xA9, data.value, data.opcode}
			case "X":
				program = []uint8{0xA2, data.value, data.opcode}
			case "Y":
				program = []uint8{0xA0, data.value, data.opcode}
			}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LD*
			s.cpu.Execute() // T**
			switch data.to {
			case "A":
				s.EqualValues(data.value, s.cpu.Registers.A)
			case "X":
				s.EqualValues(data.value, s.cpu.Registers.X)
			case "Y":
				s.EqualValues(data.value, s.cpu.Registers.Y)
			}
		})
	}
}

func (s *InstructionTAX) TestTAXTransferNegativeValueSet() {
	subtests := []struct {
		opcode       uint8
		value        uint8
		from         string
		expectedFlag bool
	}{
		{0xAA, 0x80, "A", true},
		{0xA8, 0x80, "A", true},
		{0x8A, 0x80, "X", true},
		{0x98, 0x80, "Y", true},
		{0xAA, 0xFF, "A", true},
		{0xA8, 0xFF, "A", true},
		{0x8A, 0xFF, "X", true},
		{0x98, 0xFF, "Y", true},
		{0xAA, 0x7F, "A", false},
		{0xA8, 0x7F, "A", false},
		{0x8A, 0x7F, "X", false},
		{0x98, 0x7F, "Y", false},
		{0xAA, 0x00, "A", false},
		{0xA8, 0x00, "A", false},
		{0x8A, 0x00, "X", false},
		{0x98, 0x00, "Y", false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X Value:%02X", data.opcode, data.value), func() {
			var program []uint8
			switch data.from {
			case "A":
				program = []uint8{0xA9, data.value, data.opcode}
			case "X":
				program = []uint8{0xA2, data.value, data.opcode}
			case "Y":
				program = []uint8{0xA0, data.value, data.opcode}
			}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LD*
			s.cpu.Execute() // T**
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

func (s *InstructionTAX) TestTAXTransferZeroValueSet() {
	subtests := []struct {
		opcode       uint8
		value        uint8
		from         string
		expectedFlag bool
	}{
		{0xAA, 0xFF, "A", false},
		{0xA8, 0xFF, "A", false},
		{0x8A, 0xFF, "X", false},
		{0x98, 0xFF, "Y", false},
		{0xAA, 0x00, "A", true},
		{0xA8, 0x00, "A", true},
		{0x8A, 0x00, "X", true},
		{0x98, 0x00, "Y", true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X Value:%02X", data.opcode, data.value), func() {
			var program []uint8
			switch data.from {
			case "A":
				program = []uint8{0xA9, data.value, data.opcode}
			case "X":
				program = []uint8{0xA2, data.value, data.opcode}
			case "Y":
				program = []uint8{0xA0, data.value, data.opcode}
			}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LD*
			s.cpu.Execute() // T**
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

// endregion InstructionTAX

// region InstructionTSX
type InstructionTSX struct {
	baseTest
}

func (s *InstructionTSX) TestTSXXRegisterSetCorrectly() {
	program := []uint8{0xBA}
	s.loadProgram(0x0600, program)
	stackPointer := s.cpu.Registers.SP
	s.cpu.Execute() // TSX
	s.EqualValues(stackPointer, s.cpu.Registers.X)
}

func (s *InstructionTSX) TestTSXNegativeSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x00, false},
		{0x7F, false},
		{0x80, true},
		{0xFF, true},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA2, data.value, 0x9A, 0xBA}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // TXS
			s.cpu.Execute() // TSX
			s.Equal(data.expectedFlag, s.cpu.getFlag(N))
		})
	}
}

func (s *InstructionTSX) TestTSXZeroSetCorrectly() {
	subtests := []struct {
		value        uint8
		expectedFlag bool
	}{
		{0x00, true},
		{0x01, false},
		{0xFF, false},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data.value), func() {
			program := []uint8{0xA2, data.value, 0x9A, 0xBA}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // TXS
			s.cpu.Execute() // TSX
			s.Equal(data.expectedFlag, s.cpu.getFlag(Z))
		})
	}
}

// endregion InstructionTSX

// region InstructionTXS
type InstructionTXS struct {
	baseTest
}

func (s *InstructionTXS) TestTXSStackPointerSetCorrectly() {
	subtests := []uint8{0xAA, 0x00, 0xFF}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Value:%02X", data), func() {
			program := []uint8{0xA2, data, 0x9A}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // TXS
			s.EqualValues(data, s.cpu.Registers.SP)
		})
	}
}

// endregion InstructionTXS

// region AccumulatorAddress
type AccumulatorAddress struct {
	baseTest
}

func (s *AccumulatorAddress) TestAddressCorrectForImmediateMode() {
	subtests := []uint8{0x69, 0x09, 0x29, 0x49, 0xC9, 0xE0, 0xC0, 0xA2, 0xA0, 0xA9, 0xE9}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03}
			s.loadProgram(0x0600, program)
			s.EqualValues(0x0600, s.cpu.Registers.PC)
			s.cpu.Execute()
			s.EqualValues(0x0602, s.cpu.Registers.PC)
		})
	}
}

func (s *AccumulatorAddress) TestAddressCorrectForZeroPageMode() {
	subtests := []uint8{0x65, 0x05, 0x25, 0x45, 0xC5, 0xE4, 0xC4, 0xA6, 0xA4, 0xA5, 0xE5, 0x85, 0x86, 0x84}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03}
			s.loadProgram(0x0600, program)
			s.EqualValues(0x0600, s.cpu.Registers.PC)
			s.cpu.Execute()
			s.EqualValues(0x0602, s.cpu.Registers.PC)
		})
	}
}

func (s *AccumulatorAddress) TestAddressCorrectForZeroPageXMode() {
	subtests := []uint8{0x75, 0x15, 0x35, 0x55, 0xD5, 0xB5, 0xF5, 0x95}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03}
			s.loadProgram(0x0600, program)
			s.EqualValues(0x0600, s.cpu.Registers.PC)
			s.cpu.Execute()
			s.EqualValues(0x0602, s.cpu.Registers.PC)
		})
	}
}

func (s *AccumulatorAddress) TestAddressCorrectForAbsoluteMode() {
	subtests := []uint8{0x6D, 0x0D, 0x2D, 0x4D, 0xCD, 0xEC, 0xCC, 0xAE, 0xAC, 0xAD, 0xED, 0x8D, 0x8E, 0x8C}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03, 0x05}
			s.loadProgram(0x0600, program)
			s.EqualValues(0x0600, s.cpu.Registers.PC)
			s.cpu.Execute()
			s.EqualValues(0x0603, s.cpu.Registers.PC)
		})
	}
}

func (s *AccumulatorAddress) TestAddressCorrectForAbsoluteXMode() {
	subtests := []uint8{0x7D, 0x1D, 0x3D, 0x5D, 0xDD, 0xBD, 0xFD, 0x9D}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03, 0x05}
			s.loadProgram(0x0600, program)
			s.EqualValues(0x0600, s.cpu.Registers.PC)
			s.cpu.Execute()
			s.EqualValues(0x0603, s.cpu.Registers.PC)
		})
	}
}

func (s *AccumulatorAddress) TestAddressCorrectForAbsoluteYMode() {
	subtests := []uint8{0x79, 0x19, 0x39, 0x59, 0xD9, 0xBE, 0xB9, 0xF9, 0x99}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03, 0x05}
			s.loadProgram(0x0600, program)
			s.EqualValues(0x0600, s.cpu.Registers.PC)
			s.cpu.Execute()
			s.EqualValues(0x0603, s.cpu.Registers.PC)
		})
	}
}

func (s *AccumulatorAddress) TestAddressCorrectForIndexedIndirectMode() {
	subtests := []uint8{0x61, 0x01, 0x21, 0x41, 0xC1, 0xA1, 0xE1, 0x81}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03}
			s.loadProgram(0x0600, program)
			s.EqualValues(0x0600, s.cpu.Registers.PC)
			s.cpu.Execute()
			s.EqualValues(0x0602, s.cpu.Registers.PC)
		})
	}
}

func (s *AccumulatorAddress) TestAddressCorrectForIndirectIndexedMode() {
	subtests := []uint8{0x71, 0x11, 0x31, 0x51, 0xD1, 0xB1, 0xF1, 0x91}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03}
			s.loadProgram(0x0600, program)
			s.EqualValues(0x0600, s.cpu.Registers.PC)
			s.cpu.Execute()
			s.EqualValues(0x0602, s.cpu.Registers.PC)
		})
	}
}

// endregion AccumulatorAddress

// region IndexAddress
type IndexAddress struct {
	baseTest
}

func (s *IndexAddress) TestZeroPageIndexAddressingWrapsCorrectly() {
	subtests := []struct {
		opcode   uint8
		index    uint8
		address  uint8
		expected uint8
	}{
		{0xB5, 0x01, 0xFF, 0x00},
		{0xB5, 0x03, 0xFF, 0x02},
		{0x95, 0x01, 0xFF, 0x00},
		{0x95, 0x03, 0xFF, 0x02},
		{0xB6, 0x01, 0xFF, 0x00},
		{0xB6, 0x03, 0xFF, 0x02},
		{0x96, 0x01, 0xFF, 0x00},
		{0x96, 0x03, 0xFF, 0x02},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X Index:%02X Addr:%02X", data.opcode, data.index, data.address), func() {
			s.cpu.Memory.Store(uint16(data.expected), 0x03)
			switch data.opcode {
			case 0xB5:
				// LDA $FF,X - load A from memory using X index
				program := []uint8{0xA2, data.index, data.opcode, data.address}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // LDX
				s.cpu.Execute() // LDA
				s.EqualValues(0x03, s.cpu.Registers.A)
			case 0xB6:
				// LDX $FF,Y - load X from memory using Y index
				program := []uint8{0xA0, data.index, data.opcode, data.address}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // LDY
				s.cpu.Execute() // LDX
				s.EqualValues(0x03, s.cpu.Registers.X)
			case 0x95:
				// STA $FF,X - store A to memory using X index
				program := []uint8{0xA9, 0x03, 0xA2, data.index, data.opcode, data.address}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // LDA #03
				s.cpu.Execute() // LDX
				s.cpu.Execute() // STA
				s.EqualValues(0x03, s.cpu.Memory.Fetch(uint16(data.expected)))
			default: // 0x96 - STX $FF,Y
				// STX $FF,Y - store X to memory using Y index
				program := []uint8{0xA2, 0x03, 0xA0, data.index, data.opcode, data.address}
				s.loadProgram(0x0600, program)
				s.cpu.Execute() // LDX #03
				s.cpu.Execute() // LDY
				s.cpu.Execute() // STX
				s.EqualValues(0x03, s.cpu.Memory.Fetch(uint16(data.expected)))
			}
		})
	}
}

func (s *IndexAddress) TestZeroPageAddressingCorrect() {
	subtests := []struct {
		opcode  uint8
		index   uint8
		address uint8
	}{
		{0xB5, 0x01, 0x03},
		{0x95, 0x01, 0x03},
		{0xB6, 0x01, 0x03},
		{0x96, 0x01, 0x03},
	}
	for _, data := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X Index:%02X", data.opcode, data.index), func() {
			program := []uint8{0xA2, data.index, data.opcode, data.address}
			s.loadProgram(0x0600, program)
			s.cpu.Execute() // LDX
			s.cpu.Execute() // Opcode
			s.EqualValues(0x0604, s.cpu.Registers.PC)
		})
	}
}

func (s *IndexAddress) TestAbsoluteAddressingCorrect() {
	subtests := []uint8{0xBD, 0x9D, 0xBE, 0xB9, 0x99}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03, 0x05}
			s.loadProgram(0x0600, program)
			s.cpu.Execute()
			s.EqualValues(0x0603, s.cpu.Registers.PC)
		})
	}
}

// endregion IndexAddress

// region CompareAddress
type CompareAddress struct {
	baseTest
}

func (s *CompareAddress) TestImmediateCorrect() {
	subtests := []uint8{0xC9, 0xE0, 0xC0}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03}
			s.loadProgram(0x0600, program)
			s.cpu.Execute()
			s.EqualValues(0x0602, s.cpu.Registers.PC)
		})
	}
}

func (s *CompareAddress) TestZeroPageCorrect() {
	subtests := []uint8{0xC5, 0xE4, 0xC4}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03}
			s.loadProgram(0x0600, program)
			s.cpu.Execute()
			s.EqualValues(0x0602, s.cpu.Registers.PC)
		})
	}
}

func (s *CompareAddress) TestCompareAddressAbsoluteCorrect() {
	subtests := []uint8{0xCD, 0xEC, 0xCC}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03, 0x05}
			s.loadProgram(0x0600, program)
			s.cpu.Execute()
			s.EqualValues(0x0603, s.cpu.Registers.PC)
		})
	}
}

func (s *CompareAddress) TestIndexedAddressesCorrect() {
	subtests := []uint8{0xD5, 0xDD, 0xD9, 0xC1, 0xD1}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			var program []uint8
			if opcode == 0xDD || opcode == 0xD9 {
				program = []uint8{opcode, 0x03, 0x05}
			} else {
				program = []uint8{opcode, 0x03}
			}
			s.loadProgram(0x0600, program)
			s.cpu.Execute()
			if opcode == 0xDD || opcode == 0xD9 {
				s.EqualValues(0x0603, s.cpu.Registers.PC)
			} else {
				s.EqualValues(0x0602, s.cpu.Registers.PC)
			}
		})
	}
}

// endregion CompareAddress

// region DecrementIncrementAddress
type DecrementIncrementAddress struct {
	baseTest
}

func (s *DecrementIncrementAddress) TestZeroPageCorrect() {
	subtests := []uint8{0xE6, 0xC6}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03}
			s.loadProgram(0x0600, program)
			s.cpu.Execute()
			s.EqualValues(0x0602, s.cpu.Registers.PC)
		})
	}
}

func (s *DecrementIncrementAddress) TestDecrementIncrementAddressAbsoluteCorrect() {
	subtests := []uint8{0xEE, 0xCE}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03, 0x05}
			s.loadProgram(0x0600, program)
			s.cpu.Execute()
			s.EqualValues(0x0603, s.cpu.Registers.PC)
		})
	}
}

// endregion DecrementIncrementAddress

// region StoreInMemoryAddress
type StoreInMemoryAddress struct {
	baseTest
}

func (s *StoreInMemoryAddress) TestZeroPageAddressingCorrect() {
	subtests := []uint8{0x85, 0x86, 0x84, 0x95, 0x96, 0x94}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03}
			s.loadProgram(0x0600, program)
			s.cpu.Execute()
			s.EqualValues(0x0602, s.cpu.Registers.PC)
		})
	}
}

func (s *StoreInMemoryAddress) TestAbsoluteAddressingCorrect() {
	subtests := []uint8{0x8D, 0x8E, 0x8C, 0x9D, 0x99}
	for _, opcode := range subtests {
		s.Run(fmt.Sprintf("Opcode:%02X", opcode), func() {
			program := []uint8{opcode, 0x03, 0x05}
			s.loadProgram(0x0600, program)
			s.cpu.Execute()
			s.EqualValues(0x0603, s.cpu.Registers.PC)
		})
	}
}

// endregion StoreInMemoryAddress

// region ClockCyclesTest
type ClockCyclesTest struct {
	baseTest
}

// Test number of cycles remaining is correct after operations and the
// result do not wrap
func (s *ClockCyclesTest) TestCyclesCorrectAfterOperationsThatDoNotWrap() {

	subtests := []struct {
		op       OpCode
		expected uint16
	}{
		{0x69, 2}, // ADC Immediate
		{0x65, 3}, // ADC Zero Page
		{0x75, 4}, // ADC Zero Page X
		{0x6D, 4}, // ADC Absolute
		{0x7D, 4}, // ADC Absolute X
		{0x79, 4}, // ADC Absolute Y
		{0x61, 6}, // ADC Indrect X
		{0x71, 5}, // ADC Indirect Y
		{0x29, 2}, // AND Immediate
		{0x25, 3}, // AND Zero Page
		{0x35, 4}, // AND Zero Page X
		{0x2D, 4}, // AND Absolute
		{0x3D, 4}, // AND Absolute X
		{0x39, 4}, // AND Absolute Y
		{0x21, 6}, // AND Indirect X
		{0x31, 5}, // AND Indirect Y
		{0x0A, 2}, // ASL Accumulator
		{0x06, 5}, // ASL Zero Page
		{0x16, 6}, // ASL Zero Page X
		{0x0E, 6}, // ASL Absolute
		{0x1E, 7}, // ASL Absolute X
		{0x24, 3}, // BIT Zero Page
		{0x2C, 4}, // BIT Absolute
		{0x00, 7}, // BRK Implied
		{0x18, 2}, // CLC Implied
		{0xD8, 2}, // CLD Implied
		{0x58, 2}, // CLI Implied
		{0xB8, 2}, // CLV Implied
		{0xC9, 2}, // CMP Immediate
		{0xC5, 3}, // CMP ZeroPage
		{0xD5, 4}, // CMP Zero Page X
		{0xCD, 4}, // CMP Absolute
		{0xDD, 4}, // CMP Absolute X
		{0xD9, 4}, // CMP Absolute Y
		{0xC1, 6}, // CMP Indirect X
		{0xD1, 5}, // CMP Indirect Y
		{0xE0, 2}, // CPX Immediate
		{0xE4, 3}, // CPX ZeroPage
		{0xEC, 4}, // CPX Absolute
		{0xC0, 2}, // CPY Immediate
		{0xC4, 3}, // CPY ZeroPage
		{0xCC, 4}, // CPY Absolute
		{0xC6, 5}, // DEC Zero Page
		{0xD6, 6}, // DEC Zero Page X
		{0xCE, 6}, // DEC Absolute
		{0xDE, 7}, // DEC Absolute X
		{0xCA, 2}, // DEX Implied
		{0x88, 2}, // DEY Implied
		{0x49, 2}, // EOR Immediate
		{0x45, 3}, // EOR Zero Page
		{0x55, 4}, // EOR Zero Page X
		{0x4D, 4}, // EOR Absolute
		{0x5D, 4}, // EOR Absolute X
		{0x59, 4}, // EOR Absolute Y
		{0x41, 6}, // EOR Indrect X
		{0x51, 5}, // EOR Indirect Y
		{0xE6, 5}, // INC Zero Page
		{0xF6, 6}, // INC Zero Page X
		{0xEE, 6}, // INC Absolute
		{0xFE, 7}, // INC Absolute X
		{0xE8, 2}, // INX Implied
		{0xC8, 2}, // INY Implied
		{0x4C, 3}, // JMP Absolute
		{0x6C, 5}, // JMP Indirect
		{0x20, 6}, // JSR Absolute
		{0xA9, 2}, // LDA Immediate
		{0xA5, 3}, // LDA Zero Page
		{0xB5, 4}, // LDA Zero Page X
		{0xAD, 4}, // LDA Absolute
		{0xBD, 4}, // LDA Absolute X
		{0xB9, 4}, // LDA Absolute Y
		{0xA1, 6}, // LDA Indirect X
		{0xB1, 5}, // LDA Indirect Y
		{0xA2, 2}, // LDX Immediate
		{0xA6, 3}, // LDX Zero Page
		{0xB6, 4}, // LDX Zero Page Y
		{0xAE, 4}, // LDX Absolute
		{0xBE, 4}, // LDX Absolute Y
		{0xA0, 2}, // LDY Immediate
		{0xA4, 3}, // LDY Zero Page
		{0xB4, 4}, // LDY Zero Page Y
		{0xAC, 4}, // LDY Absolute
		{0xBC, 4}, // LDY Absolute Y
		{0x4A, 2}, // LSR Accumulator
		{0x46, 5}, // LSR Zero Page
		{0x56, 6}, // LSR Zero Page X
		{0x4E, 6}, // LSR Absolute
		{0x5E, 7}, // LSR Absolute X
		{0xEA, 2}, // NOP Implied
		{0x09, 2}, // ORA Immediate
		{0x05, 3}, // ORA Zero Page
		{0x15, 4}, // ORA Zero Page X
		{0x0D, 4}, // ORA Absolute
		{0x1D, 4}, // ORA Absolute X
		{0x19, 4}, // ORA Absolute Y
		{0x01, 6}, // ORA Indirect X
		{0x11, 5}, // ORA Indirect Y
		{0x48, 3}, // PHA Implied
		{0x08, 3}, // PHP Implied
		{0x68, 4}, // PLA Implied
		{0x28, 4}, // PLP Implied
		{0x2A, 2}, // ROL Accumulator
		{0x26, 5}, // ROL Zero Page
		{0x36, 6}, // ROL Zero Page X
		{0x2E, 6}, // ROL Absolute
		{0x3E, 7}, // ROL Absolute X
		{0x6A, 2}, // ROR Accumulator
		{0x66, 5}, // ROR Zero Page
		{0x76, 6}, // ROR Zero Page X
		{0x6E, 6}, // ROR Absolute
		{0x7E, 7}, // ROR Absolute X
		{0x40, 6}, // RTI Implied
		{0x60, 6}, // RTS Implied
		{0xE9, 2}, // SBC Immediate
		{0xE5, 3}, // SBC Zero Page
		{0xF5, 4}, // SBC Zero Page X
		{0xED, 4}, // SBC Absolute
		{0xFD, 4}, // SBC Absolute X
		{0xF9, 4}, // SBC Absolute Y
		{0xE1, 6}, // SBC Indrect X
		{0xF1, 5}, // SBC Indirect Y
		{0x38, 2}, // SEC Implied
		{0xF8, 2}, // SED Implied
		{0x78, 2}, // SEI Implied
		{0x85, 3}, // STA ZeroPage
		{0x95, 4}, // STA Zero Page X
		{0x8D, 4}, // STA Absolute
		{0x9D, 5}, // STA Absolute X
		{0x99, 5}, // STA Absolute Y
		{0x81, 6}, // STA Indirect X
		{0x91, 6}, // STA Indirect Y
		{0x86, 3}, // STX Zero Page
		{0x96, 4}, // STX Zero Page Y
		{0x8E, 4}, // STX Absolute
		{0x84, 3}, // STY Zero Page
		{0x94, 4}, // STY Zero Page X
		{0x8C, 4}, // STY Absolute
		{0xAA, 2}, // TAX Implied
		{0xA8, 2}, // TAY Implied
		{0xBA, 2}, // TSX Implied
		{0x8A, 2}, // TXA Implied
		{0x9A, 2}, // TXS Implied
		{0x98, 2}, // TYA Implied
		// Illegal Instructions
		{0x0B, 2}, // AAC Immediate
		{0x87, 3}, // AAX Zero Page
		{0x97, 4}, // AAX Zero Page Y
		{0x8F, 4}, // AAX Absolute
		{0x83, 6}, // AAX Indirect X
		{0x6B, 2}, // ARR Immediate
		{0x4B, 2}, // ASR Immediate
		{0xAB, 2}, // ATX Immediate
		{0x9F, 5}, // AXA Absolute Y
		{0x93, 6}, // AXA Indirect Y
		{0xCB, 2}, // AXS Immediate
		{0xC7, 5}, // DCP Zero Page
		{0xD7, 6}, // DCP Zero Page X
		{0xCF, 6}, // DCP Absolute
		{0xDF, 7}, // DCP Absolute X
		{0xDB, 7}, // DCP Absolute Y
		{0xC3, 8}, // DCP Indirect X
		{0xD3, 8}, // DCP Indirect Y
		{0xE7, 5}, // ISC Zero Page
		{0xF7, 6}, // ISC Zero Page X
		{0xEF, 6}, // ISC Absolute
		{0xFF, 7}, // ISC Absolute X
		{0xFB, 7}, // ISC Absolute Y
		{0xE3, 8}, // ISC Indirect X
		{0xF3, 8}, // ISC Indirect Y
		// [0x, ], # KIL Immediate (Untestable)
		{0xBB, 4}, // LAR Absolute Y
		{0xA7, 3}, // LAX Zero Page
		{0xB7, 4}, // LAX Zero Page Y
		{0xAF, 4}, // LAX Absolute
		{0xBF, 4}, // LAX Absolute Y
		{0xA3, 6}, // LAX Indirect X
		{0xB3, 5}, // LAX Indirect Y
		{0x27, 5}, // RLA Zero Page
		{0x37, 6}, // RLA Zero Page X
		{0x2F, 6}, // RLA Absolute
		{0x3F, 7}, // RLA Absolute X
		{0x3B, 7}, // RLA Absolute Y
		{0x23, 8}, // RLA Indirect X
		{0x33, 8}, // RLA Indirect Y
		{0x67, 5}, // RRA Zero Page
		{0x77, 6}, // RRA Zero Page X
		{0x6F, 6}, // RRA Absolute
		{0x7F, 7}, // RRA Absolute X
		{0x7B, 7}, // RRA Absolute Y
		{0x63, 8}, // RRA Indirect X
		{0x73, 8}, // RRA Indirect Y
		{0xEB, 2}, // USBC Immediate (SBC)
		{0x07, 5}, // SLO Zero Page
		{0x17, 6}, // SLO Zero Page X
		{0x0F, 6}, // SLO Absolute
		{0x1F, 7}, // SLO Absolute X
		{0x1B, 7}, // SLO Absolute Y
		{0x03, 8}, // SLO Indirect X
		{0x13, 8}, // SLO Indirect Y
		{0x47, 5}, // SRE Zero Page
		{0x57, 6}, // SRE Zero Page X
		{0x4F, 6}, // SRE Absolute
		{0x5F, 7}, // SRE Absolute X
		{0x5B, 7}, // SRE Absolute Y
		{0x43, 8}, // SRE Indirect X
		{0x53, 8}, // SRE Indirect Y
		{0x9E, 5}, // SXA Absolute Y
		{0x9C, 5}, // SYA Absolute X
		{0x8B, 2}, // XAA Immediate
		{0x9B, 5}, // XAS Absolute Y
		{0x1A, 2}, // NOP Implied
		{0x3A, 2}, // NOP Implied
		{0x5A, 2}, // NOP Implied
		{0x7A, 2}, // NOP Implied
		{0xDA, 2}, // NOP Implied
		{0xFA, 2}, // NOP Implied
		{0x80, 2}, // NOP Immediate   (SKB)
		{0x82, 2}, // NOP Immediate   (SKB)
		{0x89, 2}, // NOP Immediate   (SKB)
		{0xC2, 2}, // NOP Immediate   (SKB)
		{0xE2, 2}, // NOP Immediate   (SKB)
		{0x04, 3}, // NOP Zero Page   (IGN)
		{0x44, 3}, // NOP Zero Page   (IGN)
		{0x64, 3}, // NOP Zero Page   (IGN)
		{0x14, 4}, // NOP Zero Page X (IGN)
		{0x34, 4}, // NOP Zero Page X (IGN)
		{0x54, 4}, // NOP Zero Page X (IGN)
		{0x74, 4}, // NOP Zero Page X (IGN)
		{0xD4, 4}, // NOP Zero Page X (IGN)
		{0xF4, 4}, // NOP Zero Page X (IGN)
		{0x0C, 4}, // NOP Absolute    (IGN)
		{0x1C, 4}, // NOP Absolute X  (IGN)
		{0x3C, 4}, // NOP Absolute X  (IGN)
		{0x5C, 4}, // NOP Absolute X  (IGN)
		{0x7C, 4}, // NOP Absolute X  (IGN)
		{0xDC, 4}, // NOP Absolute X  (IGN)
		{0xFC, 4}, // NOP Absolute X  (IGN)
	}

	r := s.Require()
	flags := s.cpu.log.Flags()
	s.cpu.log.SetFlags(log.LstdFlags | log.Lshortfile)
	s.cpu.breakError = false

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		s._instruction = data
		s.cpu.Reset()
		//s.cpu.debug |= DebugModeCycleDump
		s.Run(fmt.Sprintf("%s/%s", data.Mnemonic, data.AddressingMode), func() {
			s.cpu.Registers.PC = 0x0000
			s.cpu.Memory.Store(0x0000, uint8(data.OpCode))
			//startCycles := s.cpu.currentCycles

			cycles, err := s.cpu.Execute()
			r.NoError(err)
			s.Log(s.cpu.decode.String())
			//s.Log(fmt.Sprintf("Starting cycles: %d\ncurrent cycles: %d", startCycles, s.cpu.currentCycles))
			//actual := cycles - startCycles

			//s.Equal(data.Cycles, actual)
			s.Equal(data.Cycles, cycles)
		})
		//s.cpu.debug &= ^DebugModeCycleDump
	}
	s.cpu.breakError = true
	s.cpu.log.SetFlags(flags)
}

// Test number of cycles remaining is correct when in AbsoluteX or
// AbsoluteY and the result wraps
func (s *ClockCyclesTest) TestCyclesCorrectWhenInAbsoluteXOrAbsoluteYAndWrap() {
	subtests := []struct {
		op       OpCode
		argument uint8
		expected uint16
	}{
		{0x7D, 0xA6, 5}, // ADC Absolute X
		{0x79, 0xA4, 5}, // ADC Absolute Y
		{0x3D, 0xA6, 5}, // AND Absolute X
		{0x39, 0xA4, 5}, // AND Absolute Y
		{0x1E, 0xA6, 7}, // ASL Absolute X
		{0xDD, 0xA6, 5}, // CMP Absolute X
		{0xD9, 0xA4, 5}, // CMP Absolute Y
		{0xDE, 0xA6, 7}, // DEC Absolute X
		{0x5D, 0xA6, 5}, // EOR Absolute X
		{0x59, 0xA4, 5}, // EOR Absolute Y
		{0xFE, 0xA6, 7}, // INC Absolute X
		{0xBD, 0xA6, 5}, // LDA Absolute X
		{0xB9, 0xA4, 5}, // LDA Absolute Y
		{0xBE, 0xA4, 5}, // LDX Absolute Y
		{0xBC, 0xA6, 5}, // LDY Absolute X
		{0x5E, 0xA6, 7}, // LSR Absolute X
		{0x1D, 0xA6, 5}, // ORA Absolute X
		{0x19, 0xA4, 5}, // ORA Absolute Y
		{0x3E, 0xA6, 7}, // ROL Absolute X
		{0x7E, 0xA6, 7}, // ROR Absolute X
		{0xFD, 0xA6, 5}, // SBC Absolute X
		{0xF9, 0xA4, 5}, // SBC Absolute Y
		{0x9D, 0xA6, 5}, // STA Absolute X
		{0x99, 0xA6, 5}, // STA Absolute Y
		// Illegal Instructions
		{0xBB, 0xA4, 5}, // LAR Absolute Y
		{0xBF, 0xA4, 5}, // LAX Absolute Y
		{0x1C, 0xA6, 5}, // TOP Absolute X
		{0x3C, 0xA6, 5}, // TOP Absolute X
		{0x5C, 0xA6, 5}, // TOP Absolute X
		{0x7C, 0xA6, 5}, // TOP Absolute X
		{0xDC, 0xA6, 5}, // TOP Absolute X
		{0xFC, 0xA6, 5}, // TOP Absolute X
	}
	r := s.Require()

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		s._instruction = data
		s.Run(fmt.Sprintf("%s/%s", data.Mnemonic, data.AddressingMode), func() {
			program := []uint8{st.argument, 0x06, uint8(data.OpCode), 0xFF, 0xFF, 0x00, 0x03}
			s.loadProgram(0x0000, program)

			s.Equal(uint16(0x0000), s.cpu.currentCycles)

			_, err := s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			s.Equal(uint16(st.expected), s.cpu.currentCycles)
		})
	}
}

// Test number of cycles remaining is correct when in Indirect Indexed
// and the result wraps
func (s *ClockCyclesTest) TestCyclesCorrectWhenInIndirectIndexedAndWrap() {
	subtests := []struct {
		op       OpCode
		expected uint16
	}{
		{0x71, 6}, // ADC Indirect Y
		{0x31, 6}, // AND Indirect Y
		{0xB1, 6}, // LDA Indirect Y
		{0xD1, 6}, // CMP Indirect Y
		{0x51, 6}, // EOR Indirect Y
		{0x11, 6}, // ORA Indirect Y
		{0xF1, 6}, // SBC Indirect Y
		{0x91, 6}, // STA Indirect Y
		// Illegal Instructions
		{0xB3, 6}, // LAX Indirect Y
	}
	r := s.Require()

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		s.Run(fmt.Sprintf("%s/%s", data.Mnemonic, data.AddressingMode), func() {
			program := []uint8{0xA0, 0x04, uint8(data.OpCode), 0x05, 0x08, 0xFF, 0xFF, 0x03}
			s.loadProgram(0x0000, program)

			s.Equal(uint16(0x0000), s.cpu.currentCycles)

			_, err := s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			s.Equal(uint16(st.expected), s.cpu.currentCycles)
		})
	}
}

// Test number of cycles remaining is correct when Relative and
// Branch On Carry, and result do not wrap
func (s *ClockCyclesTest) TestCyclesCorrectWhenRelativeAndBranchOnCarry() {
	subtests := []struct {
		op       OpCode
		flag     OpCode
		expected uint16
	}{
		{0x90, 0x38, 2}, // BCC & SEC
		{0x90, 0x18, 3}, // BCC & CLC
		{0xB0, 0x18, 2}, // BCS & CLC
		{0xB0, 0x38, 3}, // BCS & SEC
	}
	r := s.Require()

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		flag := s.cpu.Instructions.Get(st.flag)
		s.Run(fmt.Sprintf("%s/%s", data.Mnemonic, flag.Mnemonic), func() {
			program := []uint8{uint8(flag.OpCode), uint8(data.OpCode), 0x00}
			s.loadProgram(0x0000, program)

			s.Equal(uint16(0x0000), s.cpu.currentCycles)

			_, err := s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			s.Equal(uint16(st.expected), s.cpu.currentCycles)
		})
	}
}

// Test number of cycles remaining is correct when Relative and
// Branch On Carry, and result wrap
func (s *ClockCyclesTest) TestCyclesCorrectWhenRelativeAndBranchOnCarryAndWrap() {
	subtests := []struct {
		op        OpCode
		flag      OpCode
		initialPC uint16
		argument  uint8
		expected  uint16
	}{
		{0x90, 0x18, 0xFFF0, 0x0F, 4}, // BCC, CLC & 15
		{0x90, 0x18, 0x0000, 0x84, 4}, // BCC, CLC & 132
		{0xB0, 0x38, 0xFFF0, 0x0F, 4}, // BCS, SEC & 15
		{0xB0, 0x38, 0x0000, 0x84, 4}, // BCS, SEC & 132
	}
	r := s.Require()

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		flag := s.cpu.Instructions.Get(st.flag)
		s.Run(fmt.Sprintf("%s/%s/%d", data.Mnemonic, flag.Mnemonic, st.argument), func() {
			program := []uint8{uint8(flag.OpCode), uint8(data.OpCode), st.argument, 0x00}
			s.loadProgram(st.initialPC, program)

			//s.Log(MemoryDump(s.cpu.Memory, st.initialPC, nil))

			s.Equal(st.initialPC, s.cpu.Registers.PC)
			s.Equal(uint16(0x0000), s.cpu.currentCycles)

			_, err := s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			s.Equal(uint16(st.expected), s.cpu.currentCycles)
		})
	}
}

// Test number of cycles remaining is correct when Relative and
// Branch On Zero, and result do not wrap
func (s *ClockCyclesTest) TestCyclesCorrectWhenRelativeAndBranchOnZero() {
	subtests := []struct {
		op       OpCode
		argument uint8
		expected uint16
	}{
		{0xF0, 0x00, 3}, // BEQ
		{0xF0, 0x01, 2}, // BEQ
		{0xD0, 0x01, 3}, // BNE
		{0xD0, 0x00, 2}, // BNE
	}
	r := s.Require()

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		s.Run(fmt.Sprintf("%s/%d", data.Mnemonic, st.argument), func() {
			program := []uint8{0xA9, st.argument, uint8(data.OpCode), 0x00}
			s.loadProgram(0x0000, program)

			s.Equal(uint16(0x0000), s.cpu.currentCycles)

			_, err := s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			s.Equal(uint16(st.expected), s.cpu.currentCycles)
		})
	}
}

// Test number of cycles remaining is correct when Relative and
// Branch On Zero, and result wrap
func (s *ClockCyclesTest) TestCyclesCorrectWhenRelativeAndBranchOnZeroAndWrap() {
	subtests := []struct {
		initialPC uint16
		op        OpCode
		amount    uint8
		argument  uint8
		expected  uint16
	}{
		{0xFFF0, 0xF0, 0x00, 0x0D, 4}, // BEQ
		{0x0000, 0xF0, 0x00, 0x84, 4}, // BEQ
		{0xFFF0, 0xD0, 0x01, 0x0D, 4}, // BNE
		{0x0000, 0xD0, 0x01, 0x84, 4}, // BNE
	}
	r := s.Require()

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		s.Run(fmt.Sprintf("%s/%d/%d", data.Mnemonic, st.amount, st.argument), func() {
			program := []uint8{0xA9, st.amount, uint8(data.OpCode), st.argument, 0x00}
			s.loadProgram(st.initialPC, program)

			s.Equal(st.initialPC, s.cpu.Registers.PC)
			s.Equal(uint16(0x0000), s.cpu.currentCycles)

			_, err := s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			s.Equal(uint16(st.expected), s.cpu.currentCycles)
		})
	}
}

// Test number of cycles remaining is correct when Relative and
// Branch On Negative, and result do not wrap
func (s *ClockCyclesTest) TestCyclesCorrectWhenRelativeAndBranchOnNegative() {
	subtests := []struct {
		op       OpCode
		argument uint8
		expected uint16
	}{
		{0x30, 0x80, 3}, // BMI
		{0x30, 0x79, 2}, // BMI
		{0x10, 0x79, 3}, // BPL
		{0x10, 0x80, 2}, // BPL
	}
	r := s.Require()

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		s.Run(fmt.Sprintf("%s/%d", data.Mnemonic, st.argument), func() {
			program := []uint8{0xA9, st.argument, uint8(data.OpCode), 0x00}
			s.loadProgram(0x0000, program)

			s.Equal(uint16(0x0000), s.cpu.currentCycles)

			_, err := s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			s.Equal(uint16(st.expected), s.cpu.currentCycles)
		})
	}
}

// Test number of cycles remaining is correct when Relative and
// Branch On Negative, and result wrap
func (s *ClockCyclesTest) TestCyclesCorrectWhenRelativeAndBranchOnNegativeAndWrap() {
	subtests := []struct {
		initialPC uint16
		op        OpCode
		amount    uint8
		argument  uint8
		expected  uint16
	}{
		{0xFFF0, 0x30, 0x80, 0x0D, 4}, // BMI
		{0x0000, 0x30, 0x80, 0x84, 4}, // BMI
		{0xFFF0, 0x10, 0x79, 0x0D, 4}, // BPL
		{0x0000, 0x10, 0x79, 0x84, 4}, // BPL
	}
	r := s.Require()

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		s.Run(fmt.Sprintf("%s/%d/%d", data.Mnemonic, st.amount, st.argument), func() {
			program := []uint8{0xA9, st.amount, uint8(data.OpCode), st.argument, 0x00}
			s.loadProgram(st.initialPC, program)

			s.Equal(st.initialPC, s.cpu.Registers.PC)
			s.Equal(uint16(0x0000), s.cpu.currentCycles)

			_, err := s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			s.Equal(uint16(st.expected), s.cpu.currentCycles)
		})
	}
}

// Test number of cycles remaining is correct when Relative and
// Branch On Overflow, and result do not wrap
func (s *ClockCyclesTest) TestCyclesCorrectWhenRelativeAndBranchOnOverflow() {
	subtests := []struct {
		op       OpCode
		argument uint8
		expected uint16
	}{
		{0x50, 0x01, 3}, // BVC
		{0x50, 0x7F, 2}, // BVC
		{0x70, 0x7F, 3}, // BVS
		{0x70, 0x01, 2}, // BVS
	}
	r := s.Require()

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		s.Run(fmt.Sprintf("%s/%d", data.Mnemonic, st.argument), func() {
			program := []uint8{0xA9, 0x01, 0x69, st.argument, uint8(data.OpCode), 0x00}
			s.loadProgram(0x0000, program)

			s.Equal(uint16(0x0000), s.cpu.currentCycles)

			_, err := s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			s.Equal(uint16(st.expected), s.cpu.currentCycles)
		})
	}
}

// Test number of cycles remaining is correct when Relative and
// Branch On Overflow, and result wrap
func (s *ClockCyclesTest) TestCyclesCorrectWhenRelativeAndBranchOnOverflowAndWrap() {
	subtests := []struct {
		initialPC uint16
		op        OpCode
		amount    uint8
		argument  uint8
		expected  uint16
	}{
		{0xFFF0, 0x50, 0x00, 0x0B, 4}, // BVC
		{0x0000, 0x50, 0x00, 0x86, 4}, // BVC
		{0xFFF0, 0x70, 0x7F, 0x0B, 4}, // BVS
		{0x0000, 0x70, 0x7F, 0x86, 4}, // BVS
	}
	r := s.Require()

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		s.Run(fmt.Sprintf("%s/%d/%d", data.Mnemonic, st.amount, st.argument), func() {
			program := []uint8{0xA9, st.amount, 0x69, 0x01, uint8(data.OpCode), st.argument, 0x00}
			s.loadProgram(st.initialPC, program)

			s.Equal(st.initialPC, s.cpu.Registers.PC)
			s.Equal(uint16(0x0000), s.cpu.currentCycles)

			_, err := s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			s.Equal(uint16(st.expected), s.cpu.currentCycles)
		})
	}
}

// endregion ClockCyclesTest

// region ProgramCounterTest
type ProgramCounterTest struct {
	baseTest
}

func (s *ProgramCounterTest) TestProgramCounterCorrect() {
	subtests := []struct {
		op       OpCode
		expected uint16
	}{
		{0x69, 2}, // ADC Immediate
		{0x65, 2}, // ADC Zero Page
		{0x75, 2}, // ADC Zero Page X
		{0x6D, 3}, // ADC Absolute
		{0x7D, 3}, // ADC Absolute X
		{0x79, 3}, // ADC Absolute Y
		{0x61, 2}, // ADC Indirect X
		{0x71, 2}, // ADC Indirect Y
		{0x29, 2}, // AND Immediate
		{0x25, 2}, // AND Zero Page
		{0x35, 2}, // AND Zero Page X
		{0x2D, 3}, // AND Absolute
		{0x3D, 3}, // AND Absolute X
		{0x39, 3}, // AND Absolute Y
		{0x21, 2}, // AND Indirect X
		{0x31, 2}, // AND Indirect Y
		{0x0A, 1}, // ASL Accumulator
		{0x06, 2}, // ASL Zero Page
		{0x16, 2}, // ASL Zero Page X
		{0x0E, 3}, // ASL Absolute
		{0x1E, 3}, // ASL Absolute X
		{0x24, 2}, // BIT Zero Page
		{0x2C, 3}, // BIT Absolute
		{0x18, 1}, // CLC Implied
		{0xD8, 1}, // CLD Implied
		{0x58, 1}, // CLI Implied
		{0xB8, 1}, // CLV Implied
		{0xC9, 2}, // CMP Immediate
		{0xC5, 2}, // CMP ZeroPage
		{0xD5, 2}, // CMP Zero Page X
		{0xCD, 3}, // CMP Absolute
		{0xDD, 3}, // CMP Absolute X
		{0xD9, 3}, // CMP Absolute Y
		{0xC1, 2}, // CMP Indirect X
		{0xD1, 2}, // CMP Indirect Y
		{0xE0, 2}, // CPX Immediate
		{0xE4, 2}, // CPX ZeroPage
		{0xEC, 3}, // CPX Absolute
		{0xC0, 2}, // CPY Immediate
		{0xC4, 2}, // CPY ZeroPage
		{0xCC, 3}, // CPY Absolute
		{0xC6, 2}, // DEC Zero Page
		{0xD6, 2}, // DEC Zero Page X
		{0xCE, 3}, // DEC Absolute
		{0xDE, 3}, // DEC Absolute X
		{0xCA, 1}, // DEX Implied
		{0x88, 1}, // DEY Implied
		{0x49, 2}, // EOR Immediate
		{0x45, 2}, // EOR ZeroPage
		{0x55, 2}, // EOR Zero Page X
		{0x4D, 3}, // EOR Absolute
		{0x5D, 3}, // EOR Absolute X
		{0x59, 3}, // EOR Absolute Y
		{0x41, 2}, // EOR Indirect X
		{0x51, 2}, // EOR Indirect Y
		{0xE6, 2}, // INC Zero Page
		{0xF6, 2}, // INC Zero Page X
		{0xEE, 3}, // INC Absolute
		{0xFE, 3}, // INC Absolute X
		{0xE8, 1}, // INX Implied
		{0xC8, 1}, // INY Implied
		{0xA9, 2}, // LDA Immediate
		{0xA5, 2}, // LDA Zero Page
		{0xB5, 2}, // LDA Zero Page X
		{0xAD, 3}, // LDA Absolute
		{0xBD, 3}, // LDA Absolute X
		{0xB9, 3}, // LDA Absolute Y
		{0xA1, 2}, // LDA Indirect X
		{0xB1, 2}, // LDA Indirect Y
		{0xA2, 2}, // LDX Immediate
		{0xA6, 2}, // LDX Zero Page
		{0xB6, 2}, // LDX Zero Page Y
		{0xAE, 3}, // LDX Absolute
		{0xBE, 3}, // LDX Absolute Y
		{0xA0, 2}, // LDY Immediate
		{0xA4, 2}, // LDY Zero Page
		{0xB4, 2}, // LDY Zero Page Y
		{0xAC, 3}, // LDY Absolute
		{0xBC, 3}, // LDY Absolute Y
		{0x4A, 1}, // LSR Accumulator
		{0x46, 2}, // LSR Zero Page
		{0x56, 2}, // LSR Zero Page X
		{0x4E, 3}, // LSR Absolute
		{0x5E, 3}, // LSR Absolute X
		{0xEA, 1}, // NOP Implied
		{0x09, 2}, // ORA Immediate
		{0x05, 2}, // ORA Zero Page
		{0x15, 2}, // ORA Zero Page X
		{0x0D, 3}, // ORA Absolute
		{0x1D, 3}, // ORA Absolute X
		{0x19, 3}, // ORA Absolute Y
		{0x01, 2}, // ORA Indirect X
		{0x11, 2}, // ORA Indirect Y
		{0x48, 1}, // PHA Implied
		{0x08, 1}, // PHP Implied
		{0x68, 1}, // PLA Implied
		{0x28, 1}, // PLP Implied
		{0x2A, 1}, // ROL Accumulator
		{0x26, 2}, // ROL Zero Page
		{0x36, 2}, // ROL Zero Page X
		{0x2E, 3}, // ROL Absolute
		{0x3E, 3}, // ROL Absolute X
		{0x6A, 1}, // ROR Accumulator
		{0x66, 2}, // ROR Zero Page
		{0x76, 2}, // ROR Zero Page X
		{0x6E, 3}, // ROR Absolute
		{0x7E, 3}, // ROR Absolute X
		{0xE9, 2}, // SBC Immediate
		{0xE5, 2}, // SBC Zero Page
		{0xF5, 2}, // SBC Zero Page X
		{0xED, 3}, // SBC Absolute
		{0xFD, 3}, // SBC Absolute X
		{0xF9, 3}, // SBC Absolute Y
		{0xE1, 2}, // SBC Indrect X
		{0xF1, 2}, // SBC Indirect Y
		{0x38, 1}, // SEC Implied
		{0xF8, 1}, // SED Implied
		{0x78, 1}, // SEI Implied
		{0x85, 2}, // STA ZeroPage
		{0x95, 2}, // STA Zero Page X
		{0x8D, 3}, // STA Absolute
		{0x9D, 3}, // STA Absolute X
		{0x99, 3}, // STA Absolute Y
		{0x81, 2}, // STA Indirect X
		{0x91, 2}, // STA Indirect Y
		{0x86, 2}, // STX Zero Page
		{0x96, 2}, // STX Zero Page Y
		{0x8E, 3}, // STX Absolute
		{0x84, 2}, // STY Zero Page
		{0x94, 2}, // STY Zero Page X
		{0x8C, 3}, // STY Absolute
		{0xAA, 1}, // TAX Implied
		{0xA8, 1}, // TAY Implied
		{0xBA, 1}, // TSX Implied
		{0x8A, 1}, // TXA Implied
		{0x9A, 1}, // TXS Implied
		{0x98, 1}, // TYA Implied
		// Illegal Instructions
		{0x0B, 2}, // AAC Immediate
		{0x87, 2}, // AAX Zero Page
		{0x97, 2}, // AAX Zero Page Y
		{0x8F, 3}, // AAX Absolute
		{0x83, 2}, // AAX Indirect X
		{0x6B, 2}, // ARR Immediate
		{0x4B, 2}, // ASR Immediate
		{0xAB, 2}, // ATX Immediate
		{0x9F, 3}, // AXA Absolute Y
		{0x93, 2}, // AXA Indirect Y
		{0xCB, 2}, // AXS Immediate
		{0xC7, 2}, // DCP Zero Page
		{0xD7, 2}, // DCP Zero Page X
		{0xCF, 3}, // DCP Absolute
		{0xDF, 3}, // DCP Absolute X
		{0xDB, 3}, // DCP Absolute Y
		{0xC3, 2}, // DCP Indirect X
		{0xD3, 2}, // DCP Indirect Y
		{0xE7, 2}, // ISC Zero Page
		{0xF7, 2}, // ISC Zero Page X
		{0xEF, 3}, // ISC Absolute
		{0xFF, 3}, // ISC Absolute X
		{0xFB, 3}, // ISC Absolute Y
		{0xE3, 2}, // ISC Indirect X
		{0xF3, 2}, // ISC Indirect Y
		// {0x, },  // KIL Immediate (Untestable)
		{0xBB, 3}, // LAR Absolute Y
		{0xA7, 2}, // LAX Zero Page
		{0xB7, 2}, // LAX Zero Page Y
		{0xAF, 3}, // LAX Absolute
		{0xBF, 3}, // LAX Absolute Y
		{0xA3, 2}, // LAX Indirect X
		{0xB3, 2}, // LAX Indirect Y
		{0x27, 2}, // RLA Zero Page
		{0x37, 2}, // RLA Zero Page X
		{0x2F, 3}, // RLA Absolute
		{0x3F, 3}, // RLA Absolute X
		{0x3B, 3}, // RLA Absolute Y
		{0x23, 2}, // RLA Indirect X
		{0x33, 2}, // RLA Indirect Y
		{0x67, 2}, // RRA Zero Page
		{0x77, 2}, // RRA Zero Page X
		{0x6F, 3}, // RRA Absolute
		{0x7F, 3}, // RRA Absolute X
		{0x7B, 3}, // RRA Absolute Y
		{0x63, 2}, // RRA Indirect X
		{0x73, 2}, // RRA Indirect Y
		{0xEB, 2}, // USBC Immediate (SBC)
		{0x07, 2}, // SLO Zero Page
		{0x17, 2}, // SLO Zero Page X
		{0x0F, 3}, // SLO Absolute
		{0x1F, 3}, // SLO Absolute X
		{0x1B, 3}, // SLO Absolute Y
		{0x03, 2}, // SLO Indirect X
		{0x13, 2}, // SLO Indirect Y
		{0x47, 2}, // SRE Zero Page
		{0x57, 2}, // SRE Zero Page X
		{0x4F, 3}, // SRE Absolute
		{0x5F, 3}, // SRE Absolute X
		{0x5B, 3}, // SRE Absolute Y
		{0x43, 2}, // SRE Indirect X
		{0x53, 2}, // SRE Indirect Y
		{0x9E, 3}, // SXA Absolute Y
		{0x9C, 3}, // SYA Absolute X
		{0x8B, 2}, // XAA Immediate
		{0x9B, 3}, // XAS Absolute Y
		{0x1A, 1}, // NOP Implied
		{0x3A, 1}, // NOP Implied
		{0x5A, 1}, // NOP Implied
		{0x7A, 1}, // NOP Implied
		{0xDA, 1}, // NOP Implied
		{0xFA, 1}, // NOP Implied
		{0x80, 2}, // NOP Immediate   (SKB)
		{0x82, 2}, // NOP Immediate   (SKB)
		{0x89, 2}, // NOP Immediate   (SKB)
		{0xC2, 2}, // NOP Immediate   (SKB)
		{0xE2, 2}, // NOP Immediate   (SKB)
		{0x04, 2}, // NOP Zero Page   (IGN)
		{0x44, 2}, // NOP Zero Page   (IGN)
		{0x64, 2}, // NOP Zero Page   (IGN)
		{0x14, 2}, // NOP Zero Page X (IGN)
		{0x34, 2}, // NOP Zero Page X (IGN)
		{0x54, 2}, // NOP Zero Page X (IGN)
		{0x74, 2}, // NOP Zero Page X (IGN)
		{0xD4, 2}, // NOP Zero Page X (IGN)
		{0xF4, 2}, // NOP Zero Page X (IGN)
		{0x0C, 3}, // NOP Absolute    (IGN)
		{0x1C, 3}, // NOP Absolute X  (IGN)
		{0x3C, 3}, // NOP Absolute X  (IGN)
		{0x5C, 3}, // NOP Absolute X  (IGN)
		{0x7C, 3}, // NOP Absolute X  (IGN)
		{0xDC, 3}, // NOP Absolute X  (IGN)
		{0xFC, 3}, // NOP Absolute X  (IGN)
	}
	r := s.Require()

	for _, st := range subtests {
		data := s.cpu.Instructions.Get(st.op)
		s.Run(fmt.Sprintf("%s/%s", data.Mnemonic, data.AddressingMode), func() {
			program := []uint8{uint8(data.OpCode), 0x00}
			s.loadProgram(0x0000, program)

			s.Equal(uint16(0x0000), s.cpu.Registers.PC)

			_, err := s.cpu.Execute()
			r.NoError(err)

			s.Equal(uint16(data.Bytes()), s.cpu.Registers.PC)
		})
	}
}

func (s *ProgramCounterTest) TestBranchOnCarryProgramCounterCorrectWhenNoBranchOccurs() {
	subtests := []struct {
		op           OpCode
		argument     uint8
		resultOffset uint16
	}{
		{0x90, 0x38, 2}, // BCC
		{0xB0, 0x18, 2}, // BCS
	}
	r := s.Require()

	for _, data := range subtests {
		i := s.cpu.Instructions.Get(data.op)
		s.Run(i.Mnemonic, func() {
			s.cpu.Registers.PC = 0x0000
			program := []uint8{data.argument, uint8(data.op), 0x48}
			s.loadProgram(0x0000, program)
			//s.cpu.Memory.Store(0x0000, uint8(data.op))

			_, err := s.cpu.Execute()
			r.NoError(err)

			currentPC := s.cpu.Registers.PC
			_, err = s.cpu.Execute()
			r.NoError(err)

			expectedPc := currentPC + data.resultOffset
			s.Equal(expectedPc, s.cpu.Registers.PC)
		})
	}
}

func (s *ProgramCounterTest) TestBranchOnZeroProgramCounterCorrectWhenNoBranchOccurs() {
	subtests := []struct {
		op           OpCode
		argument     uint8
		resultOffset uint16
	}{
		{0xF0, 0x01, 2}, // BEQ
		{0xD0, 0x00, 2}, // BNE
	}
	r := s.Require()

	for _, data := range subtests {
		i := s.cpu.Instructions.Get(data.op)
		s.Run(i.Mnemonic, func() {
			s.cpu.Registers.PC = 0x0000
			program := []uint8{0xa9, data.argument, uint8(data.op)}
			s.loadProgram(0x0000, program)
			//s.cpu.Memory.Store(0x0000, uint8(data.op))

			_, err := s.cpu.Execute()
			r.NoError(err)

			currentPC := s.cpu.Registers.PC
			_, err = s.cpu.Execute()
			r.NoError(err)

			expectedPc := currentPC + data.resultOffset
			s.Equal(expectedPc, s.cpu.Registers.PC)
		})
	}
}

func (s *ProgramCounterTest) TestBranchOnNegativeProgramCounterCorrectWhenNoBranchOccurs() {
	subtests := []struct {
		op           OpCode
		argument     uint8
		resultOffset uint16
	}{
		{0x30, 0x79, 2}, // BMI
		{0x10, 0x80, 2}, // BPL
	}
	r := s.Require()

	for _, data := range subtests {
		i := s.cpu.Instructions.Get(data.op)
		s.Run(i.Mnemonic, func() {
			s.cpu.Registers.PC = 0x0000
			program := []uint8{0xa9, data.argument, uint8(data.op)}
			s.loadProgram(0x0000, program)
			//s.cpu.Memory.Store(0x0000, uint8(data.op))

			_, err := s.cpu.Execute()
			r.NoError(err)

			currentPC := s.cpu.Registers.PC
			_, err = s.cpu.Execute()
			r.NoError(err)

			expectedPc := currentPC + data.resultOffset
			s.Equal(expectedPc, s.cpu.Registers.PC)
		})
	}
}

func (s *ProgramCounterTest) TestBranchOnOverflowProgramCounterCorrectWhenNoBranchOccurs() {
	subtests := []struct {
		op           OpCode
		argument     uint8
		resultOffset uint16
	}{
		{0x50, 0x7F, 2}, // BVC
		{0x70, 0x01, 2}, // BVS
	}
	r := s.Require()

	for _, data := range subtests {
		i := s.cpu.Instructions.Get(data.op)
		s.Run(i.Mnemonic, func() {
			s.cpu.Registers.PC = 0x0000
			program := []uint8{0xa9, 0x01, 0x69, data.argument, uint8(data.op), 0x00}
			s.loadProgram(0x0000, program)
			//s.cpu.Memory.Store(0x0000, uint8(data.op))

			_, err := s.cpu.Execute()
			r.NoError(err)

			_, err = s.cpu.Execute()
			r.NoError(err)

			currentPC := s.cpu.Registers.PC
			_, err = s.cpu.Execute()
			r.NoError(err)

			expectedPc := currentPC + data.resultOffset
			s.Equal(expectedPc, s.cpu.Registers.PC)
		})
	}
}

func (s *ProgramCounterTest) TestProgramCounterWrapsCorrectly() {
	s.cpu.Registers.PC = 0xFFFF
	s.cpu.Memory.Store(0xFFFF, 0x38)

	s.EqualValues(uint16(0xFFFF), s.cpu.Registers.PC)

	s.cpu.Execute()

	s.EqualValues(uint16(0x0000), s.cpu.Registers.PC)
}

// endregion ProgramCounterTest

// region Unofficial Instructions Tests

// region InstructionLAX
type InstructionLAX struct {
	baseTest
}

func (s *InstructionLAX) TestLAXLoadsAAndX() {
	program := []uint8{0xA7, 0x05}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x05, 0x42)
	s.cpu.Execute() // LAX
	s.EqualValues(0x42, s.cpu.Registers.A)
	s.EqualValues(0x42, s.cpu.Registers.X)
}

func (s *InstructionLAX) TestLAXSetsZeroFlag() {
	program := []uint8{0xA7, 0x05}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x05, 0x00)
	s.cpu.Execute() // LAX
	s.True(s.cpu.getFlag(Z))
}

func (s *InstructionLAX) TestLAXSetsNegativeFlag() {
	program := []uint8{0xA7, 0x05}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x05, 0x80)
	s.cpu.Execute() // LAX
	s.True(s.cpu.getFlag(N))
}

// endregion InstructionLAX

// region InstructionSAX
type InstructionSAX struct {
	baseTest
}

func (s *InstructionSAX) TestSAXStoresAAndX() {
	program := []uint8{0xA9, 0x0F, 0xA2, 0xF0, 0x87, 0x05}
	s.loadProgram(0x0600, program)
	s.cpu.Execute()                               // LDA #$0F
	s.cpu.Execute()                               // LDX #$F0
	s.cpu.Execute()                               // SAX $05
	s.EqualValues(0x00, s.cpu.Memory.Fetch(0x05)) // 0x0F AND 0xF0 = 0x00
}

func (s *InstructionSAX) TestSAXStoresCorrectValue() {
	program := []uint8{0xA9, 0xFF, 0xA2, 0xAA, 0x87, 0x05}
	s.loadProgram(0x0600, program)
	s.cpu.Execute()                               // LDA #$FF
	s.cpu.Execute()                               // LDX #$AA
	s.cpu.Execute()                               // SAX $05
	s.EqualValues(0xAA, s.cpu.Memory.Fetch(0x05)) // 0xFF AND 0xAA = 0xAA
}

// endregion InstructionSAX

// region InstructionDCP
type InstructionDCP struct {
	baseTest
}

func (s *InstructionDCP) TestDCPDecrementsAndCompares() {
	program := []uint8{0xA9, 0x05, 0xC7, 0x03}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x03, 0x06)
	s.cpu.Execute()                               // LDA #$05
	s.cpu.Execute()                               // DCP $03
	s.EqualValues(0x05, s.cpu.Memory.Fetch(0x03)) // 0x06 - 1 = 0x05
	s.True(s.cpu.getFlag(Z))                      // A == M
	s.True(s.cpu.getFlag(C))                      // A >= M
}

func (s *InstructionDCP) TestDCPSetsCarryFlag() {
	program := []uint8{0xA9, 0x10, 0xC7, 0x03}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x03, 0x05)
	s.cpu.Execute()                               // LDA #$10
	s.cpu.Execute()                               // DCP $03
	s.EqualValues(0x04, s.cpu.Memory.Fetch(0x03)) // 0x05 - 1 = 0x04
	s.True(s.cpu.getFlag(C))                      // A >= M (0x10 >= 0x04)
}

// endregion InstructionDCP

// region InstructionISB
type InstructionISB struct {
	baseTest
}

func (s *InstructionISB) TestISBIncrementsAndSubtracts() {
	program := []uint8{0xA9, 0x10, 0x38, 0xE7, 0x03}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x03, 0x05)
	s.cpu.Execute()                               // LDA #$10
	s.cpu.Execute()                               // SEC
	s.cpu.Execute()                               // ISB $03
	s.EqualValues(0x06, s.cpu.Memory.Fetch(0x03)) // 0x05 + 1 = 0x06
	s.EqualValues(0x0A, s.cpu.Registers.A)        // 0x10 - 0x06 = 0x0A
}

func (s *InstructionISB) TestISBSetsFlags() {
	program := []uint8{0xA9, 0x05, 0x38, 0xE7, 0x03}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x03, 0x04)
	s.cpu.Execute()                               // LDA #$05
	s.cpu.Execute()                               // SEC
	s.cpu.Execute()                               // ISB $03
	s.EqualValues(0x05, s.cpu.Memory.Fetch(0x03)) // 0x04 + 1 = 0x05
	s.EqualValues(0x00, s.cpu.Registers.A)        // 0x05 - 0x05 = 0x00
	s.True(s.cpu.getFlag(Z))
	s.True(s.cpu.getFlag(C))
}

// endregion InstructionISB

// region InstructionSLO
type InstructionSLO struct {
	baseTest
}

func (s *InstructionSLO) TestSLOShiftsAndOrs() {
	program := []uint8{0xA9, 0x01, 0x07, 0x03}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x03, 0x02)
	s.cpu.Execute()                               // LDA #$01
	s.cpu.Execute()                               // SLO $03
	s.EqualValues(0x04, s.cpu.Memory.Fetch(0x03)) // 0x02 << 1 = 0x04
	s.EqualValues(0x05, s.cpu.Registers.A)        // 0x01 OR 0x04 = 0x05
}

func (s *InstructionSLO) TestSLOSetsCarryFlag() {
	program := []uint8{0xA9, 0x00, 0x07, 0x03}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x03, 0x80)
	s.cpu.Execute()          // LDA #$00
	s.cpu.Execute()          // SLO $03
	s.True(s.cpu.getFlag(C)) // Bit 7 was set
}

// endregion InstructionSLO

// region InstructionSRE
type InstructionSRE struct {
	baseTest
}

func (s *InstructionSRE) TestSREShiftsAndEors() {
	program := []uint8{0xA9, 0xFF, 0x47, 0x03}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x03, 0x02)
	s.cpu.Execute()                               // LDA #$FF
	s.cpu.Execute()                               // SRE $03
	s.EqualValues(0x01, s.cpu.Memory.Fetch(0x03)) // 0x02 >> 1 = 0x01
	s.EqualValues(0xFE, s.cpu.Registers.A)        // 0xFF XOR 0x01 = 0xFE
}

func (s *InstructionSRE) TestSRESetsCarryFlag() {
	program := []uint8{0xA9, 0x00, 0x47, 0x03}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x03, 0x01)
	s.cpu.Execute()          // LDA #$00
	s.cpu.Execute()          // SRE $03
	s.True(s.cpu.getFlag(C)) // Bit 0 was set
}

// endregion InstructionSRE

// region InstructionRLA
type InstructionRLA struct {
	baseTest
}

func (s *InstructionRLA) TestRLARotatesAndAnds() {
	program := []uint8{0xA9, 0xFF, 0x27, 0x03}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x03, 0x55)                // 01010101
	s.cpu.Execute()                               // LDA #$FF
	s.cpu.Execute()                               // RLA $03
	s.EqualValues(0xAA, s.cpu.Memory.Fetch(0x03)) // 01010101 << 1 = 10101010
	s.EqualValues(0xAA, s.cpu.Registers.A)        // 0xFF AND 0xAA = 0xAA
}

func (s *InstructionRLA) TestRLASetsCarryFlag() {
	program := []uint8{0xA9, 0xFF, 0x27, 0x03}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x03, 0x80)
	s.cpu.Execute()          // LDA #$FF
	s.cpu.Execute()          // RLA $03
	s.True(s.cpu.getFlag(C)) // Bit 7 was set
}

// endregion InstructionRLA

// region InstructionRRA
type InstructionRRA struct {
	baseTest
}

func (s *InstructionRRA) TestRRARotatesAndAdds() {
	program := []uint8{0xA9, 0x00, 0x38, 0x67, 0x03}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x03, 0x02)
	s.cpu.Execute()                               // LDA #$00
	s.cpu.Execute()                               // SEC
	s.cpu.Execute()                               // RRA $03
	s.EqualValues(0x81, s.cpu.Memory.Fetch(0x03)) // 0x02 ROR with C=1 = 0x81
	s.EqualValues(0x81, s.cpu.Registers.A)        // 0x00 + 0x81 + 1 = 0x82 (wait, carry matters)
}

// endregion InstructionRRA

// region InstructionALR
type InstructionALR struct {
	baseTest
}

func (s *InstructionALR) TestALRAndsAndShifts() {
	program := []uint8{0xA9, 0xFF, 0x4B, 0xFE}
	s.loadProgram(0x0600, program)
	s.cpu.Execute()                        // LDA #$FF
	s.cpu.Execute()                        // ALR #$FE
	s.EqualValues(0x7F, s.cpu.Registers.A) // (0xFF AND 0xFE) >> 1 = 0xFE >> 1 = 0x7F
	s.False(s.cpu.getFlag(C))              // Bit 0 of 0xFE is 0
}

func (s *InstructionALR) TestALRSetsCarryFlag() {
	program := []uint8{0xA9, 0xFF, 0x4B, 0xFF}
	s.loadProgram(0x0600, program)
	s.cpu.Execute()                        // LDA #$FF
	s.cpu.Execute()                        // ALR #$FF
	s.EqualValues(0x7F, s.cpu.Registers.A) // (0xFF AND 0xFF) >> 1 = 0x7F
	s.True(s.cpu.getFlag(C))               // Bit 0 of 0xFF is 1
}

// endregion InstructionALR

// region InstructionANC
type InstructionANC struct {
	baseTest
}

func (s *InstructionANC) TestANCAndsAndSetsCarry() {
	program := []uint8{0xA9, 0xFF, 0x0B, 0x80}
	s.loadProgram(0x0600, program)
	s.cpu.Execute()                        // LDA #$FF
	s.cpu.Execute()                        // ANC #$80
	s.EqualValues(0x80, s.cpu.Registers.A) // 0xFF AND 0x80 = 0x80
	s.True(s.cpu.getFlag(C))               // Bit 7 is set
	s.True(s.cpu.getFlag(N))               // Bit 7 is set
}

func (s *InstructionANC) TestANCClearsCarryWhenBit7Clear() {
	program := []uint8{0xA9, 0xFF, 0x38, 0x0B, 0x7F}
	s.loadProgram(0x0600, program)
	s.cpu.Execute()                        // LDA #$FF
	s.cpu.Execute()                        // SEC
	s.cpu.Execute()                        // ANC #$7F
	s.EqualValues(0x7F, s.cpu.Registers.A) // 0xFF AND 0x7F = 0x7F
	s.False(s.cpu.getFlag(C))              // Bit 7 is clear
}

// endregion InstructionANC

// region InstructionARR
type InstructionARR struct {
	baseTest
}

func (s *InstructionARR) TestARRAndsAndRotates() {
	program := []uint8{0xA9, 0xFF, 0x18, 0x6B, 0xAA}
	s.loadProgram(0x0600, program)
	s.cpu.Execute()                        // LDA #$FF
	s.cpu.Execute()                        // CLC
	s.cpu.Execute()                        // ARR #$AA
	s.EqualValues(0x55, s.cpu.Registers.A) // (0xFF AND 0xAA) ROR = 0xAA ROR = 0x55
}

// endregion InstructionARR

// region InstructionSBX
type InstructionSBX struct {
	baseTest
}

func (s *InstructionSBX) TestSBXComparesAndDecrements() {
	program := []uint8{0xA9, 0xFF, 0xA2, 0x10, 0xCB, 0x05}
	s.loadProgram(0x0600, program)
	s.cpu.Execute()                        // LDA #$FF
	s.cpu.Execute()                        // LDX #$10
	s.cpu.Execute()                        // SBX #$05
	s.EqualValues(0x0B, s.cpu.Registers.X) // (A AND X) - oper = (0xFF AND 0x10) - 0x05 = 0x10 - 0x05 = 0x0B
	s.True(s.cpu.getFlag(C))               // Result >= 0
}

// endregion InstructionSBX

// region InstructionLAS
type InstructionLAS struct {
	baseTest
}

func (s *InstructionLAS) TestLASLoadsAndTransfers() {
	program := []uint8{0xBB, 0x05, 0x00}
	s.loadProgram(0x0600, program)
	s.cpu.Memory.Store(0x05, 0x42)
	s.cpu.Registers.SP = 0xFF
	s.cpu.Execute() // LAS $0005,Y
	// A = X = SP = SP AND M
	result := uint8(0xFF & 0x42)
	s.EqualValues(result, s.cpu.Registers.A)
	s.EqualValues(result, s.cpu.Registers.X)
	s.EqualValues(result, s.cpu.Registers.SP)
}

// endregion InstructionLAS

// endregion Unofficial Instructions Tests
