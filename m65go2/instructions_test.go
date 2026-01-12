package m65go2

import (
	"testing"

	s "github.com/stretchr/testify/suite"
)

func TestInstructionsSuite(t *testing.T) {
	s.Run(t, &instTest{})
}

// region instTest
type instTest struct {
	s.Suite
	cpu *M6502
	// Subtest Instruction To Test
	_instruction *Instruction
	// Previous log prefix
	_oldPrefix string
}

func (s *instTest) BeforeTest(suiteName, testName string) {
	s.cpu = NewM6502(NewBasicMemory(DefaultMemorySize))
	s.cpu.Reset()
	s.cpu.breakError = true
}

func (s *instTest) AfterTest(_, _ string) {
	s.cpu = nil
}

func (s *instTest) SetupSubTest() {
	s.cpu = NewM6502(NewBasicMemory(DefaultMemorySize))
	s.cpu.Reset()
	s.cpu.breakError = true
}

func (s *instTest) TearDownSubTest() {
}

// endregion instTest

/*
// BadOpCodeError
func (s *instTest) TestBadOpCodeError() {
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x02)

	_, err := s.cpu.Execute()
	s.ErrorContains(err, BadOpCodeError.Error())
}
*/

// BadOpCodeError

func (s *instTest) TestKil() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x02)

	_, err := s.cpu.Execute()

	s.NoError(err)
	s.True(s.cpu.freezeError)

}

// LDA

func (s *instTest) TestLdaImmediate() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa9)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestLdaZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa5)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestLdaZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb5)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestLdaAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xad)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestLdaAbsoluteX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xbd)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0xff)

	cycles, _ := s.cpu.Execute()

	s.Equal(uint16(4), cycles)
	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xbd)
	s.cpu.Memory.Store(0x0101, 0xff)
	s.cpu.Memory.Store(0x0102, 0x02)
	s.cpu.Memory.Store(0x0300, 0xff)

	cycles, _ = s.cpu.Execute()

	s.Equal(uint16(5), cycles)

}

func (s *instTest) TestLdaAbsoluteY() {

	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb9)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0xff)

	cycles, _ := s.cpu.Execute()

	s.Equal(uint16(4), cycles)
	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb9)
	s.cpu.Memory.Store(0x0101, 0xff)
	s.cpu.Memory.Store(0x0102, 0x02)
	s.cpu.Memory.Store(0x0300, 0xff)

	cycles, _ = s.cpu.Execute()

	s.Equal(uint16(5), cycles)

}

func (s *instTest) TestLdaIndirectX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa1)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x87)
	s.cpu.Memory.Store(0x0086, 0x00)
	s.cpu.Memory.Store(0x0087, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestLdaIndirectY() {

	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb1)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x86)
	s.cpu.Memory.Store(0x0085, 0x00)
	s.cpu.Memory.Store(0x0087, 0xff)

	cycles, _ := s.cpu.Execute()

	s.Equal(uint16(5), cycles)
	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb1)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)
	s.cpu.Memory.Store(0x0085, 0x02)
	s.cpu.Memory.Store(0x0300, 0xff)

	cycles, _ = s.cpu.Execute()

	s.Equal(uint16(6), cycles)

}

func (s *instTest) TestLdaZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa9)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestLdaZFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa9)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestLdaNFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa9)
	s.cpu.Memory.Store(0x0101, 0x81)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestLdaNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa9)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// LDX

func (s *instTest) TestLdxImmediate() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa2)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *instTest) TestLdxZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *instTest) TestLdxZeroPageY() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *instTest) TestLdxAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xae)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *instTest) TestLdxAbsoluteY() {

	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xbe)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *instTest) TestLdxZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa2)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestLdxZFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa2)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestLdxNFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa2)
	s.cpu.Memory.Store(0x0101, 0x81)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestLdxNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa2)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// LDY

func (s *instTest) TestLdyImmediate() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa0)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register Y is not 0xff")

}

func (s *instTest) TestLdyZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa4)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register Y is not 0xff")

}

func (s *instTest) TestLdyZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb4)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register Y is not 0xff")

}

func (s *instTest) TestLdyAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xac)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register Y is not 0xff")

}

func (s *instTest) TestLdyAbsoluteX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xbc)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register Y is not 0xff")

}

func (s *instTest) TestLdyZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa0)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestLdyZFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestLdyNFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa0)
	s.cpu.Memory.Store(0x0101, 0x81)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestLdyNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// STA

