package m65go2

import (
	s "github.com/stretchr/testify/suite"
)

// region InstructionsTest
type InstructionsTest struct {
	s.Suite
	cpu *M6502
	// Subtest Instruction To Test
	_instruction *Instruction
	// Previous log prefix
	_oldPrefix string
}

func (s *InstructionsTest) BeforeTest(suiteName, testName string) {
	s.cpu = NewM6502(NewBasicMemory(DefaultMemorySize))
	s.cpu.Reset()
	s.cpu.breakError = true
}

func (s *InstructionsTest) AfterTest(_, _ string) {
	s.cpu = nil
}

func (s *InstructionsTest) SetupSubTest() {
	s.cpu = NewM6502(NewBasicMemory(DefaultMemorySize))
	s.cpu.Reset()
	s.cpu.breakError = true
}

func (s *InstructionsTest) TearDownSubTest() {
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

func (s *InstructionsTest) TestKil() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x02)

	_, err := s.cpu.Execute()

	s.NoError(err)
	s.True(s.cpu.freezeError)

}

// LDA

func (s *InstructionsTest) TestLdaImmediate() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa9)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *InstructionsTest) TestLdaZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa5)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *InstructionsTest) TestLdaZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb5)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *InstructionsTest) TestLdaAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xad)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *InstructionsTest) TestLdaAbsoluteX() {

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

func (s *InstructionsTest) TestLdaAbsoluteY() {

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

func (s *InstructionsTest) TestLdaIndirectX() {

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

func (s *InstructionsTest) TestLdaIndirectY() {

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

func (s *InstructionsTest) TestLdaZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa9)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestLdaZFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa9)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestLdaNFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa9)
	s.cpu.Memory.Store(0x0101, 0x81)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestLdaNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa9)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// LDX

func (s *InstructionsTest) TestLdxImmediate() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa2)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *InstructionsTest) TestLdxZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *InstructionsTest) TestLdxZeroPageY() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *InstructionsTest) TestLdxAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xae)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *InstructionsTest) TestLdxAbsoluteY() {

	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xbe)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *InstructionsTest) TestLdxZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa2)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestLdxZFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa2)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestLdxNFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa2)
	s.cpu.Memory.Store(0x0101, 0x81)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestLdxNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa2)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// LDY

func (s *InstructionsTest) TestLdyImmediate() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa0)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register Y is not 0xff")

}

func (s *InstructionsTest) TestLdyZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa4)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register Y is not 0xff")

}

func (s *InstructionsTest) TestLdyZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xb4)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register Y is not 0xff")

}

func (s *InstructionsTest) TestLdyAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xac)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register Y is not 0xff")

}

func (s *InstructionsTest) TestLdyAbsoluteX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xbc)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register Y is not 0xff")

}

func (s *InstructionsTest) TestLdyZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa0)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestLdyZFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestLdyNFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa0)
	s.cpu.Memory.Store(0x0101, 0x81)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestLdyNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// STA

func (s *InstructionsTest) TestStaZeroPage() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x85)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *InstructionsTest) TestStaZeroPageX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x95)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *InstructionsTest) TestStaAbsolute() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x8d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *InstructionsTest) TestStaAbsoluteX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x9d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *InstructionsTest) TestStaAbsoluteY() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.Y = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x99)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *InstructionsTest) TestStaIndirectX() {

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

func (s *InstructionsTest) TestStaIndirectY() {

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

func (s *InstructionsTest) TestStxZeroPage() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x86)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *InstructionsTest) TestStxZeroPageY() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x96)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *InstructionsTest) TestStxAbsolute() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x8e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

// STY

func (s *InstructionsTest) TestStyZeroPage() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x84)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *InstructionsTest) TestStyZeroPageY() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x94)
	s.cpu.Memory.Store(0x0101, 0x84)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *InstructionsTest) TestStyAbsolute() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x8c)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

// TAX

func (s *InstructionsTest) TestTax() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xaa)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register is not 0xff")

}

func (s *InstructionsTest) TestTaxZFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xaa)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestTaxZFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xaa)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestTaxNFlagSet() {

	s.cpu.Registers.A = 0x81
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xaa)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestTaxNFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xaa)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// TAY

