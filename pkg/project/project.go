package project

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	git "gopkg.in/libgit2/git2go.v23"
)

const (
	// SpreadDirectory is the name of the directory that holds a Spread repository.
	SpreadDirectory = ".spread"

	// GitDirectory is the name of the directory holding the bare Git repository within the SpreadDirectory.
	GitDirectory = "git"
)

type Project struct {
	Path string
}

// InitProject creates a new Spread project including initializing a Git repository on disk.
// A target must be specified.
func InitProject(target string) (*Project, error) {
	// Check if path is specified
	if len(target) == 0 {
		return nil, errors.New("target must be specified")
	}

	// Get absolute path to directory
	target, err := filepath.Abs(target)
	if err != nil {
		return nil, fmt.Errorf("could not resolve '%s': %v", target, err)
	}

	// Check if directory exists
	if _, err = os.Stat(target); err == nil {
		return nil, fmt.Errorf("'%s' already exists", target)
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// Create .spread directory in target directory
	if err = os.MkdirAll(target, 0755); err != nil {
		return nil, fmt.Errorf("could not create repo directory: %v", err)
	}

	// Create bare Git repository in .spread directory with the directory name "git"
	gitDir := filepath.Join(target, GitDirectory)
	if _, err = git.InitRepository(gitDir, true); err != nil {
		return nil, fmt.Errorf("Could not create Object repository: %v", err)
	}

	// Create .gitignore file in directory ignoring Git repository
	ignoreName := filepath.Join(target, ".gitignore")
	ignoreData := fmt.Sprintf("/%s", GitDirectory)
	ioutil.WriteFile(ignoreName, []byte(ignoreData), 0755)
	return &Project{
		Path: target,
	}, nil
}