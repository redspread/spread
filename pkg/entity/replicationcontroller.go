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

// NewReplicationController creates a new Entity for the provided kube.ReplicationController. Must be valid.
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
		templateMeta := kubeRC.Spec.Template.ObjectMeta
		templateMeta.Name = kubeRC.Name
		rc.pod, err = NewPodFromPodSpec(templateMeta, kubeRC.Spec.Template.Spec, defaults, source)
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

// Deployment is created for RC attached with it's Pod.
func (c *ReplicationController) Deployment() (*deploy.Deployment, error) {
	deployment := new(deploy.Deployment)

	// create RC
	kubeRC, childObj, err := c.data()
	if err != nil {
		return nil, err
	}

	// add RC to deployment
	err = deployment.Add(kubeRC)
	if err != nil {
		return nil, err
	}

	// add child objects
	err = deployment.AddDeployment(childObj)
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

// Images contained by ReplicationController's Pods.
func (c *ReplicationController) Images() (images []*image.Image) {
	if c.pod != nil {
		images = c.pod.Images()
	}
	return images
}

// Attach allows Pods, Containers, and Images to be attached.
func (c *ReplicationController) Attach(e Entity) error {
	if e == nil {
		return ErrorNilEntity
	}

	if err := c.validAttach(e); err != nil {
		return err
	}

	switch e := e.(type) {
	case *Pod:
		if c.pod != nil {
			return ErrorMaxAttached
		}

		c.pod = e
		return nil
	default:
		if c.pod != nil {
			return c.pod.Attach(e)
		}

		meta := kube.ObjectMeta{Name: e.name()}
		pod, err := NewDefaultPod(meta, e.Source())
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

func (c *ReplicationController) children() []Entity {
	return []Entity{
		c.pod,
	}
}

func (c *ReplicationController) data() (*kube.ReplicationController, deploy.Deployment, error) {
	if c.pod == nil {
		return nil, deploy.Deployment{}, ErrorEntityNotReady
	}

	rc := c.rc
	pod, objects, err := c.pod.data()
	if err != nil {
		return nil, deploy.Deployment{}, err
	}

	err = objects.AddDeployment(c.objects)
	if err != nil {
		return nil, deploy.Deployment{}, err
	}

	// add selectors
	meta := pod.ObjectMeta
	meta.Labels = c.rc.Spec.Selector
	meta.Name = ""
	rc.Spec.Template = &kube.PodTemplateSpec{
		ObjectMeta: meta,
		Spec:       pod.Spec,
	}
	return rc, objects, nil
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
