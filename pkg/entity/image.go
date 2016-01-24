package entity

import (
	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	"k8s.io/kubernetes/pkg/api"
)

// Image represents a Docker image in the Redspread hierarchy. It wraps image.Image.
type Image struct {
	base
	image *image.Image
}

func NewImage(image *image.Image, defaults api.ObjectMeta, source string, objects ...deploy.KubeObject) (*Image, error) {
	base, err := newBase(EntityImage, defaults, source, objects)
	if err != nil {
		return nil, err
	} else {
		return &Image{base: base, image: image}, nil
	}
}

func (c Image) Deployment() deploy.Deployment {
	return deploy.Deployment{}
}

func (c Image) Images() []*image.Image {
	return []*image.Image{
		c.image,
	}
}

// Kubernetes representation of image
func (c Image) kube() string {
	return c.image.DockerName()
}
