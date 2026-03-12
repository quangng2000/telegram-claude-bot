package claude

import (
	"regexp"
	"strings"
)

var ansiPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\x1B\[[0-9;?]*[a-zA-Z<>]`),
	regexp.MustCompile(`\x1B\][^\x07]*\x07`),
	regexp.MustCompile(`\x1B\][^\x1B]*\x1B\\`),
	regexp.MustCompile(`\x1B[()][A-Z0-9]`),
	regexp.MustCompile(`\x1B[>=<]`),
	regexp.MustCompile(`\x1B[@-Z\\-_]`),
	regexp.MustCompile(`[\x07\x00-\x06\x08\x0E\x0F]`),
}

func StripAnsi(text string) string {
	for _, re := range ansiPatterns {
		text = re.ReplaceAllString(text, "")
	}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "")
	return strings.TrimSpace(text)
}
