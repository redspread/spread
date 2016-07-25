package cli

import (
	"errors"
	"fmt"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/entity"
	"rsprd.com/spread/pkg/input/dir"
	"rsprd.com/spread/pkg/packages"

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

			proj, err := s.project()
			if err == nil {
				if len(ref) == 0 {
					s.printf("Deploying from index...")
					docs, err := proj.Index()
					if err != nil {
						s.fatalf("Error getting index: %v", err)
					}

					dep, err = deploy.DeploymentFromDocMap(docs)
				} else {
					if docs, err := proj.ResolveCommit(ref); err == nil {
						dep, err = deploy.DeploymentFromDocMap(docs)
					} else {
						dep, err = s.globalDeploy(ref)
					}
				}
			} else {
				dep, err = s.globalDeploy(ref)
			}

			if err != nil {
				s.fatalf("Failed to assemble deployment: %v", err)
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

func (s SpreadCli) fileDeploy(srcDir string) (*deploy.Deployment, error) {
	input, err := dir.NewFileInput(srcDir)
	if err != nil {
		return nil, inputError(srcDir, err)
	}

	e, err := input.Build()
	if err != nil {
		return nil, inputError(srcDir, err)
	}

	dep, err := e.Deployment()

	// TODO: This can be removed once application (#56) is implemented
	if err == entity.ErrMissingContainer {
		// check if has pod; if not deploy objects
		pods, err := input.Entities(entity.EntityPod)
		if err != nil && len(pods) != 0 {
			return nil, fmt.Errorf("Failed to deploy: %v", err)
		}

		dep, err = objectOnlyDeploy(input)
		if err != nil {
			return nil, fmt.Errorf("Failed to deploy: %v", err)
		}

	} else if err != nil {
		return nil, inputError(srcDir, err)
	}
	return dep, nil
}

func (s SpreadCli) globalDeploy(ref string) (*deploy.Deployment, error) {
	// check if reference is local file
	dep, err := s.fileDeploy(ref)
	if err != nil {
		ref, err = packages.ExpandPackageName(ref)
		if err == nil {
			var info packages.PackageInfo
			info, err = packages.DiscoverPackage(ref, true, false)
			if err != nil {
				s.fatalf("failed to retrieve package info: %v", err)
			}

			proj, err := s.globalProject()
			if err != nil {
				s.fatalf("error setting up global project: %v", err)
			}

			remote, err := proj.Remotes().Lookup(ref)
			// if does not exist or has different URL, create new remote
			if err != nil {
				remote, err = proj.Remotes().Create(ref, info.RepoURL)
				if err != nil {
					return nil, fmt.Errorf("could not create remote: %v", err)
				}
			} else if remote.Url() != info.RepoURL {
				s.printf("changing remote URL for %s, current: '%s' new: '%s'", ref, remote.Url(), info.RepoURL)
				err = proj.Remotes().SetUrl(ref, info.RepoURL)
				if err != nil {
					return nil, fmt.Errorf("failed to change URL for %s: %v", ref, err)
				}
			}

			s.printf("pulling repo from %s", info.RepoURL)
			branch := fmt.Sprintf("%s/master", ref)
			err = proj.Fetch(remote.Name(), "master")
			if err != nil {
				return nil, fmt.Errorf("failed to fetch '%s': %v", ref, err)
			}

			docs, err := proj.Branch(branch)
			if err != nil {
				return nil, err
			}

			return deploy.DeploymentFromDocMap(docs)
		}
	}
	return dep, err
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

func inputError(srcDir string, err error) error {
	return fmt.Errorf("Error using `%s`: %v", srcDir, err)
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
