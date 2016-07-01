package project

import (
	git "gopkg.in/libgit2/git2go.v23"
)

func (p *Project) Remotes() *git.RemoteCollection {
	return &p.repo.Remotes
}
