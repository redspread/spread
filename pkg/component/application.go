package component

import (
	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"
)

// Application is the root of the Redspread hierarchy
type Application struct {
	ComponentBase
	components []Component
}

func NewApplication(source string, objects ...deploy.KubeObject) (*ReplicationController, error) {
	base, err := newComponentBase(ComponentApplication, source, objects)
	if err != nil {
		return nil, err
	}

	return &ReplicationController{ComponentBase: base}, nil
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
