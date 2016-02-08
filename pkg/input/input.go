package input

import (
	"rsprd.com/spread/pkg/entity"
)

// Input represents a source of Entities and metadata.
type Input interface {
	entity.Entity
}
