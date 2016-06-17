package project

import (
	"fmt"

	git "gopkg.in/libgit2/git2go.v23"
)

func (p *Project) Commit(refname string, author, committer Person, message string) (commitOid string, err error) {
	var parents []*git.Commit
	if head := p.headCommit(); head != nil {
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

func (p *Project) headCommit() *git.Commit {
	ref, err := p.repo.Head()
	if err != nil {
		return nil
	}

	commit, err := ref.Peel(git.ObjectCommit)
	if err != nil {
		return nil
	}
	return commit.(*git.Commit)
}
