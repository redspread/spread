package cli

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

// Pull allows references to be pulled from a remote.
func (s SpreadCli) Pull() *cli.Command {
	return &cli.Command{
		Name:        "pull",
		Usage:       "Pull changes from a remote branch",
		ArgsUsage:   "<remote> <branch>",
		Description: "Pull Spread data from a remote branch",
		Action: func(c *cli.Context) {
			remoteName := c.Args().First()
			if len(remoteName) == 0 {
				s.fatalf("a remote must be specified")
			}

			if len(c.Args()) < 2 {
				s.fatalf("a refspec must be specified")
			}

			refspec := c.Args().Get(1)
			if !strings.HasPrefix(refspec, "refs/") {
				refspec = fmt.Sprintf("refs/heads/%s", refspec)
			}

			p := s.projectOrDie()
			if err := p.Pull(remoteName, refspec); err != nil {
				s.fatalf("Failed to pull: %v", err)
			}
		},
	}
}
