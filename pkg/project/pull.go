package project

import (
	"fmt"
	"strings"

	git "gopkg.in/libgit2/git2go.v23"
)

const (
	branchRef = "refs/heads/"
)

// Pull fetches refspec from remoteName and merges them on top of HEAD.
func (p *Project) Pull(remoteName, refspec string) error {
	// fetch from remote
	if err := p.Fetch(remoteName, refspec); err != nil {
		return err
	}

	// if branch, use remote branch ref
	branchName := ""
	if strings.HasPrefix(refspec, branchRef) {
		branchName = strings.TrimPrefix(refspec, branchRef)
		refspec = fmt.Sprintf("refs/remotes/%s/%s", remoteName, branchName)
	}

	ref, err := p.repo.References.Lookup(refspec)
	if err != nil {
		return err
	}

	// get annotated commit of fetched branch's tip
	aCommit, err := p.repo.AnnotatedCommitFromRef(ref)
	if err != nil {
		return err
	}

	head, err := p.repo.Head()
	// create new branch for head if doesn't exist
	if err != nil && strings.HasSuffix(err.Error(), "not found") {
		branch := "master"
		if len(branchName) != 0 {
			branch = branchName
		}

		commit, err := ref.Peel(git.ObjectCommit)
		if err != nil {
			return err
		}

		if _, err = p.repo.CreateBranch(branch, commit.(*git.Commit), false); err != nil {
			return err
		}

		return p.repo.SetHead(fmt.Sprintf("refs/heads/%s", branch))
	} else if err != nil {
		return err
	}

	path, err := p.TempWorkdir()
	if err != nil {
		return fmt.Errorf("could not setup temporary working directory: %v", err)
	}
	defer p.CleanupWorkdir(path)

	// Perform analysis of merge
	mergeHeads := []*git.AnnotatedCommit{aCommit}
	analysis, _, err := p.repo.MergeAnalysis(mergeHeads)
	if err != nil {
		return err
	}

	switch {
	case analysis&git.MergeAnalysisUpToDate != 0:
		// no changes required
		return nil
	case analysis&git.MergeAnalysisFastForward != 0:
		return p.fastForward(ref, head)
	case analysis&git.MergeAnalysisNormal != 0:
		_, err = p.merge(mergeHeads, ref, head)
		return err
	}

	return fmt.Errorf("merge analysis failed to determine a viable strategy, result: %d", analysis)
}
