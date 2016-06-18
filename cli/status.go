package cli

import (
	"github.com/codegangsta/cli"

	"rsprd.com/spread/pkg/deploy"
)

// Status returns information about the current state of the project.
func (s SpreadCli) Status() *cli.Command {
	return &cli.Command{
		Name:        "status",
		Usage:       "spread status",
		Description: "Information about what's commited, changed, and staged.",
		Action: func(c *cli.Context) {
			proj := s.project()
			index, err := proj.Index()
			if err != nil {
				s.fatalf("Could not load Index: %v", err)
			}

			client, err := deploy.NewKubeClusterFromContext("")
			if err != nil {
				s.fatalf("Failed to connect to Kubernetes cluster: %v", err)
			}

			cluster, err := client.Deployment()
			if err != nil {
				s.fatalf("Could not load deployment from cluster: %v", err)
			}

			s.printf(index.Diff(cluster))
		},
	}
}
