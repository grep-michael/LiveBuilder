package main

import (
	"io"
	"log"
	//"os"
	"strings"
)

type IndentWriter struct {
	output    io.Writer
	indent    string
	firstLine bool
}

func NewIndentWriter(output io.Writer, indent string) *IndentWriter {
	return &IndentWriter{
		output:    output,
		indent:    indent,
		firstLine: true,
	}
}

func (w *IndentWriter) Write(p []byte) (n int, err error) {
	input := string(p)
	lines := strings.Split(input, "\n")

	var result strings.Builder

	for i, line := range lines {
		// Skip the last empty line if input ends with \n
		if i == len(lines)-1 && line == "" {
			continue
		}

		if w.firstLine && line != "" {
			result.WriteString(line)
			w.firstLine = false
		} else if line != "" {
			result.WriteString(w.indent + line)
		}

		// Add newline except for the last non-empty line
		if i < len(lines)-1 || (i == len(lines)-1 && line != "") {
			result.WriteString("\n")
		}
	}

	// Reset firstLine for next log entry when we see a complete line ending
	if strings.HasSuffix(input, "\n") {
		w.firstLine = true
	}

	written, err := w.output.Write([]byte(result.String()))
	return written, err
}

// Init sets up the global logger with indented output
func InitLoging(prefix, indent string, out io.Writer) {
	indentWriter := NewIndentWriter(out, indent)
	log.SetOutput(indentWriter)
	log.SetPrefix(prefix)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
