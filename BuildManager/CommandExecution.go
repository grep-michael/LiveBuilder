package buildmanager

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
	"unicode"
)

type OutputType string

const (
	STDERR OutputType = "STDERR"
	STDOUT OutputType = "STDOUT"
)

type CommandOut struct {
	OutType OutputType
	OutPut  string
}

func executeCommand(cmd *exec.Cmd, outputChannel chan CommandOut) error {
	log.Println("executing command")
	log.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error occured getting stdout pipe: %v\n", err)
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error occured getting stderr pipe: %v\n", err)
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return err
	}

	// Read stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			msg := scanner.Text() + "\n"
			outputChannel <- CommandOut{
				OutType: STDOUT,
				OutPut:  msg,
			}
		}
	}()

	// Read stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			msg := scanner.Text() + "\n"
			outputChannel <- CommandOut{
				OutType: STDERR,
				OutPut:  msg,
			}
		}
	}()

	return cmd.Wait()
}

func parseShellCommand(input string) []string {
	var tokens []string
	var current strings.Builder
	var inQuotes bool
	var quoteChar rune

	//Stanitize input
	input = strings.ReplaceAll(input, "\\\n", " ")
	input = strings.ReplaceAll(input, "\\\r\n", " ")
	input = strings.Join(strings.Fields(input), " ")

	for _, char := range input {
		switch {
		case char == '"' || char == '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else {
				current.WriteRune(char)
			}

		case unicode.IsSpace(char) && !inQuotes:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
