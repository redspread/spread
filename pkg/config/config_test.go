package config

import (
	"os"
	"testing"
)

func TestSetupOut(t *testing.T) {
	if Out != os.Stdout {
		t.Error("Out should have been set to use STDOUT")
	}
}
