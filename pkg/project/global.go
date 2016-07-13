package project

import (
	"errors"

	"github.com/mitchellh/go-homedir"
)

var (
	// GlobalPath is the path that holds the Spread global repository for a user. The '~' character may be used
	// to denote the home directory of a user across platforms.
	GlobalPath = "~/.spread-global"
)

// Global returns the users global project which holds data from any package downloaded.
func Global() (*Project, error) {
	return nil, nil
}

// InitGlobal initializes the global repository for this user.
func InitGlobal() (*Project, error) {
	return nil, nil
}

func GlobalLocation() (string, error) {
	return homedir.Expand(GlobalPath)
}

var (
	// ErrNoGlobal is returned when a global project does not exist for this user.
	ErrNoGlobal = errors.New("global project does not exist")
)
