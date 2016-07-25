package cli

import (
	"encoding/json"

	"github.com/codegangsta/cli"

	"rsprd.com/spread/pkg/data"
)

// Show displays data stored in Spread Documents.
func (s SpreadCli) Show() *cli.Command {
	return &cli.Command{
		Name:      "show",
		Usage:     "Display information stored in a repository",
		ArgsUsage: "<revision> <path>",
		Action: func(c *cli.Context) {
			if len(c.Args()) < 2 {
				s.fatalf("a revison and a path must be specified")
			}
			p := s.projectOrDie()

			revision := c.Args().First()
			path := c.Args().Get(1)
			doc, err := p.GetDocument(revision, path)
			if err != nil {
				s.fatalf("failed to get document: %v", err)
			}

			fields, err := data.MapFromDocument(doc)
			if err != nil {
				s.fatalf("could not get fields: %v", err)
			}

			out, err := json.MarshalIndent(&fields, "", "\t")
			if err != nil {
				s.fatalf("Couldn't create JSON: %v", err)
			}

			s.printf("%s", out)
		},
	}
}
