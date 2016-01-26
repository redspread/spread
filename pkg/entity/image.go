package entity

import (
	"errors"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	kube "k8s.io/kubernetes/pkg/api"
)

// Image represents a Docker image in the Redspread hierarchy. It wraps image.Image.
type Image struct {
	base
	image *image.Image
}

func NewImage(image *image.Image, defaults kube.ObjectMeta, source string, objects ...deploy.KubeObject) (*Image, error) {
	if image == nil {
		return nil, ErrorNilImage
	}

	base, err := newBase(EntityImage, defaults, source, objects)
	if err != nil {
		return nil, err
	} else if len(image.DockerName()) == 0 {
		return nil, ErrorEmptyImageString
	}

	return &Image{base: base, image: image}, nil
}

func (c Image) Deployment() (*deploy.Deployment, error) {
	return nil, nil
}

func (c Image) Images() []*image.Image {
	return []*image.Image{
		c.image,
	}
}

func (c Image) Attach(e Entity) error {
	return ErrorCannotAttachToImage
}

// Kubernetes representation of image
func (c Image) kube() string {
	return c.image.DockerName()
}

var (
	ErrorEmptyImageString    = errors.New("image.Image's DockerString was empty")
	ErrorNilImage            = errors.New("*image.Image cannot be nil")
	ErrorCannotAttachToImage = errors.New("No entities are allowed to attach to an Image")
)
