// Package m65go2 simulates the MOS 6502 CPU
package m65go2

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

// Flags used by P (Status) register
type Status uint8

const (
	C Status = 1 << iota // carry flag
	Z                    // zero flag
	I                    // interrupt disable
	D                    // decimal mode
	B                    // break command
	U                    // -UNUSED-
	V                    // overflow flag
	N                    // negative flag
)

// The 6502's registers, all registers are 8-bit values except for PC
// which is 16-bits.
type Registers struct {
	A  uint8  // accumulator
	X  uint8  // index register X
	Y  uint8  // index register Y
	P  Status // processor status
	SP uint8  // stack pointer
	PC uint16 // program counter
}

func (s *Status) String() string {
	var flags = []string{".", ".", ".", ".", ".", ".", ".", "."}
	if *s&N != 0 {
		flags[0] = "N"
	}
	if *s&Z != 0 {
		flags[1] = "Z"
	}
	if *s&I != 0 {
		flags[2] = "I"
	}
	if *s&D != 0 {
		flags[3] = "D"
	}
	if *s&B != 0 {
		flags[4] = "B"
	}
	if *s&U != 0 {
		flags[5] = "U"
	}
	if *s&V != 0 {
		flags[6] = "V"
	}
	if *s&C != 0 {
		flags[7] = "C"
	}

	return strings.Join(flags, "")
}

// Creates a new set of Registers.  All registers are initialized to
// 0.
func NewRegisters() (reg Registers) {
	reg = Registers{}
	reg.Reset()
	return
}

// Resets all registers.  Register P is initialized with only the I
// bit set, SP is initialized to 0xfd, PC is initialized to 0xfffc
// (the RESET vector) and all other registers are initialized to 0.
func (reg *Registers) Reset() {
	reg.A = 0
	reg.X = 0
	reg.Y = 0
	reg.P = I | U
	reg.SP = 0xfd
	reg.PC = 0xfffc
}

// Prints the values of each register to os.Stderr.
func (reg *Registers) String() string {
	return fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X (%s)", reg.A, reg.X, reg.Y, reg.P, reg.SP, reg.P.String())
}

type Interrupt uint8

const (
	Irq Interrupt = iota
	Nmi
	Rst
)

//go:generate stringer -type=Index
type Index uint8

const (
	X Index = iota
	Y
)

type decode struct {
	enabled     bool
	pc          uint16
	opcode      OpCode
	args        string
	mnemonic    string
	decodedArgs string
	registers   string
	ticks       uint64
}

func (d *decode) String() string {
	return fmt.Sprintf("%04X  %02X %-5s %4s %-26s  %25s",
		d.pc, d.opcode, d.args, d.mnemonic, d.decodedArgs, d.registers)
}

// Represents the 6502 CPU.
type M6502 struct {
	decode        decode
	debug         DebugMode
	log           *log.Logger
	Nmi           bool
	Irq           bool
	Rst           bool
	Registers     Registers        `json:"M6502Registers"`
	Memory        Memory           `json:"M6502Memory"`
	Instructions  InstructionTable `json:"-"`
	CurrentOp     *Instruction     `json:"-"`
	decimalMode   bool
	breakError    bool
	freezeError   bool
	currentCycles uint16
	magic         uint8 // Used by some unofficial opcodes
}

// Returns a pointer to a new CPU with the given Memory.
func NewM6502(mem Memory) *M6502 {
	instructions := NewInstructionTable()
	instructions.InitInstructions()
	logger := log.Default()
	//logger.SetFlags(log.LstdFlags | log.Lshortfile)

	return &M6502{
		decode:        decode{},
		log:           logger,
		debug:         0,
		Registers:     NewRegisters(),
		Memory:        mem,
		Instructions:  instructions,
		CurrentOp:     nil,
		decimalMode:   true,
		breakError:    false,
		freezeError:   false,
		Nmi:           false,
		Irq:           false,
		Rst:           false,
		currentCycles: 0,
		magic:         0x00,
	}
}

// Resets the CPU by resetting both the registers and memory.
func (cpu *M6502) Reset() {
	cpu.Registers.Reset()
	cpu.Memory.Reset()
	cpu.PerformRst()
}

// region IRQ/NMI/RST
func (cpu *M6502) Interrupt(which Interrupt, state bool) {
	switch which {
	case Irq:
		cpu.Irq = state
	case Nmi:
		cpu.Nmi = state
	case Rst:
		cpu.Rst = state
	}
}

func (cpu *M6502) InterruptLine(which Interrupt) func(state bool) {
	return func(state bool) {
		if cpu != nil {
			cpu.Interrupt(which, state)
		}
	}
}

func (cpu *M6502) GetInterrupt(which Interrupt) (state bool) {
	switch which {
	case Irq:
		state = cpu.Irq
	case Nmi:
		state = cpu.Nmi
	case Rst:
		state = cpu.Rst
	}

	return
}

func (cpu *M6502) PerformInterrupts() (cycles uint16) {
	// check interrupts
	cycles = 7

	switch {
	case cpu.Irq && cpu.Registers.P&I == 0:
		cpu.PerformIrq()
		cpu.Irq = false
	case cpu.Nmi:
		cpu.PerformNmi()
		cpu.Nmi = false
	case cpu.Rst:
		cpu.PerformRst()
		cpu.Rst = false
	default:
		cycles = 0
	}

	return
}

func (cpu *M6502) PerformIrq() {
	cpu.push16(cpu.Registers.PC)
	cpu.push(uint8((cpu.Registers.P | U) & ^B))

	cpu.Registers.P |= I

	//low := cpu.Memory.Fetch(0xfffe)
	//high := cpu.Memory.Fetch(0xffff)
	low := cpu.readAddress(0xfffe)
	high := cpu.readAddress(0xffff)

	cpu.Registers.PC = (uint16(high) << 8) | uint16(low)
}

func (cpu *M6502) PerformNmi() {
	cpu.push16(cpu.Registers.PC)
	cpu.push(uint8((cpu.Registers.P | U) & ^B))

	cpu.Registers.P |= I

	//low := cpu.Memory.Fetch(0xfffa)
	//high := cpu.Memory.Fetch(0xfffb)
	low := cpu.readAddress(0xfffa)
	high := cpu.readAddress(0xfffb)

	cpu.Registers.PC = (uint16(high) << 8) | uint16(low)
}

func (cpu *M6502) PerformRst() {
	//low := cpu.Memory.Fetch(0xfffc)
	//high := cpu.Memory.Fetch(0xfffd)
	low := cpu.readAddress(0xfffc)
	high := cpu.readAddress(0xfffd)

	cpu.Registers.PC = (uint16(high) << 8) | uint16(low)
	cpu.freezeError = false
}

// endregion IRQ/NMI/RST

func (cpu *M6502) DisableDecimalMode() {
	cpu.decimalMode = false
}

func (cpu *M6502) EnableDecode() {
	cpu.decode.enabled = true
}

func (cpu *M6502) ToggleDecode() bool {
	cpu.decode.enabled = !cpu.decode.enabled
	return cpu.decode.enabled
}

// Error type used to indicate that the CPU attempted to execute an
// invalid opcode
type BadOpCodeError OpCode

func (b BadOpCodeError) Error() string {
	return fmt.Sprintf("No such opcode %#02x", OpCode(b))
}

// Error type used to indicate that the CPU executed a BRK instruction
type BrkOpCodeError OpCode

func (b BrkOpCodeError) Error() string {
	return "Executed BRK opcode"
}

