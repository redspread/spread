package cli

import (
	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/entity"
	"rsprd.com/spread/pkg/input/dir"

	"github.com/codegangsta/cli"
)

// Build is used for rapid iteration with Kubernetes
func (s SpreadCli) Build() *cli.Command {
	return &cli.Command{
		Name:        "build",
		Usage:       "spread build PATH [kubectl context]",
		Description: "Immediately deploy objects to a remote Kubernetes cluster. In the future will also build a Dockerfile and push the resulting image",
		Action: func(c *cli.Context) {
			srcDir := c.Args().First()
			if len(srcDir) == 0 {
				s.fatalf("A directory to deploy from must be specified")
			}

			input, err := dir.NewFileInput(srcDir)
			if err != nil {
				s.fatalf(inputError(srcDir, err).Error())
			}

			e, err := input.Build()
			if err != nil {
				s.fatalf(inputError(srcDir, err).Error())
			}

			dep, err := e.Deployment()

			// TODO: This can be removed once application (#56) is implemented
			if err == entity.ErrMissingContainer {
				// check if has pod; if not deploy objects
				pods, err := input.Entities(entity.EntityPod)
				if err != nil && len(pods) != 0 {
					s.fatalf("Failed to deploy: %v", err)
				}

				dep, err = objectOnlyDeploy(input)
				if err != nil {
					s.fatalf("Failed to deploy: %v", err)
				}

			} else if err != nil {
				println("deploy")
				s.fatalf(inputError(srcDir, err).Error())
			}

			context := c.Args().Get(1)
			cluster, err := deploy.NewKubeClusterFromContext(context)
			if err != nil {
				s.fatalf("Failed to deploy: %v", err)
			}

			s.printf("Updating %d objects using the %s.", dep.Len(), displayContext(context))

			err = cluster.Deploy(dep, true, true)
			if err != nil {
				//TODO: make better error messages (one to indicate a deployment already existed; another one if a deployment did not exist but some other error was thrown
				s.fatalf("Did not deploy.: %v", err)
			}

			s.printf("Build successful!")
		},
	}
}
