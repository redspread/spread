package entity

import (
	"errors"
	"fmt"

	"rsprd.com/spread/pkg/deploy"

	kube "k8s.io/kubernetes/pkg/api"
)

// Entity is the fundamental building block of `spread`'s internal representation of state. Entities exist for both
// constructs that already exist in Kubernetes (RC, Pod, Container) as well as ones that have been added (image).
//
// An entity must at all times be capable of producing a Deployment.
//
// Entities are attached in the order: Image -> Container -> Pod -> RC. Attachments out of order are not allowed.
type Entity interface {
	deploy.Deployable
	Type() Type
	Objects() []deploy.KubeObject
	Source() string
	Attach(Entity) error
	DefaultMeta() kube.ObjectMeta

	// prevents implementation by external packages
	name() string
	children() []Entity
}

// Builder is used by input sources that create Entities.
type Builder interface {
	// Build returns an Entity based on the implementations internal logic, Errors are returned if state is invalid.
	Build() (Entity, error)
}

// base provides fields that are shared between all Entitys.
type base struct {
	entityType Type
	objects    deploy.Deployment
	source     string
	defaults   kube.ObjectMeta
}

func newBase(t Type, defaults kube.ObjectMeta, source string, objects []deploy.KubeObject) (base base, err error) {
	base.defaults = defaults

	deployment := deploy.Deployment{}
	for _, obj := range objects {
		if obj == nil {
			err = ErrorNilObject
			return
		}
		base.setDefaults(obj)
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
func (base base) DefaultMeta() kube.ObjectMeta {
	return base.defaults
}

// Type returns itself for trivial implementation of Entity
func (base base) Type() Type {
	return base.entityType
}

// validAttach checks object types to see if the attach is allowed. Objects can
// only be attached to objects higher in the hierarchy. However, to the nature of iota Application is 0, RC is 1, ...
func (base base) validAttach(e Entity) error {
	if e == nil {
		return ErrorNilEntity
	}
	if e.Type() > EntityImage {
		return ErrorInvalidAttachType
	} else if base.Type() > e.Type() {
		return ErrorBadAttachOrder
	}
	return nil
}

// setDefaults sets the bases defaults on an object
func (base base) setDefaults(obj deploy.KubeObject) {
	setMetaDefaults(obj, base.defaults)
}

// Type identifies the entity's type.
type Type int

const (
	// EntityApplication is for Application (top of tree)
	EntityApplication Type = iota
	// EntityReplicationController is for ReplicationController
	EntityReplicationController
	// EntityPod is for Pod
	EntityPod
	// EntityContainer is for Container
	EntityContainer
	// EntityImage is for Image
	EntityImage
)

// String prints the name of the type
func (t Type) String() string {
	switch t {
	case EntityApplication:
		return "Application"
	case EntityReplicationController:
		return "ReplicationController"
	case EntityPod:
		return "Pod"
	case EntityContainer:
		return "Container"
	case EntityImage:
		return "Image"
	default:
		return ""
	}
}

// metaDefaults applies a set of defaults on a KubeObject. Non-empty fields on object override defaults.
func setMetaDefaults(obj deploy.KubeObject, defaults kube.ObjectMeta) {
	meta := obj.GetObjectMeta()

	// if namespace is not set, use default
	namespace := kube.NamespaceDefault
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
	labels := map[string]string{}
	if defaults.Labels != nil {
		labels = defaults.Labels
	}
	for k, v := range meta.GetLabels() {
		labels[k] = v
	}
	meta.SetLabels(labels)

	// set default annotations
	annotations := map[string]string{}
	if defaults.Annotations != nil {
		annotations = defaults.Annotations
	}
	for k, v := range meta.GetAnnotations() {
		annotations[k] = v
	}
	meta.SetAnnotations(annotations)
}

var (
	// ErrorEntityNotReady is when Entity is not in a valid state to be deployed.
	ErrorEntityNotReady = errors.New("entity not ready to be deployed")
	// ErrorNilObject is when a deploy.KubeObject is null.
	ErrorNilObject = errors.New("an object was nil, this is not allowed")
	// ErrorInvalidAttachType is when the Type of the attached object in invalid.
	ErrorInvalidAttachType = errors.New("the entity to be attached is of an unknown type")
	// ErrorBadAttachOrder is when an attachment is attempted in an improper order
	ErrorBadAttachOrder = errors.New("entities cannot attach to entities of a lower type")
	// ErrorMaxAttached is when an Entity cannot support more Entities attaching
	ErrorMaxAttached = errors.New("no more entities can be attached")
	// ErrorNilEntity is when an entity is nil
	ErrorNilEntity = errors.New("entities cannot be nil")
)
