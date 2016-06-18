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

			head, err := proj.Head()
			if err != nil {
				head = new(deploy.Deployment)
			}

			client, err := deploy.NewKubeClusterFromContext("")
			if err != nil {
				s.fatalf("Failed to connect to Kubernetes cluster: %v", err)
			}

			cluster, err := client.Deployment()
			if err != nil {
				s.fatalf("Could not load deployment from cluster: %v", err)
			}

			stat := deploy.Stat(index, head, cluster)
			s.printStatus(stat)
		},
	}
}

func (s SpreadCli) printStatus(status deploy.DiffStat) {
	s.printf("From Index:")
	if len(status.ClusterNew) > 0 {
		s.printf("%d untracked objects:", len(status.ClusterNew))
		for _, path := range status.ClusterNew {
			s.printf("- %s", path)
		}
	}
	if len(status.ClusterModified) > 0 {
		s.printf("%d modified objects:", len(status.ClusterModified))
		for _, path := range status.ClusterModified {
			s.printf("- %s", path)
		}
	}
	if len(status.ClusterDeleted) > 0 {
		s.printf("%d deleted objects:", len(status.ClusterDeleted))
		for _, path := range status.ClusterDeleted {
			s.printf("- %s", path)
		}
	}

	s.printf("From HEAD:")
	if len(status.IndexNew) > 0 {
		s.printf("%d staged creations:", len(status.IndexNew))
		for _, path := range status.IndexNew {
			s.printf("- %s", path)
		}
	}
	if len(status.IndexModified) > 0 {
		s.printf("%d staged changes:", len(status.IndexModified))
		for _, path := range status.IndexModified {
			s.printf("- %s", path)
		}
	}
	if len(status.IndexDeleted) > 0 {
		s.printf("%d staged deletions:", len(status.IndexDeleted))
		for _, path := range status.IndexDeleted {
			s.printf("- %s", path)
		}
	}
}
