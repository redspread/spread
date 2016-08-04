package main // import "rsprd.com/spread/cmd/spread"

import (
	"os"
	"reflect"

	"rsprd.com/spread/cli"

	clilib "github.com/codegangsta/cli"
)

func main() {
	wd, _ := os.Getwd()
	spread := cli.NewSpreadCli(os.Stdin, os.Stdout, os.Stderr, Version, wd)

	// git override
	if len(os.Args) > 2 && os.Args[1] == "git" {
		spread.ExecGitCmd(os.Args[2:]...)
	}

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
		if cmdType.NumOut() == 1 && cmdType.Out(0) == reflect.TypeOf(new(clilib.Command)) {
			cmdFn := cmd.Interface().(func() *clilib.Command)
			command := cmdFn()
			if command != nil {
				cmds = append(cmds, *command)
			}
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
