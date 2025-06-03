package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// stripANSI removes ANSI escape sequences from a string.
func stripANSI(s string) string {
	ansiRegexp := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegexp.ReplaceAllString(s, "")
}

// patternFlagsWithShort allows both --pattern and -p to be used.
type patternFlagsWithShort struct {
	patternFlags
}

func (p *patternFlagsWithShort) Set(value string) error {
	return p.patternFlags.Set(value)
}

func (p *patternFlagsWithShort) String() string {
	return p.patternFlags.String()
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `ahl: highlight patterns in text streams using ANSI colors

Usage:
  ahl [OPTIONS] [PATTERN]

Options:
  --pattern, -p   Pattern in form 'REGEX=COLOR' (can be used multiple times)
  --cleanup, -c   Strip all incoming ANSI color codes before highlighting
  --help          Show this help message

If a single argument is provided (not using --pattern/-p), it is treated as a pattern with red color.

Supported colors:
  %s

Examples:
  echo "error: something failed" | ahl --pattern 'error=red' --pattern 'failed=yellow'
  echo "error: something failed" | ahl -p 'error=red' -p 'failed=yellow'
  echo "warning: disk space low" | ahl 'warning'
  echo -e "\033[31merror\033[0m: something failed" | ahl --cleanup --pattern 'error=red'
  echo -e "\033[31merror\033[0m: something failed" | ahl -c -p 'error=red'

`, strings.Join(getSupportedColors(), ", "))
}

// main is the entry point for the ahl CLI tool.
func main() {
	flag.Usage = printUsage

	var patternArgs patternFlagsWithShort
	var cleanup bool
	flag.Var(&patternArgs, "pattern", "Pattern in form 'REGEX=COLOR' (can be used multiple times)")
	flag.Var(&patternArgs, "p", "Shorthand for --pattern")
	flag.BoolVar(&cleanup, "cleanup", false, "Strip all incoming ANSI color codes before highlighting")
	flag.BoolVar(&cleanup, "c", false, "Shorthand for --cleanup")
	flag.Parse()

	// If only one non-flag argument is provided, treat it as a pattern and use red
	if len(patternArgs.patternFlags) == 0 && len(flag.Args()) == 1 {
		patternArgs.patternFlags = append(patternArgs.patternFlags, flag.Args()[0]+"=red")
	}

	if len(patternArgs.patternFlags) == 0 {
		printUsage()
		os.Exit(1)
	}

	patterns, err := parsePatterns(patternArgs.patternFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing patterns: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if cleanup {
			line = stripANSI(line)
		}
		line = highlightLine(line, patterns)
		_, _ = writer.WriteString(line + "\n")
		writer.Flush()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
