package m65go2

import (
	"reflect"
	"testing"

	s "github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	tests := []s.TestingSuite{
		&InitializationTest{},
		&CPUTest{},
		&InstructionADC{},
		&InstructionAND{},
		&InstructionASL{},
		&InstructionBCC{},
		&InstructionBCS{},
		&InstructionBEQ{},
		&InstructionBIT{},
		&InstructionBMI{},
		&InstructionBNE{},
		&InstructionBPL{},
		&InstructionBVC{},
		&InstructionBVS{},
		&InstructionCLC{},
		&InstructionCLD{},
		&InstructionCLI{},
		&InstructionCLV{},
		&InstructionCMP{},
		&InstructionCPX{},
		&InstructionCPY{},
		&InstructionDEC{},
		&InstructionDEX{},
		&InstructionDEY{},
		&InstructionEOR{},
		&InstructionINC{},
		&InstructionINX{},
		&InstructionINY{},
		&InstructionJMP{},
		&InstructionJSR{},
		&InstructionLDA{},
		&InstructionLDX{},
		&InstructionLDY{},
		&InstructionLSR{},
		&InstructionORA{},
		&InstructionPHA{},
		&InstructionPHP{},
		&InstructionPLA{},
		&InstructionPLP{},
		&InstructionROL{},
		&InstructionROR{},
		&InstructionRTI{},
		&InstructionRTS{},
		&InstructionSBC{},
		&InstructionSEC{},
		&InstructionSED{},
		&InstructionSEI{},
		&InstructionSTA{},
		&InstructionSTX{},
		&InstructionSTY{},
		&InstructionTAX{},
		&InstructionTSX{},
		&InstructionTXS{},
		&AccumulatorAddress{},
		&IndexAddress{},
		&CompareAddress{},
		&DecrementIncrementAddress{},
		&StoreInMemoryAddress{},
		&ClockCyclesTest{},
		&ProgramCounterTest{},
		// Unofficial Instructions
		&InstructionLAX{},
		&InstructionSAX{},
		&InstructionDCP{},
		&InstructionISB{},
		&InstructionSLO{},
		&InstructionSRE{},
		&InstructionRLA{},
		&InstructionRRA{},
		&InstructionALR{},
		&InstructionANC{},
		&InstructionARR{},
		&InstructionSBX{},
		&InstructionLAS{},
		// instructions test
		&InstructionsTest{},
		// memory test
		&MemoryTest{},
		// nestest
		&NesTest{},
		&FunctionalTest{},
		&InterruptTest{},
		&DecimalTest{},
	}

	for _, suite := range tests {
		var suiteName string
		if t := reflect.TypeOf(suite); t.Kind() == reflect.Pointer {
			suiteName = t.Elem().Name()
		} else {
			suiteName = t.Name()
		}

		t.Run(suiteName, func(sub *testing.T) {
			s.Run(sub, suite)
		})
	}
}

// region baseTest
type baseTest struct {
	s.Suite
	cpu *M6502
}

func (s *baseTest) SetupSuite(suiteName, testName string) {}

func (s *baseTest) TearDownSuite(suiteName, testName string) {}

func (s *baseTest) BeforeTest(suiteName, testName string) {
	s.cpu = NewM6502(NewBasicMemory(DefaultMemorySize))
	s.cpu.Reset()
	s.cpu.breakError = true
}

func (s *baseTest) AfterTest(_, _ string) {
	s.cpu = nil
}

func (s *baseTest) SetupSubTest() {
	s.cpu.Reset()
}

func (s *baseTest) TearDownSubTest() {
}

func (s *baseTest) Log(args ...any) { s.T().Log(args...) }

func (s *baseTest) Logf(format string, args ...any) { s.T().Logf(format, args...) }

// SkipShort skips the test if running in short mode
func (s *baseTest) SkipShort() {
	if testing.Short() {
		s.T().Skip("Skipping trace test in short mode")
	}
}

// endregion baseTest
