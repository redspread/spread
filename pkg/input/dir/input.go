package dir

import (
	"errors"
	"fmt"

	"rsprd.com/spread/pkg/entity"

	kube "k8s.io/kubernetes/pkg/api"
)

type fileInput struct {
	FileSource
}

func NewFileInput(path string) (*fileInput, error) {
	src, err := NewFileSource(path)
	if err != nil {
		return nil, err
	}

	return &fileInput{
		FileSource: src,
	}, nil
}

func (d fileInput) Build() (entity.Entity, error) {
	base, err := d.base()
	if err != nil {
		return nil, err
	}

	err = d.buildEntity(base)
	return base, err
}

func (d fileInput) base() (entity.Entity, error) {
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
	} else if len(rcs) == 1 {
		return pods[0], nil
	}

	return entity.NewDefaultPod(kube.ObjectMeta{}, string(d.FileSource), objects...)
}

func (d fileInput) buildEntity(parent entity.Entity) error {
	if parent == nil {
		return errors.New("parent can't be nil")
	}

	// check if any more types to attach
	if parent.Type() >= entity.EntityImage {
		return nil
	}

	// increment type number
	for typNum := int(parent.Type()) + 1; typNum <= int(entity.EntityImage); typNum++ {
		entities, err := d.Entities(entity.Type(typNum))
		if err != nil {
			return err
		}

		for _, e := range entities {
			err = parent.Attach(e)
			if err != nil {
				return fmt.Errorf("could not attach '%v' to '%v': %v", e.Type(), parent.Type(), err)
			}

			err = d.buildEntity(e)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var (
	// ErrTooManyRCs is returned when there are more than one RC in a directory
	ErrTooManyRCs = errors.New("only one RC is allowed per directory")

	// ErrTooManyPods is when there are more than one Pod in a directory
	ErrTooManyPods = errors.New("only one Pod is allowed per directory")
)
