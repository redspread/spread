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

// NewImage creates a new Entity for the image.Image it's provided with.
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

// Deployment is created with image attached to default pod. The GenerateName of the Pod
// is the name of the image.
func (c *Image) Deployment() (*deploy.Deployment, error) {
	meta := kube.ObjectMeta{
		GenerateName: c.name(),
	}

	return deployWithPod(meta, c)
}

// Images returns the associated image.Image
func (c *Image) Images() []*image.Image {
	return []*image.Image{
		c.image,
	}
}

// Attach is not allowed on Images
func (c *Image) Attach(e Entity) error {
	return ErrorMaxAttached
}

func (c *Image) name() string {
	return c.image.DockerName()
}

func (c *Image) children() []Entity {
	return []Entity{}
}

// Kubernetes representation of image
func (c *Image) data() (image string, objects deploy.Deployment, err error) {
	objects.AddDeployment(c.objects)
	return c.image.DockerName(), objects, nil
}

var (
	// ErrorEmptyImageString is when an Image is created with an empty name
	ErrorEmptyImageString = errors.New("image.Image's String was empty")
	// ErrorNilImage is when an Image is created with a nil *image.Image.
	ErrorNilImage = errors.New("*image.Image cannot be nil")
)
