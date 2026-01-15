package m65go2

import (
	s "github.com/stretchr/testify/suite"
)

type MemoryTest struct {
	s.Suite
}

func (s *MemoryTest) TestSamePage() {
	for a := uint16(0x0000); ; a += 0x0100 {
		for b := uint16(0x0000); ; b += 0x0100 {
			if a>>8 == b>>8 {
				s.True(SamePage(a, b), "Expected SamePage(%#04x, %#04x) to be true", a, b)
			} else {
				s.False(SamePage(a, b), "Expected SamePage(%#04x, %#04x) to be false", a, b)
			}

			if b == 0xff00 {
				break
			}
		}

		if a == 0xff00 {
			break
		}
	}
}
