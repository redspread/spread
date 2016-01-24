package entity

import (
	"fmt"

	"rsprd.com/spread/pkg/deploy"

	"k8s.io/kubernetes/pkg/api"
)

// A Entity is an entity (potentially containing sub-entities) that can be deployed to Kubernetes.
type Entity interface {
	deploy.Deployable
	Type() Type
	Objects() []deploy.KubeObject
	Source() string
	Attach(Entity) error
	DefaultMeta() api.ObjectMeta
}

// base provides fields that are shared between all Entitys.
type base struct {
	entityType Type
	objects    deploy.Deployment
	source     string
	defaults   api.ObjectMeta
}

func newBase(t Type, defaults api.ObjectMeta, source string, objects []deploy.KubeObject) (base base, err error) {
	base.defaults = defaults

	deployment := deploy.Deployment{}
	for _, obj := range objects {
		setMetaDefaults(obj, defaults)
		err = deployment.Add(obj)
		if err != nil {
			err = fmt.Errorf("error adding '%s': %v", source, err)
			return
		}
	}

	base.source = source
	base.entityType = t
	base.objects = deployment
	return
}

// Objects returns slice of objects attached to Entity
func (base base) Objects() []deploy.KubeObject {
	return base.objects.Objects()
}

// Source returns an import source specific identifier
func (base base) Source() string {
	return base.source
}

// DefaultMeta returns the ObjectMeta that the Entity was created with
func (base base) DefaultMeta() api.ObjectMeta {
	return base.defaults
}

// Type returns itself for trivial implementation of Entity
func (base base) Type() Type {
	return base.entityType
}

// Type identifies the entity's type.
type Type int

const (
	EntityApplication           Type = iota // Application (top of tree)
	EntityReplicationController             // Wrapper for api.ReplicationController
	EntityPod                               // Wrapper for api.Pod
	EntityContainer                         // Wrapper for api.Container
	EntityImage                             // Represented by api.Container's image field
)

// metaDefaults applies a set of defaults on a KubeObject. Non-empty fields on object override defaults.
func setMetaDefaults(obj deploy.KubeObject, defaults api.ObjectMeta) {
	meta := obj.GetObjectMeta()

	// if namespace is not set, use default
	namespace := api.NamespaceDefault
	if len(defaults.Namespace) > 0 {
		namespace = defaults.Namespace
	}

	if len(meta.GetNamespace()) == 0 {
		meta.SetNamespace(namespace)
	}

	// if name and generateName are not set use default generateName
	if len(defaults.GenerateName) > 0 && len(meta.GetName()) == 0 && len(meta.GetGenerateName()) == 0 {
		meta.SetGenerateName(defaults.GenerateName)
	}

	// set default labels
	labels := defaults.Labels
	for k, v := range meta.GetLabels() {
		labels[k] = v
	}
	meta.SetLabels(labels)

	// set default annotations
	annotations := defaults.Annotations
	for k, v := range meta.GetAnnotations() {
		annotations[k] = v
	}
	meta.SetAnnotations(annotations)
}
