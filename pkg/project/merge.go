package project

import (
	"errors"
	"time"

	git "gopkg.in/libgit2/git2go.v23"
)

const (
	mergeName  = "Merge Author"
	mergeEmail = "merge@redspread.com"
)

func (p *Project) merge(aCommits []*git.AnnotatedCommit, source, target *git.Reference) (*git.Oid, error) {
	defer p.repo.StateCleanup()

	if err := p.repo.Merge(aCommits, nil, nil); err != nil {
		return nil, err
	}

	if index, err := p.repo.Index(); err != nil {
		return nil, err
	} else if index.HasConflicts() {
		return nil, errors.New("Conflicts encountered during merge. Resolve them in the index.")
	}

	// write index to disk and get ID
	commitTree, err := p.writeIndex()
	if err != nil {
		return nil, err
	}

	lCommit, err := p.repo.LookupCommit(target.Target())
	if err != nil {
		return nil, err
	}

	rCommit, err := p.repo.LookupCommit(source.Target())
	if err != nil {
		return nil, err
	}

	return p.repo.CreateCommit("HEAD", mergeSignature(), mergeSignature(), "Auto-merged changes", commitTree, lCommit, rCommit)
}

func (p *Project) fastForward(source, target *git.Reference) error {
	_, err := target.SetTarget(source.Target(), "Fast-forward")
	return err
}

func mergeSignature() *git.Signature {
	return &git.Signature{
		Name:  mergeName,
		Email: mergeEmail,
		When:  time.Now(),
	}
}
