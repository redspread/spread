package cli

import (
	"rsprd.com/localkube/pkg/localkubectl"

	"github.com/codegangsta/cli"
)

// Cluster manages the localkube Kubernetes development environment.
func (s SpreadCli) Cluster() *cli.Command {
	return localkubectl.Command(s.out)
}
