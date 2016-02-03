package image

import (
	"strings"

	"github.com/docker/docker/reference"
)

// Image contains configuration necessary to deploy an image or if necessary, built it.
type Image struct {
	registry  bool
	namespace bool
	tag       bool
	image     reference.NamedTagged
	Build     *Build
}

// KubeImage returns a reference to the Image for use with the Image field of kube.Container.
func (i Image) KubeImage() (out string) {
	if i.image == nil {
		return ""
	}
	if i.registry {
		out += i.image.Hostname() + "/"
	}
	out += i.image.Name()
	if i.tag {
		out += ":" + i.image.Tag()
	}
	return
}

// Name returns a human readable identifier for the Image. Should be a DNS Label.
func (i Image) Name() string {
	if i.image == nil {
		return ""
	}
	if i.namespace {
		imageName := strings.Split(i.image.Name(), "/")
		return imageName[1]
	}
	return i.image.Name()
}

// FromString creates an Image using a string representation
func FromString(str string) (*Image, error) {
	named, err := reference.ParseNamed(str)
	if err != nil {
		return nil, err
	}

	// check for namespaces and tags
	hasTag := !reference.IsNameOnly(named)
	hasNamespace := strings.Contains(named.Name(), "/")

	// implicit default tag
	tagged := reference.WithDefaultTag(named)

	return &Image{
		image:     tagged.(reference.NamedTagged),
		namespace: hasNamespace,
		tag:       hasTag,
	}, nil
}
