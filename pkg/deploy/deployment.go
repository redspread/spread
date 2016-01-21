package deploy

import (
	"errors"
	"fmt"
	"reflect"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/validation"
)

// A Deployment is a collection of Kubernetes deployed. Deployment stores a slice of deployable Kubernetes objects.
// It can be used to create deployments deployments and is how the current state of a deployment is returned.
type Deployment struct {
	rcs          []*api.ReplicationController
	pods         []*api.Pod
	services     []*api.Service
	secrets      []*api.Secret
	volumes      []*api.PersistentVolume
	volumeClaims []*api.PersistentVolumeClaim
	namespaces   []*api.Namespace
}

func NewDeployment() Deployment {
	return Deployment{}
}

func (d *Deployment) Add(obj KubeObject) error {
	copy, err := deepCopy(obj)
	if err != nil {
		return err
	}

	switch t := copy.(type) {
	case *api.ReplicationController:
		errList := validation.ValidateReplicationController(t)
		if len(errList) == 0 {
			for _, v := range d.rcs {
				if err = assertUniqueName(copy, v); err != nil {
					return err
				}
			}
			d.rcs = append(d.rcs, t)
			return nil
		} else {
			return errList.ToAggregate()
		}

	case *api.Pod:
		errList := validation.ValidatePod(t)
		if len(errList) == 0 {
			for _, v := range d.pods {
				if err = assertUniqueName(copy, v); err != nil {
					return err
				}
			}
			d.pods = append(d.pods, t)
			return nil
		} else {
			return errList.ToAggregate()
		}

	case *api.Service:
		errList := validation.ValidateService(t)
		if len(errList) == 0 {
			for _, v := range d.services {
				if err = assertUniqueName(copy, v); err != nil {
					return err
				}
			}
			d.services = append(d.services, t)
			return nil
		} else {
			return errList.ToAggregate()
		}

	case *api.Secret:
		errList := validation.ValidateSecret(t)
		if len(errList) == 0 {
			for _, v := range d.secrets {
				if err = assertUniqueName(copy, v); err != nil {
					return err
				}
			}
			d.secrets = append(d.secrets, t)
			return nil
		} else {
			return errList.ToAggregate()
		}

	case *api.PersistentVolume:
		errList := validation.ValidatePersistentVolume(t)
		if len(errList) == 0 {
			for _, v := range d.volumes {
				if err = assertUniqueName(copy, v); err != nil {
					return err
				}
			}
			d.volumes = append(d.volumes, t)
			return nil
		} else {
			return errList.ToAggregate()
		}

	case *api.PersistentVolumeClaim:
		errList := validation.ValidatePersistentVolumeClaim(t)
		if len(errList) == 0 {
			for _, v := range d.volumeClaims {
				if err = assertUniqueName(copy, v); err != nil {
					return err
				}
			}
			d.volumeClaims = append(d.volumeClaims, t)
			return nil
		} else {
			return errList.ToAggregate()
		}

	case *api.Namespace:
		errList := validation.ValidateNamespace(t)
		if len(errList) == 0 {
			for _, v := range d.namespaces {
				if err = assertUniqueName(copy, v); err != nil {
					return err
				}
			}
			d.namespaces = append(d.namespaces, t)
			return nil
		} else {
			return errList.ToAggregate()
		}
	default:
		return ErrorObjectNotSupported
	}
}

func (d Deployment) Equals(other Deployment) bool {
	if !equivalent(d.rcs, other.rcs) {
		return false
	}

	if !equivalent(d.pods, other.pods) {
		return false
	}

	if !equivalent(d.services, other.services) {
		return false
	}

	if !equivalent(d.secrets, other.secrets) {
		return false
	}

	if !equivalent(d.volumes, other.volumes) {
		return false
	}

	if !equivalent(d.volumeClaims, other.volumeClaims) {
		return false
	}

	if !equivalent(d.namespaces, other.namespaces) {
		return false
	}
	return true
}

func (d Deployment) Objects() (obj []KubeObject) {
	obj = appendObjects(obj, d.rcs)
	obj = appendObjects(obj, d.pods)
	obj = appendObjects(obj, d.services)
	obj = appendObjects(obj, d.secrets)
	obj = appendObjects(obj, d.volumes)
	obj = appendObjects(obj, d.volumeClaims)
	obj = appendObjects(obj, d.namespaces)
	return

}

func appendObjects(obj []KubeObject, objectSlice interface{}) []KubeObject {
	sliceVal := reflect.ValueOf(objectSlice)
	for i := 0; i < sliceVal.Len(); i++ {
		objCopy, err := api.Scheme.DeepCopy(sliceVal.Index(i).Interface())
		if err != nil {
			panic(err)
		}
		obj = append(obj, objCopy.(KubeObject))
	}
	return obj
}

// assertUniqueName checks a slice of objects for naming collisions. It assumes that the slice is of a single type.
func assertUniqueName(a, b KubeObject) error {
	aMeta, bMeta := a.GetObjectMeta(), b.GetObjectMeta()

	if aMeta.GetName() == bMeta.GetName() && aMeta.GetNamespace() == bMeta.GetNamespace() {
		return fmt.Errorf("name/namespace combination (%s/%s) already exists for type", aMeta.GetName(), aMeta.GetNamespace())
	}

	return nil
}

func equivalent(a, b interface{}) bool {
	aSlice, bSlice := reflect.ValueOf(a), reflect.ValueOf(b)
	if aSlice.Len() != bSlice.Len() {
		return false
	}

	for i := 0; i < aSlice.Len(); i++ {
		aPtr := aSlice.Index(i).Interface()
		found := false
		for j := 0; j < bSlice.Len(); j++ {
			bPtr := bSlice.Index(j).Interface()
			if api.Semantic.DeepEqual(aPtr, bPtr) {
				found = true
			}
		}

		if !found {
			return false
		}
	}
	return true
}

// deepCopy creates a deep copy of the Kubernetes object given.
func deepCopy(obj KubeObject) (KubeObject, error) {
	copy, err := api.Scheme.DeepCopy(obj)
	if err != nil {
		return nil, err
	}
	return copy.(KubeObject), nil
}

var (
	ErrorObjectNotSupported = errors.New("could not add to deployment, object not supported")
)
