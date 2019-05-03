package main

import (
	"bufio"
	beg "github.com/u-mulder/gohbix/bitrixentitygenerator"
	"os"
	"path/filepath"
)

var scanner *bufio.Scanner
var curPath string

func main() {
	var err error
	curPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("/!\\ Error defining current path")
	}

	beg.RegisterFactories()

	scanner = bufio.NewScanner(os.Stdin)
	entityType := beg.ScanAndGetEntityType(scanner)
	entityGenerator, err := beg.GetGenerator(entityType, curPath)
	if err != nil {
		panic(err)
	}
	entityGenerator.CollectEntityParameters(scanner)
	entityGenerator.CreateEntity()
}
