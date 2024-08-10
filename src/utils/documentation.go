package utils

import (
	"strings"
)

var (
	// AlphaDisclaimer to be places at the end of description of commands in alpha release
	AlphaDisclaimer = `
		Alpha Disclaimer: this command is currently alpha.
	`

	// MacroCommandLongDescription provide a standard description for "macro" commands
	MacroCommandLongDescription = LongDesc(`
		This command is not meant to be run on its own. See list of available subcommands.
	`)
)

func LongDesc(s string) string {
	// Strip beginning and trailing space characters (including empty lines) and split the lines into a slice
	lines := strings.Split(strings.TrimSpace(s), "\n")

	output := []string{}
	paragraph := []string{}

	for _, line := range lines {
		// Remove indentation and trailing spaces from the current line
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			if len(paragraph) > 0 {
				output = append(output, strings.Join(paragraph, " ")+"\n")
				paragraph = []string{}
			}
		} else {
			// Non-empty text line, append it to the current paragraph
			paragraph = append(paragraph, trimmedLine)
			if strings.HasSuffix(line, "  ") {
				output = append(output, strings.Join(paragraph, " "))
				paragraph = []string{}
			}
		}
	}

	output = append(output, strings.Join(paragraph, " "))

	// Join all the paragraphs together with new lines in between them.
	return strings.Join(output, "\n")
}

// Examples is designed to help with producing examples for command line usage.
func Examples(s string) string {
	trimmedText := strings.TrimSpace(s)
	if trimmedText == "" {
		return ""
	}

	const indent = `  `
	inLines := strings.Split(trimmedText, "\n")
	outLines := make([]string, 0, len(inLines))

	for _, line := range inLines {
		outLines = append(outLines, indent+strings.TrimSpace(line))
	}

	return strings.Join(outLines, "\n")
}
