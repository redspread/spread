package cli

import (
	"github.com/codegangsta/cli"
)

// Push allows references to be pushed to a remote.
func (s SpreadCli) Push() *cli.Command {
	return &cli.Command{
		Name:        "push",
		Usage:       "Push references to a remote",
		ArgsUsage:   "<remote> <refspec>",
		Description: "Push Spread data to a remote",
		Action: func(c *cli.Context) {
			remoteName := c.Args().First()
			if len(remoteName) == 0 {
				s.fatalf("a remote must be specified")
			}

			if len(c.Args()) < 2 {
				s.fatalf("a refspec must be specified")
			}
			refspec := c.Args()[1:]

			p := s.projectOrDie()
			err := p.Push(remoteName, refspec...)
			if err != nil {
				s.fatalf("Failed to push: %v", err)
			}
		},
	}
}
