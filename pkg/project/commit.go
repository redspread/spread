package project

import (
	"fmt"

	git "gopkg.in/libgit2/git2go.v23"

	"rsprd.com/spread/pkg/deploy"
)

func (p *Project) Commit(refname string, author, committer Person, message string) (commitOid string, err error) {
	var parents []*git.Commit
	if head, err := p.headCommit(); err == nil {
		parents = append(parents, head)
	}

	gitAuthor, gitCommitter := git.Signature(author), git.Signature(committer)

	index, err := p.repo.Index()
	if err != nil {
		return "", fmt.Errorf("could not get index: %v", err)
	}

	treeOid, err := index.WriteTree()
	if err != nil {
		return "", fmt.Errorf("could not write index to tree: %v", err)
	}

	commitTree, err := p.repo.LookupTree(treeOid)
	if err != nil {
		return "", fmt.Errorf("could not retrieve created commit tree: %v", err)
	}

	commit, err := p.repo.CreateCommit(refname, &gitAuthor, &gitCommitter, message, commitTree, parents...)
	if err != nil {
		return "", fmt.Errorf("failed to create commit: %v", err)
	}

	return commit.String(), nil
}

func (p *Project) Head() (*deploy.Deployment, error) {
	commit, err := p.headCommit()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head: %v", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("couldn't get tree for HEAD: %v", err)
	}

	return p.deploymentFromTree(tree)
}

func (p *Project) ResolveCommit(revision string) (*deploy.Deployment, error) {
	gitObj, err := p.repo.RevparseSingle(revision)
	if err != nil {
		return nil, fmt.Errorf("couldn't resolve revspec '%s': %v", revision, err)
	}

	if gitObj.Type() != git.ObjectCommit {
		return nil, fmt.Errorf("'%s' specifies an object other than a commit", revision)
	}

	commit, err := gitObj.Peel(git.ObjectCommit)
	if err != nil {
		return nil, err
	}

	tree, err := commit.(*git.Commit).Tree()
	if err != nil {
		return nil, err
	}

	return p.deploymentFromTree(tree)
}

func (p *Project) headCommit() (*git.Commit, error) {
	ref, err := p.repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := ref.Peel(git.ObjectCommit)
	if err != nil {
		return nil, err
	}
	return commit.(*git.Commit), nil
}