func (s *InstructionsTest) TestTay() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xa8)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register is not 0xff")

}

// TXA

func (s *InstructionsTest) TestTxa() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x8a)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register is not 0xff")

}

// TYA

func (s *InstructionsTest) TestTya() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x98)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register is not 0xff")

}

// TSX

func (s *InstructionsTest) TestTsx() {

	s.cpu.Registers.SP = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xba)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register is not 0xff")

}

// TXS

func (s *InstructionsTest) TestTxs() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x9a)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.SP, "Register is not 0xff")

}

// PHA

func (s *InstructionsTest) TestPha() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x48)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.pull(), "Memory is not 0xff")

}

// PHP

func (s *InstructionsTest) TestPhp() {

	s.cpu.Registers.P = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x08)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.pull(), "Memory is not 0xff")

}

// PLA

func (s *InstructionsTest) TestPla() {

	s.cpu.Registers.PC = 0x0100
	s.cpu.push(0xff)

	s.cpu.Memory.Store(0x0100, 0x68)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Memory is not 0xff")

}

func (s *InstructionsTest) TestPlaZFlagSet() {

	s.cpu.push(0x00)
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x68)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestPlaZFlagUnset() {

	s.cpu.push(0x01)
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x68)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestPlaNFlagSet() {

	s.cpu.push(0x81)
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x68)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestPlaNFlagUnset() {

	s.cpu.push(0x01)
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x68)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// PLP

func (s *InstructionsTest) TestPlp() {

	s.cpu.Registers.PC = 0x0100
	s.cpu.push(0xff)

	s.cpu.Memory.Store(0x0100, 0x28)

	s.cpu.Execute()

	// PLP should mask out B but always set U (bit 5)
	// 0xff & ~B | U = 0xff & ~0x10 | 0x20 = 0xef | 0x20 = 0xef
	s.Equal(uint8(0xef), uint8(s.cpu.Registers.P), "Status is not 0xef")
}

// AND

func (s *InstructionsTest) TestAndImmediate() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x29)
	s.cpu.Memory.Store(0x0101, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *InstructionsTest) TestAndZeroPage() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x25)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *InstructionsTest) TestAndZeroPageX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x35)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *InstructionsTest) TestAndAbsolute() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0x0f), s.cpu.Registers.A, "Register A is not 0x0f")

}

func (s *InstructionsTest) TestAndAbsoluteX() {

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

func (s *InstructionsTest) TestAndAbsoluteY() {

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

func (s *InstructionsTest) TestAndIndirectX() {

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

func (s *InstructionsTest) TestAndIndirectY() {

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

func (s *InstructionsTest) TestAndZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x29)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestAndZFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x29)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestAndNFlagSet() {

	s.cpu.Registers.A = 0x81
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x29)
	s.cpu.Memory.Store(0x0101, 0x81)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestAndNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x29)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// EOR

func (s *InstructionsTest) TestEorImmediate() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x49)
	s.cpu.Memory.Store(0x0101, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *InstructionsTest) TestEorZeroPage() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x45)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *InstructionsTest) TestEorZeroPageX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x55)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *InstructionsTest) TestEorAbsolute() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xf0), s.cpu.Registers.A, "Register A is not 0xf0")

}

func (s *InstructionsTest) TestEorAbsoluteX() {

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

func (s *InstructionsTest) TestEorAbsoluteY() {

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

func (s *InstructionsTest) TestEorIndirectX() {

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

func (s *InstructionsTest) TestEorIndirectY() {

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

func (s *InstructionsTest) TestEorZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x49)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestEorZFlagUnset() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x49)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestEorNFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x49)
	s.cpu.Memory.Store(0x0101, 0x81)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestEorNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x49)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// ORA

func (s *InstructionsTest) TestOraImmediate() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x09)
	s.cpu.Memory.Store(0x0101, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *InstructionsTest) TestOraZeroPage() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x05)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *InstructionsTest) TestOraZeroPageX() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x15)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *InstructionsTest) TestOraAbsolute() {

	s.cpu.Registers.A = 0xf0
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x0f)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.A, "Register A is not 0xff")

}

