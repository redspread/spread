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
	repo *git.Repository
}

// InitProject creates a new Spread project including initializing a Git repository on disk.
// A target must be specified.
func InitProject(target string) (*Project, error) {
	// Check if path is specified
	if len(target) == 0 {
		return nil, ErrEmptyPath
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
	repo, err := git.InitRepository(gitDir, true)
	if err != nil {
		return nil, fmt.Errorf("Could not create Object repository: %v", err)
	}

	// Create .gitignore file in directory ignoring Git repository
	ignoreName := filepath.Join(target, ".gitignore")
	ignoreData := fmt.Sprintf("/%s", GitDirectory)
	ioutil.WriteFile(ignoreName, []byte(ignoreData), 0755)
	return &Project{
		Path: target,
		repo: repo,
	}, nil
}

// OpenProject attempts to open the project at the given path.
func OpenProject(target string) (*Project, error) {
	// Check if path is specified
	if len(target) == 0 {
		return nil, ErrEmptyPath
	}

	// check that path exists and is dir
	if fileInfo, err := os.Stat(target); err != nil {
		return nil, err
	} else if !fileInfo.IsDir() {
		return nil, ErrPathNotDir
	}

	gitDir := filepath.Join(target, GitDirectory)
	repo, err := git.OpenRepository(gitDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open Git repository: %v", err)
	}

	return &Project{
		Path: target,
		repo: repo,
	}, nil
}

var (
	// ErrEmptyPath is returned when a target string is empty.
	ErrEmptyPath = errors.New("path must be specified")

	// ErrPathNotDir is returned when a target is a file and is expected to be a directory.
	ErrPathNotDir = errors.New("a directory must be specified")
)
