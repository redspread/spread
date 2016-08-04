package cli

import (
	"rsprd.com/spread/pkg/project"

	"github.com/codegangsta/cli"
)

// Init sets up a Spread repository for versioning.
func (s SpreadCli) Init() *cli.Command {
	return &cli.Command{
		Name:        "init",
		Usage:       "spread init <path>",
		Description: "Create a new spread repository in the given directory. If none is given, the working directory will be used.",
		Action: func(c *cli.Context) {
			target := c.Args().First()
			if len(target) == 0 {
				target = project.SpreadDirectory
			}

			proj, err := project.InitProject(target)
			if err != nil {
				s.fatalf("Could not create Spread project: %v", err)
			}

			s.printf("Created Spread repository in %s.", proj.Path)
		},
	}
}
