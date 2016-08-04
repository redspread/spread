package project

import (
	"fmt"

	git "gopkg.in/libgit2/git2go.v23"

	pb "rsprd.com/spread/pkg/spreadproto"
)

func (p *Project) Commit(refname string, author, committer Person, message string) (commitOid string, err error) {
	var parents []*git.Commit
	if head, err := p.headCommit(); err == nil {
		parents = append(parents, head)
	}

	gitAuthor, gitCommitter := git.Signature(author), git.Signature(committer)

	commitTree, err := p.writeIndex()
	if err != nil {
		return "", err
	}

	commit, err := p.repo.CreateCommit(refname, &gitAuthor, &gitCommitter, message, commitTree, parents...)
	if err != nil {
		return "", fmt.Errorf("failed to create commit: %v", err)
	}

	return commit.String(), nil
}

func (p *Project) writeIndex() (*git.Tree, error) {
	index, err := p.repo.Index()
	if err != nil {
		return nil, fmt.Errorf("could not get index: %v", err)
	}

	treeOid, err := index.WriteTree()
	if err != nil {
		return nil, fmt.Errorf("could not write index to tree: %v", err)
	}

	tree, err := p.repo.LookupTree(treeOid)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve created commit tree: %v", err)
	}
	return tree, nil
}

func (p *Project) Head() (map[string]*pb.Document, error) {
	commit, err := p.headCommit()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head: %v", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("couldn't get tree for HEAD: %v", err)
	}

	return p.mapFromTree(tree)
}

func (p *Project) ResolveCommit(revision string) (map[string]*pb.Document, error) {
	gitObj, err := p.repo.RevparseSingle(revision)
	if err != nil {
		return nil, fmt.Errorf("couldn't resolve revspec '%s': %v", revision, err)
	}

	if gitObj.Type() != git.ObjectCommit {
		return nil, fmt.Errorf("'%s' specifies an object other than a commit", revision)
	}

	tree, err := gitObj.Peel(git.ObjectTree)
	if err != nil {
		return nil, err
	}

	return p.mapFromTree(tree.(*git.Tree))
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
