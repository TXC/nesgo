package rp2ago3

import (
	"fmt"

	s "github.com/stretchr/testify/suite"
)

// region apuTest
type apuTest struct {
	s.Suite
	apu *APU
}

func (s *apuTest) BeforeTest(suiteName, testName string) {
	s.apu = NewAPU(&MappedMemory{}, 0, nil)
	s.apu.Reset()
}

func (s *apuTest) AfterTest(_, _ string) {
	s.apu = nil
}

// endregion apuTest

type PulseTest struct {
	flag     PulseFlag
	index    int
	value    uint8
	expected uint8
}

func (s *apuTest) TestPulse1Channel() {
	address := uint16(0x4000)

	for i := range s.apu.Pulse1.Registers {
		s.apu.Pulse1.Registers[i] = 0x00
		s.apu.Store(address+uint16(i), 0xff)
		s.Equal(uint8(0xff), s.apu.Pulse1.Registers[i])

		s.apu.Pulse1.Registers[i] = 0xff
		s.Equal(uint8(0x00), s.apu.Fetch(address+uint16(i)))
	}

	pts := []PulseTest{
		{Duty, 0, 0xc0, 0x03},
		{Duty, 0, 0x3f, 0x00},

		{PulseEnvelopeLoopLengthCounterHalt, 0, 0x20, 0x01},
		{PulseEnvelopeLoopLengthCounterHalt, 0, 0xdf, 0x00},

		{PulseConstantVolume, 0, 0x10, 0x01},
		{PulseConstantVolume, 0, 0xef, 0x00},

		{PulseVolumeEnvelope, 0, 0x0f, 0x0f},
		{PulseVolumeEnvelope, 0, 0xf0, 0x00},

		{SweepEnabled, 1, 0x80, 0x01},
		{SweepEnabled, 1, 0x7f, 0x00},

		{SweepPeriod, 1, 0x70, 0x07},
		{SweepPeriod, 1, 0x8f, 0x00},

		{SweepNegate, 1, 0x08, 0x01},
		{SweepNegate, 1, 0xf7, 0x00},

		{SweepShift, 1, 0x07, 0x07},
		{SweepShift, 1, 0xf8, 0x00},

		{PulseTimerLow, 2, 0xff, 0xff},
		{PulseTimerLow, 2, 0x00, 0x00},

		{PulseLengthCounterLoad, 3, 0xf8, 0x1f},
		{PulseLengthCounterLoad, 3, 0x07, 0x00},

		{PulseTimerHigh, 3, 0x07, 0x07},
		{PulseTimerHigh, 3, 0xf8, 0x00},
	}

	for _, pt := range pts {
		s.Run(fmt.Sprintf("%s/%02X", pt.flag.String(), pt.value), func() {
			s.apu.Pulse1.Registers[pt.index] = pt.value

			actual := s.apu.Pulse1.registers(pt.flag)

			s.Equal(pt.expected, actual)
		})
	}
}

func (s *apuTest) TestPulse2Channel() {
	address := uint16(0x4004)

	for i := range s.apu.Pulse2.Registers {
		s.apu.Pulse2.Registers[i] = 0x00
		s.apu.Store(address+uint16(i), 0xff)

		s.Equal(uint8(0xff), s.apu.Pulse2.Registers[i])
		s.apu.Pulse2.Registers[i] = 0xff
		s.Equal(uint8(0x00), s.apu.Fetch(address+uint16(i)))
	}

	pts := []PulseTest{
		{Duty, 0, 0xc0, 0x03},
		{Duty, 0, 0x3f, 0x00},

		{PulseEnvelopeLoopLengthCounterHalt, 0, 0x20, 0x01},
		{PulseEnvelopeLoopLengthCounterHalt, 0, 0xdf, 0x00},

		{PulseConstantVolume, 0, 0x10, 0x01},
		{PulseConstantVolume, 0, 0xef, 0x00},

		{PulseVolumeEnvelope, 0, 0x0f, 0x0f},
		{PulseVolumeEnvelope, 0, 0xf0, 0x00},

		{SweepEnabled, 1, 0x80, 0x01},
		{SweepEnabled, 1, 0x7f, 0x00},

		{SweepPeriod, 1, 0x70, 0x07},
		{SweepPeriod, 1, 0x8f, 0x00},

		{SweepNegate, 1, 0x08, 0x01},
		{SweepNegate, 1, 0xf7, 0x00},

		{SweepShift, 1, 0x07, 0x07},
		{SweepShift, 1, 0xf8, 0x00},

		{PulseTimerLow, 2, 0xff, 0xff},
		{PulseTimerLow, 2, 0x00, 0x00},

		{PulseLengthCounterLoad, 3, 0xf8, 0x1f},
		{PulseLengthCounterLoad, 3, 0x07, 0x00},

		{PulseTimerHigh, 3, 0x07, 0x07},
		{PulseTimerHigh, 3, 0xf8, 0x00},
	}

	for _, pt := range pts {
		s.Run(fmt.Sprintf("%s/%02X", pt.flag.String(), pt.value), func() {
			s.apu.Pulse2.Registers[pt.index] = pt.value

			actual := s.apu.Pulse2.registers(pt.flag)

			s.Equal(pt.expected, actual)
		})
	}
}

