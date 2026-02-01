package m65go2

import (
	"fmt"
	"reflect"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Represents opcodes for the 6502 CPU
type OpCode uint8

// Represents addressing modes for the 6502 CPU
//
//go:generate stringer -type=AddressingMode -linecomment
type AddressingMode uint16

const (
	Accumulator AddressingMode = 1 << iota // Accumulator
	Immediate                              // Immediate
	Implied                                // Implied
	ZeroPage                               // ZeroPage
	ZeroPageX                              // ZeroPageIndexedX
	ZeroPageY                              // ZeroPageIndexedY
	Absolute                               // Absolute
	AbsoluteX                              // AbsoluteIndexedX
	AbsoluteY                              // AbsoluteIndexedY
	Indirect                               // Indirect
	IndirectX                              // IndexedIndirect
	IndirectY                              // IndirectIndexed
	Relative                               // Relative
)

//go:generate stringer -type=ActionMode -linecomment
type ActionMode uint8

const (
	ActionValue         ActionMode = iota // Value
	ActionAddress                         // Address
	ActionTargetAddress                   // TargetAddress
)

// Represents an instruction for the 6502 CPU.  The Exec field
// implements the instruction and returns the total clock cycles to be
// consumed by the instruction.
type Instruction struct {
	Mnemonic       string
	OpCode         OpCode
	AddressingMode AddressingMode
	ActionMode     ActionMode
	Cycles         uint16
	status         InstructionStatus
	Exec           func(*M6502) (status InstructionStatus)
}

func NewInstruction(mnemonic string, opcode OpCode, addrMode AddressingMode, actionMode ActionMode, cycles uint16) *Instruction {
	return &Instruction{
		Mnemonic:       mnemonic,
		OpCode:         opcode,
		AddressingMode: addrMode,
		ActionMode:     actionMode,
		Cycles:         cycles,
		status:         InstructionStatus(0),
	}
}

func (inst Instruction) String() string {
	prefix := fmt.Sprintf("%02X: %s ", inst.OpCode, inst.Mnemonic)

	switch inst.AddressingMode {
	case Accumulator:
		return prefix + "A "
	case Implied:
		return prefix
	case Immediate:
		return prefix + "#$BB "
	case ZeroPage:
		return prefix + "$LL "
	case ZeroPageX:
		return prefix + "$LL, X "
	case ZeroPageY:
		return prefix + "$LL, Y "
	case Absolute:
		return prefix + "$LLHH "
	case AbsoluteX:
		return prefix + "$LLHH, X "
	case AbsoluteY:
		return prefix + "$LLHH, Y "
	case Indirect:
		return prefix + "($LLHH) "
	case IndirectX:
		return prefix + "($LL, X) "
	case IndirectY:
		return prefix + "($LL), Y "
	case Relative:
		return prefix + "$BB [PC + $BB] "
	default:
		return prefix + "??? "
	}
}

// Number of bytes operation reads
func (inst *Instruction) Bytes() int {
	switch inst.AddressingMode {
	case Accumulator, Implied:
		return 1
	case Immediate, ZeroPage, ZeroPageX, ZeroPageY, IndirectX, IndirectY, Relative:
		return 2
	//case Absolute, AbsoluteX, AbsoluteY, Indirect:
	//	return 3
	default:
		return 3
	}
}

// Stores instructions understood by the 6502 CPU, indexed by opcode.
type InstructionTable struct {
	opcodes []*Instruction
}

type InstructionStatus uint16

const (
	PageCross InstructionStatus = 1 << iota
	Branched
)

// Returns a new, empty InstructionTable
func NewInstructionTable() InstructionTable {
	instructions := InstructionTable{
		opcodes: make([]*Instruction, 0x100),
	}

	return instructions
}

// Get all instructions from the InstructionTable
func (instructions InstructionTable) GetAll() []*Instruction {
	return instructions.opcodes
}

// Gets an instruction from the InstructionTable by uint8 opcode
func (instructions InstructionTable) Get(opcode uint8) *Instruction {
	return instructions.GetOp(OpCode(opcode))
}

// Gets an instruction from the InstructionTable by opcode
func (instructions InstructionTable) GetOp(opcode OpCode) *Instruction {
	return instructions.opcodes[opcode]
}

// Executes an instruction in the InstructionTable, returns number of
// cycles taken to execute
func (instructions InstructionTable) Execute(cpu *M6502) uint16 {
	inst := cpu.CurrentOp

	opName := cases.Title(language.English).String(inst.Mnemonic)
	cpuReflector := reflect.ValueOf(cpu)
	method := cpuReflector.MethodByName(opName)

	callParam := make([]reflect.Value, 0)

	switch inst.ActionMode {
	case ActionValue:
		value := cpu.getAddressValue(inst.AddressingMode, &inst.status)
		callParam = append(callParam, reflect.ValueOf(value))
	case ActionAddress:
		address := cpu.getAddress(inst.AddressingMode, &inst.status)
		callParam = append(callParam, reflect.ValueOf(address))
	}

	method.Call(callParam)

	return cpu.currentCycles
}

// Adds an instruction to the InstructionTable
func (instructions InstructionTable) AddInstruction(inst *Instruction) {
	instructions.opcodes[inst.OpCode] = inst
}

// Removes any instruction with the given opcode
func (instructions InstructionTable) RemoveInstruction(opcode OpCode) {
	instructions.opcodes[opcode] = nil
}

// Adds the 6502 CPU's instruction set to the InstructionTable.
func (i *InstructionTable) InitInstructions() {
	// region Instructions 0x0
	i.AddInstruction(NewInstruction("BRK", 0x00, Implied, ActionValue, 7))
	i.AddInstruction(NewInstruction("ORA", 0x01, IndirectX, ActionValue, 6))
	i.AddInstruction(NewInstruction("KIL", 0x02, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("SLO", 0x03, IndirectX, ActionAddress, 8))
	i.AddInstruction(NewInstruction("IGN", 0x04, ZeroPage, ActionAddress, 3))
	i.AddInstruction(NewInstruction("ORA", 0x05, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("ASL", 0x06, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("SLO", 0x07, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("PHP", 0x08, Implied, ActionAddress, 3))
	i.AddInstruction(NewInstruction("ORA", 0x09, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("ASL", 0x0A, Accumulator, ActionAddress, 2))
	i.AddInstruction(NewInstruction("ANC", 0x0B, Immediate, ActionValue, 2)) // AAC
	i.AddInstruction(NewInstruction("IGN", 0x0C, Absolute, ActionAddress, 4))
	i.AddInstruction(NewInstruction("ORA", 0x0D, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("ASL", 0x0E, Absolute, ActionAddress, 6))
	i.AddInstruction(NewInstruction("SLO", 0x0F, Absolute, ActionAddress, 6))
	// endregion Instructions 0x0

	// region Instructions 0x1
	i.AddInstruction(NewInstruction("BPL", 0x10, Relative, ActionAddress, 2))
	i.AddInstruction(NewInstruction("ORA", 0x11, IndirectY, ActionValue, 5))
	i.AddInstruction(NewInstruction("KIL", 0x12, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("SLO", 0x13, IndirectY, ActionAddress, 8))
	i.AddInstruction(NewInstruction("IGN", 0x14, ZeroPageX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("ORA", 0x15, ZeroPageX, ActionValue, 4))
	i.AddInstruction(NewInstruction("ASL", 0x16, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("SLO", 0x17, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("CLC", 0x18, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("ORA", 0x19, AbsoluteY, ActionValue, 4))
	i.AddInstruction(NewInstruction("NOP", 0x1A, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("SLO", 0x1B, AbsoluteY, ActionAddress, 7))
	i.AddInstruction(NewInstruction("IGN", 0x1C, AbsoluteX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("ORA", 0x1D, AbsoluteX, ActionValue, 4))
	i.AddInstruction(NewInstruction("ASL", 0x1E, AbsoluteX, ActionAddress, 7))
	i.AddInstruction(NewInstruction("SLO", 0x1F, AbsoluteX, ActionAddress, 7))
	// endregion Instructions 0x1

	// region Instructions 0x2
	i.AddInstruction(NewInstruction("JSR", 0x20, Absolute, ActionAddress, 6))
	i.AddInstruction(NewInstruction("AND", 0x21, IndirectX, ActionValue, 6))
	i.AddInstruction(NewInstruction("KIL", 0x22, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("RLA", 0x23, IndirectX, ActionAddress, 8))
	i.AddInstruction(NewInstruction("BIT", 0x24, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("AND", 0x25, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("ROL", 0x26, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("RLA", 0x27, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("PLP", 0x28, Implied, ActionAddress, 4))
	i.AddInstruction(NewInstruction("AND", 0x29, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("ROL", 0x2A, Accumulator, ActionAddress, 2))
	i.AddInstruction(NewInstruction("AAC", 0x2B, Immediate, ActionValue, 2)) // ANC / ANC2
	i.AddInstruction(NewInstruction("BIT", 0x2C, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("AND", 0x2D, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("ROL", 0x2E, Absolute, ActionAddress, 6))
	i.AddInstruction(NewInstruction("RLA", 0x2F, Absolute, ActionAddress, 6))
	// endregion Instructions 0x2

	// region Instructions 0x3
	i.AddInstruction(NewInstruction("BMI", 0x30, Relative, ActionAddress, 2))
	i.AddInstruction(NewInstruction("AND", 0x31, IndirectY, ActionValue, 5))
	i.AddInstruction(NewInstruction("KIL", 0x32, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("RLA", 0x33, IndirectY, ActionAddress, 8))
	i.AddInstruction(NewInstruction("IGN", 0x34, ZeroPageX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("AND", 0x35, ZeroPageX, ActionValue, 4))
	i.AddInstruction(NewInstruction("ROL", 0x36, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("RLA", 0x37, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("SEC", 0x38, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("AND", 0x39, AbsoluteY, ActionValue, 4))
	i.AddInstruction(NewInstruction("NOP", 0x3A, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("RLA", 0x3B, AbsoluteY, ActionAddress, 7))
	i.AddInstruction(NewInstruction("IGN", 0x3C, AbsoluteX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("AND", 0x3D, AbsoluteX, ActionValue, 4))
	i.AddInstruction(NewInstruction("ROL", 0x3E, AbsoluteX, ActionAddress, 7))
	i.AddInstruction(NewInstruction("RLA", 0x3F, AbsoluteX, ActionAddress, 7))
	// endregion Instructions 0x3

	// region Instructions 0x4
	i.AddInstruction(NewInstruction("RTI", 0x40, Implied, ActionAddress, 6))
	i.AddInstruction(NewInstruction("EOR", 0x41, IndirectX, ActionValue, 6))
	i.AddInstruction(NewInstruction("KIL", 0x42, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("SRE", 0x43, IndirectX, ActionAddress, 8))
	i.AddInstruction(NewInstruction("IGN", 0x44, ZeroPage, ActionAddress, 3))
	i.AddInstruction(NewInstruction("EOR", 0x45, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("LSR", 0x46, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("SRE", 0x47, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("PHA", 0x48, Implied, ActionAddress, 3))
	i.AddInstruction(NewInstruction("EOR", 0x49, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("LSR", 0x4A, Accumulator, ActionAddress, 2))
	i.AddInstruction(NewInstruction("ALR", 0x4B, Immediate, ActionValue, 2)) // ASR
	i.AddInstruction(NewInstruction("JMP", 0x4C, Absolute, ActionAddress, 3))
	i.AddInstruction(NewInstruction("EOR", 0x4D, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("LSR", 0x4E, Absolute, ActionAddress, 6))
	i.AddInstruction(NewInstruction("SRE", 0x4F, Absolute, ActionAddress, 6))
	// endregion Instructions 0x4

	// region Instructions 0x5
	i.AddInstruction(NewInstruction("BVC", 0x50, Relative, ActionAddress, 2))
	i.AddInstruction(NewInstruction("EOR", 0x51, IndirectY, ActionValue, 5))
	i.AddInstruction(NewInstruction("KIL", 0x52, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("SRE", 0x53, IndirectY, ActionAddress, 8))
	i.AddInstruction(NewInstruction("IGN", 0x54, ZeroPageX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("EOR", 0x55, ZeroPageX, ActionValue, 4))
	i.AddInstruction(NewInstruction("LSR", 0x56, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("SRE", 0x57, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("CLI", 0x58, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("EOR", 0x59, AbsoluteY, ActionValue, 4))
	i.AddInstruction(NewInstruction("NOP", 0x5A, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("SRE", 0x5B, AbsoluteY, ActionAddress, 7))
	i.AddInstruction(NewInstruction("IGN", 0x5C, AbsoluteX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("EOR", 0x5D, AbsoluteX, ActionValue, 4))
	i.AddInstruction(NewInstruction("LSR", 0x5E, AbsoluteX, ActionAddress, 7))
	i.AddInstruction(NewInstruction("SRE", 0x5F, AbsoluteX, ActionAddress, 7))
	// endregion Instructions 0x5

	// region Instructions 0x6
	i.AddInstruction(NewInstruction("RTS", 0x60, Implied, ActionAddress, 6))
	i.AddInstruction(NewInstruction("ADC", 0x61, IndirectX, ActionValue, 6))
	i.AddInstruction(NewInstruction("KIL", 0x62, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("RRA", 0x63, IndirectX, ActionAddress, 8))
	i.AddInstruction(NewInstruction("IGN", 0x64, ZeroPage, ActionAddress, 3))
	i.AddInstruction(NewInstruction("ADC", 0x65, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("ROR", 0x66, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("RRA", 0x67, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("PLA", 0x68, Implied, ActionAddress, 4))
	i.AddInstruction(NewInstruction("ADC", 0x69, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("ROR", 0x6A, Accumulator, ActionAddress, 2))
	i.AddInstruction(NewInstruction("ARR", 0x6B, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("JMP", 0x6C, Indirect, ActionAddress, 5))
	i.AddInstruction(NewInstruction("ADC", 0x6D, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("ROR", 0x6E, Absolute, ActionAddress, 6))
	i.AddInstruction(NewInstruction("RRA", 0x6F, Absolute, ActionAddress, 6))
	// endregion Instructions 0x6

	// region Instructions 0x7
	i.AddInstruction(NewInstruction("BVS", 0x70, Relative, ActionAddress, 2))
	i.AddInstruction(NewInstruction("ADC", 0x71, IndirectY, ActionValue, 5))
	i.AddInstruction(NewInstruction("KIL", 0x72, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("RRA", 0x73, IndirectY, ActionAddress, 8))
	i.AddInstruction(NewInstruction("IGN", 0x74, ZeroPageX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("ADC", 0x75, ZeroPageX, ActionValue, 4))
	i.AddInstruction(NewInstruction("ROR", 0x76, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("RRA", 0x77, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("SEI", 0x78, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("ADC", 0x79, AbsoluteY, ActionValue, 4))
	i.AddInstruction(NewInstruction("NOP", 0x7A, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("RRA", 0x7B, AbsoluteY, ActionAddress, 7))
	i.AddInstruction(NewInstruction("IGN", 0x7C, AbsoluteX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("ADC", 0x7D, AbsoluteX, ActionValue, 4))
	i.AddInstruction(NewInstruction("ROR", 0x7E, AbsoluteX, ActionAddress, 7))
	i.AddInstruction(NewInstruction("RRA", 0x7F, AbsoluteX, ActionAddress, 7))
	// endregion Instructions 0x7

	// region Instructions 0x8
	i.AddInstruction(NewInstruction("SKB", 0x80, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("STA", 0x81, IndirectX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("SKB", 0x82, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("SAX", 0x83, IndirectX, ActionAddress, 6)) // AAX
	i.AddInstruction(NewInstruction("STY", 0x84, ZeroPage, ActionAddress, 3))
	i.AddInstruction(NewInstruction("STA", 0x85, ZeroPage, ActionAddress, 3))
	i.AddInstruction(NewInstruction("STX", 0x86, ZeroPage, ActionAddress, 3))
	i.AddInstruction(NewInstruction("SAX", 0x87, ZeroPage, ActionAddress, 3)) // AAX
	i.AddInstruction(NewInstruction("DEY", 0x88, Implied, ActionAddress, 2))
	i.AddInstruction(NewInstruction("SKB", 0x89, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("TXA", 0x8A, Implied, ActionAddress, 2))
	i.AddInstruction(NewInstruction("ANE", 0x8B, Immediate, ActionValue, 2)) // XAA
	i.AddInstruction(NewInstruction("STY", 0x8C, Absolute, ActionAddress, 4))
	i.AddInstruction(NewInstruction("STA", 0x8D, Absolute, ActionAddress, 4))
	i.AddInstruction(NewInstruction("STX", 0x8E, Absolute, ActionAddress, 4))
	i.AddInstruction(NewInstruction("SAX", 0x8F, Absolute, ActionAddress, 4)) // AAX
	// endregion Instructions 0x8

	// region Instructions 0x9
	i.AddInstruction(NewInstruction("BCC", 0x90, Relative, ActionAddress, 2))
	i.AddInstruction(NewInstruction("STA", 0x91, IndirectY, ActionAddress, 6))
	i.AddInstruction(NewInstruction("KIL", 0x92, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("SHA", 0x93, IndirectY, ActionAddress, 6)) // AXA
	i.AddInstruction(NewInstruction("STY", 0x94, ZeroPageX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("STA", 0x95, ZeroPageX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("STX", 0x96, ZeroPageY, ActionAddress, 4))
	i.AddInstruction(NewInstruction("SAX", 0x97, ZeroPageY, ActionAddress, 4)) // AAX
	i.AddInstruction(NewInstruction("TYA", 0x98, Implied, ActionAddress, 2))
	i.AddInstruction(NewInstruction("STA", 0x99, AbsoluteY, ActionAddress, 5))
	i.AddInstruction(NewInstruction("TXS", 0x9A, Implied, ActionAddress, 2))
	i.AddInstruction(NewInstruction("TAS", 0x9B, AbsoluteY, ActionAddress, 5)) // SHS, XAS
	i.AddInstruction(NewInstruction("SHY", 0x9C, AbsoluteX, ActionAddress, 5)) // SHY
	i.AddInstruction(NewInstruction("STA", 0x9D, AbsoluteX, ActionAddress, 5))
	i.AddInstruction(NewInstruction("SHX", 0x9E, AbsoluteY, ActionAddress, 5))
	i.AddInstruction(NewInstruction("SHA", 0x9F, AbsoluteY, ActionAddress, 5)) // AXA
	// endregion Instructions 0x9

	// region Instructions 0xA
	i.AddInstruction(NewInstruction("LDY", 0xA0, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("LDA", 0xA1, IndirectX, ActionValue, 6))
	i.AddInstruction(NewInstruction("LDX", 0xA2, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("LAX", 0xA3, IndirectX, ActionValue, 6))
	i.AddInstruction(NewInstruction("LDY", 0xA4, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("LDA", 0xA5, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("LDX", 0xA6, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("LAX", 0xA7, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("TAY", 0xA8, Implied, ActionAddress, 2))
	i.AddInstruction(NewInstruction("LDA", 0xA9, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("TAX", 0xAA, Implied, ActionAddress, 2))
	i.AddInstruction(NewInstruction("LXA", 0xAB, Immediate, ActionValue, 2)) // LAX Immediate
	i.AddInstruction(NewInstruction("LDY", 0xAC, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("LDA", 0xAD, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("LDX", 0xAE, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("LAX", 0xAF, Absolute, ActionValue, 4))
	// endregion Instructions 0xA

	// region Instructions 0xB
	i.AddInstruction(NewInstruction("BCS", 0xB0, Relative, ActionAddress, 2))
	i.AddInstruction(NewInstruction("LDA", 0xB1, IndirectY, ActionValue, 5))
	i.AddInstruction(NewInstruction("KIL", 0xB2, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("LAX", 0xB3, IndirectY, ActionValue, 5))
	i.AddInstruction(NewInstruction("LDY", 0xB4, ZeroPageX, ActionValue, 4))
	i.AddInstruction(NewInstruction("LDA", 0xB5, ZeroPageX, ActionValue, 4))
	i.AddInstruction(NewInstruction("LDX", 0xB6, ZeroPageY, ActionValue, 4))
	i.AddInstruction(NewInstruction("LAX", 0xB7, ZeroPageY, ActionValue, 4))
	i.AddInstruction(NewInstruction("CLV", 0xB8, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("LDA", 0xB9, AbsoluteY, ActionValue, 4))
	i.AddInstruction(NewInstruction("TSX", 0xBA, Implied, ActionAddress, 2))
	i.AddInstruction(NewInstruction("LAS", 0xBB, AbsoluteY, ActionValue, 4)) // LAR
	i.AddInstruction(NewInstruction("LDY", 0xBC, AbsoluteX, ActionValue, 4))
	i.AddInstruction(NewInstruction("LDA", 0xBD, AbsoluteX, ActionValue, 4))
	i.AddInstruction(NewInstruction("LDX", 0xBE, AbsoluteY, ActionValue, 4))
	i.AddInstruction(NewInstruction("LAX", 0xBF, AbsoluteY, ActionValue, 4))
	// endregion Instructions 0xB

	// region Instructions 0xC
	i.AddInstruction(NewInstruction("CPY", 0xC0, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("CMP", 0xC1, IndirectX, ActionValue, 6))
	i.AddInstruction(NewInstruction("SKB", 0xC2, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("DCP", 0xC3, IndirectX, ActionAddress, 8))
	i.AddInstruction(NewInstruction("CPY", 0xC4, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("CMP", 0xC5, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("DEC", 0xC6, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("DCP", 0xC7, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("INY", 0xC8, Implied, ActionAddress, 2))
	i.AddInstruction(NewInstruction("CMP", 0xC9, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("DEX", 0xCA, Implied, ActionAddress, 2))
	i.AddInstruction(NewInstruction("SBX", 0xCB, Immediate, ActionValue, 2)) // AXS
	i.AddInstruction(NewInstruction("CPY", 0xCC, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("CMP", 0xCD, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("DEC", 0xCE, Absolute, ActionAddress, 6))
	i.AddInstruction(NewInstruction("DCP", 0xCF, Absolute, ActionAddress, 6))
	// endregion Instructions 0xC

	// region Instructions 0xD
	i.AddInstruction(NewInstruction("BNE", 0xD0, Relative, ActionAddress, 2))
	i.AddInstruction(NewInstruction("CMP", 0xD1, IndirectY, ActionValue, 5))
	i.AddInstruction(NewInstruction("KIL", 0xD2, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("DCP", 0xD3, IndirectY, ActionAddress, 8))
	i.AddInstruction(NewInstruction("IGN", 0xD4, ZeroPageX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("CMP", 0xD5, ZeroPageX, ActionValue, 4))
	i.AddInstruction(NewInstruction("DEC", 0xD6, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("DCP", 0xD7, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("CLD", 0xD8, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("CMP", 0xD9, AbsoluteY, ActionValue, 4))
	i.AddInstruction(NewInstruction("NOP", 0xDA, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("DCP", 0xDB, AbsoluteY, ActionAddress, 7))
	i.AddInstruction(NewInstruction("IGN", 0xDC, AbsoluteX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("CMP", 0xDD, AbsoluteX, ActionValue, 4))
	i.AddInstruction(NewInstruction("DEC", 0xDE, AbsoluteX, ActionAddress, 7))
	i.AddInstruction(NewInstruction("DCP", 0xDF, AbsoluteX, ActionAddress, 7))
	// endregion Instructions 0xD

	// region Instructions 0xE
	i.AddInstruction(NewInstruction("CPX", 0xE0, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("SBC", 0xE1, IndirectX, ActionValue, 6))
	i.AddInstruction(NewInstruction("SKB", 0xE2, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("ISB", 0xE3, IndirectX, ActionAddress, 8))
	i.AddInstruction(NewInstruction("CPX", 0xE4, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("SBC", 0xE5, ZeroPage, ActionValue, 3))
	i.AddInstruction(NewInstruction("INC", 0xE6, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("ISB", 0xE7, ZeroPage, ActionAddress, 5))
	i.AddInstruction(NewInstruction("INX", 0xE8, Implied, ActionAddress, 2))
	i.AddInstruction(NewInstruction("SBC", 0xE9, Immediate, ActionValue, 2))
	i.AddInstruction(NewInstruction("NOP", 0xEA, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("SBC", 0xEB, Immediate, ActionValue, 2)) // USBC
	i.AddInstruction(NewInstruction("CPX", 0xEC, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("SBC", 0xED, Absolute, ActionValue, 4))
	i.AddInstruction(NewInstruction("INC", 0xEE, Absolute, ActionAddress, 6))
	i.AddInstruction(NewInstruction("ISB", 0xEF, Absolute, ActionAddress, 6))
	// endregion Instructions 0xE

	// region Instructions 0xF
	i.AddInstruction(NewInstruction("BEQ", 0xF0, Relative, ActionAddress, 2))
	i.AddInstruction(NewInstruction("SBC", 0xF1, IndirectY, ActionValue, 5))
	i.AddInstruction(NewInstruction("KIL", 0xF2, Implied, ActionValue, 1))
	i.AddInstruction(NewInstruction("ISB", 0xF3, IndirectY, ActionAddress, 8))
	i.AddInstruction(NewInstruction("IGN", 0xF4, ZeroPageX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("SBC", 0xF5, ZeroPageX, ActionValue, 4))
	i.AddInstruction(NewInstruction("INC", 0xF6, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("ISB", 0xF7, ZeroPageX, ActionAddress, 6))
	i.AddInstruction(NewInstruction("SED", 0xF8, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("SBC", 0xF9, AbsoluteY, ActionValue, 4))
	i.AddInstruction(NewInstruction("NOP", 0xFA, Implied, ActionValue, 2))
	i.AddInstruction(NewInstruction("ISB", 0xFB, AbsoluteY, ActionAddress, 7))
	i.AddInstruction(NewInstruction("IGN", 0xFC, AbsoluteX, ActionAddress, 4))
	i.AddInstruction(NewInstruction("SBC", 0xFD, AbsoluteX, ActionValue, 4))
	i.AddInstruction(NewInstruction("INC", 0xFE, AbsoluteX, ActionAddress, 7))
	i.AddInstruction(NewInstruction("ISB", 0xFF, AbsoluteX, ActionAddress, 7))
	// endregion Instructions 0xF
}

/*
// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}
*/
// Helper function to create any pointers
func asPtr(v any) *any {
	return &v
}
