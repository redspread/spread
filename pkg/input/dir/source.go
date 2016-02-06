package dir

import (
	"errors"
	"fmt"
	"os"
	"path"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/entity"

	kubectl "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/resource"
)

const (
	// RCFile is the filename checked for Replication Controllers.
	RCFile = "rc.yml"
	// PodFile is the filename checked for Pods.
	PodFile = "pod.yml"
	// ContainerExtension is the file extension checked for Containers.
	ContainerExtension = ".ctr"
	// ObjectsDir is the directory checked for arbitrary Kubernetes objects.
	ObjectsDir = ".k2e"
)

type fileSource string

// NewFileSource returns a source for the path to a file or directory. Path must be valid.
func NewFileSource(path string) (fileSource, error) {
	if _, err := os.Stat(path); err != nil {
		return fileSource(""), fmt.Errorf("could not create fileSource: %v", err)
	}
	return fileSource(path), nil
}

// Entities returns the entities of the requested type from the source. Errors if any invalid entities.
func (fs fileSource) Entities(t entity.Type) ([]entity.Entity, error) {
	switch t {
	case entity.EntityReplicationController:
		return fs.getRCs()
	case entity.EntityPod:
		return fs.getPods()
	case entity.EntityContainer:
		return fs.getContainers()
	case entity.EntityImage:
		// getting images not implemented
		return []entity.Entity{}, nil
	}

	// if unknown type, return error
	return []entity.Entity{}, ErrInvalidType
}

// Objects returns the Kubernetes objects available from the source. Errors if any invalid objects.
func (fs fileSource) Objects() ([]deploy.KubeObject, error) {
	dirPath := path.Join(string(fs), ObjectsDir)

	err := walkPathForObjects(dirPath, func(info *resource.Info, err error) error {
		return nil
	})
	if err != nil {
		return nil, err
	}
	return []deploy.KubeObject{}, nil
}

func (fs fileSource) getRCs() ([]entity.Entity, error) {
	return []entity.Entity{}, nil
}

func (fs fileSource) getPods() ([]entity.Entity, error) {
	return []entity.Entity{}, nil
}

func (fs fileSource) getContainers() ([]entity.Entity, error) {
	return []entity.Entity{}, nil
}

func walkPathForObjects(path string, fn resource.VisitorFunc) error {
	f := kubectl.NewFactory(nil)

	// todo: does "cacheDir" need to be parameterizable?
	schema, err := f.Validator(false, "")
	if err != nil {
		return err
	}

	mapper, typer := f.Object()
	result := resource.NewBuilder(mapper, typer, resource.ClientMapperFunc(f.ClientForMapping), f.Decoder(true)).
		ContinueOnError().
		Schema(schema).
		FilenameParam(false, path).
		Flatten().
		Do()

	err = result.Err()
	if err != nil {
		return err
	}

	return result.Visit(fn)
}

var (
	// ErrInvalidType is returned when the entity.Type is unknown
	ErrInvalidType = errors.New("passed invalid type")
)
