package utils

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Logger struct {
	writer io.Writer
}

const indent = "    "

func NewLogger(w io.Writer) *Logger {
	return &Logger{
		writer: w,
	}
}

func FromLogger(log *log.Logger) *Logger {
	return &Logger{
		writer: log.Writer(),
	}
}

// Log formats its arguments using default formatting, analogous to [fmt.Println],
// and records the text in the error log.
// For tests, the text will be printed only if the test fails or the -test.v flag is
// set.
// For benchmarks, the text is always printed to avoid having performance
// depend on the value of the -test.v flag.
// It is an error to call Log after a test or benchmark returns.
func (c *Logger) Log(args ...any) { c.log(fmt.Sprintln(args...)) }

// Logf formats its arguments according to the format, analogous to [fmt.Printf], and
// records the text in the error log. A final newline is added if not provided.
// For tests, the text will be printed only if the test fails or the -test.v flag is
// set.
// For benchmarks, the text is always printed to avoid having performance
// depend on the value of the -test.v flag.
// It is an error to call Logf after a test or benchmark returns.
func (c *Logger) Logf(format string, args ...any) {
	c.log(fmt.Sprintf(format, args...))
}

// log generates the output. It is always at the same stack depth. log inserts
// indentation and the final newline if necessary.
func (c *Logger) log(s string) {
	if c.writer == nil {
		return
	}

	s = strings.TrimSuffix(s, "\n")

	// Second and subsequent lines are indented 4 spaces. This is in addition to
	// the indentation provided by outputWriter.
	s = strings.ReplaceAll(s, "\n", "\n"+indent)
	s += "\n"

	c.writer.Write([]byte(s))
}
