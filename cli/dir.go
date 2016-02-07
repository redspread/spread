package cli

import (
	"fmt"
	"os"

	"rsprd.com/spread/pkg/input/dir"

	"github.com/codegangsta/cli"

	"rsprd.com/spread/pkg/entity"
)

// Version returns the current spread version
func (spread SpreadCli) Dir() *cli.Command {
	return &cli.Command{
		Name:  "file",
		Usage: "File info",
		Action: func(c *cli.Context) {
			wd, _ := os.Getwd()
			fmt.Fprintf(spread.out, "Current directory is %s\n", wd)
			fs, err := dir.NewFileSource(c.Args().First())
			if err != nil {
				fmt.Errorf("FSError: %v", err)
			}
			obj, err := fs.Objects()
			if err != nil {
				fmt.Errorf("Objects Error: %v", err)
			}
			for _, v := range obj {
				fmt.Printf("I see %s\n\n", v.GetObjectMeta().GetName())
			}

			entities, err := fs.Entities(entity.EntityReplicationController)
			if err != nil {
				fmt.Errorf("RC Error: %v", err)
			}
			for _, v := range entities {
				println(v.Source())
			}

		},
	}
}
