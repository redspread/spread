package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	git "gopkg.in/libgit2/git2go.v23"
)

const (
	testDir = "../../test/repos/"
)

func init() {
	os.RemoveAll(testDir)
}

func TestInitProjectNoPath(t *testing.T) {
	target := ""
	proj, err := InitProject(target)
	assert.Error(t, err, "projects require paths")
	assert.Nil(t, proj, "invalid projects should not be returned")
}

func TestInitProjectRelative(t *testing.T) {
	target := filepath.Join(testDir, "relativeTest")
	defer os.RemoveAll(target)

	proj, err := InitProject(target)
	assert.NoError(t, err)
	assert.NotNil(t, proj)

	checkProject(t, proj.Path)
}

func TestInitProjectAbsolute(t *testing.T) {
	target := filepath.Join(testDir, "absoluteTest")
	target, err := filepath.Abs(target)
	assert.NoError(t, err)

	defer os.RemoveAll(target)

	proj, err := InitProject(target)
	assert.NoError(t, err)
	assert.NotNil(t, proj)

	checkProject(t, proj.Path)
}

func checkProject(t *testing.T, target string) {
	// check that project directory was created
	checkPath(t, target, true)

	// check .gitignore
	ignorePath := filepath.Join(target, ".gitignore")
	checkPath(t, ignorePath, false)

	// check git directory
	gitDir := filepath.Join(target, GitDirectory)
	checkPath(t, gitDir, true)
	repo, err := git.OpenRepository(gitDir)
	assert.NoError(t, err)

	assert.True(t, filepath.Clean(repo.Path()) == gitDir)
}

// checkPath determines if a path has been created.
// The test errors if it doesn't exist or doesn't match dir.
func checkPath(t *testing.T, target string, dir bool) {
	fileInfo, err := os.Stat(target)
	if assert.NoError(t, err) {
		assert.True(t, fileInfo.IsDir() == dir, "should have created directory at %s", target)
	}
}
