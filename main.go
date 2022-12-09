package main

import (
	"gitlab.com/everactive/everactive-cli/cmd"
	"gitlab.com/everactive/everactive-cli/lib"
)

func main() {
	lib.InitConfiguration()
	cmd.Execute()
}


