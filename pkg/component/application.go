package component

import (
	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	"k8s.io/kubernetes/pkg/api"
)

// Application is the root of the Redspread hierarchy
type Application struct {
	Base
	components []Component
}

func NewApplication(source string, defaults api.ObjectMeta, objects ...deploy.KubeObject) (*ReplicationController, error) {
	base, err := newBase(ComponentApplication, defaults, source, objects)
	if err != nil {
		return nil, err
	}

	return &ReplicationController{Base: base}, nil
}

func (c Application) Deployment() deploy.Deployment {
	return deploy.Deployment{}
}

func (c Application) Images() (images []*image.Image) {
	for _, v := range c.components {
		images = append(images, v.Images()...)
	}
	return
}
