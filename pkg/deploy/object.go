package deploy

import (
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/runtime"
)

// A KubeObject is an alias for Kubernetes objects.
type KubeObject interface {
	meta.ObjectMetaAccessor
	runtime.Object
}
