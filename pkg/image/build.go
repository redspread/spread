package image

import (
	"reflect"

	docker "github.com/fsouza/go-dockerclient"
)

// Build holds the configuration required to build a Docker context
type Build struct {
	ContextPath string
	Config      docker.BuildImageOptions
}

// Equal returns whether the build is the same as the other
func (b Build) Equal(other *Build) bool {
	if other == nil {
		return false
	}
	return reflect.DeepEqual(b, other)
}
