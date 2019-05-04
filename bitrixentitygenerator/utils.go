package bitrixentitygenerator

import (
	"fmt"
	"os"
	"strings"
)

var ansiColorCodes = map[string]string{
	"red":     "\x1b[31m",
	"green":   "\x1b[32m",
	"yellow":  "\x1b[33m",
	"blue":    "\x1b[34m",
	"magenta": "\x1b[35m",
	"cyan":    "\x1b[36m",
	"reset":   "\x1b[0m",
}

const (
	langFileLine = iota
	openTag
	prologCheckLine
)

var codeContents = map[int]string{
	0: "<?php\n$MESS[''] = '';\n",
	1: "<?php\n",
	2: "<?php\nif (!defined(\"B_PROLOG_INCLUDED\") || B_PROLOG_INCLUDED!==true) die();\n",
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

func getContentByType(contType int) string {
	content, ok := codeContents[contType]
	if ok {
		return content
	}

	return ""
}

func isYes(str string) bool {
	return (str == "y" || str == "Y" || str == "yes")
}

func mkDir(name string) error {
	if err := os.Mkdir(name, defaultDirPermissions); err != nil {
		fmt.Println("Error creating path '"+colorize(name, "green")+"': ", colorize(err.Error(), "red"))
		return err
	}

	return nil
}

func mkFileWithContents(name string, contType int) error {
	fh, err := os.Create(name)

	if err != nil {
		fmt.Println("Error creating file '"+colorize(name, "green")+"': ", colorize(err.Error(), "red"))
		return err
	}

	_ = os.Chmod(name, defaultFilePermissions)
	_, err = fh.Write([]byte(getContentByType(contType)))

	if err != nil {
		fmt.Println("Error writing contents to file '"+colorize(name, "green")+"': ", colorize(err.Error(), "red"))
		return err
	}

	return nil
}
