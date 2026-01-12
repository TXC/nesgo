package rp2ago3

import (
	"testing"

	s "github.com/stretchr/testify/suite"
)

func TestApuSuite(t *testing.T) {
	s.Run(t, &baseApuTest{})
}

// region baseApuTest
type baseApuTest struct {
	s.Suite
	apu *APU
}

func (s *baseApuTest) BeforeTest(suiteName, testName string) {
	s.apu = NewAPU(&MappedMemory{}, 0, nil)
	s.apu.Reset()
}

func (s *baseApuTest) AfterTest(_, _ string) {
	s.apu = nil
}

// endregion baseApuTest

type PulseTest struct {
	flag     PulseFlag
	index    int
	value    uint8
	expected uint8
}

func (s *baseApuTest) TestPulse1Channel() {
	address := uint16(0x4000)

	for i := range s.apu.Pulse1.Registers {
		s.apu.Pulse1.Registers[i] = 0x00
		s.apu.Store(address+uint16(i), 0xff)
		s.Equal(uint8(0xff), s.apu.Pulse1.Registers[i])

		s.apu.Pulse1.Registers[i] = 0xff
		s.Equal(uint8(0x00), s.apu.Fetch(address+uint16(i)))
	}

	pts := []PulseTest{}

	pts = append(pts, PulseTest{Duty, 0, 0xc0, 0x03})
	pts = append(pts, PulseTest{Duty, 0, 0x3f, 0x00})

	pts = append(pts, PulseTest{PulseEnvelopeLoopLengthCounterHalt, 0, 0x20, 0x01})
	pts = append(pts, PulseTest{PulseEnvelopeLoopLengthCounterHalt, 0, 0xdf, 0x00})

	pts = append(pts, PulseTest{PulseConstantVolume, 0, 0x10, 0x01})
	pts = append(pts, PulseTest{PulseConstantVolume, 0, 0xef, 0x00})

	pts = append(pts, PulseTest{PulseVolumeEnvelope, 0, 0x0f, 0x0f})
	pts = append(pts, PulseTest{PulseVolumeEnvelope, 0, 0xf0, 0x00})

	pts = append(pts, PulseTest{SweepEnabled, 1, 0x80, 0x01})
	pts = append(pts, PulseTest{SweepEnabled, 1, 0x7f, 0x00})

	pts = append(pts, PulseTest{SweepPeriod, 1, 0x70, 0x07})
	pts = append(pts, PulseTest{SweepPeriod, 1, 0x8f, 0x00})

	pts = append(pts, PulseTest{SweepNegate, 1, 0x08, 0x01})
	pts = append(pts, PulseTest{SweepNegate, 1, 0xf7, 0x00})

	pts = append(pts, PulseTest{SweepShift, 1, 0x07, 0x07})
	pts = append(pts, PulseTest{SweepShift, 1, 0xf8, 0x00})

	pts = append(pts, PulseTest{PulseTimerLow, 2, 0xff, 0xff})
	pts = append(pts, PulseTest{PulseTimerLow, 2, 0x00, 0x00})

	pts = append(pts, PulseTest{PulseLengthCounterLoad, 3, 0xf8, 0x1f})
	pts = append(pts, PulseTest{PulseLengthCounterLoad, 3, 0x07, 0x00})

	pts = append(pts, PulseTest{PulseTimerHigh, 3, 0x07, 0x07})
	pts = append(pts, PulseTest{PulseTimerHigh, 3, 0xf8, 0x00})

	for _, pt := range pts {
		s.apu.Pulse1.Registers[pt.index] = pt.value

		actual := s.apu.Pulse1.registers(pt.flag)
		expected := pt.expected

		s.Equal(expected, actual)
	}
}

func (s *baseApuTest) TestPulse2Channel() {
	address := uint16(0x4004)

	for i := range s.apu.Pulse2.Registers {
		s.apu.Pulse2.Registers[i] = 0x00
		s.apu.Store(address+uint16(i), 0xff)

		s.Equal(uint8(0xff), s.apu.Pulse2.Registers[i])
		s.apu.Pulse2.Registers[i] = 0xff
		s.Equal(uint8(0x00), s.apu.Fetch(address+uint16(i)))
	}
}

type TriangleTest struct {
	flag     TriangleFlag
	index    int
	value    uint8
	expected uint8
}

