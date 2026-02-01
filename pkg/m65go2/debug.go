// Package m65go2 simulates the MOS 6502 CPU
package m65go2

import (
	"encoding/hex"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type DebugMode uint16

const (
	DebugModeTraceCPU DebugMode = 1 << iota
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
	var end uint16
	if stop == nil {
		end = (start & 0xFFF0) | 0x000F
		start = start & 0xFFF0
	} else {
		end = *stop
	}

	return hex.Dump(mem.(*BasicMemory).Raw()[start : end+1])
}

type DecodeBytes struct {
	Instruction *Instruction
	PC          uint16
	Hi          uint8
	Lo          uint8
}

func NewDecodeBytes(inst *Instruction, pc uint16, hi uint8, lo uint8) *DecodeBytes {
	db := &DecodeBytes{
		Instruction: inst,
	}
	if inst.Bytes() >= 2 {
		db.Lo = lo
	}
	if inst.Bytes() == 3 {
		db.Hi = hi
	}
	return db
}

func NewDecodeBytesFromCPU(cpu *M6502) *DecodeBytes {
	return NewDecodeBytes(cpu.Instructions.GetOp(cpu.CurrentOp.OpCode),
		cpu.Registers.PC,
		cpu.Memory.Fetch(cpu.Registers.PC+2),
		cpu.Memory.Fetch(cpu.Registers.PC+1))
}

// This is the disassembly struct.
// Its workings are not required for emulation.
//
// It is merely a convenience function to turn the
// binary instruction code into human readable form.
// Its included as part of the emulator because it can
// take advantage of many of the CPUs internal
// operations to do
type Assembly struct {
	cpu          *M6502
	instructions *InstructionTable
}

func NewAssembly() *Assembly {
	instructions := NewInstructionTable()
	instructions.InitInstructions()

	return &Assembly{
		instructions: &instructions,
	}
}

func (a *Assembly) GetInstructions() *InstructionTable { return a.instructions }

func (a *Assembly) ModeString(mode AddressingMode) string {
	switch mode {
	case Accumulator:
		return "A"
	case Immediate:
		return "#$BB"
	case Implied:
		return ""
	case ZeroPage:
		return "$LL"
	case ZeroPageX:
		return "$LL, X"
	case ZeroPageY:
		return "$LL, Y"
	case Absolute:
		return "$LLHH"
	case AbsoluteX:
		return "$LLHH, X"
	case AbsoluteY:
		return "$LLHH, Y"
	case Indirect:
		return "($LLHH)"
	case IndirectX:
		return "($LL, X)"
	case IndirectY:
		return "($LL), Y"
	case Relative:
		return "$BB [PC + $BB]"
	}
	return ""
}

func (a *Assembly) DisassembleByte(db DecodeBytes) *Decode {
	args := ""
	switch db.Instruction.Bytes() {
	case 2:
		args = fmt.Sprintf("%02X", db.Lo)
	case 3:
		args = fmt.Sprintf("%02X %02X", db.Lo, db.Hi)
	}
	return &Decode{
		pc:          db.PC,
		opcode:      db.Instruction.OpCode,
		args:        args,
		mnemonic:    db.Instruction.Mnemonic,
		decodedArgs: a.ModeString(db.Instruction.AddressingMode),
		registers:   "",
		//registers: a.cpu.Registers.String(),
	}

}

func (a *Assembly) DisassembleAddr(cpu *M6502, addr uint16) *Decode {
	b := a.cpu.Memory.Fetch(addr)
	op := a.instructions.Get(b)
	asm := &Decode{
		pc:          addr,
		opcode:      op.OpCode,
		args:        "",
		mnemonic:    op.Mnemonic,
		decodedArgs: "",
		registers:   "",
		//registers: a.cpu.Registers.String(),
	}
	var lo, hi uint8

	switch op.Bytes() {
	case 2:
		lo = cpu.Memory.Fetch(cpu.Registers.PC + 1)
	case 3:
		lo = cpu.Memory.Fetch(cpu.Registers.PC + 1)
		hi = cpu.Memory.Fetch(cpu.Registers.PC + 2)
	}

	switch op.AddressingMode {
	case Accumulator:
		asm.args = "A"
		asm.decodedArgs = "A"
	case Immediate:
		asm.args = fmt.Sprintf("%02X", lo)
		asm.decodedArgs = "#$"
	case Implied:
		//asm.args = fmt.Sprintf("A")
		//asm.decodedArgs = fmt.Sprintf("A")
		break
	case ZeroPage:
		asm.args = fmt.Sprintf("%02X", lo)
		asm.decodedArgs = fmt.Sprintf("$%02X", lo)
	case ZeroPageX:
		result := (lo + cpu.Registers.X)
		asm.args = fmt.Sprintf("%02X", lo)
		asm.decodedArgs = fmt.Sprintf("$%02X,X @ %02X", lo, result)
	case ZeroPageY:
		result := (lo + cpu.Registers.Y)
		asm.args = fmt.Sprintf("%02X", lo)
		asm.decodedArgs = fmt.Sprintf("$%02X,Y @ %02X", lo, result)
	case Absolute:
		asm.args = fmt.Sprintf("%02X %02X", lo, hi)
		asm.decodedArgs = fmt.Sprintf("$%02X%02X = ", lo, hi)
	case AbsoluteX:
		result := (uint16(hi)<<8 | uint16(lo)) + uint16(cpu.Registers.X)
		asm.args = fmt.Sprintf("%02X %02X", lo, hi)
		asm.decodedArgs = fmt.Sprintf("$%04X = ", result)
	case AbsoluteY:
		result := (uint16(hi)<<8 | uint16(lo)) + uint16(cpu.Registers.Y)
		asm.args = fmt.Sprintf("%02X %02X", lo, hi)
		asm.decodedArgs = fmt.Sprintf("$%04X = ", result)
	case Indirect:
		asm.args = fmt.Sprintf("%02X %02X", lo, hi)
		addrL := (uint16(hi) << 8) | uint16(lo)
		var addrH uint16
		if lo == 0xFF {
			addrH = addrL & 0xff00
		} else {
			addrH = addrL + 1
		}
		resultL := cpu.Memory.Fetch(addrL)
		resultH := cpu.Memory.Fetch(addrH)

		asm.decodedArgs = fmt.Sprintf("(%02X%02X)", resultL, resultH)
	case IndirectX:
		asm.args = fmt.Sprintf("%02X", lo)
		address := uint16(cpu.Memory.Fetch(uint16(lo + cpu.Registers.X)))
		resultL := cpu.Memory.Fetch(address)
		resultH := cpu.Memory.Fetch(address + 1)
		result := uint16(resultH)<<8 | uint16(resultL)

		asm.decodedArgs = fmt.Sprintf("($%02X,X) @ %02X = %04X = ", lo, address, result)
	case IndirectY:
		asm.args = fmt.Sprintf("%02X", lo)
		addrL := uint16(cpu.Memory.Fetch(uint16(lo)))
		addrH := uint16(cpu.Memory.Fetch(uint16((lo + 1) & 0x00FF)))
		address := (uint16(addrH)<<8 | uint16(addrL))
		result := address + uint16(cpu.Registers.Y)

		asm.decodedArgs = fmt.Sprintf("($%02X),Y = %04X @ %04X = ", lo, address, result)
	case Relative:
		var offset uint16
		if lo > 0x7f {
			offset = -(0x0100 - uint16(lo))
		} else {
			offset = uint16(lo)
		}
		asm.args = fmt.Sprintf("%02X", lo)
		asm.decodedArgs = fmt.Sprintf("$%04X", addr+offset)
	}

	return asm
}