func (s *InstructionsTest) TestOraAbsoluteX() {

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

func (s *InstructionsTest) TestOraAbsoluteY() {

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

func (s *InstructionsTest) TestOraIndirectX() {

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

func (s *InstructionsTest) TestOraIndirectY() {

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

func (s *InstructionsTest) TestOraZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x09)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestOraZFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x09)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestOraNFlagSet() {

	s.cpu.Registers.A = 0x81
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x09)
	s.cpu.Memory.Store(0x0101, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestOraNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x09)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// BIT

func (s *InstructionsTest) TestBitZeroPage() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x7f)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *InstructionsTest) TestBitAbsolute() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2c)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x7f)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *InstructionsTest) TestBitNFlagSet() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestBitNFlagUnset() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x7f)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *InstructionsTest) TestBitVFlagSet() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&V, "V flag is not set")

}

func (s *InstructionsTest) TestBitVFlagUnset() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x3f)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&V, "V flag is set")

}

func (s *InstructionsTest) TestBitZFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestBitZFlagUnset() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x24)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x3f)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

// ADC

func (s *InstructionsTest) TestAdcImmediate() {

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

func (s *InstructionsTest) TestAdcZeroPage() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x65)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x03), s.cpu.Registers.A, "Register A is not 0x03")

}

func (s *InstructionsTest) TestAdcZeroPageX() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x75)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x03), s.cpu.Registers.A, "Register A is not 0x03")

}

func (s *InstructionsTest) TestAdcAbsolute() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6d)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x03), s.cpu.Registers.A, "Register A is not 0x03")

}

func (s *InstructionsTest) TestAdcAbsoluteX() {

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

func (s *InstructionsTest) TestAdcAbsoluteY() {

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

func (s *InstructionsTest) TestAdcIndirectX() {

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

func (s *InstructionsTest) TestAdcIndirectY() {

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

func (s *InstructionsTest) TestAdcCFlagSet() {

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

func (s *InstructionsTest) TestAdcCFlagUnset() {

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

func (s *InstructionsTest) TestAdcZFlagSet() {

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

func (s *InstructionsTest) TestAdcZFlagUnset() {

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

func (s *InstructionsTest) TestAdcVFlagSet() {

	s.cpu.Registers.A = 0x7f // +127
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&V, "V flag is not set")

}

func (s *InstructionsTest) TestAdcVFlagUnset() {

	s.cpu.Registers.A = 0x01 // +1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&V, "V flag is set")

}

func (s *InstructionsTest) TestAdcNFlagSet() {

	s.cpu.Registers.A = 0x01 // +1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0xfe) // -2

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestAdcNFlagUnset() {

	s.cpu.Registers.A = 0x01 // +1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x69)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// SBC

func (s *InstructionsTest) TestSbcImmediate() {

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

func (s *InstructionsTest) TestSbcZeroPage() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe5)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x01)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.A, "Register A is not 0x01")

}

func (s *InstructionsTest) TestSbcZeroPageX() {

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

func (s *InstructionsTest) TestSbcAbsolute() {

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

func (s *InstructionsTest) TestSbcAbsoluteX() {

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

func (s *InstructionsTest) TestSbcAbsoluteY() {

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

func (s *InstructionsTest) TestSbcIndirectX() {

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

func (s *InstructionsTest) TestSbcIndirectY() {

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

func (s *InstructionsTest) TestSbcCFlagSet() {

	s.cpu.Registers.A = 0xc4 // -60
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x3c) // +60

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *InstructionsTest) TestSbcCFlagUnset() {

	s.cpu.Registers.A = 0x02 // +2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x04) // +4

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

func (s *InstructionsTest) TestSbcZFlagSet() {

	s.cpu.Registers.A = 0x02 // +2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestSbcZFlagUnset() {

	s.cpu.Registers.A = 0x02 // +2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x02) // +2

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestSbcVFlagSet() {

	s.cpu.Registers.A = 0x80 // -128
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&V, "V flag is not set")

}

func (s *InstructionsTest) TestSbcVFlagUnset() {

	s.cpu.Registers.A = 0x01 // +1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&V, "V flag is set")

}

func (s *InstructionsTest) TestSbcNFlagSet() {

	s.cpu.Registers.A = 0xfd // -3
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestSbcNFlagUnset() {

	s.cpu.Registers.A = 0x02 // +2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe9)
	s.cpu.Memory.Store(0x0101, 0x01) // +1

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// CMP

func (s *InstructionsTest) TestCmpImmediate() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCmpZeroPage() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc5)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCmpZeroPageX() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xd5)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCmpAbsolute() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xcd)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCmpAbsoluteX() {

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

func (s *InstructionsTest) TestCmpAbsoluteY() {

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

func (s *InstructionsTest) TestCmpIndirectX() {

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

func (s *InstructionsTest) TestCmpIndirectY() {

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

func (s *InstructionsTest) TestCmpNFlagSet() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestCmpNFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc9)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *InstructionsTest) TestCmpZFlagSet() {

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

func (s *InstructionsTest) TestCmpZFlagUnset() {

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

func (s *InstructionsTest) TestCmpCFlagSet() {

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

func (s *InstructionsTest) TestCmpCFlagUnset() {

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

func (s *InstructionsTest) TestCpxImmediate() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCpxZeroPage() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe4)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCpxAbsolute() {

	s.cpu.Registers.X = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xec)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCpxNFlagSet() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestCpxNFlagUnset() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *InstructionsTest) TestCpxZFlagSet() {

	s.cpu.Registers.X = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCpxZFlagUnset() {

	s.cpu.Registers.X = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestCpxCFlagSet() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *InstructionsTest) TestCpxCFlagUnset() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

// CPY

func (s *InstructionsTest) TestCpyImmediate() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCpyZeroPage() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc4)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCpyAbsolute() {

	s.cpu.Registers.Y = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xcc)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xff)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCpyNFlagSet() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestCpyNFlagUnset() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

func (s *InstructionsTest) TestCpyZFlagSet() {

	s.cpu.Registers.Y = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestCpyZFlagUnset() {

	s.cpu.Registers.Y = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestCpyCFlagSet() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x01)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *InstructionsTest) TestCpyCFlagUnset() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc0)
	s.cpu.Memory.Store(0x0101, 0x02)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

// INC

func (s *InstructionsTest) TestIncZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xfe)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *InstructionsTest) TestIncZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xf6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0xfe)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *InstructionsTest) TestIncAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xee)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0xfe)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0084), "Memory is not 0xff")

}

func (s *InstructionsTest) TestIncAbsoluteX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xfe)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0xfe)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Memory.Fetch(0x0085), "Memory is not 0xff")

}

