package rp2ago3

import (
	"testing"

	s "github.com/stretchr/testify/suite"
	"github.com/txc/nesgo/pkg/utils"
)

func TestSuite(t *testing.T) {
	tests := []s.TestingSuite{
		&apuTest{},
		&cpuTest{},
	}
	utils.RunSuiteTests(t, tests)
}
