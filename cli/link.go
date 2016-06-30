package cli

import (
	"github.com/codegangsta/cli"

	"rsprd.com/spread/pkg/data"
	"rsprd.com/spread/pkg/deploy"
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
			dep, err := proj.Index()
			if err != nil {
				s.fatalf("Error retrieving index: %v", err)
			}

			kubeObj, err := dep.Get(attach.Path)
			if err != nil {
				s.fatalf("Could not get object: %v", err)
			}

			path, err := deploy.ObjectPath(kubeObj)
			if err != nil {
				s.fatalf("Failed to determine path to save object: %v", err)
			}

			obj, err := data.CreateObject(kubeObj.GetObjectMeta().GetName(), path, kubeObj)
			if err != nil {
				s.fatalf("failed to encode object: %v", err)
			}

			link := data.NewLink("test", target, false)
			if err = data.CreateLinkInObject(obj, link, attach); err != nil {
				s.fatalf("Could not create link: %v", err)
			}

			if err = proj.AddObjectToIndex(obj); err != nil {
				s.fatalf("Failed to add object to Git index: %v", err)
			}
		},
	}
}
