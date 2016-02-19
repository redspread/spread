package cli

import (
	"fmt"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/entity"
	"rsprd.com/spread/pkg/input/dir"

	"github.com/codegangsta/cli"
)

// Version returns the current spread version
func (s SpreadCli) Deploy() *cli.Command {
	return &cli.Command{
		Name:      "deploy",
		Usage:     "spread deploy [-u] PATH [kubectl context]",
		ArgsUsage: "-u attempts to update objects (otherwise fails)",
		Action: func(c *cli.Context) {
			srcDir := c.Args().First()
			if len(srcDir) == 0 {
				s.fatal("A directory to deploy from must be specified")
			}

			input, err := dir.NewFileInput(srcDir)
			if err != nil {
				s.fatal(inputError(srcDir, err))
			}

			e, err := input.Build()
			if err != nil {
				println("build")
				s.fatal(inputError(srcDir, err))
			}

			dep, err := e.Deployment()
			if err == entity.ErrMissingContainer {
				// check if has pod; if not deploy objects
				pods, _ := input.Entities(entity.EntityPod)
				if len(pods) != 0 {
					s.fatal("Failed to deploy: %v", err)
				}

				objects, err := input.Objects()
				if err != nil {
					s.fatal(inputError(srcDir, err))
				} else if len(objects) == 0 {
					s.fatal("Couldn't find objects to deploy in '%s'", srcDir)
				}

				dep = new(deploy.Deployment)
				for _, obj := range objects {
					err = dep.Add(obj)
					if err != nil {
						s.fatal(inputError(srcDir, err))
					}
				}
			} else if err != nil {
				println("deploy")
				s.fatal(inputError(srcDir, err))
			}

			context := c.Args().Get(1)
			cluster, err := deploy.NewKubeClusterFromContext(context)
			if err != nil {
				s.fatal("Failed to deploy: %v", err)
			}

			s.printf("Deploying %d objects using the %s.", dep.Len(), displayContext(context))

			update := c.Bool("u")
			err = cluster.Deploy(dep, update)
			if err != nil {
				s.fatal("Failed to deploy: %v", err)
			}

			s.printf("Deployment successful!")
		},
	}
}

func inputError(srcDir string, err error) string {
	return fmt.Sprintf("Error using `%s`: %v", srcDir, err)
}

func displayContext(name string) string {
	if name == deploy.DefaultContext {
		return "default context"
	}
	return fmt.Sprintf("context '%s'", name)
}
