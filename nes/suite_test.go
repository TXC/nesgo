package nes

import (
	"testing"

	s "github.com/stretchr/testify/suite"
	"github.com/txc/nesgo/utils"
)

func TestSuite(t *testing.T) {
	tests := []s.TestingSuite{
		&ControllersTestSuite{},
		&mirroringTest{},
		&ROMTestSuite{},
	}
	utils.RunSuiteTests(t, tests)
}
