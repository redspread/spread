package dir

import (
	"os"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/entity"
)

type dirSource string

// NewDirSource returns a source for the path to a file or directory. Path must be valid.
func NewDirSource(path string) (dirSource, error) {
	if _, err := os.Stat(path); err != nil {
		return dirSource(""), err
	}
	return dirSource(path), nil
}

// Entities returns the entities of the requested type from the source. Errors if any invalid entities.
func (ds *dirSource) Entities(t entity.Type) ([]entity.Entity, error) {
	return []entity.Entity{}, nil
}

// Objects returns the Kubernetes objects available from the source. Errors if any invalid objects.
func (ds *dirSource) Objects() ([]deploy.KubeObject, error) {
	return []deploy.KubeObject{}, nil
}

func getPods(filename string) []entity.Entity {
	return []entity.Entity{}
}