func (s *instTest) TestStaZeroPage() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x85)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *instTest) TestStaZeroPageX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x95)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *instTest) TestStaAbsolute() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x8d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *instTest) TestStaAbsoluteX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x9d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *instTest) TestStaAbsoluteY() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x99)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *instTest) TestStaIndirectX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x81)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x87)
	s.cpu.Memory.Store(0x0086, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0087), "Memory is not 0xff")

}

func (s *instTest) TestStaIndirectY() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x91)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x86)
	s.cpu.Memory.Store(0x0085, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0087), "Memory is not 0xff")

}

// STX

func (s *instTest) TestStxZeroPage() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x86)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *instTest) TestStxZeroPageY() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x96)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *instTest) TestStxAbsolute() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x8e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

// STY

func (s *instTest) TestStyZeroPage() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x84)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *instTest) TestStyZeroPageY() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x94)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *instTest) TestStyAbsolute() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x8c)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

// TAX

func (s *instTest) TestTax() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xaa)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register is not 0xff")

}

func (s *instTest) TestTaxZFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xaa)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestTaxZFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xaa)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestTaxNFlagSet() {

	s.cpu.Registers.A = 0x81
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xaa)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestTaxNFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xaa)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// TAY

func (s *instTest) TestTay() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa8)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register is not 0xff")

}

// TXA

func (s *instTest) TestTxa() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x8a)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register is not 0xff")

}

// TYA

func (s *instTest) TestTya() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x98)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register is not 0xff")

}

// TSX

func (s *instTest) TestTsx() {

	s.cpu.Registers.SP = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xba)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register is not 0xff")

}

// TXS

func (s *instTest) TestTxs() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x9a)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.SP, "Register is not 0xff")

}

// PHA

func (s *instTest) TestPha() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x48)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.pull(), "Memory is not 0xff")

}

// PHP

func (s *instTest) TestPhp() {

	s.cpu.Registers.P = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x08)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.pull(), "Memory is not 0xff")

}

// PLA

func (s *instTest) TestPla() {

	s.cpu.Registers.PC = 0x0100
	s.cpu.push(0xff)

	s.cpu.Memory.Store(0x0100, 0x68)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Memory is not 0xff")

}

func (s *instTest) TestPlaZFlagSet() {

	s.cpu.push(0x00)
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x68)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestPlaZFlagUnset() {

	s.cpu.push(0x01)
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x68)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestPlaNFlagSet() {

	s.cpu.push(0x81)
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x68)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestPlaNFlagUnset() {

	s.cpu.push(0x01)
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x68)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// PLP

func (s *instTest) TestPlp() {

	s.cpu.Registers.PC = 0x0100
	s.cpu.push(0xff)

	s.cpu.Memory.Store(0x0100, 0x28)

	s.cpu.Execute()

	// PLP should mask out B but always set U (bit 5)
	// 0xff & ~B | U = 0xff & ~0x10 | 0x20 = 0xef | 0x20 = 0xef
	s.Equal(uint8(0xef), uint8(s.cpu.Registers.P), "Status is not 0xef")
}

// AND

func (s *instTest) TestAndImmediate() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x29)
	s.cpu.Memory.Store(0x0101, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *instTest) TestAndZeroPage() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x25)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *instTest) TestAndZeroPageX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x35)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *instTest) TestAndAbsolute() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *instTest) TestAndAbsoluteX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x3d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *instTest) TestAndAbsoluteY() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x39)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *instTest) TestAndIndirectX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x21)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x87)
	s.cpu.Memory.Store(0x0086, 0x00)
	s.cpu.Memory.Store(0x0087, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *instTest) TestAndIndirectY() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x31)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x86)
	s.cpu.Memory.Store(0x0085, 0x00)
	s.cpu.Memory.Store(0x0087, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *instTest) TestAndZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x29)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestAndZFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x29)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestAndNFlagSet() {

	s.cpu.Registers.A = 0x81
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x29)
	s.cpu.Memory.Store(0x0101, 0x81)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestAndNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x29)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// EOR

func (s *instTest) TestEorImmediate() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x49)
	s.cpu.Memory.Store(0x0101, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *instTest) TestEorZeroPage() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x45)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *instTest) TestEorZeroPageX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x55)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *instTest) TestEorAbsolute() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *instTest) TestEorAbsoluteX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x5d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *instTest) TestEorAbsoluteY() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x59)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *instTest) TestEorIndirectX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x41)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x87)
	s.cpu.Memory.Store(0x0086, 0x00)
	s.cpu.Memory.Store(0x0087, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *instTest) TestEorIndirectY() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x51)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x86)
	s.cpu.Memory.Store(0x0085, 0x00)
	s.cpu.Memory.Store(0x0087, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *instTest) TestEorZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x49)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestEorZFlagUnset() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x49)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestEorNFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x49)
	s.cpu.Memory.Store(0x0101, 0x81)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestEorNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x49)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// ORA

