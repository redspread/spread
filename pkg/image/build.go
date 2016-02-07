package image

import (
	docker "github.com/fsouza/go-dockerclient"
)

// Build holds the configuration required to build a Docker context
type Build struct {
	ContextPath string
	Config      docker.BuildImageOptions
}