func (s *InstructionsTest) TestIncZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xff) // -1

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestIncZFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x00)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestIncNFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0xfe) // -2

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestIncNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x00)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// INX

func (s *InstructionsTest) TestInx() {

	s.cpu.Registers.X = 0xfe
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe8)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.X, "Register X is not 0xff")

}

func (s *InstructionsTest) TestInxZFlagSet() {

	s.cpu.Registers.X = 0xff // -1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe8)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestInxZFlagUnset() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestInxNFlagSet() {

	s.cpu.Registers.X = 0xfe // -2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe8)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestInxNFlagUnset() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xe8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// INY

func (s *InstructionsTest) TestIny() {

	s.cpu.Registers.Y = 0xfe // -2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc8)

	s.cpu.Execute()

	s.Equal(uint8(0xff), s.cpu.Registers.Y, "Register X is not 0xff")

}

func (s *InstructionsTest) TestInyZFlagSet() {

	s.cpu.Registers.Y = 0xff // -1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc8)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestInyZFlagUnset() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestInyNFlagSet() {

	s.cpu.Registers.Y = 0xfe // -2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc8)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestInyNFlagUnset() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc8)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// DEC

func (s *InstructionsTest) TestDecZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x01")

}

func (s *InstructionsTest) TestDecZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xd6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x01")

}

func (s *InstructionsTest) TestDecAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xce)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x01")

}

func (s *InstructionsTest) TestDecAbsoluteX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xde)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x01")

}

func (s *InstructionsTest) TestDecZFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x01)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestDecZFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestDecNFlagSet() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x00)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestDecNFlagUnset() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xc6)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x01)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// DEX