func (s *instTest) TestOraImmediate() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x09)
	s.cpu.Memory.Store(0x0101, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestOraZeroPage() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x05)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestOraZeroPageX() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x15)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestOraAbsolute() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestOraAbsoluteX() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x1d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestOraAbsoluteY() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x19)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestOraIndirectX() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x01)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x87)
	s.cpu.Memory.Store(0x0086, 0x00)
	s.cpu.Memory.Store(0x0087, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestOraIndirectY() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x11)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x86)
	s.cpu.Memory.Store(0x0085, 0x00)
	s.cpu.Memory.Store(0x0087, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *instTest) TestOraZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x09)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestOraZFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x09)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestOraNFlagSet() {

	s.cpu.Registers.A = 0x81
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x09)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestOraNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x09)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// BIT

func (s *instTest) TestBitZeroPage() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x7f)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *instTest) TestBitAbsolute() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2c)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x7f)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *instTest) TestBitNFlagSet() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestBitNFlagUnset() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x7f)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *instTest) TestBitVFlagSet() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&V, "V flag is not set")

}

func (s *instTest) TestBitVFlagUnset() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x3f)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&V, "V flag is set")

}

func (s *instTest) TestBitZFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestBitZFlagUnset() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x3f)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

// ADC

func (s *instTest) TestAdcImmediate() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x03), s.cpu.Registers.A, "Register A is not 0x03")

	s.cpu.Registers.P |= D
	s.cpu.Registers.A = 0x29 // BCD
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x11) // BCD

	s.cpu.Execute()

	s.Equal(uint8(0x40), s.cpu.Registers.A, "Register A is not 0x40") // BCD

	s.cpu.Registers.P |= D
	s.cpu.Registers.A = 0x29 | uint8(N) // BCD
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x29) // BCD

	s.cpu.Execute()

	s.Equal(uint8(0x38), s.cpu.Registers.A, "Register A is not 0x38") // BCD

	s.cpu.Registers.P |= D
	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x58 // BCD
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x46) // BCD

	s.cpu.Execute()

	s.Equal(uint8(0x05), s.cpu.Registers.A, "Register A is not 0x05") // BCD

}

func (s *instTest) TestAdcZeroPage() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x65)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x03), s.cpu.Registers.A, "Register A is not 0x03")

}

func (s *instTest) TestAdcZeroPageX() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x75)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x03), s.cpu.Registers.A, "Register A is not 0x03")

}

func (s *instTest) TestAdcAbsolute() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x03), s.cpu.Registers.A, "Register A is not 0x03")

}

func (s *instTest) TestAdcAbsoluteX() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x7d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x03), s.cpu.Registers.A, "Register A is not 0x03")

}

func (s *instTest) TestAdcAbsoluteY() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x79)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x03), s.cpu.Registers.A, "Register A is not 0x03")

}

func (s *instTest) TestAdcIndirectX() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x61)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x87)
	s.cpu.Memory.Store(0x0086, 0x00)
	s.cpu.Memory.Store(0x0087, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x03), s.cpu.Registers.A, "Register A is not 0x03")

}

func (s *instTest) TestAdcIndirectY() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x71)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x86)
	s.cpu.Memory.Store(0x0085, 0x00)
	s.cpu.Memory.Store(0x0087, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x03), s.cpu.Registers.A, "Register A is not 0x03")

}

func (s *instTest) TestAdcCFlagSet() {

	s.cpu.Registers.A = 0xff // -1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0xff // -1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x00) // +0

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *instTest) TestAdcCFlagUnset() {

	s.cpu.Registers.A = 0x00 // +0
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

	s.cpu.Registers.P &^= C
	s.cpu.Registers.A = 0xff // -1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x00) // +0

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

func (s *instTest) TestAdcZFlagSet() {

	s.cpu.Registers.A = 0x00 // +0
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x00) // +0

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0xfe // -2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestAdcZFlagUnset() {

	s.cpu.Registers.A = 0x00 // +0
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0xff) // -1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

	s.cpu.Registers.A = 0xfe // -2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestAdcVFlagSet() {

	s.cpu.Registers.A = 0x7f // +127
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&V, "V flag is not set")

}

