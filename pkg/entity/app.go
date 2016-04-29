package entity

import (
	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	kube "k8s.io/kubernetes/pkg/api"
)

// ensure implements Entity
var _ Entity = new(App)

// App is a new Spread construct that can group together sets of Entities.
type App struct {
	base
	entities []Entity
}

// NewApp creates a new App entity.
func NewApp(entities []Entity, defaults kube.ObjectMeta, source string, objects ...deploy.KubeObject) (*App, error) {
	if entities == nil {
		entities = []Entity{}
	}

	base, err := newBase(EntityApplication, defaults, source, objects)
	if err != nil {
		return nil, err
	}

	return &App{base: base, entities: entities}, nil
}

// Deployment is created using Deployments of child Entities
func (c *App) Deployment() (*deploy.Deployment, error) {
	d := new(deploy.Deployment)
	for _, entity := range c.children() {
		deploy, err := entity.Deployment()
		if err != nil {
			return nil, err
		}

		err = d.AddDeployment(*deploy)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

// Images returns images of child entities
func (c *App) Images() (images []*image.Image) {
	for _, e := range c.children() {
		images = append(images, e.Images()...)
	}
	return
}

// Attach is allowed on any valid Entity
func (c *App) Attach(e Entity) error {
	if err := c.validAttach(e); err != nil {
		return err
	}

	// add to entities
	c.entities = append(c.entities, e)
	return nil
}

func (c *App) name() string {
	return "app"
}

func (c *App) children() []Entity {
	return c.entities
}
