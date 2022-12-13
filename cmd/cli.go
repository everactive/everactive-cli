package cmd

import (
	"fmt"
	"github.com/fatih/color"
)

func Tui_info(message string) {
	color.Cyan(message)
}

func Tui_error(message string) {
	color.Red(message)
}

func Tui_warning(message string) {
	color.Yellow(message)
}

func Tui_message(message string) {
	fmt.Println(message)
}

func Tui_debug(message string) {
	if DebugEnabled {
		color.White(message)
	}
}

func Tui_bell() {
	fmt.Print("\a")
}