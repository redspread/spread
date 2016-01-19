package component

import (
	"fmt"

	"rsprd.com/spread/pkg/deploy"
)

// A Component is an entity (potentially containing sub-components) that can be deployed to Kubernetes.
type Component interface {
	deploy.Deployable
	Type() ComponentType
	Objects() []*deploy.KubeObject
	Source() string
}

// ComponentBase provides fields that are shared between all Components.
type ComponentBase struct {
	ComponentType
	objects deploy.Deployment
	source  string
}

func newComponentBase(t ComponentType, source string, objects []deploy.KubeObject) (base ComponentBase, err error) {
	deployment := deploy.NewDeployment()
	for _, v := range objects {
		err = deployment.Add(v)
		if err != nil {
			err = fmt.Errorf("error adding '%s': %v", source, err)
			return
		}
	}

	base.source = source
	base.ComponentType = t
	base.objects = deployment
	return
}

// Objects returns slice of objects attached to Component
func (base ComponentBase) Objects() []*deploy.KubeObject {
	// TODO: Implement
	return []*deploy.KubeObject{}
}

// Source returns an import source specific identifier
func (base ComponentBase) Source() string {
	return base.source
}

// ComponentType identifies the component's type.
type ComponentType int

// Type returns itself for trivial implementation of Component
func (t ComponentType) Type() ComponentType {
	return t
}

const (
	ComponentApplication           ComponentType = iota // Application (top of tree)
	ComponentReplicationController                      // Wrapper for api.ReplicationController
	ComponentPod                                        // Wrapper for api.Pod
	ComponentContainer                                  // Wrapper for api.Container
	ComponentImage                                      // Represented by api.Container's image field
)
