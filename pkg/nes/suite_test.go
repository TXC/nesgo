package nes

import (
	"testing"

	s "github.com/stretchr/testify/suite"
	"github.com/txc/nesgo/pkg/utils"
)

func TestSuite(t *testing.T) {
	tests := []s.TestingSuite{
		&ControllersTestSuite{},
		&mirroringTest{},
		&iNESRomTestSuite{},
		&NES20RomTestSuite{},
	}
	utils.RunSuiteTests(t, tests)
}
