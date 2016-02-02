package entity

import (
	"fmt"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/validation"
)

// ReplicationController represents kube.ReplicationController in the Redspread hierarchy.
type ReplicationController struct {
	base
	rc  *kube.ReplicationController
	pod *Pod
}

func NewReplicationController(kubeRC *kube.ReplicationController, defaults kube.ObjectMeta, source string, objects ...deploy.KubeObject) (*ReplicationController, error) {
	if kubeRC == nil {
		return nil, fmt.Errorf("cannot create ReplicationController from nil `%s`", source)
	}

	base, err := newBase(EntityReplicationController, defaults, source, objects)
	if err != nil {
		return nil, err
	}

	// deep copy
	kubeRC, err = copyRC(kubeRC)
	if err != nil {
		return nil, err
	}

	rc := ReplicationController{base: base}
	if kubeRC.Spec.Template != nil {
		rc.pod, err = NewPodFromPodSpec(kubeRC.Spec.Template.ObjectMeta, kubeRC.Spec.Template.Spec, defaults, source)
		if err != nil {
			return nil, err
		}
		kubeRC.Spec.Template = nil
	}

	base.setDefaults(kubeRC)
	if err = validateRC(kubeRC); err != nil {
		return nil, err
	}
	rc.rc = kubeRC

	return &rc, nil
}

func (c *ReplicationController) Deployment() (*deploy.Deployment, error) {
	return nil, nil
}

func (c *ReplicationController) Images() (images []*image.Image) {
	if c.pod != nil {
		images = c.pod.Images()
	}
	return images
}

func (c *ReplicationController) Attach(e Entity) error {
	if c.pod != nil {
		return ErrorMaxAttached
	}

	if err := c.validAttach(e); err != nil {
		return err
	}

	switch e := e.(type) {
	case *Pod:
		c.pod = e
		return nil
	default:
		meta := kube.ObjectMeta{Name: e.name()}
		pod, err := newDefaultPod(meta, e.Source())
		if err != nil {
			return err
		}

		err = pod.Attach(e)
		if err != nil {
			return err
		}

		return c.Attach(pod)
	}
}

func (c *ReplicationController) name() string {
	return c.rc.ObjectMeta.Name
}

func copyRC(rc *kube.ReplicationController) (*kube.ReplicationController, error) {
	copy, err := kube.Scheme.DeepCopy(rc)
	if err != nil {
		return nil, err
	}

	return copy.(*kube.ReplicationController), nil
}

func validateRC(rc *kube.ReplicationController) error {
	errList := validation.ValidateReplicationController(rc).Filter(
		// remove errors about missing template
		func(e error) bool {
			return e.Error() == "spec.template: Required value"
		},
	)
	return errList.ToAggregate()
}
