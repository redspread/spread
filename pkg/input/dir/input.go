package dir

import (
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

	rcs, err := d.Entities(entity.EntityReplicationController)
	if err != nil {
		return nil, err
	}

	for _, rc := range rcs {
		err = base.Attach(rc)
		if err != nil {
			return nil, err
		}
	}

	pods, err := d.Entities(entity.EntityPod)

	for _, pod := range pods {
		err = base.Attach(pod)
		if err != nil {
			return nil, err
		}
	}
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

	// create app
	return entity.NewApp(nil, kube.ObjectMeta{}, string(d.FileSource), objects...)
}
