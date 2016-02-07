package cli

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// Version returns the current spread version
func (spread SpreadCli) Version() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "prints spread version",
		Action: func(c *cli.Context) {
			fmt.Fprintf(spread.out, "Spread version is %s\n", spread.version)
		},
	}
}