func (s *instTest) TestAdcVFlagUnset() {

	s.cpu.Registers.A = 0x01 // +1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&V, "V flag is set")

}

func (s *instTest) TestAdcNFlagSet() {

	s.cpu.Registers.A = 0x01 // +1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0xfe) // -2

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestAdcNFlagUnset() {

	s.cpu.Registers.A = 0x01 // +1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// SBC

func (s *instTest) TestSbcImmediate() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.A, "Register A is not 0x01")

	s.cpu.Registers.P |= D
	s.cpu.Registers.A = 0x29 // BCD
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x11) // BCD

	s.cpu.Execute()

	s.Equal(uint8(0x18), s.cpu.Registers.A, "Register A is not 0x18") // BCD

}

func (s *instTest) TestSbcZeroPage() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe5)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x01)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.A, "Register A is not 0x01")

}

func (s *instTest) TestSbcZeroPageX() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x02
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xf5)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x01)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.A, "Register A is not 0x01")

}

func (s *instTest) TestSbcAbsolute() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xed)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x01)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.A, "Register A is not 0x01")

}

func (s *instTest) TestSbcAbsoluteX() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x02
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xfd)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x01)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.A, "Register A is not 0x01")

}

func (s *instTest) TestSbcAbsoluteY() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x02
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xf9)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x01)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.A, "Register A is not 0x01")

}

func (s *instTest) TestSbcIndirectX() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x02
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe1)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x87)
	s.cpu.Memory.Store(0x0086, 0x00)
	s.cpu.Memory.Store(0x0087, 0x01)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.A, "Register A is not 0x01")

}

func (s *instTest) TestSbcIndirectY() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x02
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xf1)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x86)
	s.cpu.Memory.Store(0x0085, 0x00)
	s.cpu.Memory.Store(0x0087, 0x01)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.A, "Register A is not 0x01")

}

func (s *instTest) TestSbcCFlagSet() {

	s.cpu.Registers.A = 0xc4 // -60
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x3c) // +60

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *instTest) TestSbcCFlagUnset() {

	s.cpu.Registers.A = 0x02 // +2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x04) // +4

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

func (s *instTest) TestSbcZFlagSet() {

	s.cpu.Registers.A = 0x02 // +2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestSbcZFlagUnset() {

	s.cpu.Registers.A = 0x02 // +2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x02) // +2

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestSbcVFlagSet() {

	s.cpu.Registers.A = 0x80 // -128
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&V, "V flag is not set")

}

func (s *instTest) TestSbcVFlagUnset() {

	s.cpu.Registers.A = 0x01 // +1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&V, "V flag is set")

}

func (s *instTest) TestSbcNFlagSet() {

	s.cpu.Registers.A = 0xfd // -3
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestSbcNFlagUnset() {

	s.cpu.Registers.A = 0x02 // +2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// CMP

func (s *instTest) TestCmpImmediate() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCmpZeroPage() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc5)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCmpZeroPageX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xd5)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCmpAbsolute() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xcd)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCmpAbsoluteX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xdd)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCmpAbsoluteY() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xd9)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCmpIndirectX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc1)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x87)
	s.cpu.Memory.Store(0x0086, 0x00)
	s.cpu.Memory.Store(0x0087, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCmpIndirectY() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xd1)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x86)
	s.cpu.Memory.Store(0x0085, 0x00)
	s.cpu.Memory.Store(0x0087, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCmpNFlagSet() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestCmpNFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *instTest) TestCmpZFlagSet() {

	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

	s.cpu.Registers.A = 0xfe // -2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0xfe) // -2

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCmpZFlagUnset() {

	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

	s.cpu.Registers.A = 0xfe // -2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0xff) // -1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestCmpCFlagSet() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

	s.cpu.Registers.A = 0xfe // -2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0xfd) // -3

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *instTest) TestCmpCFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

	s.cpu.Registers.A = 0xfd // -3
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0xfe) // -2

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

// CPX