func (s *InstructionsTest) TestDex() {

	s.cpu.Registers.X = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xca)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.X, "Register X is not 0x01")

}

func (s *InstructionsTest) TestDexZFlagSet() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xca)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestDexZFlagUnset() {

	s.cpu.Registers.X = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xca)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestDexNFlagSet() {

	s.cpu.Registers.X = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xca)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestDexNFlagUnset() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0xca)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// DEY

func (s *InstructionsTest) TestDey() {

	s.cpu.Registers.Y = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x88)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.Y, "Register X is not 0x01")

}

func (s *InstructionsTest) TestDeyZFlagSet() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x88)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestDeyZFlagUnset() {

	s.cpu.Registers.Y = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x88)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestDeyNFlagSet() {

	s.cpu.Registers.Y = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x88)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestDeyNFlagUnset() {

	s.cpu.Registers.Y = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x88)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// ASL

func (s *InstructionsTest) TestAslAccumulator() {

	s.cpu.Registers.A = 0x2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.Equal(uint8(0x04), s.cpu.Registers.A, "Register A is not 0x04")

}

func (s *InstructionsTest) TestAslZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x06)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x04), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x04")

}

func (s *InstructionsTest) TestAslZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x16)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x04), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x04")

}

func (s *InstructionsTest) TestAslAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x04), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x04")

}

func (s *InstructionsTest) TestAslAbsoluteX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x1e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x04), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x04")

}

func (s *InstructionsTest) TestAslCFlagSet() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *InstructionsTest) TestAslCFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

func (s *InstructionsTest) TestAslZFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestAslZFlagUnset() {

	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestAslNFlagSet() {

	s.cpu.Registers.A = 0xfe
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestAslNFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x0a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// LSR

func (s *InstructionsTest) TestLsrAccumulator() {

	s.cpu.Registers.A = 0x2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Registers.A, "Register A is not 0x01")

}

func (s *InstructionsTest) TestLsrZeroPage() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x46)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x01")

}

func (s *InstructionsTest) TestLsrZeroPageX() {

	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x56)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x01")

}

func (s *InstructionsTest) TestLsrAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x01")

}

func (s *InstructionsTest) TestLsrAbsoluteX() {

	s.cpu.Registers.X = 1
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x5e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x01), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x01")

}

func (s *InstructionsTest) TestLsrCFlagSet() {

	s.cpu.Registers.A = 0xff
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *InstructionsTest) TestLsrCFlagUnset() {

	s.cpu.Registers.A = 0x10
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

func (s *InstructionsTest) TestLsrZFlagSet() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestLsrZFlagUnset() {

	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

// func (s *instTest) TestLsrNFlagSet() { }
// not tested, N bit always set to 0

func (s *InstructionsTest) TestLsrNFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// ROL

func (s *InstructionsTest) TestRolAccumulator() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x2
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.Equal(uint8(0x05), s.cpu.Registers.A, "Register A is not 0x05")

}

func (s *InstructionsTest) TestRolZeroPage() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x26)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x05), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x05")

}

func (s *InstructionsTest) TestRolZeroPageX() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x36)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x05), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x05")

}

func (s *InstructionsTest) TestRolAbsolute() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x02)

	s.cpu.Execute()

	s.Equal(uint8(0x05), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x05")

}

func (s *InstructionsTest) TestRolAbsoluteX() {

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

func (s *InstructionsTest) TestRolCFlagSet() {

	s.cpu.Registers.A = 0x80
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *InstructionsTest) TestRolCFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

func (s *InstructionsTest) TestRolZFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestRolZFlagUnset() {

	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestRolNFlagSet() {

	s.cpu.Registers.A = 0xfe
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestRolNFlagUnset() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x2a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// ROR

func (s *InstructionsTest) TestRorAccumulator() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0x08
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.Equal(uint8(0x84), s.cpu.Registers.A, "Register A is not 0x84")

}

func (s *InstructionsTest) TestRorZeroPage() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x66)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0084, 0x08)

	s.cpu.Execute()

	s.Equal(uint8(0x84), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x84")

}

