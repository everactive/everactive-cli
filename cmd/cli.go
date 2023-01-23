package cmd

import (
	"fmt"
	"github.com/fatih/color"
)

func Tui_info(message string) {
	color.Set(color.FgCyan)
	fmt.Println(message)
	color.Unset()
}

func Tui_error(message string) {
	color.Set(color.FgRed)
	fmt.Println(message)
	color.Unset()
}

func Tui_warning(message string) {
	color.Set(color.FgYellow)
	fmt.Println(message)
	color.Unset()
}

func Tui_message(message string) {
	fmt.Println(message)
}

func Tui_debug(message string) {
	if DebugEnabled {
		fmt.Println(message)
	}
}

func Tui_bell() {
	fmt.Print("\a")
}