package main

import (
	"github.com/codegangsta/cli"
	"github.com/maleck13/gogen/cmd"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gogen"
	app.Usage = "Generates a new rest service in golang"
	app.Commands = []cli.Command{
		cmd.GenerateCommand(),
	}

	app.Run(os.Args)
}