// TestPulseChannelSweepDifference tests the key difference between Pulse1 and Pulse2:
// Pulse1 uses one's complement (MinusOne=true) while Pulse2 uses two's complement (MinusOne=false)
// This affects the target period calculation when sweep negation is enabled.
func (s *apuTest) TestPulseChannelSweepDifference() {
	// Test case: Period = 256, Shift = 1, Negate enabled
	// Expected behavior:
	// - Delta = 256 >> 1 = 128
	// - Pulse1 (one's complement): TargetPeriod = 256 + (-128 - 1) = 127
	// - Pulse2 (two's complement): TargetPeriod = 256 + (-128) = 128

	//   76543210
	//   EPPPNSSS
	// 0b10001001 (0x89)
	//   ||||||||
	//   |||||+++  S: Shift
	//   ||||+---  N: Negate=1
	//   |+++----  P: Period
	//   +-------  E: Enable=1
	sweep := uint8(0x89) // Negate=1, Shift=1
	period := uint16(256)
	result := int16(128)

	s.apu.Pulse1.Store(1, sweep)
	s.apu.Pulse1.Store(2, uint8(period&0xFF))
	s.apu.Pulse1.Store(3, uint8(period>>8))

	s.Equal(uint8(1), s.apu.Pulse1.registers(SweepNegate))
	s.Equal(uint8(1), s.apu.Pulse1.registers(SweepShift))
	s.Equal(int16(period), s.apu.Pulse1.Divider.Period)

	target1 := s.apu.Pulse1.TargetPeriod()

	// Setup Pulse2
	s.apu.Pulse2.Store(1, sweep)
	s.apu.Pulse2.Store(2, uint8(period&0xFF))
	s.apu.Pulse2.Store(3, uint8(period>>8))

	s.Equal(uint8(1), s.apu.Pulse2.registers(SweepNegate))
	s.Equal(uint8(1), s.apu.Pulse2.registers(SweepShift))
	s.Equal(int16(period), s.apu.Pulse2.Divider.Period)

	target2 := s.apu.Pulse2.TargetPeriod()

	// Evaluate results
	s.Equal(result-1, target1, "Pulse1 should use one's complement (subtract 1)")
	s.Equal(result, target2, "Pulse2 should use two's complement (no extra subtraction)")
	s.NotEqual(target1, target2, "Pulse1 and Pulse2 sweep calculations should differ")
}

// TestPulseChannelSweepDifferenceMultipleCases tests various sweep configurations
// to verify the consistent difference between Pulse1 and Pulse2 behavior.
func (s *apuTest) TestPulseChannelSweepDifferenceMultipleCases() {
	testCases := []struct {
		name           string
		sweep          uint8
		period         uint16
		expectedPulse1 int16
		expectedPulse2 int16
	}{
		// 512 + (-(512>>2)) = 512 - 128 = 384
		{"Period=512,Shift=2", 0x0A, 0x0200, 0x017F, 0x0180},
		// 1024 + (-(1024>>3)) = 1024 - 128 = 896
		{"Period=1024,Shift=3", 0x0B, 0x0400, 0x037F, 0x0380},
		// 100 + (-(100>>0)) = 100 - 100 = 0 for Pulse2, -1 for Pulse1
		{"Period=100,Shift=0", 0x08, 0x0064, -1, 0},
		// 800 + (-(800>>4)) = 800 - 50 = 750
		{"Period=800,Shift=4", 0x0C, 0x0320, 0x02ed, 0x02ee},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup Pulse1
			s.apu.Pulse1.Store(1, tc.sweep)
			s.apu.Pulse1.Store(2, uint8(tc.period&0xFF))
			s.apu.Pulse1.Store(3, uint8(tc.period>>8))

			s.Equal(int16(tc.period), s.apu.Pulse1.Divider.Period)
			target1 := s.apu.Pulse1.TargetPeriod()

			// Setup Pulse2
			s.apu.Pulse2.Store(1, tc.sweep)
			s.apu.Pulse2.Store(2, uint8(tc.period&0xFF))
			s.apu.Pulse2.Store(3, uint8(tc.period>>8))

			s.Equal(int16(tc.period), s.apu.Pulse2.Divider.Period)
			target2 := s.apu.Pulse2.TargetPeriod()

			s.Equal(tc.expectedPulse1, target1, "Pulse1 target period mismatch")
			s.Equal(tc.expectedPulse2, target2, "Pulse2 target period mismatch")

			// The difference should always be exactly 1 when negate is enabled
			s.Equal(int16(1), target2-target1, "Pulse2 should always be 1 higher than Pulse1 in negate mode")
		})
	}
}

