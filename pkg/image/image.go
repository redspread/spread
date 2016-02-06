package image

import (
	"io"
	"strings"

	"github.com/docker/docker/reference"
	docker "rsprd.com/spread/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
)

// Image contains configuration necessary to deploy an image or if necessary, built it.
type Image struct {
	tag   bool
	image reference.NamedTagged
	Build *Build
}

// KubeImage returns a reference to the Image for use with the Image field of kube.Container.
func (i Image) KubeImage() (out string) {
	if i.image == nil {
		return ""
	}
	imageStr := i.image.String()
	if !i.tag {
		// remove default latest tag
		return strings.TrimSuffix(imageStr, ":latest")
	}
	return imageStr
}

// Name returns a human readable identifier for the Image. Should be a DNS Label.
func (i Image) Name() string {
	if i.image == nil {
		return ""
	}
	userName := strings.Split(i.image.RemoteName(), "/")
	return userName[1]
}

// PushOptions returns the parameters needed to push an image.
func (i Image) PushOptions(out io.Writer, json bool) docker.PushImageOptions {
	if i.image == nil {
		return docker.PushImageOptions{}
	}
	name := strings.TrimPrefix(i.image.FullName(), "docker.io/")
	opts := docker.PushImageOptions{
		Name:          name,
		Registry:      i.image.Hostname(),
		OutputStream:  out,
		RawJSONStream: json,
	}

	if i.tag {
		opts.Tag = i.image.Tag()
	}
	return opts
}

// FromString creates an Image using a string representation
func FromString(str string) (*Image, error) {
	named, err := reference.ParseNamed(str)
	if err != nil {
		return nil, err
	}

	// check for tag
	hasTag := !reference.IsNameOnly(named)

	// implicit default tag
	tagged := reference.WithDefaultTag(named)

	return &Image{
		image: tagged.(reference.NamedTagged),
		tag:   hasTag,
	}, nil
}