// Error type used to indicate that the CPU is in a freezed state
// The processor is trapped infinitely in T1 phase with $FF on
// the data bus.
// **Reset required**
type CpuFreezeError error

// Executes the instruction pointed to by the PC register in the
// number of cycles as returned by the instruction's Exec function.
// Returns the number of cycles executed and any error (such as
// BadOpCodeError).
func (cpu *M6502) Execute() (cycles uint16, error error) {
	if cpu.freezeError {
		err := CpuFreezeError(fmt.Errorf("CPU is in a freezed state"))
		return 0, err
	}
	// Reset current cycle count
	cpu.currentCycles = 0

	// check interrupts
	//cycles += cpu.PerformInterrupts()
	cpu.PerformInterrupts()

	// fetch
	//opcode := OpCode(cpu.Memory.Fetch(cpu.Registers.PC))
	opcode := OpCode(cpu.readAddress(cpu.Registers.PC))
	inst := cpu.Instructions.opcodes[opcode]

	if inst == nil {
		return 0, BadOpCodeError(opcode)
	}

	cpu.CurrentOp = inst
	inst = nil

	// execute
	if cpu.decode.enabled {
		cpu.decode.pc = cpu.Registers.PC
		cpu.decode.opcode = opcode
		cpu.decode.args = ""
		cpu.decode.mnemonic = cpu.CurrentOp.Mnemonic
		cpu.decode.decodedArgs = ""
		cpu.decode.registers = cpu.Registers.String()
	}

	cpu.Registers.PC++
	cycles += cpu.Instructions.Execute(cpu)

	if cpu.debug&DebugModeVerbose != 0 && cpu.decode.enabled {
		cpu.log.Println(cpu.decode.String())
	}

	cpu.CurrentOp = nil

	if cpu.breakError && opcode == 0x00 {
		return cycles, BrkOpCodeError(opcode)
	}

	return cycles, nil
}

// Executes instruction until Execute() returns an error.
func (cpu *M6502) Run() (err error) {
	for {
		if _, err = cpu.Execute(); err != nil {
			return
		}
	}
}

func (cpu *M6502) IncrementCycle() {
	if cpu.debug&DebugModeCycleDump != 0 {
		frames := GetStackFrames()
		mnemonic := ""
		if cpu.CurrentOp != nil {
			mnemonic = cpu.CurrentOp.Mnemonic
		}
		if len(frames) > 0 {
			f := cpu.log.Flags()
			cpu.log.SetFlags(log.LstdFlags)
			cpu.log.Printf("%s: cycle incremented:\n", mnemonic)
			for i, frame := range frames {
				cpu.log.Printf("%s: %s %s (%s:%d)\n",
					mnemonic,
					strings.Repeat("  ", i),
					frame.Function, frame.File, frame.Line)
			}
			cpu.log.SetFlags(f)
		}
	}
	cpu.currentCycles++
}

// region Read/Write Memory
func (cpu *M6502) readAddress(address uint16) uint8 {
	value := cpu.Memory.Fetch(address)
	cpu.IncrementCycle()
	return value
}
func (cpu *M6502) readNext() uint8 {
	value := cpu.readAddress(cpu.Registers.PC)
	cpu.Registers.PC++
	return value
}

func (cpu *M6502) store(address uint16, value uint8) {
	oldValue := cpu.Memory.Store(address, value)
	cpu.IncrementCycle()

	if cpu.decode.enabled {
		if !strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", oldValue)
	}
}

// endregion Read/Write Memory

// region Flags
func (cpu *M6502) getFlag(value Status) bool {
	return cpu.Registers.P&value != 0
}

func (cpu *M6502) setFlag(flag Status, value bool) {
	if value {
		cpu.Registers.P |= flag
	} else {
		cpu.Registers.P &= ^flag
	}
}

func (cpu *M6502) setZFlag(value uint8) uint8 {
	cpu.setFlag(Z, value == 0)
	return value
}

func (cpu *M6502) setNFlag(value uint8) uint8 {
	cpu.setFlag(N, value&0x80 != 0)
	return value
}

func (cpu *M6502) setZNFlags(value uint8) uint8 {
	cpu.setZFlag(value)
	cpu.setNFlag(value)
	return value
}

func (cpu *M6502) setCFlagAddition(value uint16) uint16 {
	cpu.Registers.P = (cpu.Registers.P & ^C) | Status(value>>8&uint16(C))
	return value
}

func (cpu *M6502) setVFlagAddition(term1 uint16, term2 uint16, result uint16) uint16 {
	cpu.Registers.P = (cpu.Registers.P & ^V) | Status((^(term1^term2)&(term1^result)&uint16(N))>>1)
	return result
}

// endregion Flags

// region Addressing Modes
func (cpu *M6502) getAddress(mode AddressingMode, status *InstructionStatus) uint16 {
	var addr uint16
	switch mode {
	case Accumulator, Implied:
		//cpu.readNext()
		cpu.IncrementCycle()
	//case Accumulator:
	//	cpu.readNext()
	case Immediate:
		addr = cpu.immediateAddress()
	//case Implied:
	//	cpu.readNext()
	case ZeroPage:
		addr = cpu.zeroPageAddress()
	case ZeroPageX:
		addr = cpu.zeroPageIndexedAddress(X)
	case ZeroPageY:
		addr = cpu.zeroPageIndexedAddress(Y)
	case Absolute:
		addr = cpu.absoluteAddress()
	case AbsoluteX:
		addr = cpu.absoluteIndexedAddress(X, status)
	case AbsoluteY:
		addr = cpu.absoluteIndexedAddress(Y, status)
	case Indirect:
		addr = cpu.indirectAddress()
	case IndirectX:
		addr = cpu.indexedIndirectAddress()
	case IndirectY:
		addr = cpu.indirectIndexedAddress(status)
	case Relative:
		addr = cpu.relativeAddress()
	default:
		panic("unknown addressing mode: " + mode.String())
	}
	return addr
}

func (cpu *M6502) getAddressValue(mode AddressingMode, status *InstructionStatus) uint8 {
	addr := cpu.getAddress(mode, status)
	if mode == Implied {
		return uint8(0)
	}
	//if mode == Immediate {
	//	return uint8(addr)
	//}
	//return cpu.Memory.Fetch(addr)
	return cpu.readAddress(addr)
}

func (cpu *M6502) containsOpCode(ops []OpCode, opcode OpCode) bool {
	return slices.Contains(ops, opcode)
}

// Immediate Addressing Mode
func (cpu *M6502) immediateAddress() (result uint16) {
	result = cpu.Registers.PC
	cpu.Registers.PC++

	if cpu.decode.enabled {
		value := cpu.Memory.Fetch(result)
		cpu.decode.args = fmt.Sprintf("%02X", value)
		cpu.decode.decodedArgs = "#$"
	}

	return
}

// Zero Page Addressing Mode
func (cpu *M6502) zeroPageAddress() (result uint16) {
	//result = uint16(cpu.Memory.Fetch(cpu.Registers.PC))
	//cpu.Registers.PC++
	result = uint16(cpu.readNext())

	if cpu.decode.enabled {
		cpu.decode.args = fmt.Sprintf("%02X", result)
		cpu.decode.decodedArgs = fmt.Sprintf("$%02X", result)
	}

	return
}

// Helper - Index to Register
func (cpu *M6502) IndexToRegister(which Index) uint8 {
	var index uint8

	switch which {
	case X:
		index = cpu.Registers.X
	case Y:
		index = cpu.Registers.Y
	}

	return index
}

