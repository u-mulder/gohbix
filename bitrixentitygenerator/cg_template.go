package bitrixentitygenerator

import (
	"bufio"
	"fmt"
)

type templateData struct {
	Name         string
	AddResModfr  bool
	AddCmpEpilog bool
	Langs        []string
	EntityPathData
}

func NewTemplateGenerator(path string) (EntityGenerator, error) {
	td := templateData{}
	td.Path = path

	return &td, nil
}

func (td *templateData) CollectEntityParameters(scanner *bufio.Scanner) {
	fmt.Println("Template will be created from current path '" + colorize(td.Path, "green") + "'")
	fmt.Println("Enter template name. Leave blank to create '" + colorize(defaultTemplateName, "green") + "' template")
	scanner.Scan()
	td.Name = clearTextValue(scanner.Text())
	if "" == td.Name {
		td.Name = defaultTemplateName
	}

	fmt.Println("Add result_modifier file? (y/n):")
	scanner.Scan()
	td.AddResModfr = isYes(clearTextValue(scanner.Text()))

	fmt.Println("Add component_epilog file? (y/n):")
	scanner.Scan()
	td.AddCmpEpilog = isYes(clearTextValue(scanner.Text()))

	fmt.Println("Add language support? (y/n):")
	scanner.Scan()
	if isYes(clearTextValue(scanner.Text())) {
		fmt.Println("Default lang '" + colorize(defaultLanguage, "green") + "' already added")
		td.Langs = append(td.Langs, defaultLanguage)
		langName := defaultLanguage
		for "" != langName {
			fmt.Println("Enter lang. Leave blank to stop entering langs: ")
			scanner.Scan()
			langName = clearTextValue(scanner.Text())
			// TODO - languages can repeat
			if "" != langName && defaultLanguage != langName {
				td.Langs = append(td.Langs, langName)
			}
		}
	}
}

// TODO - intellect guessing what folder is it (component, templates, certain template)
func (td *templateData) CreateEntity() {
	var err error

	fmt.Println("Creating template '" + colorize(td.Name, "green") + "'")

	tplPath := td.Path + "/" + td.Name
	err = mkDir(tplPath)

	if err == nil {
		mkFileWithContents(tplPath+"/template.php", prologCheckLine)

		if td.AddResModfr {
			mkFileWithContents(tplPath+"/result_modifier.php", prologCheckLine)
		}

		if td.AddCmpEpilog {
			mkFileWithContents(tplPath+"/component_epilog.php", prologCheckLine)
		}

		if 0 < len(td.Langs) {
			langPath := tplPath + "/lang"
			err = mkDir(langPath)

			if err == nil {
				for _, v := range td.Langs {
					err := mkDir(langPath + "/" + v)

					if err == nil {
						mkFileWithContents(langPath+"/"+v+"/template.php", langFileLine)
					}
				}
			}
		}
	}

	fmt.Println(colorize("Creating lang '"+td.Name+"' completed", "green"))
}
