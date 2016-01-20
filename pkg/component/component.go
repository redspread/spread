package component

import (
	"fmt"

	"rsprd.com/spread/pkg/deploy"
)

// A Component is an entity (potentially containing sub-components) that can be deployed to Kubernetes.
type Component interface {
	deploy.Deployable
	Type() Type
	Objects() []*deploy.KubeObject
	Source() string
}

// Base provides fields that are shared between all Components.
type Base struct {
	componentType Type
	objects       deploy.Deployment
	source        string
}

func newBase(t Type, source string, objects []deploy.KubeObject) (base Base, err error) {
	deployment := deploy.NewDeployment()
	for _, v := range objects {
		err = deployment.Add(v)
		if err != nil {
			err = fmt.Errorf("error adding '%s': %v", source, err)
			return
		}
	}

	base.source = source
	base.componentType = t
	base.objects = deployment
	return
}

// Objects returns slice of objects attached to Component
func (base Base) Objects() []*deploy.KubeObject {
	// TODO: Implement
	return []*deploy.KubeObject{}
}

// Source returns an import source specific identifier
func (base Base) Source() string {
	return base.source
}

// Type returns itself for trivial implementation of Component
func (base Base) Type() Type {
	return base.componentType
}

// Type identifies the component's type.
type Type int

const (
	ComponentApplication           Type = iota // Application (top of tree)
	ComponentReplicationController             // Wrapper for api.ReplicationController
	ComponentPod                               // Wrapper for api.Pod
	ComponentContainer                         // Wrapper for api.Container
	ComponentImage                             // Represented by api.Container's image field
)