// Zero Page Indexed Addressing Mode
func (cpu *M6502) zeroPageIndexedAddress(index Index) (result uint16) {
	//value := cpu.Memory.Fetch(cpu.Registers.PC)
	value := cpu.readNext()
	cpu.readAddress(uint16(value))

	effectiveValue := value + cpu.IndexToRegister(index)
	result = uint16(effectiveValue)

	if cpu.decode.enabled {
		cpu.decode.args = fmt.Sprintf("%02X", value)
		cpu.decode.decodedArgs = fmt.Sprintf("$%02X,%s @ %02X",
			value, index.String(), result)
	}

	return
}

// Relative Addressing Mode
func (cpu *M6502) relativeAddress() (result uint16) {
	//value := uint16(cpu.Memory.Fetch(cpu.Registers.PC))
	//cpu.Registers.PC++
	value := uint16(cpu.readNext())

	var offset uint16
	if value > 0x7f {
		offset = -(0x0100 - value)
	} else {
		offset = value
	}
	//offset := uint16(cpu.fromTwosCom(value))

	result = cpu.Registers.PC + offset

	if cpu.decode.enabled {
		cpu.decode.args = fmt.Sprintf("%02X", value)
		cpu.decode.decodedArgs = fmt.Sprintf("$%04X", result)
	}

	return
}

// Absolute Addressing Mode
func (cpu *M6502) absoluteAddress() (result uint16) {
	//low := cpu.Memory.Fetch(cpu.Registers.PC)
	//high := cpu.Memory.Fetch(cpu.Registers.PC + 1)
	//cpu.Registers.PC += 2
	low := cpu.readNext()
	high := cpu.readNext()

	result = (uint16(high) << 8) | uint16(low)

	if cpu.decode.enabled {
		cpu.decode.args = fmt.Sprintf("%02X %02X", low, high)
		cpu.decode.decodedArgs = fmt.Sprintf("$%04X = ", result)
	}

	return
}

// Indirect Addressing Mode
func (cpu *M6502) indirectAddress() (result uint16) {
	//low := cpu.Memory.Fetch(cpu.Registers.PC)
	//high := cpu.Memory.Fetch(cpu.Registers.PC + 1)
	//cpu.Registers.PC += 2
	low := cpu.readNext()
	high := cpu.readNext()

	if cpu.decode.enabled {
		cpu.decode.args = fmt.Sprintf("%02X %02X", low, high)
	}

	// XXX: The 6502 had a bug in which it incremented only the
	// high byte instead of the whole 16-bit address when
	// computing the address.
	//
	// See http://www.obelisk.demon.co.uk/6502/reference.html#JMP
	// and http://www.6502.org/tutorials/6502opcodes.html#JMP for
	// details
	aLow := (uint16(high) << 8) | uint16(low)
	var aHigh uint16
	if low == 0xff {
		// Page boundary wrap bug: MSB comes from start of same page
		aHigh = aLow & 0xff00
	} else {
		aHigh = aLow + 1
	}
	low = cpu.readAddress(aLow)
	high = cpu.readAddress(aHigh)

	/*
		aHigh := (uint16(high) << 8) | uint16(low+1)
		aLow := (uint16(high) << 8) | uint16(low)

		//low = cpu.Memory.Fetch(aLow)
		//high = cpu.Memory.Fetch(aHigh)
		low = cpu.readAddress(aLow)
		high = cpu.readAddress(aHigh)
	*/

	result = (uint16(high) << 8) | uint16(low)
	badResult := (uint16(cpu.Memory.Fetch(aLow+1)) << 8) | uint16(low)

	if cpu.decode.enabled {
		cpu.decode.decodedArgs = fmt.Sprintf("($%04X) = %04X", aLow, badResult)
	}
	return
}

// Absolute Indexed Addressing Mode
func (cpu *M6502) absoluteIndexedAddress(index Index, _ *InstructionStatus) (result uint16) {
	low := cpu.readNext()
	high := cpu.readNext()
	indexValue := uint16(cpu.IndexToRegister(index))

	// Calculate effective address with potentially invalid high byte
	lowWithIndex := uint16(low) + indexValue
	address := (uint16(high) << 8) | uint16(low)

	// Effective address with potentially invalid high byte (may be $100 smaller if page crossed)
	effectiveAddrWithPotentiallyInvalidHigh := (uint16(high) << 8) | (lowWithIndex & 0xFF)

	// Cycle 4: read from effective address (potentially invalid), fix the high byte
	cpu.Memory.Fetch(effectiveAddrWithPotentiallyInvalidHigh)
	//cpu.readAddress(effectiveAddrWithPotentiallyInvalidHigh)

	// Calculate the correct effective address
	result = address + indexValue

	var ops []OpCode
	if index == X {
		ops = []OpCode{
			0x1E, 0x1F, 0x3E, 0x3F, 0x5E, 0x5F, 0x7E,
			0x7F, 0x9D, 0x9F, 0xDE, 0xDF, 0xFE, 0xFF,
		}
	} else {
		ops = []OpCode{
			0x1B, 0x3B, 0x5B, 0x7B, 0x99, 0xDB, 0xFB,
		}
	}

	// Cycle 5 (optional): if page boundary was crossed, re-read from corrected effective address
	if !SamePage(address, result) {
		cpu.IncrementCycle()
	} else if cpu.CurrentOp != nil && cpu.containsOpCode(ops, cpu.CurrentOp.OpCode) {
		// Some opcodes always perform the read on cycle 5 even if no page crossing occurred
		cpu.IncrementCycle()
		//cpu.Memory.Fetch(result)
	}

	if cpu.decode.enabled {
		cpu.decode.args = fmt.Sprintf("%02X %02X", low, high)
		cpu.decode.decodedArgs = fmt.Sprintf("$%04X,%s @ %04X = ", address, index.String(), result)
	}

	return
}

// X-Indexed Indirect Addressing Mode
func (cpu *M6502) indexedIndirectAddress() (result uint16) {
	value := cpu.readNext()
	cpu.readAddress(uint16(value))
	address := uint16(value + cpu.Registers.X)

	low := cpu.readAddress(address)
	high := cpu.readAddress((address + 1) & 0x00ff)

	result = (uint16(high) << 8) | uint16(low)

	if cpu.decode.enabled {
		cpu.decode.args = fmt.Sprintf("%02X", value)
		cpu.decode.decodedArgs = fmt.Sprintf("($%02X,X) @ %02X = %04X = ", value, address, result)
	}

	return
}

// Indirect Y-Indexed Addressing Mode
func (cpu *M6502) indirectIndexedAddress(_ *InstructionStatus) (result uint16) {
	value := cpu.readNext()
	address := uint16(value)

	low := cpu.readAddress(address)
	high := cpu.readAddress((address + 1) & 0x00ff)

	address = (uint16(high) << 8) | uint16(low)

	result = address + uint16(cpu.Registers.Y)

	ops := []OpCode{
		0x13, 0x33, 0x53, 0x73, 0x91, 0xD3, 0xF3,
	}

	if !SamePage(address, result) {
		cpu.IncrementCycle()
	} else if cpu.CurrentOp != nil && cpu.containsOpCode(ops, cpu.CurrentOp.OpCode) {
		// Some opcodes always perform the read on cycle 5 even if no page crossing occurred
		cpu.IncrementCycle()
		//cpu.Memory.Fetch(result)
	}

	if cpu.decode.enabled {
		cpu.decode.args = fmt.Sprintf("%02X", value)
		cpu.decode.decodedArgs = fmt.Sprintf("($%02X),Y = %04X @ %04X = ", value, address, result)
	}

	return
}

