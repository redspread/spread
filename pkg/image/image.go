package image

import (
	"strings"

	"github.com/docker/docker/reference"
)

// Image contains configuration necessary to deploy an image or if necessary, built it.
type Image struct {
	tag   bool
	image reference.NamedTagged
	Build *Build
}

// KubeImage returns a reference to the Image for use with the Image field of kube.Container.
func (i Image) KubeImage() (out string) {
	imageStr := i.image.String()
	if !i.tag {
		// remove default latest tag
		return strings.TrimRight(imageStr, ":latest")
	}
	return imageStr
}

// Name returns a human readable identifier for the Image. Should be a DNS Label.
func (i Image) Name() string {
	userName := strings.Split(i.image.RemoteName(), "/")
	return userName[1]
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
