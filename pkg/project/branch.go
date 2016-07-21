package project

import (
	"fmt"

	git "gopkg.in/libgit2/git2go.v23"

	"rsprd.com/spread/pkg/deploy"
)

func (p *Project) Branch(name string) (*deploy.Deployment, error) {
	br, err := p.repo.LookupBranch(name, git.BranchRemote)
	if err != nil {
		return nil, fmt.Errorf("unable to locate branch: %v", err)
	}

	tree, err := br.Peel(git.ObjectTree)
	if err != nil {
		return nil, err
	}

	return p.deploymentFromTree(tree.(*git.Tree))
}
