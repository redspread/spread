package project

import (
	"os"
	"testing"
)

// returns global location after cleaning up directory
func cleanupGlobal(t *testing.T) string {
	// set global path to test path
	GlobalPath = "~/test-spread-global-path"

	// delete contents
	if path, err := GlobalLocation(); err == nil {
		if err = os.RemoveAll(path); err != nil {
			t.Fatalf("could not remove global location: %v", err)
			return ""
		}
		return path
	} else {
		t.Fatalf("could not resolve global location: %v", err)
		return ""
	}
}

func TestNoInitGlobal(t *testing.T) {
	cleanupGlobal(t)
	_, err := Global()
	if err != ErrNoGlobal {
		t.Error("did not return error about getting global while none exists")
	}
}

func TestInitGlobal(t *testing.T) {
	globalPath := cleanupGlobal(t)
	proj, err := InitGlobal()
	if err != nil {
		t.Errorf("failed to init global: %v", err)
	} else if proj == nil {
		t.Error("the returned project was nil")
	} else if proj.Path != globalPath {
		t.Errorf("the path the created project ('%s') does not match the global path ('%s')", proj.Path, globalPath)
	}
}

func TestGetGlobal(t *testing.T) {
	globalPath := cleanupGlobal(t)
	_, err := InitGlobal()
	if err != nil {
		t.Errorf("failed to init global: %v", err)
	}

	proj, err := Global()
	if err != nil {
		t.Errorf("could not get global: %v", err)
	} else if proj == nil {
		t.Error("the returned project was nil")
	} else if proj.Path != globalPath {
		t.Errorf("the path the created project ('%s') does not match the global path ('%s')", proj.Path, globalPath)
	}
}

func TestDoubleInitGlobal(t *testing.T) {
	cleanupGlobal(t)
	_, err := InitGlobal()
	if err != nil {
		t.Errorf("failed to init global: %v", err)
	}

	_, err = InitGlobal()
	if err == nil {
		t.Error("did not throw error for double initialization of global repo")
	}
}