func (s *instTest) TestCpxImmediate() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCpxZeroPage() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe4)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCpxAbsolute() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xec)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCpxNFlagSet() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestCpxNFlagUnset() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *instTest) TestCpxZFlagSet() {

	s.cpu.Registers.X = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCpxZFlagUnset() {

	s.cpu.Registers.X = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestCpxCFlagSet() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *instTest) TestCpxCFlagUnset() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

// CPY

func (s *instTest) TestCpyImmediate() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCpyZeroPage() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc4)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCpyAbsolute() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xcc)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCpyNFlagSet() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestCpyNFlagUnset() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *instTest) TestCpyZFlagSet() {

	s.cpu.Registers.Y = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestCpyZFlagUnset() {

	s.cpu.Registers.Y = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestCpyCFlagSet() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *instTest) TestCpyCFlagUnset() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

// INC

func (s *instTest) TestIncZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xfe)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *instTest) TestIncZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xf6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0xfe)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *instTest) TestIncAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xee)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xfe)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *instTest) TestIncAbsoluteX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xfe)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0xfe)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *instTest) TestIncZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff) // -1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestIncZFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x00)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestIncNFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xfe) // -2

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestIncNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x00)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// INX

func (s *instTest) TestInx() {

	s.cpu.Registers.X = 0xfe
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe8)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *instTest) TestInxZFlagSet() {

	s.cpu.Registers.X = 0xff // -1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe8)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestInxZFlagUnset() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestInxNFlagSet() {

	s.cpu.Registers.X = 0xfe // -2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe8)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestInxNFlagUnset() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// INY

func (s *instTest) TestIny() {

	s.cpu.Registers.Y = 0xfe // -2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc8)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register X is not 0xff")

}

func (s *instTest) TestInyZFlagSet() {

	s.cpu.Registers.Y = 0xff // -1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc8)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestInyZFlagUnset() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestInyNFlagSet() {

	s.cpu.Registers.Y = 0xfe // -2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc8)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestInyNFlagUnset() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// DEC

func (s *instTest) TestDecZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x01")

}

func (s *instTest) TestDecZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xd6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x01")

}

func (s *instTest) TestDecAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xce)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x01")

}

func (s *instTest) TestDecAbsoluteX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xde)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x01")

}

func (s *instTest) TestDecZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x01)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestDecZFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestDecNFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestDecNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// DEX

func (s *instTest) TestDex() {

	s.cpu.Registers.X = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xca)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.X, "Register X is not 0x01")

}

func (s *instTest) TestDexZFlagSet() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xca)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestDexZFlagUnset() {

	s.cpu.Registers.X = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xca)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestDexNFlagSet() {

	s.cpu.Registers.X = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xca)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestDexNFlagUnset() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xca)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// DEY

func (s *instTest) TestDey() {

	s.cpu.Registers.Y = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x88)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.Y, "Register X is not 0x01")

}

func (s *instTest) TestDeyZFlagSet() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x88)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestDeyZFlagUnset() {

	s.cpu.Registers.Y = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x88)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestDeyNFlagSet() {

	s.cpu.Registers.Y = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x88)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestDeyNFlagUnset() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x88)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// ASL

func (s *instTest) TestAslAccumulator() {

	s.cpu.Registers.A = 0x2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.Equal(uint8(0x04), s.cpu.Registers.A, "Register A is not 0x04")

}

func (s *instTest) TestAslZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x06)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x04), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x04")

}

func (s *instTest) TestAslZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x16)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x04), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x04")

}

func (s *instTest) TestAslAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x04), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x04")

}

func (s *instTest) TestAslAbsoluteX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x1e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x04), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x04")

}

func (s *instTest) TestAslCFlagSet() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *instTest) TestAslCFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

func (s *instTest) TestAslZFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestAslZFlagUnset() {

	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestAslNFlagSet() {

	s.cpu.Registers.A = 0xfe
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestAslNFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// LSR

func (s *instTest) TestLsrAccumulator() {

	s.cpu.Registers.A = 0x2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.A, "Register A is not 0x01")

}

func (s *instTest) TestLsrZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x46)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x01")

}

func (s *instTest) TestLsrZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x56)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x01")

}

func (s *instTest) TestLsrAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x01")

}

func (s *instTest) TestLsrAbsoluteX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x5e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x01")

}

func (s *instTest) TestLsrCFlagSet() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *instTest) TestLsrCFlagUnset() {

	s.cpu.Registers.A = 0x10
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

func (s *instTest) TestLsrZFlagSet() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestLsrZFlagUnset() {

	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

// func (s *instTest) TestLsrNFlagSet() { }
// not tested, N bit always set to 0

func (s *instTest) TestLsrNFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// ROL

func (s *instTest) TestRolAccumulator() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.Equal(uint8(0x05), s.cpu.Registers.A, "Register A is not 0x05")

}

func (s *instTest) TestRolZeroPage() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x26)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x05), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x05")

}

