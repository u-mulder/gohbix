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
    // TODO - intellect gues what current folder is
    fmt.Println("Creating component '" + cd.Namespace + ":" + cd.Name + "'")

    var path string
    var hasLangs bool

    path = curPath + "/" + cd.Namespace
    // check if exists
    if mkDir(path) {
        path = path + "/" + cd.Name
        if mkDir(path) {
            mkFileWithContents(path + "/component.php", FL_TYPE_TPL)
            if 0 < len(cd.Langs) {
                hasLangs = true
            }

            path = path + "/templates"
            if mkDir(path) {
                for _, v := range cd.TplNames {
                    if mkDir (path + "/" + v) {
                        mkFileWithContents(path + "/" + v + "/template.php", FL_TYPE_TPL)

                        if hasLangs {
                            if mkDir (path + "/" + v + "/lang") {
                                for _, lv := range cd.Langs {
                                    if mkDir(path + "/" + v + "/lang/" + lv) {
                                        mkFileWithContents(path + "/" + v + "/lang/" + lv + "/template.php", FL_TYPE_LANG)
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    fmt.Println("Creating component '" + cd.Namespace + ":" + cd.Name + "' completed")
}

type LangFileData struct {
    //CurPosition string
    //TemplateName string
    Langs []string
}

func (lfd *LangFileData) Create() {
    fmt.Println("Creating lang files")

    langPath := curPath + "/lang"
    if mkDir(langPath) {
        for _, v := range lfd.Langs {
            if mkDir(langPath + "/" + v) {
                mkFileWithContents(langPath + "/" + v + "/template.php", FL_TYPE_LANG)
            }
        }
    }

    fmt.Println("Creating lang files completed")
}

type ModuleData struct {
    Name string
    AddOptions bool
    AddAdminFolder bool
    InstallFolders []string
    Langs []string
}

func (md *ModuleData) Create() {
    fmt.Println("Creating module '" + md.Name + "'")

    path := curPath + "/" + md.Name
    var ipath string
    if mkDir(path) {
        mkFileWithContents(path + "/include.php", FL_TYPE_OTAG)

        ipath = path + "/install"
        if mkDir(ipath) {
            mkFileWithContents(ipath + "/index.php", FL_TYPE_OTAG)
            mkFileWithContents(ipath + "/version.php", FL_TYPE_OTAG)
            for _, v := range md.InstallFolders  {
                mkDir(ipath + "/" + v)
            }
        }

        if md.AddOptions {
            mkFileWithContents(path + "/options.php", FL_TYPE_OTAG)
        }

        if md.AddAdminFolder {
            mkDir(path + "/admin")
        }

        if 0 < len(md.Langs) {
            // TODO
        }

    }

    fmt.Println("Creating module '" + md.Name + "' completed")
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

    defMdlName = "my.module.name"
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
    FL_TYPE_OTAG

    LF_POS_COMPONENT
    LF_POS_TEMPLATES
    LF_POS_TEMPLATE
)

func init() {
    tpls = make(map[int]string)
    tpls[FL_TYPE_TPL] = "<?php\nif (!defined(\"B_PROLOG_INCLUDED\") || B_PROLOG_INCLUDED!==true) die();\n"
    tpls[FL_TYPE_LANG] = "<?php\n$MESS[''] = '';\n"
    tpls[FL_TYPE_OTAG] = "<?php\n"
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

func createLangFile() {
    // TODO - intellect guessing what folder is it - component, templates, certain template
    langName := "_"
    newLfd := new(LangFileData)

    fmt.Println("Lang file(s) and related folder(s) will be created in current path '" + curPath + "'")
    for "" != langName {
        fmt.Println("Enter lang. Leave blank to stop entering langs: ")
        scanner.Scan()
        langName = clearTextValue(scanner.Text())
        if "" != langName {
            newLfd.Langs = append(newLfd.Langs, langName)
        }
    }

    newLfd.Create()
}

func createModule() {
    var mdlName, addOpts, addAfldr, addIfldrs, fldName, addLangs, langName string
    newMd := new(ModuleData)

    fmt.Println("Module will be created from current path '" + curPath + "'")
    fmt.Println("Enter <fg=blue;bg=red>module</> name. Leave blank to create '" + defMdlName + "' template")
    scanner.Scan()
    mdlName = clearTextValue(scanner.Text())
    if "" == mdlName {
        mdlName = defMdlName
    }
    newMd.Name = mdlName

    fmt.Println("Add options.php file? (y/n):")
    scanner.Scan()
    addOpts = clearTextValue(scanner.Text())
    newMd.AddOptions = (addOpts == "y" || addOpts == "yes")

    fmt.Println("Add admin folder? (y/n):")
    scanner.Scan()
    addAfldr = clearTextValue(scanner.Text())
    newMd.AddAdminFolder = (addAfldr == "y" || addAfldr == "yes")

    fmt.Println("Add install folders? (y/n):")
    scanner.Scan()
    addIfldrs = clearTextValue(scanner.Text())
    if "y" == addIfldrs || "yes" == addIfldrs {
        fldName = "_"
        for "" != fldName {
            fmt.Println("Enter folder name. Leave blank to stop entering folders: ")
            scanner.Scan()
            fldName = clearTextValue(scanner.Text())
            if "" != fldName  {
                newMd.InstallFolders = append(newMd.InstallFolders, fldName)
            }
        }
    }

    fmt.Println("Add language support? (y/n):")
    scanner.Scan()
    addLangs = clearTextValue(scanner.Text())
    if addLangs == "y" || addLangs == "yes" {
        fmt.Println("Default lang '" + defLang + "' already added")
        newMd.Langs = append(newMd.Langs, defLang)
        langName = defLang
        for "" != langName {
            fmt.Println("Enter lang. Leave blank to stop entering langs: ")
            scanner.Scan()
            langName = clearTextValue(scanner.Text())
            if "" != langName && defLang != langName {
                newMd.Langs = append(newMd.Langs, langName)
            }
        }
    }

    newMd.Create()
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
