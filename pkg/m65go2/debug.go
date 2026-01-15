// Package m65go2 simulates the MOS 6502 CPU
package m65go2

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type DebugMode uint16

const (
	DebugModeInstruction DebugMode = 1 << iota
	DebugModeCycleDump
	DebugModeVerbose
)

func GetStackFrames() []runtime.Frame {
	// Ask runtime.Callers for up to 10 PCs, including runtime.Callers itself.
	pc := make([]uintptr, 10)
	n := runtime.Callers(3, pc)
	if n == 0 {
		// No PCs available. This can happen if the first argument to
		// runtime.Callers is large.
		//
		// Return now to avoid processing the zero Frame that would
		// otherwise be returned by frames.Next below.
		return []runtime.Frame{}
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	var ret = make([]runtime.Frame, 0, n)
	// Loop to get frames.
	// A fixed number of PCs can expand to an indefinite number of Frames.
	for {
		frame, more := frames.Next()

		function := strings.ReplaceAll(frame.Function, "github.com/txc/nesgo/m65go2.", "")
		ret = append(ret, runtime.Frame{
			PC:       frame.PC,
			Function: function,
			File:     filepath.Base(frame.File),
			Line:     frame.Line,
		})

		// Check whether there are more frames to process after this one.
		if !more {
			break
		}
	}
	return ret
}

// Prints the memory data between `start` and `stop` parameters.
// If `stop` parameter is left out, it will dump the 16 byte segment where
// the `start` address exists.
// (`start = start & 0xFFF0` and `stop = start | 0x000F`)
//
// This method reads over blocks.
func MemoryDump(mem Memory, start uint16, stop *uint16) string {
	var output strings.Builder

	var end uint16
	if stop == nil {
		end = (start & 0xFFF0) | 0x000F
		start = start & 0xFFF0
	} else {
		end = *stop
	}

	fmt.Fprintf(&output, "MEMORY DUMP FOR: $%04X - $%04X\n", start, end)
	output.WriteString("ADDR 00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F\n")

	modulus := (((end + 1) - start) + 1) % 16

	for n := range modulus {
		addr := start + uint16(n*16)
		fmt.Fprintf(&output, "%04X", addr)
		for offset := range uint16(16) {
			if addr+offset > end {
				break
			}
			value := mem.Fetch(addr + offset)
			fmt.Fprintf(&output, " %02X", value)
		}
		output.WriteString("\n")
	}

	return output.String()
}
