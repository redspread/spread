package dir

import (
	"errors"
	"fmt"
	"os"

	"rsprd.com/spread/pkg/entity"

	kube "k8s.io/kubernetes/pkg/api"
)

// FileInput produces Entities from objects stored at a given path of the filesystem using the Redspread convention.
type FileInput struct {
	FileSource
}

// NewFileInput returns an Input based on a file system
func NewFileInput(path string) (*FileInput, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	return &FileInput{
		FileSource: FileSource(path),
	}, nil
}

// Build creates an Entity by selecting the entity with the lowest type number and attaching higher objects recursively
func (d *FileInput) Build() (entity.Entity, error) {
	base, err := d.base()
	if err != nil {
		return nil, err
	}

	err = d.buildEntity(base)
	return base, err
}

// Path returns the location the fileInput was created
func (d FileInput) Path() string {
	return string(d.FileSource)
}

func (d *FileInput) base() (entity.Entity, error) {
	objects, err := d.Objects()
	if err != nil {
		return nil, err
	}

	// try as rc
	rcs, err := d.Entities(entity.EntityReplicationController, objects...)
	if err != nil {
		return nil, err
	} else if len(rcs) > 1 {
		return nil, ErrTooManyRCs
	} else if len(rcs) == 1 {
		return rcs[0], nil
	}

	// try as pod
	pods, err := d.Entities(entity.EntityPod, objects...)
	if err != nil {
		return nil, err
	} else if len(pods) > 1 {
		return nil, ErrTooManyPods
	} else if len(pods) == 1 {
		return pods[0], nil
	}

	return entity.NewDefaultPod(kube.ObjectMeta{GenerateName: "spread"}, string(d.FileSource), objects...)
}

func (d *FileInput) buildEntity(parent entity.Entity) error {
	if parent == nil {
		return errors.New("parent can't be nil")
	}

	// check if any more types to attach
	if parent.Type() >= entity.EntityImage {
		return nil
	}

	childTypeNum := int(parent.Type()) + 1

	// increment type number
	for typNum := childTypeNum; typNum <= int(entity.EntityImage); typNum++ {
		typ := entity.Type(typNum)
		entities, err := d.Entities(typ)
		if err != nil {
			return err
		}

		// if none, check next Entity type
		if len(entities) == 0 {
			continue
		}

		// attach any Entities matching type
		for _, e := range entities {
			err = parent.Attach(e)
			if err != nil {
				return fmt.Errorf("could not attach '%s' to '%s': %v", e.Type().String(), parent.Type().String(), err)
			}

			err = d.buildEntity(e)
			if err != nil {
				return err
			}
		}

		// stop on Entity attach
		return nil
	}
	return nil
}

var (
	// ErrTooManyRCs is returned when there are more than one RC in a directory
	ErrTooManyRCs = errors.New("only one RC is allowed per directory")

	// ErrTooManyPods is when there are more than one Pod in a directory
	ErrTooManyPods = errors.New("only one Pod is allowed per directory")
)
