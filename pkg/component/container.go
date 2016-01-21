package component

import (
	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	"k8s.io/kubernetes/pkg/api"
)

// Container represents api.Container in the Redspread hierarchy.
type Container struct {
	Base
	container api.Container
	image     *Image
}

func NewContainer(container api.Container, source string, objects ...deploy.KubeObject) (*Container, error) {
	base, err := newBase(ComponentContainer, source, objects)
	if err != nil {
		return nil, err
	}

	newContainer := Container{Base: base}
	if len(container.Image) != 0 {
		image, err := image.FromString(container.Image)
		if err != nil {
			return nil, err
		} else {
			newContainer.image, err = NewImage(image, source)
			if err != nil {
				return nil, err
			}
			container.Image = "placeholder"
		}
	}

	newContainer.container = container
	return &newContainer, nil
}

func (c Container) Deployment() deploy.Deployment {
	return deploy.Deployment{}
}

func (c Container) Images() []*image.Image {
	return c.image.Images()
}

func (c Container) kube() api.Container {
	c.container.Image = c.image.kube()
	return c.container
}