// endregion Addressing Modes

// region Helpers

// Converts unsigned int value to signed int value using two's complement
func (cpu *M6502) fromTwosCom(value uint8) int8 { return int8(value) }

// Converts signed int value to unsigned int value using two's complement
func (cpu *M6502) toTwosCom(value int8) uint8 { return uint8(value) }

// Converts BCD (Binary Coded Decimal) to uint8
func (cpu *M6502) FromBCD(value uint8) uint8 {
	return ((value >> 4) * 10) + (value & 0x0F)
}

// Converts int to BCD (Binary Coded Decimal)
func (cpu *M6502) ToBCD(value uint8) uint8 {
	return ((value / 10) << 4) | (value % 10)
}

func (cpu *M6502) GetHighLow(value uint16) (high, low uint8) {
	high = uint8(value >> 8)
	low = uint8(value)
	return
}

// region Stack push/pull

func (cpu *M6502) push(value uint8) {
	//cpu.Memory.Store(0x0100|uint16(cpu.Registers.SP), value)
	cpu.store(0x0100|uint16(cpu.Registers.SP), value)
	cpu.Registers.SP--
}

func (cpu *M6502) push16(value uint16) {
	cpu.push(uint8(value >> 8))
	cpu.push(uint8(value))
}

func (cpu *M6502) pull() (value uint8) {
	cpu.Registers.SP++
	//value = cpu.Memory.Fetch(0x0100 | uint16(cpu.Registers.SP))
	value = cpu.readAddress(0x0100 | uint16(cpu.Registers.SP))
	return
}

func (cpu *M6502) pull16() (value uint16) {
	low := cpu.pull()
	high := cpu.pull()

	value = (uint16(high) << 8) | uint16(low)
	return
}

// endregion Stack Operations

// region Transfer
func (cpu *M6502) transfer(from uint8, to *uint8) {
	*to = cpu.setZNFlags(from)
}

// endregion Transfer

// region Arithmetic
func (cpu *M6502) addition(input uint8) {
	value := uint16(input)
	orig := uint16(cpu.Registers.A)

	if !cpu.decimalMode || cpu.Registers.P&D == 0 {
		result := cpu.setCFlagAddition(orig + value + uint16(cpu.Registers.P&C))
		cpu.Registers.A = cpu.setZNFlags(uint8(cpu.setVFlagAddition(orig, value, result)))
	} else {
		low := uint16(orig&0x000f) + uint16(value&0x000f) + uint16(cpu.Registers.P&C)
		high := uint16(orig&0x00f0) + uint16(value&0x00f0)

		if low >= 0x000a {
			low -= 0x000a
			high += 0x0010
		}

		// Check for BCD carry before adjusting
		carry := uint16(0)
		if high >= 0x00a0 {
			high -= 0x00a0
			carry = 0x0100 // Set carry flag
		}

		result := high | (low & 0x000f) | carry
		cpu.setCFlagAddition(result)
		cpu.Registers.A = cpu.setZNFlags(uint8(cpu.setVFlagAddition(orig, value, result)))
	}
}

// endregion Arithmetic

// region Conditional Branch
func (cpu *M6502) branch(address uint16, condition func() bool) {
	v := condition()
	if !v {
		return
	}
	cpu.IncrementCycle()
	cpu.CurrentOp.status |= Branched

	if !SamePage(cpu.Registers.PC, address) {
		cpu.IncrementCycle()
		cpu.CurrentOp.status |= PageCross
	}

	cpu.Registers.PC = address
}

// endregion Conditional Branch

// region Comparisons
func (cpu *M6502) compare(value uint16, register uint8) {
	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	value = value ^ 0xff + 1
	cpu.setZNFlags(uint8(cpu.setCFlagAddition(uint16(register) + value)))
}

// endregion Comparisons

// region Decrements & Increments
func (cpu *M6502) increment(register *uint8) {
	*register = cpu.setZNFlags(*register + 1)
}

func (cpu *M6502) decrement(register *uint8) {
	*register = cpu.setZNFlags(*register - 1)
}

// endregion Decrements & Increments

// region Shift & Rotate
type direction int

const (
	left direction = iota
	right
)

func (cpu *M6502) shift(direction direction, value uint8, store func(uint8)) {
	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	c := Status(0)

	switch direction {
	case left:
		c = Status((value & uint8(N)) >> 7)
		value <<= 1
	case right:
		c = Status(value & uint8(C))
		value >>= 1
	}

	cpu.Registers.P &= ^C
	cpu.Registers.P |= c

	store(cpu.setZNFlags(value))
}

func (cpu *M6502) rotate(direction direction, value uint8, store func(uint8)) {
	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	c := Status(0)

	switch direction {
	case left:
		c = Status(value & uint8(N) >> 7)
		value = ((value << 1) & uint8(^C)) | uint8(cpu.Registers.P&C)
	case right:
		c = Status(value & uint8(C))
		value = ((value >> 1) & uint8(^N)) | uint8((cpu.Registers.P&C)<<7)
	}

	cpu.Registers.P &= ^C
	cpu.Registers.P |= c

	store(cpu.setZNFlags(value))
}

// endregion Shift & Rotate

// endregion Helpers

// region Operations

// region Transfer Instructions

// Loads a byte of memory into the accumulator setting the zero and
// negative flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of A is set
func (cpu *M6502) Lda(value uint8) {
	ret := cpu.setZNFlags(value)
	cpu.Registers.A = ret
	//cpu.load(address, &cpu.Registers.A)
}

// Unofficial
//
// Loads a byte of memory into the accumulator and X setting the zero
// and negative flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of A is set
func (cpu *M6502) Lax(value uint8) {
	cpu.Registers.A = value
	cpu.Registers.X = value
	cpu.setZNFlags(value)
}

// Loads a byte of memory into the X register setting the zero and
// negative flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if X = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of X is set
func (cpu *M6502) Ldx(value uint8) {
	ret := cpu.setZNFlags(value)
	cpu.Registers.X = ret
	//cpu.load(value, &cpu.Registers.X)
}

// Loads a byte of memory into the Y register setting the zero and
// negative flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if Y = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of Y is set
func (cpu *M6502) Ldy(value uint8) {
	ret := cpu.setZNFlags(value)
	cpu.Registers.Y = ret
	//cpu.load(address, &cpu.Registers.Y)
}

// Unofficial
//
// A and X are put on the bus at the same time (resulting effectively
// in an AND operation) and stored in M
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Sax(address uint16) {
	cpu.store(address, cpu.Registers.A&cpu.Registers.X)
}

// Stores the contents of the accumulator into memory.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Sta(address uint16) {
	cpu.store(address, cpu.Registers.A)
}

// Stores the contents of the X register into memory.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Stx(address uint16) {
	cpu.store(address, cpu.Registers.X)
}

// Stores the contents of the Y register into memory.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Sty(address uint16) {
	cpu.store(address, cpu.Registers.Y)
}

// Copies the current contents of the accumulator into the X register
// and sets the zero and negative flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if X = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of X is set
func (cpu *M6502) Tax(_ uint16) {
	cpu.transfer(cpu.Registers.A, &cpu.Registers.X)
}

