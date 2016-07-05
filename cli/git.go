package cli

import (
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/codegangsta/cli"
)

func (s SpreadCli) Git() *cli.Command {
	return &cli.Command{
		Name:  "git",
		Usage: "Allows access to git commands while Spread is build out",
		Action: func(c *cli.Context) {
			cli.ShowSubcommandHelp(c)
		},
	}
}

func (s SpreadCli) ExecGitCmd(args ...string) {
	git, err := exec.LookPath("git")
	if err != nil {
		s.fatalf("Could not locate git: %v", err)
	}

	proj := s.projectOrDie()
	gitDir := filepath.Join(proj.Path, "git")

	gitArgs := []string{git, "--git-dir=" + gitDir}
	gitArgs = append(gitArgs, args...)

	err = syscall.Exec(git, gitArgs, []string{})
	if err != nil {
		s.fatalf("could not run git: %v", err)
	}
}
