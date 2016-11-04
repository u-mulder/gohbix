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

// ComponentData is a struct which stores new component info
type ComponentData struct {
    Namespace string
    Name string
    TplNames []string
    Langs []string
}

// Create folders and file structure for nw component
// TODO - intellect gues what current folder it is
// TODO - check if component path already exists
func (cd *ComponentData) Create() {
    fmt.Println("Creating component '" + cd.Namespace + ":" + cd.Name + "'")

    var path string

    path = curPath + "/" + cd.Namespace
    if mkDir(path) {
        path = path + "/" + cd.Name
        if mkDir(path) {
            mkFileWithContents(path + "/component.php", FL_TYPE_TPL)

            path = path + "/templates"
            if mkDir(path) {
                for _, v := range cd.TplNames {
                    if mkDir (path + "/" + v) {
                        mkFileWithContents(path + "/" + v + "/template.php", FL_TYPE_TPL)

                        if 0 < len(cd.Langs) {
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

// TODO - intellect guessing what folder is it (component, templates, certain template)
func (lfd *LangFileData) Create() {
    fmt.Println("Creating lang files")

    if 0 < len(lfd.Langs) {
        langPath := curPath + "/lang"
        if mkDir(langPath) {
            for _, v := range lfd.Langs {
                if mkDir(langPath + "/" + v) {
                    mkFileWithContents(langPath + "/" + v + "/template.php", FL_TYPE_LANG)
                }
            }
        }
    } else {
        fmt.Println("/!\\ Langs list empty")
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
    fmt.Printf("Creating module %s%s%s \n", ansiColorGreen, md.Name, ansiColorReset)

    path := curPath + "/" + md.Name
    var ipath string
    if mkDir(path) {
        mkFileWithContents(path + "/include.php", FL_TYPE_OTAG)

        ipath = path + "/install"
        langFiles := map[string][]string{}

        if mkDir(ipath) {
            mkFileWithContents(ipath + "/index.php", FL_TYPE_OTAG)
            mkFileWithContents(ipath + "/version.php", FL_TYPE_OTAG)
            for _, v := range md.InstallFolders  {
                mkDir(ipath + "/" + v)
            }

            langFiles["/install"] = []string{"/version.php", "/index.php"}
        }

        if md.AddOptions {
            mkFileWithContents(path + "/options.php", FL_TYPE_OTAG)
            langFiles["/"] = []string{"options.php"}
        }

        if md.AddAdminFolder {
            mkDir(path + "/admin")

            langFiles["/admin"] = []string{"/menu.php"}
        }

        if 0 < len(md.Langs) && 0 < len(langFiles) {
            langPath := path + "/lang"
            var curLangPath string
            var hasPath bool
            if mkDir(langPath) {
                for _, v := range md.Langs {
                    curLangPath = langPath + "/" + v
                    if mkDir(curLangPath) {
                        for lfPath, lfFiles := range langFiles {
                            hasPath = true
                            if "/" != lfPath {
                                hasPath = mkDir(curLangPath + lfPath)
                            }

                            if hasPath {
                                for _, lfFile := range lfFiles  {
                                    mkFileWithContents(curLangPath + lfPath + lfFile, FL_TYPE_LANG)
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    fmt.Printf("Creating module %s%s%s completed\n", ansiColorGreen, md.Name, ansiColorReset)
}

// TemplateData is a struct which stores new template info
type TemplateData struct {
    Name string
    AddResModfr bool
    AddCmpEpilog bool
    Langs []string
}

// TODO - intellect guessing what folder is it (component, templates, certain template)
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
    entCmp = "c"
    entLfl = "l"
    entMdl = "m"
    entTpl = "t"

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

    LfPosComponent
    LfPosTemplates
    LfPosTemplate
)

const (
    ansiColorRed        = "\x1b[31m"
    ansiColorGreen      = "\x1b[32m"
    ansiColorYellow     = "\x1b[33m"
    ansiColorBlue       = "\x1b[34m"
    ansiColorMagenta    = "\x1b[35m"
    ansiColorCyan       = "\x1b[36m"
    ansiColorReset      = "\x1b[0m"
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
        case entCmp:
            createComponent()

        case entLfl:
            createLangFile()

        case entMdl:
            createModule()

        case entTpl:
            createTemplate()

        default:
            fmt.Printf("Type %s not supported \n", entityType)
    }
}

// TODO - colorize messages

func createComponent() {
    var moreTpls, tplName, addLangs, langName string
    newCp := new(ComponentData)

    fmt.Println("Component will be created from current path '" + curPath + "'")
    fmt.Print("Enter namespace ('" + defNamespace + "' is default), if folder not exists - it will be created: ")
    scanner.Scan()
    newCp.Namespace = clearTextValue(scanner.Text())
    if "" == newCp.Namespace {
        newCp.Namespace = defNamespace
    }

    fmt.Println("Enter component name ('" + defCmpName + "' is default), if component exists - you'll have to additionally confirm overwrite: ")
    scanner.Scan()
    newCp.Name = clearTextValue(scanner.Text())
    if "" == newCp.Name {
        newCp.Name = defCmpName
    }

    fmt.Println("Create other templates except '" + defTplName + "'? (y/n):")
    scanner.Scan()
    moreTpls = clearTextValue(scanner.Text())
    if isYes(moreTpls) {
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
    if isYes(addLangs) {
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
    var addIfldrs, fldName, langName string
    newMd := new(ModuleData)

    fmt.Println("Module will be created from current path '" + curPath + "'")
    fmt.Printf("Enter module name. Leave blank to create %s%s%s template\n", ansiColorGreen, defMdlName, ansiColorReset)
    scanner.Scan()
    newMd.Name = clearTextValue(scanner.Text())
    if "" == newMd.Name {
        newMd.Name = defMdlName
    }

    fmt.Printf("Add %soptions.php%s file? (y/n):\n", ansiColorGreen, ansiColorReset)
    scanner.Scan()
    newMd.AddOptions = isYes(clearTextValue(scanner.Text()))

    fmt.Printf("Add %sadmin%s folder? (y/n):\n", ansiColorGreen, ansiColorReset)
    scanner.Scan()
    newMd.AddAdminFolder = isYes(clearTextValue(scanner.Text()))

    fmt.Printf("Add %sinstall%s folder(s)? (y/n):\n", ansiColorGreen, ansiColorReset)
    scanner.Scan()
    addIfldrs = clearTextValue(scanner.Text())
    if isYes(addIfldrs) {
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
    if isYes(clearTextValue(scanner.Text())) {
        fmt.Printf("Default lang %s%s%s already added\n", ansiColorBlue, defLang, ansiColorReset)
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
    var addResModfr, addCmpEpilog, addLangs, langName string
    newTp := new(TemplateData)

    fmt.Println("Template will be created from current path '" + curPath + "'")
    fmt.Println("Enter template name. Leave blank to create '" + defTplName + "' template")
    scanner.Scan()
    newTp.Name = clearTextValue(scanner.Text())
    if "" == newTp.Name {
        newTp.Name = defTplName
    }

    fmt.Println("Add result_modifier file? (y/n):")
    scanner.Scan()
    addResModfr = clearTextValue(scanner.Text())
    if isYes(addResModfr) {
        newTp.AddResModfr = true
    }

    fmt.Println("Add component_epilog file? (y/n):")
    scanner.Scan()
    addCmpEpilog = clearTextValue(scanner.Text())
    if isYes(addCmpEpilog) {
        newTp.AddCmpEpilog = true
    }

    fmt.Println("Add language support? (y/n):")
    scanner.Scan()
    addLangs = clearTextValue(scanner.Text())
    if isYes(addLangs) {
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

func isYes(str string) bool {
    var b bool
    if str == "y" || str == "Y" || str == "yes" {
        b = true
    }

    return b
}
