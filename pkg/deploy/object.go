package deploy

import (
	"errors"
	"fmt"

	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	types "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/runtime"
)

// A KubeObject is an alias for Kubernetes objects.
type KubeObject interface {
	meta.ObjectMetaAccessor
	runtime.Object
}

// ObjectPath returns the full path of an object.
// This uses the format "<apiVersion>/namespaces/<namespace>/<kind>/<name>"
func ObjectPath(obj KubeObject) (string, error) {
	// attempt to determine ObjectKind
	gkv, err := objectKind(obj)
	if err != nil {
		return "", fmt.Errorf("could not get object path: %v", err)
	}

	meta := obj.GetObjectMeta()
	return fmt.Sprintf("%s/namespaces/%s/%s/%s", obj.GetObjectKind().GroupVersionKind().Version, meta.GetNamespace(), gkv.Kind, meta.GetName()), nil
}

// objectKind is a helper function which determines type information from given KubeObject.
// An error is returned if the GroupVersionKind is empty or cannot be determined.
func objectKind(obj KubeObject) (types.GroupVersionKind, error) {
	gkv, err := kube.Scheme.ObjectKind(obj)
	if err != nil {
		return types.GroupVersionKind{}, err
	} else if gkv.IsEmpty() {
		return types.GroupVersionKind{}, errors.New("empty ObjectKind")
	}
	return gkv, nil
}
