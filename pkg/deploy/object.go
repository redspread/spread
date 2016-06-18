package deploy

import (
	"errors"
	"fmt"
	"strings"

	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	types "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/runtime"

	"rsprd.com/spread/pkg/data"
	pb "rsprd.com/spread/pkg/spreadproto"
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
	path := fmt.Sprintf("namespaces/%s/%s/%s", meta.GetNamespace(), gkv.Kind, meta.GetName())
	return strings.ToLower(path), nil
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

func KubeObjectFromObject(kind string, obj *pb.Object) (KubeObject, error) {
	base := BaseObject(kind)
	if base == nil {
		return nil, fmt.Errorf("unable to find Kind for '%s'", kind)
	}

	err := data.Unmarshal(obj, &base)
	if err != nil {
		return nil, err
	}
	return base, nil
}
