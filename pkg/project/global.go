package project

import (
	"github.com/mitchellh/go-homedir"
)

var (
	// GlobalPath is the path that holds the Spread global repository for a user. The '~' character may be used
	// to denote the home directory of a user across platforms.
	GlobalPath = "~/.spread-global"
)

// Global returns the users global project which holds data from any package downloaded.
func Global() (*Project, error) {
	path, err := GlobalLocation()
	if err != nil {
		return nil, err
	}

	return OpenProject(path)
}

// InitGlobal initializes the global repository for this user.
func InitGlobal() (*Project, error) {
	path, err := GlobalLocation()
	if err != nil {
		return nil, err
	}

	return InitProject(path)
}

// GlobalLocation returns the path of the global project. An error is returned if the path doesn't exist.
func GlobalLocation() (string, error) {
	return homedir.Expand(GlobalPath)
}
