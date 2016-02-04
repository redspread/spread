package input

import (
	"rsprd.com/spread/pkg/entity"
)

// EntityBuilder is used by input sources that create Entities.
type EntityBuilder interface {
	// Build returns an Entity based on the implementations internal logic, Errors are returned if state is invalid.
	Build() (entity.Entity, error)
}
