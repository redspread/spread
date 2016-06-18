package cli

import (
	"time"

	"github.com/codegangsta/cli"

	"rsprd.com/spread/pkg/project"
)

// Commit sets up a Spread repository for versioning.
func (s SpreadCli) Commit() *cli.Command {
	return &cli.Command{
		Name:        "commit",
		Usage:       "spread commit -m <msg>",
		Description: "Create new commit based on the current index",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "m",
				Usage: "Message to store the commit with",
			},
		},
		Action: func(c *cli.Context) {
			msg := c.String("m")
			if len(msg) == 0 {
				s.fatalf("All commits must have a message. Specify one with 'spread commit -m \"message\"'")
			}

			proj := s.projectOrDie()
			notImplemented := project.Person{
				Name:  "not implemented",
				Email: "not@implemented.com",
				When:  time.Now(),
			}

			oid, err := proj.Commit("HEAD", notImplemented, notImplemented, msg)
			if err != nil {
				s.fatalf("Could not commit: %v", err)
			}

			s.printf("New commit: [%s] %s", oid, msg)
		},
	}
}
