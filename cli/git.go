package cli

import (
	"os/exec"
	"path/filepath"

	"github.com/codegangsta/cli"
)

func (s SpreadCli) Git() *cli.Command {
	return &cli.Command{
		Name:  "git",
		Usage: "Allows access to git commands while Spread is build out",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "context",
				Value: "",
				Usage: "kubectl context to use for requests",
			},
		},
		Action: func(c *cli.Context) {
			proj := s.project()
			gitDir := filepath.Join(proj.Path, "git")

			gitArgs := []string{"--git-dir=" + gitDir}
			gitArgs = append(gitArgs, c.Args()...)

			cmd := exec.Command("git", gitArgs...)
			cmd.Stdin = s.in
			cmd.Stdout = s.out
			cmd.Stderr = s.err
			err := cmd.Run()
			if err != nil {
				s.fatalf("could not run git: %v", err)
			}
		},
	}
}
