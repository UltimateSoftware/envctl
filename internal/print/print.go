package print

import (
	"fmt"

	"github.com/fatih/color"
)

var msgOK = color.GreenString("OK")
var msgWarn = color.YellowString("WARN")
var msgError = color.RedString("ERROR")
var tpl = "[ %v ]\n"

func OK() {
	fmt.Printf(tpl, msgOK)
}

func Warn() {
	fmt.Printf(tpl, msgWarn)
}

func Error() {
	fmt.Printf(tpl, msgError)
}
