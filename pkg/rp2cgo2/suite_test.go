package rp2cgo2

import (
	"testing"

	s "github.com/stretchr/testify/suite"
	"github.com/txc/nesgo/pkg/utils"
)

func TestSuite(t *testing.T) {
	tests := []s.TestingSuite{
		&oamTest{},
		&ntscTest{},
		&palTest{},
	}
	utils.RunSuiteTests(t, tests)
}
