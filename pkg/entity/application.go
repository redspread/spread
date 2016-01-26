package entity

import (
	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	kube "k8s.io/kubernetes/pkg/api"
)

// Application is the root of the Redspread hierarchy
type Application struct {
	base
	entities []Entity
}

func NewApplication(source string, defaults kube.ObjectMeta, objects ...deploy.KubeObject) (*ReplicationController, error) {
	base, err := newBase(EntityApplication, defaults, source, objects)
	if err != nil {
		return nil, err
	}

	return &ReplicationController{base: base}, nil
}

func (c Application) Deployment() (*deploy.Deployment, error) {
	return nil, nil
}

func (c Application) Images() (images []*image.Image) {
	for _, v := range c.entities {
		images = append(images, v.Images()...)
	}
	return
}

func (c Application) Attach(e Entity) error {
	return nil
}
