package deploy

import (
	"rsprd.com/spread/pkg/image"
)

// A Deployable can produce a Deployment
type Deployable interface {
	// Deployment creates a new Deployment based on the types current state. Errors are returned if not possible.
	Deployment() (*Deployment, error)
	// Images returns the images required for deployment
	Images() []*image.Image
}