func (s *instTest) TestRolZeroPageX() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x36)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x05), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x05")

}

func (s *instTest) TestRolAbsolute() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x05), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x05")

}

func (s *instTest) TestRolAbsoluteX() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x3e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x05), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x05")

}

func (s *instTest) TestRolCFlagSet() {

	s.cpu.Registers.A = 0x80
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *instTest) TestRolCFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

func (s *instTest) TestRolZFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestRolZFlagUnset() {

	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestRolNFlagSet() {

	s.cpu.Registers.A = 0xfe
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestRolNFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// ROR

func (s *instTest) TestRorAccumulator() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x08
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.Equal(uint8(0x84), s.cpu.Registers.A, "Register A is not 0x84")

}

func (s *instTest) TestRorZeroPage() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x66)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x08)

	s.cpu.Execute()

	s.Equal(uint8(0x84), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x84")

}

func (s *instTest) TestRorZeroPageX() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x76)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x08)

	s.cpu.Execute()

	s.Equal(uint8(0x84), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x84")

}

func (s *instTest) TestRorAbsolute() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x08)

	s.cpu.Execute()

	s.Equal(uint8(0x84), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x84")

}

func (s *instTest) TestRorAbsoluteX() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x7e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x08)

	s.cpu.Execute()

	s.Equal(uint8(0x84), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x84")

}

func (s *instTest) TestRorCFlagSet() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *instTest) TestRorCFlagUnset() {

	s.cpu.Registers.A = 0x10
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

func (s *instTest) TestRorZFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *instTest) TestRorZFlagUnset() {

	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *instTest) TestRorNFlagSet() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0xfe
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *instTest) TestRorNFlagUnset() {

	s.cpu.Registers.P &^= C
	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// JMP

func (s *instTest) TestJmpAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4c)
	s.cpu.Memory.Store(0x0101, 0xff)
	s.cpu.Memory.Store(0x0102, 0x01)

	s.cpu.Execute()

	s.Equal(uint16(0x01ff), s.cpu.Registers.PC, "Register PC is not 0x01ff")

}

func (s *instTest) TestJmpIndirect() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6c)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x01)
	s.cpu.Memory.Store(0x0184, 0xff)
	s.cpu.Memory.Store(0x0185, 0xff)

	s.cpu.Execute()

	s.Equal(uint16(0xffff), s.cpu.Registers.PC, "Register PC is not 0xffff")

}

// JSR

func (s *instTest) TestJsr() {

	s.Run("JSR Test 1", func() {
		s.cpu.Registers.PC = 0x0100

		s.cpu.Memory.Store(0x0100, 0x20)
		s.cpu.Memory.Store(0x0101, 0xff)
		s.cpu.Memory.Store(0x0102, 0x01)

		s.cpu.Execute()

		s.Equal(uint16(0x01ff), s.cpu.Registers.PC, "Register PC is not 0x01ff")
		s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x01fd), "Memory is not 0x01")
		s.Equal(uint8(0x02), s.cpu.Memory.Fetch(0x01fc), "Memory is not 0x02")
	})

	s.Run("JSR Test 2", func() {
		s.cpu.Registers.PC = 0x0100

		s.cpu.Memory.Store(0x0100, 0x20) // JSR
		s.cpu.Memory.Store(0x0101, 0x84)
		s.cpu.Memory.Store(0x0102, 0x00)
		s.cpu.Memory.Store(0x0084, 0x60) // RTS

		s.cpu.Execute()
		s.cpu.Execute()

		s.Equal(uint16(0x0103), s.cpu.Registers.PC, "Register PC is not 0x0103")
		s.Equal(uint8(0xfd), s.cpu.Registers.SP, "Register SP is not 0xfd")
	})

	s.Run("JSR Test 3", func() {
		s.cpu.Registers.PC = 0x0100

		s.cpu.Memory.Store(0x0100, 0x20) // JSR $0084
		s.cpu.Memory.Store(0x0101, 0x84)
		s.cpu.Memory.Store(0x0102, 0x00)
		s.cpu.Memory.Store(0x0103, 0xa9) // LDA #$ff
		s.cpu.Memory.Store(0x0104, 0xff)
		s.cpu.Memory.Store(0x0105, 0x02) // illegal opcode
		s.cpu.Memory.Store(0x0084, 0x60) // RTS

		s.cpu.Run()

		s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")
	})
}

// RTS