// TestPulseChannelPositiveSweep tests that both channels behave identically
// when sweep negation is disabled (positive sweep).
func (s *apuTest) TestPulseChannelPositiveSweep() {
	sweep := uint8(0xf2) // Negate=0, Shift=1
	period := uint16(256)
	result := int16(320)

	s.apu.Pulse1.Store(1, sweep)
	s.apu.Pulse1.Store(2, uint8(period&0xFF))
	s.apu.Pulse1.Store(3, uint8(period>>8))

	s.apu.Pulse2.Store(1, sweep)
	s.apu.Pulse2.Store(2, uint8(period&0xFF))
	s.apu.Pulse2.Store(3, uint8(period>>8))

	target1 := s.apu.Pulse1.TargetPeriod()
	target2 := s.apu.Pulse2.TargetPeriod()

	// Both should be 256 + (256 >> 2) = 256 + 64 = 320
	s.Equal(result, target1, "Pulse1 positive sweep")
	s.Equal(result, target2, "Pulse2 positive sweep")
	s.Equal(target1, target2, "Both channels should behave identically with positive sweep")

	s.apu.Pulse1.Store(1, 0x92)
	s.apu.Pulse1.Store(2, 0x05)
	s.apu.Pulse1.Store(3, 0x00)

	s.apu.Pulse2.Store(1, 0x92)
	s.apu.Pulse2.Store(2, 0x05)
	s.apu.Pulse2.Store(3, 0x00)

	result2 := int16(6)
	s.Equal(result2, s.apu.Pulse1.TargetPeriod())
	s.Equal(result2, s.apu.Pulse2.TargetPeriod())

}

type TriangleTest struct {
	flag     TriangleFlag
	index    int
	value    uint8
	expected uint8
}

func (s *apuTest) TestTriangleChannel() {

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

	pts := []TriangleTest{
		{LengthCounterHaltLinearCounterControl, 0, 0x80, 0x01},
		{LengthCounterHaltLinearCounterControl, 0, 0x7f, 0x00},
		{LinearCounterLoad, 0, 0x7f, 0x7f},
		{LinearCounterLoad, 0, 0x80, 0x00},
		{TriangleTimerLow, 1, 0xff, 0xff},
		{TriangleTimerLow, 1, 0x00, 0x00},
		{TriangleLengthCounterLoad, 2, 0xf8, 0x1f},
		{TriangleLengthCounterLoad, 2, 0x07, 0x00},
		{TriangleTimerHigh, 2, 0x07, 0x07},
		{TriangleTimerHigh, 2, 0xf8, 0x00},
	}

	for _, pt := range pts {
		s.Run(fmt.Sprintf("%s/%02X", pt.flag.String(), pt.value), func() {
			s.apu.Triangle.Registers[pt.index] = pt.value

			actual := s.apu.Triangle.registers(pt.flag)
			s.Equal(pt.expected, actual)
		})
	}
}

type NoiseTest struct {
	flag     NoiseFlag
	index    int
	value    uint8
	expected uint8
}