// Copies the current contents of the accumulator into the Y register
// and sets the zero and negative flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if Y = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of Y is set
func (cpu *M6502) Tay(_ uint16) {
	cpu.transfer(cpu.Registers.A, &cpu.Registers.Y)
}

// Copies the current contents of the X register into the accumulator
// and sets the zero and negative flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of A is set
func (cpu *M6502) Txa(_ uint16) {
	cpu.transfer(cpu.Registers.X, &cpu.Registers.A)
}

// Copies the current contents of the Y register into the accumulator
// and sets the zero and negative flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of A is set
func (cpu *M6502) Tya(_ uint16) {
	cpu.transfer(cpu.Registers.Y, &cpu.Registers.A)
}

// Copies the current contents of the stack register into the X
// register and sets the zero and negative flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if X = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of X is set
func (cpu *M6502) Tsx(_ uint16) {
	cpu.transfer(cpu.Registers.SP, &cpu.Registers.X)
}

// Copies the current contents of the X register into the stack
// register.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Txs(_ uint16) {
	cpu.Registers.SP = cpu.Registers.X
}

// endregion Transfer Instructions

// region Stack Instructions

// Pushes a copy of the accumulator on to the stack.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Pha(_ uint16) {
	//cpu.readNext()
	cpu.push(cpu.Registers.A)
}

// Pushes a copy of the status flags on to the stack.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Php(_ uint16) {
	//cpu.readNext()
	cpu.push(uint8(cpu.Registers.P | B | U))
}

// Pulls an 8 bit value from the stack and into the accumulator. The
// zero and negative flags are set as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of A is set
func (cpu *M6502) Pla(_ uint16) {
	//cpu.readNext()
	cpu.IncrementCycle()
	cpu.Registers.A = cpu.setZNFlags(cpu.pull())
}

// Pulls an 8 bit value from the stack and into the processor
// flags. The flags will take on new states as determined by the value
// pulled.
//
//	C 	Carry Flag 	  Set from stack
//	Z 	Zero Flag 	  Set from stack
//	I 	Interrupt Disable Set from stack
//	D 	Decimal Mode Flag Set from stack
//	B 	Break Command 	  Set from stack
//	V 	Overflow Flag 	  Set from stack
//	N 	Negative Flag 	  Set from stack
func (cpu *M6502) Plp(_ uint16) {
	//cpu.readNext()
	cpu.IncrementCycle()
	cpu.Registers.P = (Status(cpu.pull()) & ^B) | U
}

// endregion Stack Instructions

// region Logical Operations

// A logical AND is performed, bit by bit, on the accumulator contents
// using the contents of a byte of memory.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 set
func (cpu *M6502) And(value uint8) {
	//value := cpu.Memory.Fetch(address)
	//value := cpu.readAddress(address)

	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	cpu.Registers.A = cpu.setZNFlags(cpu.Registers.A & value)
}

// An exclusive OR is performed, bit by bit, on the accumulator
// contents using the contents of a byte of memory.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 set
func (cpu *M6502) Eor(value uint8) {
	//value := cpu.Memory.Fetch(address)
	//value := cpu.readAddress(address)

	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	cpu.Registers.A = cpu.setZNFlags(cpu.Registers.A ^ value)
}

// An inclusive OR is performed, bit by bit, on the accumulator
// contents using the contents of a byte of memory.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 set
func (cpu *M6502) Ora(value uint8) {
	//value := cpu.Memory.Fetch(address)
	//value := cpu.readAddress(address)

	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	cpu.Registers.A = cpu.setZNFlags(cpu.Registers.A | value)
}

// endregion Logical Operations

// region Bit Test

// This instructions is used to test if one or more bits are set in a
// target memory location. The mask pattern in A is ANDed with the
// value in memory to set or clear the zero flag, but the result is
// not kept. Bits 7 and 6 of the value from memory are copied into the
// N and V flags.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if the result if the AND is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Set to bit 6 of the memory value
//	N 	Negative Flag 	  Set to bit 7 of the memory value
func (cpu *M6502) Bit(value uint8) {
	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	cpu.setZFlag(value & cpu.Registers.A)
	cpu.Registers.P = (cpu.Registers.P & ^(N | V)) | Status(value&uint8(N|V))
}

// endregion Bit Test

// region Arithmetic Operations

// This instruction adds the contents of a memory location to the
// accumulator together with the carry bit. If overflow occurs the
// carry bit is set, this enables multiple byte addition to be
// performed.
//
//	C 	Carry Flag 	  Set if overflow in bit 7
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Set if sign bit is incorrect
//	N 	Negative Flag 	  Set if bit 7 set
func (cpu *M6502) Adc(value uint8) {
	//value := uint16(cpu.Memory.Fetch(address))
	//value := uint16(cpu.readAddress(address))

	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	cpu.addition(value)
}

// This instruction subtracts the contents of a memory location to the
// accumulator together with the not of the carry bit. If overflow
// occurs the carry bit is clear, this enables multiple byte
// subtraction to be performed.
//
//	C 	Carry Flag 	  Clear if overflow in bit 7
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Set if sign bit is incorrect
//	N 	Negative Flag 	  Set if bit 7 set
func (cpu *M6502) Sbc(value uint8) {
	//value := uint16(cpu.Memory.Fetch(address))
	//value := uint16(cpu.readAddress(address))

	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	if !cpu.decimalMode || cpu.Registers.P&D == 0 {
		value ^= 0xff
	} else {
		value = 0x99 - value
	}

	cpu.addition(value)
}

// endregion Arithmetic Operations

// region Comparisons

// This instruction compares the contents of the accumulator with
// another memory held value and sets the zero and carry flags as
// appropriate.
//
//	C 	Carry Flag 	  Set if A >= M
//	Z 	Zero Flag 	  Set if A = M
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Cmp(value uint8) {
	//value := uint16(cpu.Memory.Fetch(address))
	//value := uint16(cpu.readAddress(address))
	cpu.compare(uint16(value), cpu.Registers.A)
}

// This instruction compares the contents of the X register with
// another memory held value and sets the zero and carry flags as
// appropriate.
//
//	C 	Carry Flag 	  Set if X >= M
//	Z 	Zero Flag 	  Set if X = M
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Cpx(value uint8) {
	//value := uint16(cpu.Memory.Fetch(address))
	//value := uint16(cpu.readAddress(address))
	cpu.compare(uint16(value), cpu.Registers.X)
}

// This instruction compares the contents of the Y register with
// another memory held value and sets the zero and carry flags as
// appropriate.
//
//	C 	Carry Flag 	  Set if Y >= M
//	Z 	Zero Flag 	  Set if Y = M
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Cpy(value uint8) {
	//value := uint16(cpu.Memory.Fetch(address))
	//value := uint16(cpu.readAddress(address))
	cpu.compare(uint16(value), cpu.Registers.Y)
}

// endregion Comparisons

// region Decrements & Increments

// Adds one to the value held at a specified memory location setting
// the zero and negative flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Inc(address uint16) {
	//value := cpu.Memory.Fetch(address)
	value := cpu.readAddress(address)
	cpu.store(address, value)

	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	//cpu.Memory.Store(address, cpu.setZNFlags(value+1))
	cpu.store(address, cpu.setZNFlags(value+1))
}

// Adds one to the X register setting the zero and negative flags as
// appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if X is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of X is set
func (cpu *M6502) Inx(_ uint16) {
	cpu.increment(&cpu.Registers.X)
}

