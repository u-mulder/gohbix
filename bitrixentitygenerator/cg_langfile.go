package bitrixentitygenerator

import (
	"bufio"
	"fmt"
)

type LangFileData struct {
    //CurPosition string
    //TemplateName string
    Langs []string
    EntityPathData
}

func NewLangFileGenerator(path string) (EntityGenerator, error) {
    lfd := LangFileData{}
    lfd.Path = path

    return &lfd, nil
}

func (lfd *LangFileData) CollectEntityParameters(scanner *bufio.Scanner) {
    langName := "_"

    fmt.Println("Lang file(s) and related folder(s) will be created in current path '" + colorize(lfd.Path, "green") + "'")
    for "" != langName {
        fmt.Println("Enter lang. Leave blank to stop entering langs: ")
        scanner.Scan()
        langName = clearTextValue(scanner.Text())
        if "" != langName {
            lfd.Langs = append(lfd.Langs, langName)
        }
    }
}

// TODO - intellect guessing what folder is it (component, templates, certain template)
func (lfd *LangFileData) CreateEntity() {
    var err error

    if 0 < len(lfd.Langs) {
        langPath := lfd.Path + "/lang"
        err = mkDir(langPath)

        if err {
            for _, v := range lfd.Langs {
                err := mkDir(langPath + "/" + v)
                if err != nil {
                    // TODO
                    //err = mkFileWithContents(langPath + "/" + v + "/template.php", FL_TYPE_LANG)
                    if err {
                        // TODO
                    }
                } else {
                    // TODO
                    //fmt.Println(colorize("/!\\ Langs list empty", "red") 
                }
            }
        } else {
            // TODO
            //fmt.Println(colorize("/!\\ Langs list empty", "red")
        }
    } else {
        fmt.Println(colorize("/!\\ Language list is empty", "yellow")
    }

    fmt.Println(colorize("Creating lang files completed", "green"))
}
