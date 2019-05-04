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

func isYes(str string) bool {
    return (str == "y" || str == "Y" || str == "yes")
}

func mkDir(name string) error {
    /*var err error
    if err := os.Mkdir(name, defDirMode); err == nil {
        r = true
    } else {
        fmt.Println("Error creating path '" + name + "': ", err.Error())
    }

    return r*/
}

func mkFileWithContents(name string, contType int) error {
    //var err error

    /*if fh, err := os.Create(name); err == nil {
        r = true
        _ = os.Chmod(name, defFileMode)

        _, err = fh.Write([]byte(tpls[contType]))

        if err != nil {
            fmt.Println("Error writing contents to file '" + name + "': ", err.Error())
        }

    } else {
        fmt.Println("Error creating file '" + name + "': ", err.Error())
    }*/

    // return r
}