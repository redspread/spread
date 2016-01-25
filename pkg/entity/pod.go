package entity

import (
	"fmt"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/validation"
)

// Pod represents api.Pod in the Redspread hierarchy.
type Pod struct {
	base
	pod        *api.Pod
	containers []*Container
}

func NewPod(kubePod *api.Pod, defaults api.ObjectMeta, source string, objects ...deploy.KubeObject) (*Pod, error) {
	if kubePod == nil {
		return nil, fmt.Errorf("cannot create Pod from nil `%s`", source)
	}

	// deep copy
	kubePod, err := copyPod(kubePod)
	if err != nil {
		return nil, err
	}

	base, err := newBase(EntityPod, defaults, source, objects)
	if err != nil {
		return nil, fmt.Errorf("could not create Pod from `%s`: %v", source, err)
	}

	pod := Pod{base: base}
	for _, v := range kubePod.Spec.Containers {
		container, err := NewContainer(v, defaults, source)
		if err != nil {
			return nil, err
		} else {
			pod.containers = append(pod.containers, container)
		}
	}
	kubePod.Spec.Containers = []api.Container{}

	base.setDefaults(kubePod)
	if err := validatePod(kubePod, true); err != nil {
		return nil, fmt.Errorf("could not create Pod from `%s`: %v", source, err)
	}

	pod.pod = kubePod
	return &pod, nil
}

func NewPodFromPodSpec(name string, podSpec api.PodSpec, defaults api.ObjectMeta, source string, objects ...deploy.KubeObject) (*Pod, error) {
	pod := api.Pod{
		ObjectMeta: api.ObjectMeta{
			Name: name,
		},
		Spec: podSpec,
	}
	return NewPod(&pod, defaults, source, objects...)
}

func (c Pod) Deployment() (*deploy.Deployment, error) {
	return nil, nil
}

func (c Pod) Images() (images []*image.Image) {
	for _, v := range c.containers {
		images = append(images, v.Images()...)
	}
	return
}

func (c Pod) Attach(e Entity) error {
	return nil
}

func (c Pod) kube() (*api.Pod, error) {
	containers := []api.Container{}
	for _, container := range c.containers {
		kubeContainer, err := container.kube()
		if err != nil {
			return nil, err
		}
		containers = append(containers, kubeContainer)
	}

	if len(containers) == 0 {
		return nil, ErrorEntityNotReady
	}

	pod, err := copyPod(c.pod)
	if err != nil {
		return nil, err
	}

	pod.Spec.Containers = containers
	return pod, nil
}

func copyPod(pod *api.Pod) (*api.Pod, error) {
	copy, err := api.Scheme.DeepCopy(pod)
	if err != nil {
		return nil, err
	}

	return copy.(*api.Pod), nil
}

func validatePod(pod *api.Pod, ignoreContainers bool) error {
	errList := validation.ValidatePod(pod)

	// remove error for no containers if requested
	if ignoreContainers {
		errList = errList.Filter(func(e error) bool {
			return e.Error() == "spec.containers: Required value"
		})
	}
	return errList.ToAggregate()
}
