package cli

import (
	"github.com/codegangsta/cli"

	"rsprd.com/spread/pkg/deploy"
)

// Diff shows the difference bettwen the cluster and the index.
func (s SpreadCli) Diff() *cli.Command {
	return &cli.Command{
		Name:        "diff",
		Usage:       "spread diff",
		Description: "Diffs index against state of cluster",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "context",
				Value: "",
				Usage: "kubectl context to use for requests",
			},
		},
		Action: func(c *cli.Context) {
			proj := s.projectOrDie()
			index, err := proj.Index()
			if err != nil {
				s.fatalf("Could not load Index: %v", err)
			}

			context := c.String("context")
			client, err := deploy.NewKubeClusterFromContext(context)
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
