package cli

import (
	"encoding/json"

	"github.com/codegangsta/cli"

	"rsprd.com/spread/pkg/data"
	pb "rsprd.com/spread/pkg/spreadproto"
)

// Show displays data stored in Spread Documents.
func (s SpreadCli) Show() *cli.Command {
	return &cli.Command{
		Name:      "show",
		Usage:     "Display information stored in a repository",
		ArgsUsage: "<revision> <path>",
		Action: func(c *cli.Context) {
			var doc *pb.Document
			p := s.projectOrDie()
			switch len(c.Args()) {
			case 1: // from index
				docs, err := p.Index()
				if err != nil {
					s.fatalf("Failed to get index: %v", err)
				}
				path := c.Args().First()
				var ok bool
				if doc, ok = docs[path]; !ok {
					s.fatalf("Path '%s' not found in index.", path)
				}
			case 2: // from revision
				revision := c.Args().First()
				path := c.Args().Get(1)
				var err error
				doc, err = p.GetDocument(revision, path)
				if err != nil {
					s.fatalf("failed to get document: %v", err)
				}
			default:
				s.fatalf("a path OR path and revision must be specified")
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
