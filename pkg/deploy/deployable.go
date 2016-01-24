package deploy

import (
	"rsprd.com/spread/pkg/image"
)

// A Deployable can produce a Deployment
type Deployable interface {
	Deployment() Deployment
	Images() []*image.Image
}
