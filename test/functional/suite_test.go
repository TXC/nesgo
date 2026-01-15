package functional

import (
	"testing"

	s "github.com/stretchr/testify/suite"
	"github.com/txc/nesgo/pkg/m65go2"
	"github.com/txc/nesgo/pkg/utils"
)

const (
	romPath = "../roms"
)

func TestSuite(t *testing.T) {
	tests := []s.TestingSuite{
		// nestest
		&NesTest{},
		&FunctionalTest{},
		&InterruptTest{},
		&DecimalTest{},
	}
	utils.RunSuiteTests(t, tests)
}

// region baseTest
type baseTest struct {
	utils.CommonTest
	cpu *m65go2.M6502
}

//func (s *baseTest) SetupSuite(suiteName, testName string) {}

//func (s *baseTest) TearDownSuite(suiteName, testName string) {}

func (s *baseTest) BeforeTest(suiteName, testName string) {
	s.cpu = m65go2.NewM6502(m65go2.NewBasicMemory(m65go2.DefaultMemorySize))
	s.cpu.Reset()
}

func (s *baseTest) AfterTest(_, _ string) {
	s.cpu = nil
}

func (s *baseTest) SetupSubTest() {
	s.cpu.Reset()
}

//func (s *baseTest) TearDownSubTest() {}

// endregion baseTest
