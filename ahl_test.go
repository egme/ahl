package main

import (
	"regexp"
	"testing"
)

func TestParsePatterns_Valid(t *testing.T) {
	args := []string{"foo=red", "bar=green"}
	patterns, err := parsePatterns(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(patterns) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(patterns))
	}
	if patterns[0].color != ansiColors["red"] {
		t.Errorf("expected color red, got %s", patterns[0].color)
	}
	if patterns[1].color != ansiColors["green"] {
		t.Errorf("expected color green, got %s", patterns[1].color)
	}
}

func TestParsePatterns_InvalidColor(t *testing.T) {
	args := []string{"foo=notacolor"}
	_, err := parsePatterns(args)
	if err == nil {
		t.Fatal("expected error for invalid color, got nil")
	}
}

func TestParsePatterns_InvalidFormat(t *testing.T) {
	args := []string{"foo"}
	_, err := parsePatterns(args)
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestParsePatterns_InvalidRegex(t *testing.T) {
	args := []string{"[=red"}
	_, err := parsePatterns(args)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestHighlightLine(t *testing.T) {
	patterns := []patternColor{
		{regex: regexp.MustCompile("foo"), color: ansiColors["red"]},
		{regex: regexp.MustCompile("bar"), color: ansiColors["green"]},
	}
	line := "foo and bar"
	highlighted := highlightLine(line, patterns)
	if highlighted == line {
		t.Errorf("expected highlighted output, got unchanged line")
	}
	if !regexp.MustCompile("\\x1b\\[31mfoo\\x1b\\[0m").MatchString(highlighted) {
		t.Errorf("expected foo to be highlighted in red")
	}
	if !regexp.MustCompile("\\x1b\\[32mbar\\x1b\\[0m").MatchString(highlighted) {
		t.Errorf("expected bar to be highlighted in green")
	}
}

func TestStripANSI(t *testing.T) {
	input := "\x1b[31mfoo\x1b[0m and \x1b[32mbar\x1b[0m"
	expected := "foo and bar"
	stripped := stripANSI(input)
	if stripped != expected {
		t.Errorf("expected '%s', got '%s'", expected, stripped)
	}
}

// Note: --cleanup flag logic should be tested via integration tests or manually, as it depends on the runtime environment.

// ---
// Integration test samples (for manual testing or future automation):
//
// 1. Pattern highlighting with long and short flags:
//    echo "error: something failed" | ./ahl --pattern 'error=red' --pattern 'failed=yellow'
//    echo "error: something failed" | ./ahl -p 'error=red' -p 'failed=yellow'
//
// 2. Single pattern (default red):
//    echo "warning: disk space low" | ./ahl 'warning'
//
// 3. --cleanup/-c flag with colored input:
//    echo -e "\033[31merror\033[0m: something failed" | ./ahl --cleanup --pattern 'error=red'
//    echo -e "\033[31merror\033[0m: something failed" | ./ahl -c -p 'error=red'
//
// 4. --cleanup/-c with no patterns (should error):
//    echo -e "\033[31merror\033[0m: something failed" | ./ahl --cleanup
//
