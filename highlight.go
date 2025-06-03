// Package main provides highlighting logic for ahl.
package main

import "fmt"

// ansiColors maps color names to their ANSI color codes.
var ansiColors = map[string]string{
	"black": "30", "red": "31", "green": "32", "yellow": "33",
	"blue": "34", "magenta": "35", "cyan": "36", "white": "37",
	"gray": "90", "brightred": "91", "brightgreen": "92", "brightyellow": "93",
	"brightblue": "94", "brightmagenta": "95", "brightcyan": "96", "brightwhite": "97",
}

// highlightLine applies all pattern highlights to a single line.
func highlightLine(line string, patterns []patternColor) string {
	for _, pc := range patterns {
		line = pc.regex.ReplaceAllStringFunc(line, func(match string) string {
			return fmt.Sprintf("\x1b[%sm%s\x1b[0m", pc.color, match)
		})
	}
	return line
}
