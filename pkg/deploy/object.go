package deploy

import (
	"errors"
	"fmt"

	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/runtime"
)

// A KubeObject is an alias for Kubernetes objects.
type KubeObject interface {
	meta.ObjectMetaAccessor
	runtime.Object
}

// FullObjectPath returns the full path of an object.
// This uses the format "<apiVersion>/namespaces/<namespace>/<kind>/<name>"
func FullObjectPath(obj KubeObject) (string, error) {
	errText := "could not get object path"

	// attempt to determine ObjectKind
	gkv, err := kube.Scheme.ObjectKind(obj)
	if err != nil {
		return "", fmt.Errorf(errText+": %v", err)
	} else if gkv.IsEmpty() {
		return "", errors.New(errText + ": empty ObjectKind")
	}

	meta := obj.GetObjectMeta()
	return fmt.Sprintf("%s/namespaces/%s/%s/%s", obj.GetObjectKind().GroupVersionKind().Version, meta.GetNamespace(), gkv.Kind, meta.GetName()), nil
}
