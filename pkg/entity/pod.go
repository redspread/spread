package entity

import (
	"fmt"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/validation"
)

var DefaultPodSpec = kube.PodSpec{
	RestartPolicy: kube.RestartPolicyAlways,
	DNSPolicy:     kube.DNSDefault,
}

// Pod represents kube.Pod in the Redspread hierarchy.
type Pod struct {
	base
	pod        *kube.Pod
	containers []*Container
}

func NewPod(kubePod *kube.Pod, defaults kube.ObjectMeta, source string, objects ...deploy.KubeObject) (*Pod, error) {
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
	kubePod.Spec.Containers = []kube.Container{}

	base.setDefaults(kubePod)
	if err := validatePod(kubePod, true); err != nil {
		return nil, fmt.Errorf("could not create Pod from `%s`: %v", source, err)
	}

	pod.pod = kubePod
	return &pod, nil
}

func NewPodFromPodSpec(meta kube.ObjectMeta, podSpec kube.PodSpec, defaults kube.ObjectMeta, source string, objects ...deploy.KubeObject) (*Pod, error) {
	pod := kube.Pod{
		ObjectMeta: meta,
		Spec:       podSpec,
	}
	return NewPod(&pod, defaults, source, objects...)
}

func newDefaultPod(meta kube.ObjectMeta, source string) (*Pod, error) {
	return NewPodFromPodSpec(meta, DefaultPodSpec, kube.ObjectMeta{}, source)
}

func (c *Pod) Deployment() (*deploy.Deployment, error) {
	return nil, nil
}

func (c *Pod) Images() (images []*image.Image) {
	for _, v := range c.containers {
		images = append(images, v.Images()...)
	}
	return
}

func (c *Pod) Attach(e Entity) error {
	return nil
}

func (c *Pod) kube() (*kube.Pod, error) {
	containers := []kube.Container{}
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

func copyPod(pod *kube.Pod) (*kube.Pod, error) {
	copy, err := kube.Scheme.DeepCopy(pod)
	if err != nil {
		return nil, err
	}

	return copy.(*kube.Pod), nil
}

func validatePod(pod *kube.Pod, ignoreContainers bool) error {
	errList := validation.ValidatePod(pod)

	// remove error for no containers if requested
	if ignoreContainers {
		errList = errList.Filter(func(e error) bool {
			return e.Error() == "spec.containers: Required value"
		})
	}
	return errList.ToAggregate()
}

func deployWithPod(meta kube.ObjectMeta, attached Entity) (*deploy.Deployment, error) {
	pod, err := newDefaultPod(meta, attached.Source())
	if err != nil {
		return nil, err
	}

	err = pod.Attach(attached)
	if err != nil {
		return nil, err
	}

	return pod.Deployment()
}
