package image

import (
	docker "rsprd.com/spread/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
)

// Build holds the configuration required to build a Docker context
type Build struct {
	ContextPath string
	Config      docker.BuildImageOptions
}