func (s *InstructionsTest) TestRorZeroPageX() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.X = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x76)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0085, 0x08)

	s.cpu.Execute()

	s.Equal(uint8(0x84), s.cpu.Memory.Fetch(0x0085), "Memory is not 0x84")

}

func (s *InstructionsTest) TestRorAbsolute() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6e)
	s.cpu.Memory.Store(0x0101, 0x84)
	s.cpu.Memory.Store(0x0102, 0x00)
	s.cpu.Memory.Store(0x0084, 0x08)

	s.cpu.Execute()

	s.Equal(uint8(0x84), s.cpu.Memory.Fetch(0x0084), "Memory is not 0x84")

}

func (s *InstructionsTest) TestRorAbsoluteX() {

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

func (s *InstructionsTest) TestRorCFlagSet() {

	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&C, "C flag is not set")

}

func (s *InstructionsTest) TestRorCFlagUnset() {

	s.cpu.Registers.A = 0x10
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&C, "C flag is set")

}

func (s *InstructionsTest) TestRorZFlagSet() {

	s.cpu.Registers.A = 0x00
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&Z, "Z flag is not set")

}

func (s *InstructionsTest) TestRorZFlagUnset() {

	s.cpu.Registers.A = 0x02
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&Z, "Z flag is set")

}

func (s *InstructionsTest) TestRorNFlagSet() {

	s.cpu.Registers.P |= C
	s.cpu.Registers.A = 0xfe
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.NotZero(s.cpu.Registers.P&N, "N flag is not set")

}

func (s *InstructionsTest) TestRorNFlagUnset() {

	s.cpu.Registers.P &^= C
	s.cpu.Registers.A = 0x01
	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x6a)

	s.cpu.Execute()

	s.Zero(s.cpu.Registers.P&N, "N flag is set")

}

// JMP

func (s *InstructionsTest) TestJmpAbsolute() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Memory.Store(0x0100, 0x4c)
	s.cpu.Memory.Store(0x0101, 0xff)
	s.cpu.Memory.Store(0x0102, 0x01)

	s.cpu.Execute()

	s.Equal(uint16(0x01ff), s.cpu.Registers.PC, "Register PC is not 0x01ff")

}

func (s *InstructionsTest) TestJmpIndirect() {

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

func (s *InstructionsTest) TestJsr() {

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

func (s *InstructionsTest) TestRts() {

	s.cpu.Registers.PC = 0x0100
	s.cpu.push16(0x0102)

	s.cpu.Memory.Store(0x0100, 0x60)

	s.cpu.Execute()

	s.Equal(uint16(0x0103), s.cpu.Registers.PC, "Register PC is not 0x0103")

}

// BCC

func (s *InstructionsTest) TestBcc() {

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

func (s *InstructionsTest) TestBcs() {

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

func (s *InstructionsTest) TestBeq() {

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

func (s *InstructionsTest) TestBmi() {

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

func (s *InstructionsTest) TestBne() {

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

func (s *InstructionsTest) TestBpl() {

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

func (s *InstructionsTest) TestBvc() {

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

func (s *InstructionsTest) TestBvs() {

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

func (s *InstructionsTest) TestClc() {

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

func (s *InstructionsTest) TestCld() {

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

func (s *InstructionsTest) TestCli() {

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

func (s *InstructionsTest) TestClv() {

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

func (s *InstructionsTest) TestSec() {

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

func (s *InstructionsTest) TestSed() {

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

func (s *InstructionsTest) TestSei() {

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

func (s *InstructionsTest) TestBrk() {

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

func (s *InstructionsTest) TestRti() {

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

func (s *InstructionsTest) TestIrq() {

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

func (s *InstructionsTest) TestNmi() {

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

func (s *InstructionsTest) TestRst() {

	s.cpu.Registers.PC = 0x0100

	s.cpu.Interrupt(Rst, true)
	s.cpu.Memory.Store(0xfffc, 0x40)
	s.cpu.Memory.Store(0xfffd, 0x01)

	s.cpu.PerformInterrupts()

	s.Equal(uint16(0x0140), s.cpu.Registers.PC, "Register PC is not 0x0140")
	s.False(s.cpu.GetInterrupt(Rst), "Interrupt is set")
}
