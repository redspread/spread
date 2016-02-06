// +build solaris

package logrus

import (
	"os"

	"rsprd.com/spread/Godeps/_workspace/src/github.com/fsouza/go-dockerclient/external/golang.org/x/sys/unix"
)

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal() bool {
	_, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TCGETA)
	return err == nil
}
