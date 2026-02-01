package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/txc/nesgo/pkg/m65go2"
	"github.com/txc/nesgo/pkg/nes"
)

type positionalArg struct {
	name        string
	funcName    func(string) error
	description string
}

var (
	progname    = filepath.Base(os.Args[0])
	allowedArgs = []positionalArg{
		{
			name:        "tags",
			funcName:    func(f string) error { return checkTags(f) },
			description: "Check tags in ROM",
		},
		{
			name:        "rom",
			funcName:    func(f string) error { return checkRom(f) },
			description: "Check ROM",
		},
		{
			name:        "banks",
			funcName:    func(f string) error { return bankingTest(f) },
			description: "Test banking",
		},
		{
			name:        "dump",
			funcName:    func(f string) error { return dumpRom(f) },
			description: "Disassemble and dump ROM",
		},
		{
			name:        "wram",
			funcName:    func(f string) error { return testWram(f) },
			description: "Test WRAM functionality",
		},
	}
)

func main() {
	var romFile string
	//optOne := flag.Int("one", 1, "option one controls Thing one")
	flag.StringVar(&romFile, "file", "", "path to ROM file")
	flag.StringVar(&romFile, "f", "", "path to ROM file (shorthand)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Do something with Thing one and Thing two
Usage of %s:
  %s [flags] [argument]

Argument must be one of the following:
%s

Flags:
`, progname, progname, strings.Join(getAllowedArgs(), "\n"))
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	argsLen := len(args)
	if argsLen == 0 {
		flag.Usage()
		os.Exit(1)
	} else if argsLen > 1 {
		flag.Usage()
		log.Fatal("only one positional argument is allowed")
	}
	if !isAllowedArg(args[0]) {
		flag.Usage()
		log.Fatal("invalid argument")
	}
	f := getFuncByName(args[0])
	err := f(romFile)
	if err != nil {
		log.Fatal(err)
	}

}

// region internal utils
func getAllowedArgs() []string {
	keys := make([]string, 0, len(allowedArgs))
	for _, k := range allowedArgs {
		val := fmt.Sprintf(" %-8s%s", k.name, k.description)
		keys = append(keys, val)
	}
	return keys
}

func isAllowedArg(arg string) bool {
	for _, k := range allowedArgs {
		if k.name == arg {
			return true
		}
	}
	return false
}

func getFuncByName(name string) func(string) error {
	for _, arg := range allowedArgs {
		if arg.name == name {
			return arg.funcName
		}
	}
	return nil
}

// endregion internal utils

func checkTags(file string) error {
	buf, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	prgData := buf[16 : 16+65536] // Skip header, get 64KB PRG

	fmt.Println("Bank tags at offset $FF8 of each 4KB bank:")
	for bank := range 16 {
		offset := bank*0x1000 + 0xFF8
		tag := prgData[offset]
		fmt.Printf("4KB Bank %2d (offset $%04X): tag = 0x%02X (%d)\n", bank, offset, tag, tag)
	}
	return nil
}

func checkRom(file string) error {
	buf, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	// Skip 16-byte iNES header
	prgData := buf[16:]
	fmt.Printf("First 10 bytes of PRG (bank 0, $8000): ")
	for i := 0; i < 10; i++ {
		fmt.Printf("%02X ", prgData[i])
	}
	fmt.Printf("\nFirst 10 bytes of PRG (bank 1, $8000): ")
	for i := 0x8000; i < 0x8000+10; i++ {
		fmt.Printf("%02X ", prgData[i])
	}
	fmt.Println()
	return nil
}

func bankingTest(file string) error {
	buf, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	romf, err := nes.NewROMFile(buf)
	if err != nil {
		return err
	}

	// Create mapper
	mapper, err := romf.GetMapper()
	if err != nil {
		return err
	}
	cd := mapper.(*nes.ColorDreams)

	fmt.Printf("Initial BankSelect: 0x%02X\n", cd.Registers.BankSelect)
	fmt.Printf("PRG ROM Size: %d (0x%X)\n", cd.PRGRomSize, cd.PRGRomSize)
	fmt.Printf("Number of 32KB banks: %d\n", cd.PRGRomSize/0x8000)

	// Test reading reset vector
	lo := cd.Fetch(0xFFFC)
	hi := cd.Fetch(0xFFFD)
	resetVec := uint16(lo) | (uint16(hi) << 8)
	fmt.Printf("\nReset vector at $FFFC: $%04X\n", resetVec)

	// Read first 10 bytes at reset vector
	if resetVec >= 0x8000 {
		fmt.Printf("First 10 bytes at $%04X: ", resetVec)
		for i := 0; i < 10; i++ {
			b := cd.Fetch(resetVec + uint16(i))
			fmt.Printf("%02X ", b)
		}
		fmt.Println()
	}

	// Test bank switching
	fmt.Println("\n--- Testing bank switching ---")
	for bank := uint8(0); bank < 2; bank++ {
		cd.Store(0x8000, bank) // Switch to bank
		lo := cd.Fetch(0xFFFC)
		hi := cd.Fetch(0xFFFD)
		vec := uint16(lo) | (uint16(hi) << 8)
		fmt.Printf("Bank %d: BankSelect=0x%02X, Reset vector=$%04X\n", bank, cd.Registers.BankSelect, vec)
	}

	// Test specific addresses
	fmt.Println("\n--- Testing address mapping ---")
	cd.Registers.BankSelect = 0 // Bank 0
	addr := uint16(0x8000)
	val := cd.Fetch(addr)
	fmt.Printf("Bank 0, addr $%04X: $%02X\n", addr, val)

	cd.Registers.BankSelect = 1 // Bank 1
	val = cd.Fetch(addr)
	fmt.Printf("Bank 1, addr $%04X: $%02X\n", addr, val)
	return nil
}

func dumpRom(file string) error {
	buf, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	rom, err := nes.NewROMFile(buf)
	if err != nil {
		return err
	}

	romInfo := strings.SplitSeq(rom.String(), "\n")
	for line := range romInfo {
		fmt.Printf("; %s\n", line)
	}

	bank := 1
	asm := m65go2.NewAssembly()
	fmt.Printf("; Disassembling ROM: %s\n"+
		"; Start: 0x%08X\n\n", file, 0x8000+(bank*0x4000))

	for bank := 0; bank < len(rom.ROMBanks); bank++ {
		for i := 0; i < len(rom.ROMBanks[bank]); i++ {
			b := rom.ROMBanks[bank][i]
			if b == 0xff {
				continue
			}
			op := asm.GetInstructions().GetOp(m65go2.OpCode(b))
			db := m65go2.DecodeBytes{
				PC:          uint16(0x8000 + (bank * 0x4000) + i),
				Instruction: op,
				Lo:          0,
				Hi:          0,
			}
			switch op.Bytes() {
			case 2:
				if i+1 < len(rom.ROMBanks[bank]) {
					db.Lo = rom.ROMBanks[bank][i+1]
					i += 1
				}
			case 3:
				if i+1 < len(rom.ROMBanks[bank]) {
					db.Lo = rom.ROMBanks[bank][i+1]
					i += 1
				}
				if i+1 < len(rom.ROMBanks[bank]) {
					db.Hi = rom.ROMBanks[bank][i+1]
					i += 1
				}
			}
			a := asm.DisassembleByte(db)
			//if op.Bytes() > 1 {
			//	i += int(op.Bytes() - 1)
			//}
			fmt.Println(a.String())
		}
	}
	return nil
}

func testWram(file string) error {
	options := &nes.Options{
		Region: "NTSC",
	}

	ns, err := nes.NewNES(file, options)
	if err != nil {
		return err
	}

	ns.Reset()

	fmt.Println("Testing WRAM writes and reads:")

	// Test writing to WRAM at $6000-$7FFF
	testData := []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0}

	// Write to WRAM
	for i, val := range testData {
		ns.ROM.Store(0x6000+uint16(i), val)
	}

	// Read back from WRAM
	allGood := true
	for i, expected := range testData {
		val := ns.ROM.Fetch(0x6000 + uint16(i))
		status := "✓"
		if val != expected {
			status = "✗"
			allGood = false
		}
		fmt.Printf("  $%04X = 0x%02X (expected 0x%02X) %s\n", 0x6000+i, val, expected, status)
	}

	if allGood {
		fmt.Println("\n✓ WRAM is working!")
	} else {
		fmt.Println("\n✗ WRAM is not working properly!")
	}

	// Test CPU RAM at $0000-$07FF
	fmt.Println("\nTesting CPU RAM writes and reads:")

	// The NES instance should have CPU with RAM
	// But we need to go through the CPU's memory interface
	// Let's test through the memory interface
	for i, val := range testData {
		ns.CPU.M6502.Memory.Store(0x0400+uint16(i), val)
	}

	allGood = true
	for i, expected := range testData {
		val := ns.CPU.M6502.Memory.Fetch(0x0400 + uint16(i))
		status := "✓"
		if val != expected {
			status = "✗"
			allGood = false
		}
		fmt.Printf("  $%04X = 0x%02X (expected 0x%02X) %s\n", 0x0400+i, val, expected, status)
	}

	if allGood {
		fmt.Println("\n✓ CPU RAM is working!")
	} else {
		fmt.Println("\n✗ CPU RAM is not working properly!")
	}
	return nil
}
