package project

import (
	"io/ioutil"
	"os"

	git "gopkg.in/libgit2/git2go.v23"
)

// TempWorkdir creates a temporary directory and configures the Repository to use it as a work dir.
// HEAD is checked out to the temporary working directory. The path of the working directory is returned as a string.
func (p *Project) TempWorkdir() (string, error) {
	name, err := ioutil.TempDir("", "spread-workdir")
	if err != nil {
		return "", err
	}

	opts := &git.CheckoutOpts{
		TargetDirectory: name,
	}

	if err = p.repo.CheckoutHead(opts); err != nil {
		os.RemoveAll(name)
		return "", err
	}

	return name, p.repo.SetWorkdir(name, false)
}

// CleanupWorkdir removes the given directory and sets the repositories Workdir to an empty string.
func (p *Project) CleanupWorkdir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	return p.repo.SetWorkdir("", false)
}