func (s *instTest) TestRts() {

	s.cpu.Registers.PC = 0x0100
	s.cpu.push16(0x0102)

	s.cpu.Memory.Store(0x0100, 0x60)

	s.cpu.Execute()

	s.Equal(uint16(0x0103), s.cpu.Registers.PC, "Register PC is not 0x0103")

}

// BCC

func (s *instTest) TestBcc() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x90)

	cycles, _ := s.cpu.Execute()

	s.Equal(uint16(2), cycles)
	s.Equal(uint16(0x0102), s.cpu.Registers.PC)

	s.cpu.Registers.P &^= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x90)
	s.cpu.Memory.Store(0x0101, 0x02) // +2

	cycles, _ = s.cpu.Execute()

	s.Equal(uint16(3), cycles)
	s.Equal(uint16(0x0104), s.cpu.Registers.PC)

	s.cpu.Registers.P &^= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x90)
	s.cpu.Memory.Store(0x0101, 0xfd) // -3

	cycles, _ = s.cpu.Execute()

	s.Equal(uint16(4), cycles)
	s.Equal(uint16(0x00ff), s.cpu.Registers.PC)

}

// BCS

func (s *instTest) TestBcs() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb0)
	s.cpu.Memory.Store(0x0101, 0x02) // +2

	s.cpu.Execute()

	s.Equal(uint16(0x0104), s.cpu.Registers.PC)

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb0)
	s.cpu.Memory.Store(0x0101, 0xfe) // -2

	s.cpu.Execute()

	s.Equal(uint16(0x0100), s.cpu.Registers.PC)

}

// BEQ

func (s *instTest) TestBeq() {

	s.cpu.Registers.P |= Z
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xf0)
	s.cpu.Memory.Store(0x0101, 0x02) // +2

	s.cpu.Execute()

	s.Equal(uint16(0x0104), s.cpu.Registers.PC)

	s.cpu.Registers.P |= Z
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xf0)
	s.cpu.Memory.Store(0x0101, 0xfe) // -2

	s.cpu.Execute()

	s.Equal(uint16(0x0100), s.cpu.Registers.PC)

}

// BMI

func (s *instTest) TestBmi() {

	s.cpu.Registers.P |= N
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x30)
	s.cpu.Memory.Store(0x0101, 0x02) // +2

	s.cpu.Execute()

	s.Equal(uint16(0x0104), s.cpu.Registers.PC)

	s.cpu.Registers.P |= N
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x30)
	s.cpu.Memory.Store(0x0101, 0xfe) // -2

	s.cpu.Execute()

	s.Equal(uint16(0x0100), s.cpu.Registers.PC)

}

// BNE

func (s *instTest) TestBne() {

	s.cpu.Registers.P &^= Z
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xd0)
	s.cpu.Memory.Store(0x0101, 0x02) // +2

	s.cpu.Execute()

	s.Equal(uint16(0x0104), s.cpu.Registers.PC)

	s.cpu.Registers.P &^= Z
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xd0)
	s.cpu.Memory.Store(0x0101, 0xfe) // -2

	s.cpu.Execute()

	s.Equal(uint16(0x0100), s.cpu.Registers.PC)

}

// BPL

func (s *instTest) TestBpl() {

	s.cpu.Registers.P &^= N
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x10)
	s.cpu.Memory.Store(0x0101, 0x02) // +2

	s.cpu.Execute()

	s.Equal(uint16(0x0104), s.cpu.Registers.PC)

	s.cpu.Registers.P &^= N
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x10)
	s.cpu.Memory.Store(0x0101, 0xfe) // -2

	s.cpu.Execute()

	s.Equal(uint16(0x0100), s.cpu.Registers.PC)

}

// BVC

func (s *instTest) TestBvc() {

	s.cpu.Registers.P &^= V
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x50)
	s.cpu.Memory.Store(0x0101, 0x02) // +2

	s.cpu.Execute()

	s.Equal(uint16(0x0104), s.cpu.Registers.PC)

	s.cpu.Registers.P &^= V
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x50)
	s.cpu.Memory.Store(0x0101, 0xfe) // -2

	s.cpu.Execute()

	s.Equal(uint16(0x0100), s.cpu.Registers.PC)

}

// BVS

func (s *instTest) TestBvs() {

	s.cpu.Registers.P |= V
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x70)
	s.cpu.Memory.Store(0x0101, 0x02) // +2

	s.cpu.Execute()

	s.Equal(uint16(0x0104), s.cpu.Registers.PC)

	s.cpu.Registers.P |= V
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x70)
	s.cpu.Memory.Store(0x0101, 0xfe) // -2

	s.cpu.Execute()

	s.Equal(uint16(0x0100), s.cpu.Registers.PC)

}

