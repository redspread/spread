package main // import "rsprd.com/spread/cmd/spread"

import (
	"os"
	"reflect"

	"rsprd.com/spread/cli"

	clilib "rsprd.com/spread/Godeps/_workspace/src/github.com/codegangsta/cli"
)

func main() {
	spread := cli.NewSpreadCli(os.Stdin, os.Stdout, os.Stderr, Version)

	app := app()
	app.Commands = commands(spread)
	app.Run(os.Args)
}

// Get commands provided by spread
func commands(spread *cli.SpreadCli) []clilib.Command {
	cmds := []clilib.Command{}

	cliType := reflect.ValueOf(spread)
	for i := 0; i < cliType.NumMethod(); i++ {
		cmd := cliType.Method(i)

		// check that is returning command
		cmdType := cmd.Type()
		if cmdType.NumOut() == 1 && cmdType.Out(0) == reflect.TypeOf(clilib.Command{}) {
			command := cmd.Interface().(func() clilib.Command)
			cmds = append(cmds, command())
		}
	}
	return cmds
}

func app() *clilib.App {
	app := clilib.NewApp()
	app.Usage = Usage
	app.Version = Version
	return app
}