func (s *baseApuTest) TestTriangleChannel() {

	var address uint16

	for i := range s.apu.Triangle.Registers {
		switch i {
		case 0:
			address = 0x4008
		case 1:
			address = 0x400a
		case 2:
			address = 0x400b
		}

		s.apu.Triangle.Registers[i] = 0x00
		s.apu.Store(address, 0xff)

		s.Equal(uint8(0xff), s.apu.Triangle.Registers[i])

		s.apu.Triangle.Registers[i] = 0xff
		s.Equal(uint8(0x00), s.apu.Fetch(address))
	}

	pts := []TriangleTest{}

	pts = append(pts, TriangleTest{LengthCounterHaltLinearCounterControl, 0, 0x80, 0x01})
	pts = append(pts, TriangleTest{LengthCounterHaltLinearCounterControl, 0, 0x7f, 0x00})

	pts = append(pts, TriangleTest{LinearCounterLoad, 0, 0x7f, 0x7f})
	pts = append(pts, TriangleTest{LinearCounterLoad, 0, 0x80, 0x00})

	pts = append(pts, TriangleTest{TriangleTimerLow, 1, 0xff, 0xff})
	pts = append(pts, TriangleTest{TriangleTimerLow, 1, 0x00, 0x00})

	pts = append(pts, TriangleTest{TriangleLengthCounterLoad, 2, 0xf8, 0x1f})
	pts = append(pts, TriangleTest{TriangleLengthCounterLoad, 2, 0x07, 0x00})

	pts = append(pts, TriangleTest{TriangleTimerHigh, 2, 0x07, 0x07})
	pts = append(pts, TriangleTest{TriangleTimerHigh, 2, 0xf8, 0x00})

	for _, pt := range pts {
		s.apu.Triangle.Registers[pt.index] = pt.value

		actual := s.apu.Triangle.registers(pt.flag)
		s.Equal(pt.expected, actual)
	}
}

type NoiseTest struct {
	flag     NoiseFlag
	index    int
	value    uint8
	expected uint8
}

func (s *baseApuTest) TestNoiseChannel() {
	var address uint16

	for i := range s.apu.Noise.Registers {
		switch i {
		case 0:
			address = 0x400c
		case 1:
			address = 0x400e
		case 2:
			address = 0x400f
		}

		s.apu.Noise.Registers[i] = 0x00
		s.apu.Store(address, 0xff)

		s.Equal(uint8(0xff), s.apu.Noise.Registers[i])

		s.apu.Noise.Registers[i] = 0xff
		s.Equal(uint8(0x00), s.apu.Fetch(address))
	}

	pts := []NoiseTest{}

	pts = append(pts, NoiseTest{NoiseEnvelopeLoopLengthCounterHalt, 0, 0x20, 0x01})
	pts = append(pts, NoiseTest{NoiseEnvelopeLoopLengthCounterHalt, 0, 0xdf, 0x00})

	pts = append(pts, NoiseTest{NoiseConstantVolume, 0, 0x10, 0x01})
	pts = append(pts, NoiseTest{NoiseConstantVolume, 0, 0xef, 0x00})

	pts = append(pts, NoiseTest{NoiseVolumeEnvelope, 0, 0x0f, 0x0f})
	pts = append(pts, NoiseTest{NoiseVolumeEnvelope, 0, 0xf0, 0x00})

	pts = append(pts, NoiseTest{LoopNoise, 1, 0x80, 0x01})
	pts = append(pts, NoiseTest{LoopNoise, 1, 0x7f, 0x00})

	pts = append(pts, NoiseTest{NoisePeriod, 1, 0x0f, 0x0f})
	pts = append(pts, NoiseTest{NoisePeriod, 1, 0xf0, 0x00})

	pts = append(pts, NoiseTest{NoiseLengthCounterLoad, 2, 0xf8, 0x1f})
	pts = append(pts, NoiseTest{NoiseLengthCounterLoad, 2, 0x07, 0x00})

	for _, pt := range pts {
		s.apu.Noise.Registers[pt.index] = pt.value

		actual := s.apu.Noise.registers(pt.flag)
		s.Equal(pt.expected, actual)
	}
}

type DMCTest struct {
	flag     DMCFlag
	index    int
	value    uint8
	expected uint8
}

