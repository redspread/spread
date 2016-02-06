package deploy

import (
	"rsprd.com/spread/Godeps/_workspace/src/k8s.io/kubernetes/pkg/api/meta"
	"rsprd.com/spread/Godeps/_workspace/src/k8s.io/kubernetes/pkg/runtime"
)

// A KubeObject is an alias for Kubernetes objects.
type KubeObject interface {
	meta.ObjectMetaAccessor
	runtime.Object
}
