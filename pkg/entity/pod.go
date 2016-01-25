package entity

import (
	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	"k8s.io/kubernetes/pkg/api"
)

// Pod represents api.Pod in the Redspread hierarchy.
type Pod struct {
	base
	pod        *api.Pod
	containers []*Container
}

func NewPod(kubePod *api.Pod, defaults api.ObjectMeta, source string, objects ...deploy.KubeObject) (*Pod, error) {
	base, err := newBase(EntityPod, defaults, source, objects)
	if err != nil {
		return nil, err
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

	pod.pod = kubePod
	return &pod, nil
}

func NewPodFromPodSpec(podSpec api.PodSpec, defaults api.ObjectMeta, source string, objects ...deploy.KubeObject) (*Pod, error) {
	pod := api.Pod{
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

func (c Pod) kube() *api.Pod {
	containers := []api.Container{}
	for _, container := range c.containers {
		containers = append(containers, container.kube())
	}

	c.pod.Spec.Containers = containers
	return c.pod
}
