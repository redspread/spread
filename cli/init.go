package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	git "gopkg.in/libgit2/git2go.v23"
)

// Init sets up a Spread repository for versioning.
func (s SpreadCli) Init() *cli.Command {
	return &cli.Command{
		Name:        "init",
		Usage:       "spread init <path>",
		Description: "Create a new spread repository in the given directory. If none is given, the working directory will be used.",
		Action: func(c *cli.Context) {
			target := SpreadDirectory
			// Check if path is specified
			if len(c.Args().First()) != 0 {
				target = c.Args().First()
			}

			// Get absolute path to directory
			target, err := filepath.Abs(target)
			if err != nil {
				s.fatalf("Could not resolve '%s': %v", target, err)
			}

			// Check if directory exists
			if _, err = os.Stat(target); err == nil {
				s.fatalf("Could not initialize repo: '%s' already exists", target)
			} else if !os.IsNotExist(err) {
				s.fatalf("Could not stat repo directory: %v", err)
			}

			// Create .spread directory in target directory
			if err = os.MkdirAll(target, 0755); err != nil {
				s.fatalf("Could not create repo directory: %v", err)
			}

			// Create bare Git repository in .spread directory with the directory name "git"
			gitDir := filepath.Join(target, GitDirectory)
			if _, err = git.InitRepository(gitDir, true); err != nil {
				s.fatalf("Could not create Object repository: %v", err)
			}

			// Create .gitignore file in directory ignoring Git repository
			ignoreName := filepath.Join(SpreadDirectory, ".gitignore")
			ignoreData := fmt.Sprintf("/%s", GitDirectory)
			ioutil.WriteFile(ignoreName, []byte(ignoreData), 0755)
		},
	}
}
