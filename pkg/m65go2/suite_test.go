package m65go2

import (
	"testing"

	s "github.com/stretchr/testify/suite"
	"github.com/txc/nesgo/pkg/utils"
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
	}
	utils.RunSuiteTests(t, tests)
}

// region baseTest
type baseTest struct {
	utils.CommonTest
	cpu *M6502
}

//func (s *baseTest) SetupSuite(suiteName, testName string) {}

//func (s *baseTest) TearDownSuite(suiteName, testName string) {}

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

//func (s *baseTest) TearDownSubTest() {}

// endregion baseTest
