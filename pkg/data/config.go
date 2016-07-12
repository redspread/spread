package data

import (
	"io"
	"os"
)

// Out is the Writer that debugging information is written to.
var Out io.Writer

func init() {
	if Out == nil {
		Out = os.Stdout
	}
}
