package project

import (
	"fmt"

	git "gopkg.in/libgit2/git2go.v23"
	pb "rsprd.com/spread/pkg/spreadproto"
)

func (p *Project) Branch(name string) (map[string]*pb.Document, error) {
	br, err := p.repo.LookupBranch(name, git.BranchRemote)
	if err != nil {
		return nil, fmt.Errorf("unable to locate branch: %v", err)
	}

	tree, err := br.Peel(git.ObjectTree)
	if err != nil {
		return nil, err
	}

	return p.mapFromTree(tree.(*git.Tree))
}
