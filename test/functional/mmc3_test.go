package functional

import (
	"strings"

	//"github.com/txc/nesgo/pkg/rp2ago3"
	"github.com/txc/nesgo/pkg/nes"
	"github.com/txc/nesgo/pkg/utils"
)

type MMC3Test struct {
	utils.CommonTest
	//nes *nes.NES
}

func (s *MMC3Test) FetchString(ns *nes.NES, address uint16) (result string) {
	var sb strings.Builder
	b := uint8(0xFF)
	offset := uint16(0x6004)
	for b != uint8(0x00) {
		b = ns.CPU.Memory.Fetch(offset)
		sb.WriteByte(b)
		offset++
	}
	return sb.String()
}

func (s *MMC3Test) TestMMC3() {
	tests := []struct {
		filename string
	}{
		{"1-clocking.nes"},
		{"2-details.nes"},
		{"3-A12_clocking.nes"},
		{"4-scanline_timing.nes"},
		{"5-MMC3.nes"},
		{"6-MMC6.nes"},
	}

	//testPath := romPath + "/mmc3_test/"
	testPath := romPath + "/mmc3_test_2/rom_singles/"
	for _, tt := range tests {
		// Run each test as a sub-test
		s.Run(tt.filename, func() {
			// Create NES instance with test ROM
			options := &nes.Options{
				Region: "NTSC",
				//Recorder: "jpeg",
			}

			ns, err := nes.NewNES(testPath+tt.filename, options)
			if !s.NoError(err) || !s.NotNil(ns) {
				return
			}

			// At 60 FPS, 600 frames = 10 seconds should be enough
			// At 60 FPS, 900 frames = 15 seconds should be enough
			//maxFrames := 600
			//maxFrames := 900
			maxFrames := 1800

			isLoop := false

			// Reset the NES
			ns.Reset()
			//ns.CPU.EnableDecode()

			frameCount := 0
			// Step through one frame worth of cycles
			// NTSC NES runs at ~29780 cycles per frame
			cyclesPerFrame := 29780
			if ns.ROM.Region() == nes.PAL {
				cyclesPerFrame = 33247
			}
			cyclesRun := 0

			for frameCount < maxFrames {
				cyclesRun = 0

				for cyclesRun < cyclesPerFrame {
					cycles, err := ns.Execute()
					if err != nil {
						s.FailNow("CPU execution error", "Error at frame %d: %v", frameCount, err)
						return
					}

					cyclesRun += int(cycles)
					check := []byte{0x00, 0x00, 0x00, 0x00}

					check[0] = ns.CPU.Memory.Fetch(0x6000)
					check[1] = ns.CPU.Memory.Fetch(0x6001)
					check[2] = ns.CPU.Memory.Fetch(0x6002)
					check[3] = ns.CPU.Memory.Fetch(0x6003)

					/*
						if check[1] == 0xDE && check[2] == 0xB0 {
							s.Logf("PC=0x%04X\t"+
								"0x6000=0x%02X, "+
								"0x6001=0x%02X, "+
								"0x6002=0x%02X, "+
								"0x6003=0x%02X",
								ns.CPU.M6502.Registers.PC,
								check[0], check[1], check[2], check[3])
						}
					*/
					if check[0] == 0x81 {
						s.Log("CPU Needs Reset!")
						ns.CPU.Reset()
					}
					if check[0] < 0x80 && check[1] == 0xDE && check[2] == 0xB0 && check[3] == 0x61 {
						s.Logf("Found $DE, $B0, $61 - Test Complete @ Frame %d", frameCount)
						isLoop = true
						break
					}

					if dec := ns.CPU.M6502.Decoder(); dec != nil {
						s.Log(dec.String())
					}
				}

				frameCount++
				if isLoop {
					break
				}
			}

			s.Equal(uint8(0x00), ns.CPU.Memory.Fetch(0x6000), strings.Trim(s.FetchString(ns, 0x6004), " \n\t\r"))
		})
	}

}
