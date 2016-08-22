package main

import (
    "fmt"
    "flag"
    "os"
    "path/filepath"
    "bufio"
    "strings"
)

var scanner *bufio.Scanner
var curPath string
var tpls map[int]string

type EntityCreator interface {
    Create()
}

type ComponentData struct {
    Namespace string
    Name string
    TplNames []string
    Langs []string
}

func (cd *ComponentData) Create() {

    fmt.Println( "cd is ", cd )     // TODO

}

type LangFileData struct {

}

func (lfd *LangFileData) Create() {

    fmt.Println( "lfd is ", lfd )     // TODO

}

type ModuleData struct {
    Name string
}

func (md *ModuleData) Create() {

    // TODO
    fmt.Println("md is ", md)

}

type TemplateData struct {
    Name string
    AddResModfr bool
    AddCmpEpilog bool
    Langs []string
}

func (td *TemplateData) Create() {
    fmt.Println("Creating template " + td.Name)
    tplPath := curPath + "/" + td.Name
    if mkDir(tplPath) {
        mkFileWithContents(tplPath + "/template.php", FL_TYPE_TPL)

        if td.AddResModfr {
            mkFileWithContents(tplPath + "/result_modifier.php", FL_TYPE_TPL)
        }

        if td.AddCmpEpilog {
            mkFileWithContents(tplPath + "/component_epilog.php", FL_TYPE_TPL)
        }

        if 0 < len(td.Langs) {
            tplPath := tplPath + "/lang"
            if mkDir(tplPath) {
                for _, v := range td.Langs {
                    langPath := tplPath + "/" + v
                    if mkDir(langPath) {
                        mkFileWithContents(langPath + "/template.php", FL_TYPE_LANG)
                    }
                }
            }
        }
    }

    fmt.Println("Creating template " + td.Name + " completed.")
}

const (
    EN_C = "c"
    EN_L = "l"
    EN_M = "m"
    EN_T = "t"

    defNamespace = "myns"
    defCmpName = "my.component"
    defLang = "ru"
    defTplName = ".default"
    defDirMode = 0755
    defFileMode = 0644
)

// reset iota
const (
    FL_TYPE_TPL = iota
    FL_TYPE_LANG
)

func init() {
    tpls = make(map[int]string)
    tpls[FL_TYPE_TPL] = "<?php\nif (!defined(\"B_PROLOG_INCLUDED\") || B_PROLOG_INCLUDED!==true) die();\n"
    tpls[FL_TYPE_LANG] = "<?php\n$MESS[''] = '';\n"
}

func main() {
    var err error
    curPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        panic("/!\\ Error defining current path")
    }

    var entityType string
    flag.StringVar(&entityType, "config", "c", "Entity type to create, default value is 'c' which stands for 'component'")
    flag.Parse()

    scanner = bufio.NewScanner(os.Stdin)

    switch entityType {
        case EN_C:
            createComponent()

        case EN_L:
            createLangFile()

        case EN_M:
            createModule()

        case EN_T:
            createTemplate()
    }
}