// Adds one to the Y register setting the zero and negative flags as
// appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if Y is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of Y is set
func (cpu *M6502) Iny(_ uint16) {
	cpu.increment(&cpu.Registers.Y)
}

// Subtracts one from the value held at a specified memory location
// setting the zero and negative flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Dec(address uint16) {
	//value := cpu.Memory.Fetch(address)
	value := cpu.readAddress(address)
	cpu.store(address, value)

	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	//cpu.Memory.Store(address, cpu.setZNFlags(value-1))
	cpu.store(address, cpu.setZNFlags(value-1))
}

// Subtracts one from the X register setting the zero and negative
// flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if X is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of X is set
func (cpu *M6502) Dex(_ uint16) {
	cpu.decrement(&cpu.Registers.X)
}

// Subtracts one from the Y register setting the zero and negative
// flags as appropriate.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if Y is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of Y is set
func (cpu *M6502) Dey(_ uint16) {
	cpu.decrement(&cpu.Registers.Y)
}

// endregion Decrements & Increments

// region Shift & Rotate Instructions

// This operation shifts all the bits of the accumulator one bit
// left. Bit 0 is set to 0 and bit 7 is placed in the carry flag. The
// effect of this operation is to multiply the memory contents by 2
// (ignoring 2's complement considerations), setting the carry if the
// result will not fit in 8 bits.
//
//	C 	Carry Flag 	  Set to contents of old bit 7
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) AslA() {
	cpu.shift(left, cpu.Registers.A, func(value uint8) { cpu.Registers.A = value })

	if cpu.decode.enabled {
		cpu.decode.decodedArgs = "A"
	}
}

// This operation shifts all the bits of the memory contents one bit
// left. Bit 0 is set to 0 and bit 7 is placed in the carry flag. The
// effect of this operation is to multiply the memory contents by 2
// (ignoring 2's complement considerations), setting the carry if the
// result will not fit in 8 bits.
//
//	C 	Carry Flag 	  Set to contents of old bit 7
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Asl(address uint16) {
	if cpu.CurrentOp.AddressingMode == Accumulator {
		cpu.AslA()
		return
	}
	value := cpu.readAddress(address)
	cpu.store(address, value)
	cpu.shift(left, value, func(value uint8) { cpu.store(address, value) })
}

// Each of the bits in A is shift one place to the right. The bit that
// was in bit 0 is shifted into the carry flag. Bit 7 is set to zero.
//
//	C 	Carry Flag 	  Set to contents of old bit 0
//	Z 	Zero Flag 	  Set if result = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) LsrA() {
	cpu.shift(right, cpu.Registers.A, func(value uint8) { cpu.Registers.A = value })

	if cpu.decode.enabled {
		cpu.decode.decodedArgs = "A"
	}
}

// Each of the bits in M is shift one place to the right. The bit that
// was in bit 0 is shifted into the carry flag. Bit 7 is set to zero.
//
//	C 	Carry Flag 	  Set to contents of old bit 0
//	Z 	Zero Flag 	  Set if result = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Lsr(address uint16) {
	if cpu.CurrentOp.AddressingMode == Accumulator {
		cpu.LsrA()
		return
	}
	value := cpu.readAddress(address)
	cpu.store(address, value)
	cpu.shift(right, value, func(value uint8) { cpu.store(address, value) })
}

// Move each of the bits in A one place to the left. Bit 0 is filled
// with the current value of the carry flag whilst the old bit 7
// becomes the new carry flag value.
//
//	C 	Carry Flag 	  Set to contents of old bit 7
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) RolA() {
	cpu.rotate(left, cpu.Registers.A, func(value uint8) { cpu.Registers.A = value })

	if cpu.decode.enabled {
		cpu.decode.decodedArgs = "A"
	}
}

// Move each of the bits in A one place to the left. Bit 0 is filled
// with the current value of the carry flag whilst the old bit 7
// becomes the new carry flag value.
//
//	C 	Carry Flag 	  Set to contents of old bit 7
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Rol(address uint16) {
	if cpu.CurrentOp.AddressingMode == Accumulator {
		cpu.RolA()
		return
	}
	value := cpu.readAddress(address)
	cpu.store(address, value)
	cpu.rotate(left, value, func(value uint8) { cpu.store(address, value) })
}

// Move each of the bits in A one place to the right. Bit 7 is filled
// with the current value of the carry flag whilst the old bit 0
// becomes the new carry flag value.
//
//	C 	Carry Flag 	  Set to contents of old bit 0
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) RorA() {
	cpu.rotate(right, cpu.Registers.A, func(value uint8) { cpu.Registers.A = value })

	if cpu.decode.enabled {
		cpu.decode.decodedArgs = "A"
	}
}

// Move each of the bits in M one place to the right. Bit 7 is filled
// with the current value of the carry flag whilst the old bit 0
// becomes the new carry flag value.
//
//	C 	Carry Flag 	  Set to contents of old bit 0
//	Z 	Zero Flag 	  Set if A = 0
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Ror(address uint16) {
	if cpu.CurrentOp.AddressingMode == Accumulator {
		cpu.RorA()
		return
	}
	value := cpu.readAddress(address)
	cpu.store(address, value)
	cpu.rotate(right, value, func(value uint8) { cpu.store(address, value) })
}

// endregion Shift & Rotate Instructions

// region Jumps & Subroutines

// Sets the program counter to the address specified by the operand.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Jmp(address uint16) {
	if cpu.decode.enabled {
		if strings.HasPrefix(cpu.decode.decodedArgs, "$") {
			// delete ' = '
			cpu.decode.decodedArgs = cpu.decode.decodedArgs[:len(cpu.decode.decodedArgs)-3]
		}
	}

	cpu.Registers.PC = address
}

// The JSR instruction pushes the address (minus one) of the return
// point on to the stack and then sets the program counter to the
// target memory address.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Jsr(address uint16) {
	if cpu.decode.enabled {
		cpu.decode.decodedArgs = fmt.Sprintf("$%04X", address)
	}

	value := cpu.Registers.PC - 1

	cpu.push16(value)

	cpu.IncrementCycle()
	cpu.Registers.PC = address
}

// The RTS instruction is used at the end of a subroutine to return to
// the calling routine. It pulls the program counter (minus one) from
// the stack.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Rts(_ uint16) {
	cpu.IncrementCycle()
	cpu.Registers.PC = cpu.pull16() + 1
	cpu.IncrementCycle()
}

// endregion Jumps & Subroutines

// region Conditional Branch Instructions

// If the carry flag is clear then add the relative displacement to
// the program counter to cause a branch to a new location.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Bcc(address uint16) {
	cpu.branch(address, func() bool { return cpu.Registers.P&C == 0 })
}

// If the carry flag is set then add the relative displacement to the
// program counter to cause a branch to a new location.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Bcs(address uint16) {
	cpu.branch(address, func() bool { return cpu.Registers.P&C != 0 })
}

// If the zero flag is set then add the relative displacement to the
// program counter to cause a branch to a new location.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Beq(address uint16) {
	cpu.branch(address, func() bool { return cpu.Registers.P&Z != 0 })
}

// If the negative flag is set then add the relative displacement to
// the program counter to cause a branch to a new location.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Bmi(address uint16) {
	cpu.branch(address, func() bool { return cpu.Registers.P&N != 0 })
}

// If the zero flag is clear then add the relative displacement to the
// program counter to cause a branch to a new location.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Bne(address uint16) {
	cpu.branch(address, func() bool { return cpu.Registers.P&Z == 0 })
}

