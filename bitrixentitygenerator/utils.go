package bitrixentitygenerator 

import "strings"

var ansiColorCodes = map[string]string{
	"red": "\x1b[31m",
    "green": "\x1b[32m",
    "yellow": "\x1b[33m",
    "blue": "\x1b[34m",
    "magenta": "\x1b[35m",
    "cyan": "\x1b[36m",
    "reset": "\x1b[0m",
}

func isYes(str string) bool {
    return (str == "y" || str == "Y" || str == "yes")
}

func clearTextValue(txt string) string {
    return strings.TrimSpace(txt)
}

func colorize(str string, ansiColorCodeKey string) string {
	ansiColorCodeValue, ok := ansiColorCodes[ansiColorCodeKey]
	if !ok {
		return str
	}

	return ansiColorCodeValue + str + ansiColorCodes["reset"]
}
