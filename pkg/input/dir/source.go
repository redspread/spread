package dir

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/entity"

	kube "k8s.io/kubernetes/pkg/api"
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
func (fs fileSource) Objects() (objects []deploy.KubeObject, err error) {
	dirPath := path.Join(string(fs), ObjectsDir)

	err = walkPathForObjects(dirPath, func(info *resource.Info, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		obj, ok := info.Object.(deploy.KubeObject)
		if !ok {
			return ErrTypeMismatch
		}

		objects = append(objects, obj)
		return nil
	})

	// don't throw error if simply didn't find anything
	if err != nil && !checkErrNotFound(err) {
		return nil, err
	}
	return objects, nil
}

// getRCs returns entities for the rcs in the RCFile
func (fs fileSource) getRCs() (rcs []entity.Entity, err error) {
	filePath := path.Join(string(fs), RCFile)

	err = walkPathForObjects(filePath, func(info *resource.Info, err error) error {
		kubeRC, ok := info.Object.(*kube.ReplicationController)
		if !ok {
			return fmt.Errorf("expected type `ReplicationController` but found `%s`", info.Object.GetObjectKind().GroupVersionKind().Kind)
		}

		rc, err := entity.NewReplicationController(kubeRC, kube.ObjectMeta{}, info.Source)
		if err != nil {
			return err
		}

		rcs = append(rcs, rc)
		return nil
	})
	return
}

// getPods returns Pods for the rcs in the PodFile
func (fs fileSource) getPods() (pods []entity.Entity, err error) {
	filePath := path.Join(string(fs), PodFile)

	err = walkPathForObjects(filePath, func(info *resource.Info, err error) error {
		kubePod, ok := info.Object.(*kube.Pod)
		if !ok {
			return fmt.Errorf("expected type `Pod` but found `%s`", info.Object.GetObjectKind().GroupVersionKind().Kind)
		}

		pod, err := entity.NewPod(kubePod, kube.ObjectMeta{}, info.Source)
		if err != nil {
			return err
		}

		pods = append(pods, pod)
		return nil
	})
	return
}

// getContainers creates entities from files with the ContainerExtension
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
	result := resource.NewBuilder(mapper, typer, resource.DisabledClientForMapping{}, f.Decoder(true)).
		ContinueOnError().
		Schema(schema).
		FilenameParam(false, path).
		Flatten().
		Do()

	err = result.Err()
	if err != nil && !checkErrNotFound(err) {
		return err
	}

	err = result.Visit(fn)
	if err != nil {
		return err
	}
	return nil
}

func checkErrNotFound(err error) bool {
	return strings.HasPrefix(err.Error(), "you must provide one or more resources")
}

var (
	// ErrInvalidType is returned when the entity.Type is unknown
	ErrInvalidType        = errors.New("passed invalid type")
	ErrTypeMismatch       = errors.New("was expecting a KubeObject")
	ErrEntityPathNotFound = errors.New("the path being searched doesn't exist")
)
