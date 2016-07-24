package cli

import (
	"github.com/codegangsta/cli"
)

// Remote manages the Git repositories remotes.
func (s SpreadCli) Remote() *cli.Command {
	return &cli.Command{
		Name:        "remote",
		Usage:       "View/modify versioning remotes",
		Description: "Manages repository remotes.",
		Action: func(c *cli.Context) {
			p := s.projectOrDie()
			remotes, err := p.Remotes().List()
			if err != nil {
				s.fatalf("couldn't retrieve remotes")
			}

			s.printf("Remotes:")
			for _, r := range remotes {
				s.printf("- %s", r)
			}

			println()
			cli.ShowSubcommandHelp(c)
		},
		HideHelp: true,
		Subcommands: []cli.Command{
			{
				Name:      "add",
				Usage:     "add new remote",
				ArgsUsage: "<name> <url>",
				Action: func(c *cli.Context) {
					name := c.Args().First()
					if len(name) == 0 {
						s.fatalf("a name must be specified")
					}

					url := c.Args().Get(1)
					if len(url) == 0 {
						s.fatalf("a url must be specified")
					}

					p := s.projectOrDie()
					_, err := p.Remotes().Create(name, url)
					if err != nil {
						s.fatalf("Could not create new remote: %v", err)
					}

					s.printf("Created remote '%s'", name)
				},
			},
			{
				Name:      "remove",
				Usage:     "delete remote",
				ArgsUsage: "<name>",
				Action: func(c *cli.Context) {
					name := c.Args().First()
					if len(name) == 0 {
						s.fatalf("a name must be specified")
					}

					p := s.projectOrDie()
					if err := p.Remotes().Delete(name); err != nil {
						s.fatalf("Could not delete remote: %v", err)
					}

					s.printf("Removed remote '%s'", name)
				},
			},
			{
				Name:      "set-url",
				Usage:     "change remotes url",
				ArgsUsage: "<name> <url>",
				Action: func(c *cli.Context) {
					name := c.Args().First()
					if len(name) == 0 {
						s.fatalf("a name must be specified")
					}

					url := c.Args().Get(1)
					if len(url) == 0 {
						s.fatalf("a url must be specified")
					}

					p := s.projectOrDie()
					if err := p.Remotes().SetUrl(name, url); err != nil {
						s.fatalf("Could not change url for remote: %v", err)
					}

					s.printf("Set url for remote %s to '%s'", name, url)
				},
			},
		},
	}
}
