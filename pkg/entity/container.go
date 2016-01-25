package entity

import (
	"fmt"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/validation"
)

// Container represents api.Container in the Redspread hierarchy.
type Container struct {
	base
	container api.Container
	image     *Image
}

func NewContainer(container api.Container, defaults api.ObjectMeta, source string, objects ...deploy.KubeObject) (*Container, error) {
	err := validateContainer(container)
	if err != nil {
		return nil, fmt.Errorf("could not create Container from `%s`: %v", source, err)
	}

	base, err := newBase(EntityContainer, defaults, source, objects)
	if err != nil {
		return nil, err
	}

	newContainer := Container{base: base}
	if len(container.Image) != 0 {
		image, err := image.FromString(container.Image)
		if err != nil {
			return nil, err
		} else {
			newContainer.image, err = NewImage(image, defaults, source)
			if err != nil {
				return nil, err
			}
			container.Image = "placeholder"
		}
	}

	newContainer.container = container
	return &newContainer, nil
}

func (c Container) Deployment() (*deploy.Deployment, error) {
	return nil, nil
}

func (c Container) Images() []*image.Image {
	return c.image.Images()
}

func (c Container) Attach(e Entity) error {
	return nil
}

func (c Container) kube() api.Container {
	c.container.Image = c.image.kube()
	return c.container
}

func validateContainer(c api.Container) error {
	validMeta := api.ObjectMeta{
		Name:      "valid",
		Namespace: "object",
	}

	pod := api.Pod{
		ObjectMeta: validMeta,
		Spec: api.PodSpec{
			Containers:    []api.Container{c},
			RestartPolicy: api.RestartPolicyAlways,
			DNSPolicy:     api.DNSClusterFirst,
		},
	}

	return validation.ValidatePod(&pod).ToAggregate()
}
