package entity

import (
	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	"k8s.io/kubernetes/pkg/api"
)

// ReplicationController represents api.ReplicationController in the Redspread hierarchy.
type ReplicationController struct {
	base
	rc  *api.ReplicationController
	pod *Pod
}

func NewReplicationController(kubeRC *api.ReplicationController, defaults api.ObjectMeta, source string, objects ...deploy.KubeObject) (*ReplicationController, error) {
	base, err := newBase(EntityReplicationController, defaults, source, objects)
	if err != nil {
		return nil, err
	}

	rc := ReplicationController{base: base}
	if kubeRC.Spec.Template != nil {
		rc.pod, err = NewPodFromPodSpec(kubeRC.Spec.Template.Spec, defaults, source)
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

func (c ReplicationController) kube() (*api.ReplicationController, error) {
	if c.pod == nil {
		return nil, ErrorEntityNotReady
	}

	pod, err := c.pod.kube()
	if err != nil {
		return nil, err
	}

	c.rc.Spec.Template.Spec = pod.Spec
	return c.rc, nil
}
