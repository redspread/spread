package deploy

import (
	kube "k8s.io/kubernetes/pkg/api"
)

var shortForms = map[string]string{
	"cs":     "componentstatuses",
	"ds":     "daemonsets",
	"ep":     "endpoints",
	"ev":     "events",
	"hpa":    "horizontalpodautoscalers",
	"ing":    "ingresses",
	"limits": "limitranges",
	"no":     "nodes",
	"ns":     "namespaces",
	"po":     "pods",
	"psp":    "podSecurityPolicies",
	"pvc":    "persistentvolumeclaims",
	"pv":     "persistentvolumes",
	"quota":  "resourcequotas",
	"rc":     "replicationcontrollers",
	"rs":     "replicasets",
	"svc":    "services",
}

func KubeShortForm(resource string) string {
	if long, ok := shortForms[resource]; ok {
		return long
	}
	return resource
}

// TODO: ensure all Kinds are supported
var kinds = map[string]KubeObject{
	"componentstatus": &kube.ComponentStatus{},
	"endpoint":        &kube.Endpoints{},
	"event":           &kube.Event{},
	"limitrange":      &kube.LimitRange{},
	"node":            &kube.Node{},
	"namespace":       &kube.Namespace{},
	"pod":             &kube.Pod{},
	"persistentvolumeclaim": &kube.PersistentVolumeClaim{},
	"persistentvolume":      &kube.PersistentVolume{},
	"resourcequota":         &kube.ResourceQuota{},
	"replicationcontroller": &kube.ReplicationController{},
	"service":               &kube.Service{},
	"secret":                &kube.Secret{},
}

// BaseObject returns a Kubernetes object of the given kind to be used to populate.
// Nil is returned if the Kind is unknown
func BaseObject(kind string) KubeObject {
	obj, ok := kinds[kind]
	if !ok {
		return nil
	}

	copy, err := kube.Scheme.DeepCopy(obj)
	if err != nil {
		return nil
	}
	return copy.(KubeObject)
}
