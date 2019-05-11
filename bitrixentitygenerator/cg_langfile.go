package bitrixentitygenerator

import (
	"bufio"
	"fmt"
)

type langFileData struct {
	//CurPosition string
	//TemplateName string
	Langs []string
	EntityPathData
}

// NewLangFileGenerator returns new instance of `langFileData`
func NewLangFileGenerator(path string) (EntityGenerator, error) {
	lfd := langFileData{}
	lfd.Path = path

	return &lfd, nil
}

func (lfd *langFileData) CollectEntityParameters(scanner *bufio.Scanner) {
	langName := "_"

	fmt.Println("Lang file(s) and related folder(s) will be created in current path '" + colorize(lfd.Path, "green") + "'")
	for "" != langName {
		fmt.Println("Enter language. Leave blank to stop entering langs: ")
		scanner.Scan()
		// TODO - language usually has 2 symbols (ru, en, ua, de)
		langName = clearTextValue(scanner.Text())
		if "" != langName {
			lfd.Langs = append(lfd.Langs, langName)
		}
	}
}

// TODO - intellect guessing what folder is it (component, templates, certain template)
func (lfd *langFileData) CreateEntity() {
	var err error

	if 0 < len(lfd.Langs) {
		langPath := lfd.Path + "/lang"
		err = mkDir(langPath)

		if err == nil {
			for _, v := range lfd.Langs {
				err := mkDir(langPath + "/" + v)
				if err == nil {
					mkFileWithContents(langPath+"/"+v+"/template.php", langFileLine)
				}
			}
		}
	} else {
		fmt.Println(colorize("/!\\ Language list is empty", "yellow"))
	}

	fmt.Println(colorize("Creating lang files completed", "green"))
}