func createComponent() {
    var cmpNs, cmpName, moreTpls, tplName, addLangs, langName string
    newCp := new(ComponentData)

    fmt.Println("Component will be created from current path '" + curPath + "'")
    fmt.Print("Enter namespace ('" + defNamespace + "' is default), if folder not exists - it will be created: ")
    scanner.Scan()
    cmpNs = clearTextValue(scanner.Text())
    if "" == cmpNs {
        cmpNs = defNamespace
    }
    newCp.Namespace = cmpNs

    fmt.Println("Enter component name ('" + defCmpName + "' is default), if component exists - you'll have to additionally confirm overwrite: ")
    scanner.Scan()
    cmpName = clearTextValue(scanner.Text())
    if "" == cmpName {
        cmpName = defCmpName
    }
    newCp.Name = cmpName

    fmt.Println("Create other templates except '" + defTplName + "'? (y/n):")
    scanner.Scan()
    moreTpls = clearTextValue(scanner.Text())
    if moreTpls == "y" || moreTpls == "yes" {
        tplName = "y"
        for "" != tplName {
            fmt.Println("Enter template name ('" + defTplName + "' will be added automatically). Leave blank to stop entering template names: ")
            scanner.Scan()
            tplName = clearTextValue(scanner.Text())
            if "" != tplName && tplName != defTplName {
                newCp.TplNames = append(newCp.TplNames, tplName)
            }
        }
    }
    newCp.TplNames = append(newCp.TplNames, ".default")

    fmt.Println("Add language support? (y/n):")
    scanner.Scan()
    addLangs = clearTextValue(scanner.Text())
    if addLangs == "y" || addLangs == "yes" {
        fmt.Println("Default lang '" + defLang + "' already added")
        newCp.Langs = append(newCp.Langs, defLang)
        langName = defLang
        for "" != langName {
            fmt.Println("Enter lang. Leave blank to stop entering langs: ")
            scanner.Scan()
            langName = clearTextValue(scanner.Text())
            if "" != langName && defLang != langName {
                newCp.Langs = append(newCp.Langs, langName)
            }
        }
    }

    newCp.Create()
}

func createLangFile() {     // TODO
    // TODO
    fmt.Println("Lang file and it's folder will be created in current path")
    fmt.Println("Enter lang (two symbols, 'ru' is default)")
    // scan

    //


}

func createModule() {  // TODO
    fmt.Println("This option is under development!")
    // TODO
}

func createTemplate() {
    var tplName, addLangs, langName, addResModfr, addCmpEpilog string
    newTp := new(TemplateData)

    fmt.Println("Template will be created from current path '" + curPath + "'")
    fmt.Println("Enter template name. Leave blank to create '" + defTplName + "' template")
    scanner.Scan()
    tplName = clearTextValue(scanner.Text())
    if "" == tplName {
        tplName = defTplName
    }
    newTp.Name = tplName

    fmt.Println("Add result_modifier file? (y/n):")
    scanner.Scan()
    addResModfr = clearTextValue(scanner.Text())
    if addResModfr == "y" || addResModfr == "yes" {
        newTp.AddResModfr = true
    }

    fmt.Println("Add component_epilog file? (y/n):")
    scanner.Scan()
    addCmpEpilog = clearTextValue(scanner.Text())
    if addCmpEpilog == "y" || addCmpEpilog == "yes" {
        newTp.AddCmpEpilog = true
    }

    fmt.Println("Add language support? (y/n):")
    scanner.Scan()
    addLangs = clearTextValue(scanner.Text())
    if addLangs == "y" || addLangs == "yes" {
        fmt.Println("Default lang '" + defLang + "' already added")
        newTp.Langs = append(newTp.Langs, defLang)
        langName = defLang
        for "" != langName {
            fmt.Println("Enter lang. Leave blank to stop entering langs: ")
            scanner.Scan()
            langName = clearTextValue(scanner.Text())
            if "" != langName && defLang != langName {
                newTp.Langs = append(newTp.Langs, langName)
            }
        }
    }

    newTp.Create()
}

func clearTextValue(txt string) string {
    return strings.TrimSpace(txt)
}

func mkDir(name string) bool {
    var r bool
    if err := os.Mkdir(name, defDirMode); err == nil {
        r = true
    } else {
        fmt.Println("Error creating path '" + name + "': ", err.Error())
    }

    return r
}

func mkFileWithContents(name string, contType int) bool {
    var r bool

    if fh, err := os.Create(name); err == nil {
        r = true
        _ = os.Chmod(name, defFileMode)

        _, err = fh.Write([]byte(tpls[contType]))

        if err != nil {
            fmt.Println("Error writing contents to file '" + name + "': ", err.Error())
        }

    } else {
        fmt.Println("Error creating file '" + name + "': ", err.Error())
    }

    return r
}