// CLC

func (s *instTest) TestClc() {

	s.cpu.Registers.P &^= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x18)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x18)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

// CLD

func (s *instTest) TestCld() {

	s.cpu.Registers.P &^= D
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xd8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&D, "D flag is set")

	s.cpu.Registers.P |= D
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xd8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&D, "D flag is set")

}

// CLI

func (s *instTest) TestCli() {

	s.cpu.Registers.P &^= I
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x58)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&I, "I flag is set")

	s.cpu.Registers.P |= I
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x58)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&I, "I flag is set")

}

// CLV

func (s *instTest) TestClv() {

	s.cpu.Registers.P &^= V
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&V, "V flag is set")

	s.cpu.Registers.P |= V
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&V, "V flag is set")

}

// SEC

func (s *instTest) TestSec() {

	s.cpu.Registers.P &^= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x38)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x38)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

// SED

func (s *instTest) TestSed() {

	s.cpu.Registers.P &^= D
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xf8)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&D, "D flag is not set")

	s.cpu.Registers.P |= D
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xf8)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&D, "D flag is not set")

}

// SEI

func (s *instTest) TestSei() {

	s.cpu.Registers.P &^= I
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x78)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&I, "I flag is not set")

	s.cpu.Registers.P |= I
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x78)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&I, "I flag is not set")

}

// BRK

func (s *instTest) TestBrk() {

	s.cpu.Registers.P = 0xff & (^B)
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x00)
	s.cpu.Memory.Store(0xfffe, 0xff)
	s.cpu.Memory.Store(0xffff, 0x01)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.pull(), "Memory is not 0xff")

	s.Equal(uint16(0x0102), s.cpu.pull16(), "Memory is not 0x0102")

	s.Equal(uint16(0x01ff), s.cpu.Registers.PC, "Register PC is not 0x01ff")

}

// RTI

func (s *instTest) TestRti() {

	s.cpu.Registers.PC = 0x0100
	s.cpu.push16(0x0102)
	s.cpu.push(0x03)

	s.cpu.Memory.Store(0x0100, 0x40)

	s.cpu.Execute()

	// RTI should mask out B but always set U (bit 5)
	// 0x03 & ~B | U = 0x03 & ~0x10 | 0x20 = 0x03 | 0x20 = 0x23
	s.Equal(Status(0x23), s.cpu.Registers.P, "Register P is not 0x23")
	s.Equal(uint16(0x0102), s.cpu.Registers.PC, "Register PC is not 0x0102")
}

// Irq

func (s *instTest) TestIrq() {

	s.cpu.Registers.P = 0xfb
	s.cpu.Registers.PC = 0x0100

	s.cpu.Interrupt(Irq, true)
	s.cpu.Memory.Store(0xfffe, 0x40)
	s.cpu.Memory.Store(0xffff, 0x01)

	s.cpu.PerformInterrupts()

	s.Equal(uint8(0xeb), s.cpu.pull(), "Memory is not 0xeb")
	s.Equal(uint16(0x0100), s.cpu.pull16(), "Memory is not 0x0100")
	s.Equal(uint16(0x0140), s.cpu.Registers.PC, "Register PC is not 0x0140")
	s.False(s.cpu.GetInterrupt(Irq), "Interrupt is set")
}

// Nmi

func (s *instTest) TestNmi() {

	s.cpu.Registers.P = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Interrupt(Nmi, true)
	s.cpu.Memory.Store(0xfffa, 0x40)
	s.cpu.Memory.Store(0xfffb, 0x01)

	s.cpu.PerformInterrupts()

	s.Equal(uint8(0xef), s.cpu.pull(), "Memory is not 0xef")
	s.Equal(uint16(0x0100), s.cpu.pull16(), "Memory is not 0x0100")
	s.Equal(uint16(0x0140), s.cpu.Registers.PC, "Register PC is not 0x0140")
	s.False(s.cpu.GetInterrupt(Nmi), "Interrupt is set")
}

// Rst

func (s *instTest) TestRst() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Interrupt(Rst, true)
	s.cpu.Memory.Store(0xfffc, 0x40)
	s.cpu.Memory.Store(0xfffd, 0x01)

	s.cpu.PerformInterrupts()

	s.Equal(uint16(0x0140), s.cpu.Registers.PC, "Register PC is not 0x0140")
	s.False(s.cpu.GetInterrupt(Rst), "Interrupt is set")
}