func (s *apuTest) TestNoiseChannel() {
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

	pts := []NoiseTest{
		{NoiseEnvelopeLoopLengthCounterHalt, 0, 0x20, 0x01},
		{NoiseEnvelopeLoopLengthCounterHalt, 0, 0xdf, 0x00},
		{NoiseConstantVolume, 0, 0x10, 0x01},
		{NoiseConstantVolume, 0, 0xef, 0x00},
		{NoiseVolumeEnvelope, 0, 0x0f, 0x0f},
		{NoiseVolumeEnvelope, 0, 0xf0, 0x00},
		{LoopNoise, 1, 0x80, 0x01},
		{LoopNoise, 1, 0x7f, 0x00},
		{NoisePeriod, 1, 0x0f, 0x0f},
		{NoisePeriod, 1, 0xf0, 0x00},
		{NoiseLengthCounterLoad, 2, 0xf8, 0x1f},
		{NoiseLengthCounterLoad, 2, 0x07, 0x00},
	}

	for _, pt := range pts {
		s.Run(fmt.Sprintf("%s/%02X", pt.flag.String(), pt.value), func() {
			s.apu.Noise.Registers[pt.index] = pt.value

			actual := s.apu.Noise.registers(pt.flag)
			s.Equal(pt.expected, actual)
		})
	}
}

type DMCTest struct {
	flag     DMCFlag
	index    int
	value    uint8
	expected uint8
}

func (s *apuTest) TestDMCChannel() {
	address := uint16(0x4010)

	for i := range s.apu.DMC.Registers {
		s.apu.DMC.Registers[i] = 0x00
		s.apu.Store(address+uint16(i), 0xff)

		s.Equal(uint8(0xff), s.apu.DMC.Registers[i])

		s.apu.DMC.Registers[i] = 0xff
		s.Equal(uint8(0x00), s.apu.Fetch(address+uint16(i)))
	}

	pts := []DMCTest{
		{IRQEnable, 0, 0x80, 0x01},
		{IRQEnable, 0, 0x7f, 0x00},
		{Loop, 0, 0x40, 0x01},
		{Loop, 0, 0xbf, 0x00},
		{Frequency, 0, 0x0f, 0x0f},
		{Frequency, 0, 0xf0, 0x00},
		{LoadCounter, 1, 0x7f, 0x7f},
		{LoadCounter, 1, 0x80, 0x00},
		{SampleAddress, 2, 0xff, 0xff},
		{SampleAddress, 2, 0x00, 0x00},
		{SampleLength, 3, 0xff, 0xff},
		{SampleLength, 3, 0x00, 0x00},
	}

	for _, pt := range pts {
		s.Run(fmt.Sprintf("%s/%02X", pt.flag.String(), pt.value), func() {
			s.apu.DMC.Registers[pt.index] = pt.value

			actual := s.apu.DMC.registers(pt.flag)

			s.Equal(pt.expected, actual)
		})
	}
}

type ControlTest struct {
	flag     ControlFlag
	value    uint8
	expected bool
}

func (s *apuTest) TestControl() {
	address := uint16(0x4015)

	s.apu.Registers.Control = 0x00
	s.apu.Store(address, 0xef)

	s.Equal(Control(0xef), s.apu.Registers.Control)

	s.apu.Registers.Control = 0x00
	s.apu.Registers.Status = 0x00

	s.Equal(uint8(0x00), s.apu.Fetch(address))

	pts := []ControlTest{
		{EnableDMC, 0x10, true},
		{EnableDMC, 0xef, false},
		{EnableNoise, 0x08, true},
		{EnableNoise, 0xf7, false},
		{EnableTriangle, 0x04, true},
		{EnableTriangle, 0xfb, false},
		{EnablePulseChannel2, 0x02, true},
		{EnablePulseChannel2, 0xfd, false},
		{EnablePulseChannel1, 0x01, true},
		{EnablePulseChannel1, 0xfe, false},
	}

	for _, pt := range pts {
		s.Run(fmt.Sprintf("%s/%02X", pt.flag.String(), pt.value), func() {
			s.apu.Registers.Control = Control(pt.value)

			actual := s.apu.control(pt.flag)

			s.Equal(pt.expected, actual)
		})
	}
}

type StatusTest struct {
	flag     StatusFlag
	value    uint8
	expected bool
}

func (s *apuTest) TestStatus() {
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

func (s *apuTest) TestFrameCounter() {
	address := uint16(0x4017)

	s.apu.FrameCounter.Register = 0xff

	s.apu.Store(address, 0x00)

	s.Equal(uint8(0x00), s.apu.FrameCounter.Register)

	pts := []FrameCounterTest{
		{Mode, 0x80, 5},
		{Mode, 0x7f, 4},
		{IRQInhibit, 0x40, 0x01},
		{IRQInhibit, 0xbf, 0x00},
	}

	for _, pt := range pts {
		s.Run(fmt.Sprintf("%s/%02X", pt.flag.String(), pt.value), func() {
			s.apu.FrameCounter.Register = pt.value

			actual := s.apu.FrameCounter.register(pt.flag)
			s.Equal(pt.expected, actual)
		})
	}
}
