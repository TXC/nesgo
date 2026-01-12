# Klaus Dormann's 65o2 Test Suites

Included assembler code, listings and binary files

Compiler is available in (Klaus Dormann's Repository)[https://github.com/Klaus2m5/6502_65C02_functional_tests]

Compile with:
```
as65 -l -m -s2 -w -h0
      |  |  |   |  no page headers in listing
      |  |  |   wide listing (133 char/col)
      |  |  write intel hex file instead of binary
      |  expand macros in listing
      generate pass2 listing
```