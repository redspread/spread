package entity

import (
	"errors"
	"strings"

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

func (c *Image) Deployment() (*deploy.Deployment, error) {
	podName := strings.Join([]string{c.name(), "image"}, "-")
	meta := kube.ObjectMeta{
		Name: podName,
	}

	return deployWithPod(meta, c)
}

func (c *Image) Images() []*image.Image {
	return []*image.Image{
		c.image,
	}
}

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
	ErrorEmptyImageString = errors.New("image.Image's DockerString was empty")
	ErrorNilImage         = errors.New("*image.Image cannot be nil")
)
