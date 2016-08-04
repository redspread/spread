package cli

import (
	"github.com/codegangsta/cli"

	"rsprd.com/spread/pkg/data"
)

// Link allows the links to be created on the Index
func (s SpreadCli) Link() *cli.Command {
	return &cli.Command{
		Name:        "link",
		Usage:       "spread link <target-url> <attach-point>",
		Description: "Create/remove links on Index",
		Action: func(c *cli.Context) {
			targetUrl := c.Args().First()
			if len(targetUrl) == 0 {
				s.fatalf("A target URL must be specified")
			}

			target, err := data.ParseSRI(targetUrl)
			if err != nil {
				s.fatalf("Error using target: %v", err)
			}

			attachPoint := c.Args().Get(1)
			if len(attachPoint) == 0 {
				s.fatalf("An attach point must be specified")
			}

			attach, err := data.ParseSRI(attachPoint)
			if err != nil {
				s.fatalf("Error using attach-point: %v", err)
			}

			proj := s.projectOrDie()
			index, err := proj.Index()
			if err != nil {
				s.fatalf("Error retrieving index: %v", err)
			}

			doc, ok := index[attach.Path]
			if !ok {
				s.fatalf("Path '%s' not found", attach.Path)
			}

			link := data.NewLink("test", target, false)
			if err = data.CreateLinkInDocument(doc, link, attach); err != nil {
				s.fatalf("Could not create link: %v", err)
			}

			if err = proj.AddDocumentToIndex(doc); err != nil {
				s.fatalf("Failed to add object to Git index: %v", err)
			}
		},
	}
}
