package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"rsprd.com/spread/pkg/project"
)

// SpreadCli is the spread command line client.
type SpreadCli struct {
	// input stream (ie. stdin)
	in io.ReadCloser

	// output stream (ie. stdout)
	out io.Writer

	// error stream (ie. stderr)
	err io.Writer

	version string
	workDir string
}

// NewSpreadCli returns a spread command line interface (CLI) client.NewSpreadCli.
// All functionality accessible from the command line is attached to this struct.
func NewSpreadCli(in io.ReadCloser, out, err io.Writer, version, workDir string) *SpreadCli {
	return &SpreadCli{
		in:      in,
		out:     out,
		err:     err,
		version: version,
		workDir: workDir,
	}
}

func (c SpreadCli) project() (*project.Project, error) {
	if len(c.workDir) == 0 {
		return nil, ErrNoWorkDir
	}

	root, found := findPath(c.workDir, project.SpreadDirectory, true)
	if !found {
		return nil, fmt.Errorf("could not find spread directory from %s", c.workDir)
	}

	return project.OpenProject(root)
}

func (c SpreadCli) printf(message string, data ...interface{}) {
	// add newline if doesn't have one
	if !strings.HasSuffix(message, "\n") {
		message = message + "\n"
	}
	fmt.Fprintf(c.out, message, data...)
}

func (c SpreadCli) fatalf(message string, data ...interface{}) {
	c.printf(message, data...)
	os.Exit(1)
}

func findPath(leafDir, targetFile string, dir bool) (string, bool) {
	if len(leafDir) == 0 {
		return "", false
	}
	spread := filepath.Join(leafDir, targetFile)
	if exists, err := pathExists(spread, dir); err == nil && exists {
		return spread, true
	} else {
		if leafDir == "/" {
			return "", false
		}
		parent := filepath.Dir(leafDir)
		return findPath(parent, targetFile, dir)
	}
}

func pathExists(path string, dir bool) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir() == dir, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

var (
	// ErrNoWorkDir is returned when the CLI was started without a working directory set.
	ErrNoWorkDir = errors.New("no working directory was set")
)