// If the negative flag is clear then add the relative displacement to
// the program counter to cause a branch to a new location.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Bpl(address uint16) {
	cpu.branch(address, func() bool { return cpu.Registers.P&N == 0 })
}

// If the overflow flag is clear then add the relative displacement to
// the program counter to cause a branch to a new location.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Bvc(address uint16) {
	cpu.branch(address, func() bool { return cpu.Registers.P&V == 0 })
}

// If the overflow flag is set then add the relative displacement to
// the program counter to cause a branch to a new location.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Bvs(address uint16) {
	cpu.branch(address, func() bool { return cpu.Registers.P&V != 0 })
}

// endregion Conditional Branch Instructions

// region Flag Instructions

// Set the carry flag to zero.
//
//	C 	Carry Flag 	  Set to 0
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Clc(_ uint8) {
	cpu.Registers.P &^= C
}

// Set the decimal mode flag to zero.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Set to 0
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Cld(_ uint8) {
	cpu.Registers.P &^= D
}

// Clears the interrupt disable flag allowing normal interrupt
// requests to be serviced.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Set to 0
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Cli(_ uint8) {
	cpu.Registers.P &^= I
}

// Clears the interrupt disable flag allowing normal interrupt
// requests to be serviced.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Set to 0
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Clv(_ uint8) {
	cpu.Registers.P &^= V
}

// Set the carry flag to one.
//
//	C 	Carry Flag 	  Set to 1
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Sec(_ uint8) {
	cpu.Registers.P |= C
}

// Set the decimal mode flag to one.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Set to 1
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Sed(_ uint8) {
	cpu.Registers.P |= D
}

// Set the interrupt disable flag to one.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Set to 1
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Sei(_ uint8) {
	cpu.Registers.P |= I
}

// endregion Flag Instructions

// region Interrupts

// The BRK instruction forces the generation of an interrupt
// request. The program counter and processor status are pushed on the
// stack then the IRQ interrupt vector at $FFFE/F is loaded into the
// PC and the break flag in the status set to one.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Set to 1
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Brk(_ uint8) {
	cpu.Registers.PC++

	cpu.push16(cpu.Registers.PC)
	cpu.push(uint8(cpu.Registers.P | B | U))

	cpu.Registers.P |= I

	low := cpu.readAddress(0xfffe)
	high := cpu.readAddress(0xffff)

	cpu.Registers.PC = (uint16(high) << 8) | uint16(low)
}

// The RTI instruction is used at the end of an interrupt processing
// routine. It pulls the processor flags from the stack followed by
// the program counter.
//
//	C 	Carry Flag 	  Set from stack
//	Z 	Zero Flag 	  Set from stack
//	I 	Interrupt Disable Set from stack
//	D 	Decimal Mode Flag Set from stack
//	B 	Break Command 	  Set from stack
//	V 	Overflow Flag 	  Set from stack
//	N 	Negative Flag 	  Set from stack
func (cpu *M6502) Rti(_ uint16) {
	cpu.readNext()
	//cpu.IncrementCycle()
	cpu.Registers.P = (Status(cpu.pull()) & ^B) | U
	cpu.Registers.PC = cpu.pull16()
}

// endregion Interrupts

// region Other Instructions

// The NOP instruction causes no changes to the processor other than
// the normal incrementing of the program counter to the next
// instruction.
//
// Sadly not all NOPs are equal, Ive added a few here
// based on `CPU_unofficial_opcodes` and will add more
// based on game compatibility, and ultimately I'd like
// to cover all illegal opcodes too
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Nop(value uint8) {
	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}
}

// Unofficial
//
// Skip Byte - (NOP)
//
// See [cpu.Nop]
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Skb(value uint8) {
	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}
}

// Unofficial
//
// Ignore Byte - (NOP)
//
// See [cpu.Nop]
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Ign(address uint16) {
	v := cpu.readAddress(address)

	if cpu.decode.enabled {
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", v)
	}
}

// endregion Other Instructions

// region Unofficial Instructions

// Unofficial
//
// Description
//
//	C 	Carry Flag 	  Set if A >= result
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Alr(value uint8) {
	cpu.And(value)
	cpu.LsrA()
}

// Unofficial
//
// AND oper + set C as [cpu.Asl]
//
//	C 	Carry Flag 	  Set if A >= result
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Anc(value uint8) {
	cpu.And(value)
	cpu.Registers.P = (cpu.Registers.P & ^C) | cpu.Registers.P>>7
}

// Unofficial
//
// AND oper + set C as [cpu.Rol]
// effectively the same as instr. 0B ([cpu.Anc])
//
//	C 	Carry Flag 	  Set if A >= result
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Aac(value uint8) {
	cpu.And(value)
	cpu.Registers.P = (cpu.Registers.P & ^C) | cpu.Registers.P>>7
}

// Unofficial
//
// * OR X + AND oper
//
// # Highly unstable, do not use.
//
// A base value in A is determined based on the contets of A and
// a constant, which may be typically $00, $ff, $ee, etc.
// The value of this constant depends on temerature, the chip series,
// and maybe other factors, as well.
// In order to eliminate these uncertaincies from the equation,
// use either 0 as the operand or a value of $FF in the accumulator.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Ane(value uint8) {
	cpu.Registers.A = (cpu.Registers.A | cpu.magic) & cpu.Registers.X & value
	cpu.setZNFlags(cpu.Registers.A)
}

// Unofficial
//
// AND oper + ROR
//
// This operation involves the adder:
// V-flag is set according to (A AND oper) + oper
// The carry is not set, but bit 7 (sign) is exchanged with the carry
//
//	C 	Carry Flag 	  Set to bit 6 of result
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Set to bit 6 XOR bit 5 of result
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Arr(value uint8) {
	// AND with accumulator
	cpu.Registers.A &= value

	// ROR with current carry
	carry := uint8(0)
	if cpu.getFlag(C) {
		carry = 0x80
	}
	cpu.Registers.A = (cpu.Registers.A >> 1) | carry

	// ARR has special flag behavior:
	// C flag = bit 6 of result
	cpu.setFlag(C, cpu.Registers.A&0x40 != 0)

	// V flag = bit 6 XOR bit 5 of result
	bit6 := (cpu.Registers.A >> 6) & 1
	bit5 := (cpu.Registers.A >> 5) & 1
	cpu.setFlag(V, bit6^bit5 != 0)

	// Z and N flags set normally
	cpu.setZNFlags(cpu.Registers.A)
}

