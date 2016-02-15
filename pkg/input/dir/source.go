package dir

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/entity"

	"github.com/ghodss/yaml"
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

// FileSource provides access to Entities stored according to the Redspread file convention.
// Under this convention RC's are looked for in RCFile, Pods are looked for in PodFile, and anything with the extension
// ContainerExtension is considered a container.
type FileSource string

// NewFileSource returns a source for the path to a file or directory. Path must be valid.
func NewFileSource(path string) (FileSource, error) {
	if _, err := os.Stat(path); err != nil {
		return FileSource(""), fmt.Errorf("could not create fileSource: %v", err)
	}
	return FileSource(path), nil
}

// Entities returns the entities of the requested type from the source. Errors if any invalid entities.
func (fs FileSource) Entities(t entity.Type, objects ...deploy.KubeObject) ([]entity.Entity, error) {
	switch t {
	case entity.EntityReplicationController:
		return fs.rcs(objects)
	case entity.EntityPod:
		return fs.pods(objects)
	case entity.EntityContainer:
		return fs.containers(objects)
	case entity.EntityImage:
		// getting images not implemented
		return []entity.Entity{}, nil
	}

	// if unknown type, return error
	return []entity.Entity{}, ErrInvalidType
}

// Objects returns the Kubernetes objects available from the source. Errors if any invalid objects.
func (fs FileSource) Objects() (objects []deploy.KubeObject, err error) {
	dirPath := path.Join(string(fs), ObjectsDir)

	err = walkPathForObjects(dirPath, func(info *resource.Info, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		obj := info.Object.(deploy.KubeObject)

		objects = append(objects, obj)
		return nil
	})

	// don't throw error if simply didn't find anything
	if err != nil && !checkErrNoResources(err) && !checkErrPathDoesNotExist(err) {
		return nil, err
	}
	return objects, nil
}

// rcs returns entities for the rcs in the RCFile
func (fs FileSource) rcs(objects []deploy.KubeObject) (rcs []entity.Entity, err error) {
	filePath := path.Join(string(fs), RCFile)

	err = walkPathForObjects(filePath, func(info *resource.Info, err error) error {
		kubeRC, ok := info.Object.(*kube.ReplicationController)
		if !ok {
			return fmt.Errorf("expected type `ReplicationController` but found `%s`", info.Object.GetObjectKind().GroupVersionKind().Kind)
		}

		rc, err := entity.NewReplicationController(kubeRC, kube.ObjectMeta{}, info.Source, objects...)
		if err != nil {
			return err
		}

		rcs = append(rcs, rc)
		return nil
	})

	if checkErrPathDoesNotExist(err) {
		// it's okay if directory doesn't exit
		err = nil
	}
	return
}

// pods returns Pods for the rcs in the PodFile
func (fs FileSource) pods(objects []deploy.KubeObject) (pods []entity.Entity, err error) {
	filePath := path.Join(string(fs), PodFile)

	err = walkPathForObjects(filePath, func(info *resource.Info, err error) error {
		kubePod, ok := info.Object.(*kube.Pod)
		if !ok {
			return fmt.Errorf("expected type `Pod` but found `%s`", info.Object.GetObjectKind().GroupVersionKind().Kind)
		}

		pod, err := entity.NewPod(kubePod, kube.ObjectMeta{}, info.Source, objects...)
		if err != nil {
			return err
		}

		pods = append(pods, pod)
		return nil
	})

	if checkErrPathDoesNotExist(err) {
		// it's okay if directory doesn't exit
		err = nil
	}
	return
}

// containers creates entities from files with the ContainerExtension
func (fs FileSource) containers(objects []deploy.KubeObject) (containers []entity.Entity, err error) {
	info, err := os.Stat(string(fs))
	if err != nil {
		return
	}

	// check if file
	if !info.IsDir() {
		kubeCtr, err := unmarshalContainer(string(fs))
		if err != nil {
			return nil, err
		}

		ctr, err := entity.NewContainer(kubeCtr, kube.ObjectMeta{}, string(fs), objects...)
		if err != nil {
			return nil, err
		}

		return []entity.Entity{ctr}, nil
	}

	dir, err := os.Open(string(fs))
	if err != nil {
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ContainerExtension {
				filename := path.Join(string(fs), file.Name())

				kubeCtr, err := unmarshalContainer(filename)
				if err != nil {
					return nil, err
				}

				ctr, err := entity.NewContainer(kubeCtr, kube.ObjectMeta{}, filename)
				if err != nil {
					return nil, err
				}

				containers = append(containers, ctr)
			}
		}
	}
	return containers, nil
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
	if err != nil && !checkErrNoResources(err) {
		return err
	}

	err = result.Visit(fn)
	if err != nil {
		return err
	}
	return nil
}

func unmarshalContainer(path string) (ctr kube.Container, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &ctr)
	return
}

func checkErrNoResources(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "you must provide one or more resources")
}

func checkErrPathDoesNotExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasSuffix(err.Error(), " does not exist")
}

var (
	// ErrInvalidType is returned when the entity.Type is unknown
	ErrInvalidType = errors.New("passed invalid type")
)
