package utils

import (
	"reflect"
	"testing"

	s "github.com/stretchr/testify/suite"
)

func RunSuiteTests(t *testing.T, tests []s.TestingSuite) {
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

// region CommonTest
type CommonTest struct {
	s.Suite
}

func (s *CommonTest) SetupSuite(suiteName, testName string) {}

func (s *CommonTest) TearDownSuite(suiteName, testName string) {}

func (s *CommonTest) BeforeTest(suiteName, testName string) {}

func (s *CommonTest) AfterTest(_, _ string) {}

func (s *CommonTest) SetupSubTest() {}

func (s *CommonTest) TearDownSubTest() {}

func (s *CommonTest) Log(args ...any) { s.T().Log(args...) }

func (s *CommonTest) Logf(format string, args ...any) { s.T().Logf(format, args...) }

// SkipShort skips the test if running in short mode
func (s *CommonTest) SkipShort() {
	if testing.Short() {
		s.T().Skip("Skipping trace test in short mode")
	}
}

// ValidateRom validates the NES ROM format and returns the PRG ROM size
func (s *CommonTest) ValidateRom(romData []byte) (prgRomSize int) {
	// NES ROM format: 16-byte header + PRG ROM + CHR ROM
	// Header format: NES\x1a + 16KB PRG banks + 8KB CHR banks + flags...
	s.Greater(len(romData), 16)
	s.Equal([]uint8{0x4E, 0x45, 0x53, 0x1A}, romData[0:4])

	prgRomSize = int(romData[4]) * 16384 // Number of 16KB banks
	s.Greater(len(romData), 16+prgRomSize)
	return
}

// endregion CommonTest