// Unofficial
//
// Decrements the operand and then compares the result to the
// accumulator.
//
//	C 	Carry Flag 	  Set if A >= result
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Dcp(address uint16) {
	if cpu.decode.enabled {
		value := cpu.Memory.Fetch(address)
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	enabled := cpu.decode.enabled
	cpu.decode.enabled = false
	cpu.Dec(address)
	//value := cpu.readAddress(address)
	value := cpu.Memory.Fetch(address)
	cpu.Cmp(value)
	cpu.decode.enabled = enabled
}

// Unofficial
//
// INC oper + SBC oper
//
//	C 	Carry Flag 	  Clear if overflow in bit 7
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Set if sign bit is incorrect
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Isb(address uint16) {
	if cpu.decode.enabled {
		value := cpu.Memory.Fetch(address)
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	enabled := cpu.decode.enabled
	cpu.decode.enabled = false
	cpu.Inc(address)
	//value := cpu.readAddress(address)
	value := cpu.Memory.Fetch(address)
	cpu.Sbc(value)
	cpu.decode.enabled = enabled
}

// Unofficial
//
// LDA/TSX oper
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Las(value uint8) {
	res := cpu.setZNFlags(cpu.Registers.SP & value)

	cpu.Registers.A = res
	cpu.Registers.X = res
	cpu.Registers.SP = res
}

// Unofficial
//
// Store * AND oper in A and X
//
// Highly unstable, involves a 'magic' constant, see [cpu.Ane]
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Lxa(value uint8) {
	res := cpu.setZNFlags((cpu.Registers.A ^ cpu.magic) & value)
	cpu.Registers.A = res
	cpu.Registers.X = res
}

// Unofficial
//
// Description
//
//	C 	Carry Flag 	  Set if A >= result
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Rla(address uint16) {
	if cpu.decode.enabled {
		value := cpu.Memory.Fetch(address)
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	enabled := cpu.decode.enabled
	cpu.decode.enabled = false
	cpu.Rol(address)
	//cpu.And(address)
	//value := cpu.readAddress(address)
	value := cpu.Memory.Fetch(address)
	cpu.And(value)
	cpu.decode.enabled = enabled
}

// Unofficial
//
// ROR oper + ADC oper
//
//	C 	Carry Flag 	  Set if A >= result
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Set if sign bit is incorrect
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Rra(address uint16) {
	if cpu.decode.enabled {
		value := cpu.Memory.Fetch(address)
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	enabled := cpu.decode.enabled
	cpu.decode.enabled = false
	cpu.Ror(address)
	//cpu.Adc(address)
	//value := cpu.readAddress(address)
	value := cpu.Memory.Fetch(address)
	cpu.Adc(value)
	cpu.decode.enabled = enabled
}

// Unofficial
//
// CMP and DEX at once, sets flags like CMP
//
//	C 	Carry Flag 	  Set if A >= result
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Sbx(value uint8) {
	//value := cpu.Memory.Fetch(address)
	//value := cpu.readAddress(address)
	cpu.Registers.X &= cpu.Registers.A
	cpu.compare(uint16(value), cpu.Registers.X)
	cpu.Registers.X -= value
}

// Unofficial
//
// Stores A AND X AND (high-byte of addr. + 1) at addr.
//
// Unstable: sometimes 'AND (H+1)' is dropped, page boundary crossings
// may not work (with the high-byte of the value used as the high-byte
// of the address)
//
// There are a few illegal opcodes which and the high bit of the
// address with registers and write the values back into that
// address. These operations are particularly screwy.
// These posts are used as reference but I am unsure whether they
// are correct.
// - http://forums.nesdev.com/viewtopic.php?f=3&t=3831&start=30#p113343
// - http://forums.nesdev.com/viewtopic.php?f=3&t=10698
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Sha(address uint16) {
	offset := (address - uint16(cpu.Registers.Y))
	low := uint8(offset & 0x00FF)
	high := uint8(offset >> 8)

	cpu.readAddress(offset)
	var result uint16
	if !SamePage(address, offset) {
		result = (uint16(high&cpu.Registers.X) << 8) | uint16(low+cpu.Registers.Y)
	} else {
		result = (uint16(high) << 8) | uint16(low+cpu.Registers.Y)
	}
	value := cpu.Registers.A & cpu.Registers.X & (high + 1)
	cpu.store(result, value)
}

// Unofficial
//
// Stores X AND (high-byte of addr. + 1) at addr.
//
// See [cpu.Sha] for details.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Shx(address uint16) {
	offset := (address - uint16(cpu.Registers.Y))
	low := uint8(offset & 0x00FF)
	high := uint8(offset >> 8)

	cpu.readAddress(offset)
	var result uint16
	if !SamePage(address, offset) {
		result = (uint16(high&cpu.Registers.X) << 8) | uint16(low+cpu.Registers.Y)
	} else {
		result = (uint16(high) << 8) | uint16(low+cpu.Registers.Y)
	}
	value := cpu.Registers.X & (high + 1)
	cpu.store(result, value)
}

// Unofficial
//
// Stores Y AND (high-byte of addr. + 1) at addr.
//
// See [cpu.Sha] for details.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Shy(address uint16) {
	offset := (address - uint16(cpu.Registers.X))
	low := uint8(offset & 0x00FF)
	high := uint8(offset >> 8)

	cpu.readAddress(offset)
	var result uint16
	if !SamePage(address, offset) {
		result = (uint16(high&cpu.Registers.Y) << 8) | uint16(low+cpu.Registers.X)
	} else {
		result = (uint16(high) << 8) | uint16(low+cpu.Registers.X)
	}
	value := cpu.Registers.Y & (high + 1)
	cpu.store(result, value)
}

// Unofficial
//
// ASL oper + ORA oper
//
//	C 	Carry Flag 	  Set if A >= result
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Slo(address uint16) {
	if cpu.decode.enabled {
		value := cpu.Memory.Fetch(address)
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	enabled := cpu.decode.enabled
	cpu.decode.enabled = false
	cpu.Asl(address)
	//value := cpu.readAddress(address)
	value := cpu.Memory.Fetch(address)
	cpu.Ora(value)
	cpu.decode.enabled = enabled
}

// Unofficial
//
// LSR oper + EOR oper
//
//	C 	Carry Flag 	  Set if A >= result
//	Z 	Zero Flag 	  Set if result is zero
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Set if bit 7 of the result is set
func (cpu *M6502) Sre(address uint16) {
	if cpu.decode.enabled {
		value := cpu.Memory.Fetch(address)
		if !strings.HasPrefix(cpu.decode.decodedArgs, "#") &&
			!strings.HasSuffix(cpu.decode.decodedArgs, " = ") {
			cpu.decode.decodedArgs += " = "
		}

		cpu.decode.decodedArgs += fmt.Sprintf("%02X", value)
	}

	enabled := cpu.decode.enabled
	cpu.decode.enabled = false
	cpu.Lsr(address)
	//value := cpu.readAddress(address)
	value := cpu.Memory.Fetch(address)
	cpu.Eor(value)
	cpu.decode.enabled = enabled
}

// Unofficial
//
// Puts A AND X in SP and stores A AND X AND (high-byte of addr. + 1) at addr.
//
// See [cpu.Sha] for details.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Tas(address uint16) {
	cpu.Registers.SP = cpu.Registers.A & cpu.Registers.X
	cpu.Sha(address)
	/*
		# First set the stack pointer's value
		self.r.s = self.r.a & self.r.x

		# Then write to memory using the new value of the stack pointer
		o = (a - self.r.y) & 0xFFFF
		low = o & 0xFF
		high = o >> 8
		self.increment_cycle_count()
		if low + self.r.y > 0xFF:  # crossed page
			a = ((high & self.r.s) << 8) + low + self.r.y
		else:
			a = (high << 8) + low + self.r.y

		v = self.r.s & (high + 1)
		self.writeByte(a, v)
	*/
}

// Unofficial
//
// This instruction freezes the CPU.
//
//	C 	Carry Flag 	  Not affected
//	Z 	Zero Flag 	  Not affected
//	I 	Interrupt Disable Not affected
//	D 	Decimal Mode Flag Not affected
//	B 	Break Command 	  Not affected
//	V 	Overflow Flag 	  Not affected
//	N 	Negative Flag 	  Not affected
func (cpu *M6502) Kil(value uint8) {
	cpu.freezeError = true
}

// endregion Unofficial Instructions

// endregion Operations
