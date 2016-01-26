package entity

import (
	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	kube "k8s.io/kubernetes/pkg/api"
)

// ReplicationController represents kube.ReplicationController in the Redspread hierarchy.
type ReplicationController struct {
	base
	rc  *kube.ReplicationController
	pod *Pod
}

func NewReplicationController(kubeRC *kube.ReplicationController, defaults kube.ObjectMeta, source string, objects ...deploy.KubeObject) (*ReplicationController, error) {
	base, err := newBase(EntityReplicationController, defaults, source, objects)
	if err != nil {
		return nil, err
	}

	rc := ReplicationController{base: base}
	if kubeRC.Spec.Template != nil {
		rc.pod, err = NewPodFromPodSpec("unamed", kubeRC.Spec.Template.Spec, defaults, source)
		if err != nil {
			return nil, err
		}
		kubeRC.Spec.Template = nil
	}
	return &rc, nil
}

func (c ReplicationController) Deployment() (*deploy.Deployment, error) {
	return nil, nil
}

func (c ReplicationController) Images() (images []*image.Image) {
	return c.pod.Images()
}

func (c ReplicationController) Attach(e Entity) error {
	return nil
}
