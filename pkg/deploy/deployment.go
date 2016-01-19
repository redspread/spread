package deploy

import (
	"errors"

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
	// TODO: implement deep-copy
	switch t := obj.(type) {
	case *api.ReplicationController:
		errList := validation.ValidateReplicationController(t)
		if len(errList) == 0 {
			d.rcs = append(d.rcs, t)
			return nil
		} else {
			return errList.ToAggregate()
		}
	case *api.Pod:
		errList := validation.ValidatePod(t)
		if len(errList) == 0 {
			d.pods = append(d.pods, t)
			return nil
		} else {
			return errList.ToAggregate()
		}
	case *api.Service:
		errList := validation.ValidateService(t)
		if len(errList) == 0 {
			d.services = append(d.services, t)
			return nil
		} else {
			return errList.ToAggregate()
		}
	case *api.Secret:
		errList := validation.ValidateSecret(t)
		if len(errList) == 0 {
			d.secrets = append(d.secrets, t)
			return nil
		} else {
			return errList.ToAggregate()
		}
	case *api.PersistentVolume:
		errList := validation.ValidatePersistentVolume(t)
		if len(errList) == 0 {
			d.volumes = append(d.volumes, t)
			return nil
		} else {
			return errList.ToAggregate()
		}
	case *api.PersistentVolumeClaim:
		errList := validation.ValidatePersistentVolumeClaim(t)
		if len(errList) == 0 {
			d.volumeClaims = append(d.volumeClaims, t)
			return nil
		} else {
			return errList.ToAggregate()
		}
	case *api.Namespace:
		errList := validation.ValidateNamespace(t)
		if len(errList) == 0 {
			d.namespaces = append(d.namespaces, t)
			return nil
		} else {
			return errList.ToAggregate()
		}
	default:
		return ErrorObjectNotSupported
	}
}

var (
	ErrorObjectNotSupported = errors.New("could not add to deployment, object not supported")
)
