package project

import (
	"fmt"

	git "gopkg.in/libgit2/git2go.v23"

	"rsprd.com/spread/pkg/deploy"
)

func (p *Project) deploymentFromTree(tree *git.Tree) (*deploy.Deployment, error) {
	deployment := new(deploy.Deployment)
	var walkErr error
	err := tree.Walk(func(path string, entry *git.TreeEntry) int {
		// add objects to deployment
		if entry.Type == git.ObjectBlob {
			kubeObj, err := p.getKubeObject(entry.Id, path)
			if err != nil {
				walkErr = err
				return -1
			}

			walkErr = deployment.Add(kubeObj)
			if walkErr != nil {
				return -1
			}
		}
		return 0
	})

	if err != nil {
		return nil, fmt.Errorf("error starting walk: %v", err)
	}

	if walkErr != nil {
		return nil, fmt.Errorf("error during walk: %v", err)
	}

	return deployment, nil
}
