package cli

import (
	"errors"
	"fmt"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/entity"
	"rsprd.com/spread/pkg/input/dir"

	"github.com/codegangsta/cli"
)

// Deploy allows the creation of deploy.Deployments remotely
func (s SpreadCli) Deploy() *cli.Command {
	return &cli.Command{
		Name:        "deploy",
		Usage:       "spread deploy [-s] PATH | COMMIT [kubectl context]",
		Description: "Deploys objects to a remote Kubernetes cluster.",
		ArgsUsage:   "-s will deploy only if no other deployment found (otherwise fails)",
		Action: func(c *cli.Context) {
			ref := c.Args().First()
			var dep *deploy.Deployment

			if proj, err := s.project(); err == nil {
				if len(ref) == 0 {
					s.printf("Deploying from index...")
					index, err := proj.Index()
					if err != nil {
						s.fatalf("Error reading index: %v", err)
					}
					dep = index
				} else {
					if commit, err := proj.ResolveCommit(ref); err == nil {
						dep = commit
					} else {
						dep = s.fileDeploy(ref)
					}
				}
			} else {
				dep = s.fileDeploy(ref)
			}

			context := c.Args().Get(1)
			cluster, err := deploy.NewKubeClusterFromContext(context)
			if err != nil {
				s.fatalf("Failed to deploy: %v", err)
			}

			s.printf("Deploying %d objects using the %s.", dep.Len(), displayContext(context))

			update := !c.Bool("s")
			err = cluster.Deploy(dep, update, false)
			if err != nil {
				//TODO: make better error messages (one to indicate a deployment already existed; another one if a deployment did not exist but some other error was thrown
				s.fatalf("Did not deploy.: %v", err)
			}

			s.printf("Deployment successful!")
		},
	}
}

func (s SpreadCli) fileDeploy(srcDir string) *deploy.Deployment {
	input, err := dir.NewFileInput(srcDir)
	if err != nil {
		s.fatalf(inputError(srcDir, err))
	}

	e, err := input.Build()
	if err != nil {
		println("build")
		s.fatalf(inputError(srcDir, err))
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
		s.fatalf(inputError(srcDir, err))
	}
	return dep
}

func objectOnlyDeploy(input *dir.FileInput) (*deploy.Deployment, error) {
	objects, err := input.Objects()
	if err != nil {
		return nil, err
	} else if len(objects) == 0 {
		return nil, ErrNothingDeployable
	}

	deployment := new(deploy.Deployment)
	for _, obj := range objects {
		err = deployment.Add(obj)
		if err != nil {
			return nil, err
		}
	}
	return deployment, nil
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

var (
	ErrNothingDeployable = errors.New("there is nothing deployable")
)
