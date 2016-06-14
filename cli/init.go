package cli

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

const (
	// SpreadDirectory is the name of the directory that holds a Spread repository.
	SpreadDirectory = ".spread"
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

			// Create .spread directory in target directory
			if err = os.MkdirAll(target, 0755); os.IsNotExist(err) {
				s.fatalf("Could not initialize repo: '%s' already exists", target)
			} else if err != nil {
				s.fatalf("Could not create repo directory: %v", err)
			}

			// Create bare Git repository in .spread directory with the directory name "git"
			// Create .gitignore file in directory ignoring Git repository
		},
	}
}
