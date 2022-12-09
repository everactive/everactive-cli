package cmd

import (
	"fmt"
	"github.com/fatih/color"
)

func TUI_Info(message string) {
	color.Cyan(message)
}

func TUI_Error(message string) {
	color.Red(message)
}

func TUI_Warning(message string) {
	color.Yellow(message)
}

func TUI_Message(message string) {
	fmt.Println(message)
}

func TUI_Debug(message string) {
	if DebugEnabled {
		color.White(message)
	}
}
