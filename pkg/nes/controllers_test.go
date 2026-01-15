package nes

import (
	s "github.com/stretchr/testify/suite"
)

type ControllersTestSuite struct {
	s.Suite
	ctrls *Controllers
}

func (s *ControllersTestSuite) BeforeTest(suiteName, testName string) {
	s.ctrls = NewControllers()
	s.ctrls.Reset()
}

func (s *ControllersTestSuite) AfterTest(_, _ string) {
	s.ctrls = nil
}

// endregion baseApuTest

func (s *ControllersTestSuite) TestControllers() {
	s.ctrls.controllers[0].buttons = 0x01

	s.Equal(uint8(0x41), s.ctrls.Fetch(0x4016))

	s.ctrls.controllers[0].buttons = 0x00

	s.Equal(uint8(0x40), s.ctrls.Fetch(0x4016))

	s.ctrls.Store(0x4016, 1)
	s.ctrls.Store(0x4016, 0)

	s.ctrls.controllers[0].buttons = 0xff
	s.ctrls.controllers[1].buttons = 0xff

	for range 20 {
		s.Equal(uint8(0x41), s.ctrls.Fetch(0x4016))
	}

	for range 20 {
		s.Equal(uint8(0x41), s.ctrls.Fetch(0x4017))
	}

	s.ctrls.Store(0x4016, 1)
	s.ctrls.Store(0x4016, 0)
	s.ctrls.controllers[0].buttons = 0x00
	s.ctrls.controllers[1].buttons = 0x00

	for i := range 20 {
		if i < 8 {
			s.Equal(uint8(0x40), s.ctrls.Fetch(0x4016))
		} else {
			s.Equal(uint8(0x41), s.ctrls.Fetch(0x4016))
		}
	}

	for i := range 20 {
		if i < 8 {
			s.Equal(uint8(0x40), s.ctrls.Fetch(0x4017))
		} else {
			s.Equal(uint8(0x41), s.ctrls.Fetch(0x4017))
		}
	}
}