func (s *baseApuTest) TestDMCChannel() {
	address := uint16(0x4010)

	for i := range s.apu.DMC.Registers {
		s.apu.DMC.Registers[i] = 0x00
		s.apu.Store(address+uint16(i), 0xff)

		s.Equal(uint8(0xff), s.apu.DMC.Registers[i])

		s.apu.DMC.Registers[i] = 0xff
		s.Equal(uint8(0x00), s.apu.Fetch(address+uint16(i)))
	}

	pts := []DMCTest{}

	pts = append(pts, DMCTest{IRQEnable, 0, 0x80, 0x01})
	pts = append(pts, DMCTest{IRQEnable, 0, 0x7f, 0x00})

	pts = append(pts, DMCTest{Loop, 0, 0x40, 0x01})
	pts = append(pts, DMCTest{Loop, 0, 0xbf, 0x00})

	pts = append(pts, DMCTest{Frequency, 0, 0x0f, 0x0f})
	pts = append(pts, DMCTest{Frequency, 0, 0xf0, 0x00})

	pts = append(pts, DMCTest{LoadCounter, 1, 0x7f, 0x7f})
	pts = append(pts, DMCTest{LoadCounter, 1, 0x80, 0x00})

	pts = append(pts, DMCTest{SampleAddress, 2, 0xff, 0xff})
	pts = append(pts, DMCTest{SampleAddress, 2, 0x00, 0x00})

	pts = append(pts, DMCTest{SampleLength, 3, 0xff, 0xff})
	pts = append(pts, DMCTest{SampleLength, 3, 0x00, 0x00})

	for _, pt := range pts {
		s.apu.DMC.Registers[pt.index] = pt.value

		actual := s.apu.DMC.registers(pt.flag)

		s.Equal(pt.expected, actual)
	}
}

type ControlTest struct {
	flag     ControlFlag
	value    uint8
	expected bool
}

func (s *baseApuTest) TestControl() {
	address := uint16(0x4015)

	s.apu.Registers.Control = 0x00
	s.apu.Store(address, 0xef)

	s.Equal(Control(0xef), s.apu.Registers.Control)

	s.apu.Registers.Control = 0x00
	s.apu.Registers.Status = 0x00

	s.Equal(uint8(0x00), s.apu.Fetch(address))

	pts := []ControlTest{}

	pts = append(pts, ControlTest{EnableDMC, 0x10, true})
	pts = append(pts, ControlTest{EnableDMC, 0xef, false})

	pts = append(pts, ControlTest{EnableNoise, 0x08, true})
	pts = append(pts, ControlTest{EnableNoise, 0xf7, false})

	pts = append(pts, ControlTest{EnableTriangle, 0x04, true})
	pts = append(pts, ControlTest{EnableTriangle, 0xfb, false})

	pts = append(pts, ControlTest{EnablePulseChannel2, 0x02, true})
	pts = append(pts, ControlTest{EnablePulseChannel2, 0xfd, false})

	pts = append(pts, ControlTest{EnablePulseChannel1, 0x01, true})
	pts = append(pts, ControlTest{EnablePulseChannel1, 0xfe, false})

	for _, pt := range pts {
		s.apu.Registers.Control = Control(pt.value)

		actual := s.apu.control(pt.flag)

		s.Equal(pt.expected, actual)
	}
}

type StatusTest struct {
	flag     StatusFlag
	value    uint8
	expected bool
}

func (s *baseApuTest) TestStatus() {
	address := uint16(0x4015)

	s.apu.Registers.Status = 0x00
	value := s.apu.Fetch(address)

	s.apu.Noise.LengthCounter = 0

	s.Equal(uint8(0x00), value)
}

type FrameCounterTest struct {
	flag     FrameCounterFlag
	value    uint8
	expected uint8
}

func (s *baseApuTest) TestFrameCounter() {
	address := uint16(0x4017)

	s.apu.FrameCounter.Register = 0xff

	s.apu.Store(address, 0x00)

	s.Equal(uint8(0x00), s.apu.FrameCounter.Register)

	pts := []FrameCounterTest{}

	pts = append(pts, FrameCounterTest{Mode, 0x80, 5})
	pts = append(pts, FrameCounterTest{Mode, 0x7f, 4})

	pts = append(pts, FrameCounterTest{IRQInhibit, 0x40, 0x01})
	pts = append(pts, FrameCounterTest{IRQInhibit, 0xbf, 0x00})

	for _, pt := range pts {
		s.apu.FrameCounter.Register = pt.value

		actual := s.apu.FrameCounter.register(pt.flag)
		s.Equal(pt.expected, actual)
	}
}
