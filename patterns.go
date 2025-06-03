// Package main provides pattern parsing and color support for ahl.
package main

import (
	"fmt"
	"regexp"
	"strings"
)

// patternColor associates a compiled regex with an ANSI color code.
type patternColor struct {
	regex *regexp.Regexp
	color string
}

// patternFlags is a slice of pattern strings for use with flag.Var.
type patternFlags []string

// String returns the string representation of patternFlags.
func (p *patternFlags) String() string {
	return strings.Join(*p, ", ")
}

// Set appends a new pattern to patternFlags.
func (p *patternFlags) Set(value string) error {
	*p = append(*p, value)
	return nil
}

// parsePatterns parses pattern arguments in the form REGEX=COLOR and returns compiled patterns.
func parsePatterns(patternArgs []string) ([]patternColor, error) {
	var patterns []patternColor
	for _, arg := range patternArgs {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid pattern (expected REGEX=COLOR): %s", arg)
		}
		pattern, colorName := parts[0], strings.ToLower(parts[1])
		colorCode, ok := ansiColors[colorName]
		if !ok {
			return nil, fmt.Errorf("invalid color name '%s' in pattern '%s'", colorName, arg)
		}
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex '%s' in pattern '%s': %v", pattern, arg, err)
		}
		patterns = append(patterns, patternColor{
			regex: regex,
			color: colorCode,
		})
	}
	return patterns, nil
}

// getSupportedColors returns a list of supported color names.
func getSupportedColors() []string {
	var colors []string
	for color := range ansiColors {
		colors = append(colors, color)
	}
	return colors
}
