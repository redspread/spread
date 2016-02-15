package entity

import (
	"fmt"
	"sort"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/validation"
)

// DefaultPodSpec is the default for newly created Pods
var DefaultPodSpec = kube.PodSpec{
	RestartPolicy: kube.RestartPolicyAlways,
	DNSPolicy:     kube.DNSDefault,
}

type containers []*Container

func (c containers) Len() int {
	return len(c)
}

func (c containers) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c containers) Less(i, j int) bool {
	return c[i].name() < c[j].name()
}

// Pod represents kube.Pod in the Redspread hierarchy.
type Pod struct {
	base
	pod *kube.Pod
	containers
}

// NewPod creates a Entity for a corresponding *kube.Pod. Pod must be valid.
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
		}
		pod.containers = append(pod.containers, container)
	}
	kubePod.Spec.Containers = []kube.Container{}

	base.setDefaults(kubePod)
	if err := validatePod(kubePod, true); err != nil {
		return nil, fmt.Errorf("could not create Pod from `%s`: %v", source, err)
	}

	sort.Sort(pod.containers)

	pod.pod = kubePod
	return &pod, nil
}

// NewPodFromPodSpec creates a new Pod using ObjectMeta and a PodSpec.
func NewPodFromPodSpec(meta kube.ObjectMeta, podSpec kube.PodSpec, defaults kube.ObjectMeta, source string, objects ...deploy.KubeObject) (*Pod, error) {
	pod := kube.Pod{
		ObjectMeta: meta,
		Spec:       podSpec,
	}
	return NewPod(&pod, defaults, source, objects...)
}

// NewDefaultPod creates a Pod using spreads defaults without any containers attached. Containers must be attached
// before this Pod can be deployed.
func NewDefaultPod(meta kube.ObjectMeta, source string, objects ...deploy.KubeObject) (*Pod, error) {
	return NewPodFromPodSpec(meta, DefaultPodSpec, kube.ObjectMeta{}, source, objects...)
}

// Deployment is created containing Pod with attached Containers.
func (c *Pod) Deployment() (*deploy.Deployment, error) {
	deployment := deploy.Deployment{}

	// create Pod from tree of Entities
	kubePod, childObj, err := c.data()
	if err != nil {
		return nil, err
	}

	// add Pod to deployment
	err = deployment.Add(kubePod)
	if err != nil {
		return nil, err
	}

	// add child objects
	err = deployment.AddDeployment(childObj)
	if err != nil {
		return nil, err
	}
	return &deployment, nil
}

// Images for all Containers in Pod
func (c *Pod) Images() (images []*image.Image) {
	for _, v := range c.containers {
		images = append(images, v.Images()...)
	}
	return
}

// Attach appends Images and Containers.
func (c *Pod) Attach(curEntity Entity) error {
	if curEntity == nil {
		return ErrorNilEntity
	}

	if err := c.validAttach(curEntity); err != nil {
		return err
	}
	for {
		switch e := curEntity.(type) {
		case *Image:
			container, err := newDefaultContainer(e.name(), e.Source())
			if err != nil {
				return err
			}

			err = container.Attach(e)
			curEntity = container
			break
		case *Container:
			c.containers = append(c.containers, e)
			sort.Sort(c.containers)
			return nil
		default:
			panic("Unexpected type")
		}
	}
}

func (c *Pod) name() string {
	return c.pod.ObjectMeta.Name
}

func (c *Pod) children() (children []Entity) {
	for _, v := range c.containers {
		children = append(children, v)
	}
	return
}

func (c *Pod) data() (pod *kube.Pod, objects deploy.Deployment, err error) {
	containers := []kube.Container{}
	for _, container := range c.containers {
		kubeContainer, cObj, err := container.data()
		if err != nil {
			return nil, objects, err
		}
		containers = append(containers, kubeContainer)
		// add containers objects
		objects.AddDeployment(cObj)
	}

	if len(containers) == 0 {
		return nil, objects, ErrorEntityNotReady
	}

	// add own objects
	objects.AddDeployment(c.objects)

	pod, err = copyPod(c.pod)
	if err != nil {
		return nil, objects, err
	}

	pod.Spec.Containers = containers
	return pod, objects, nil
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

	meta := pod.GetObjectMeta()
	if len(meta.GetName()) == 0 && len(meta.GetGenerateName()) > 0 {
		errList = errList.Filter(func(e error) bool {
			return e.Error() == "metadata.name: Required value: name or generateName is required"
		})
	}

	return errList.ToAggregate()
}

func deployWithPod(meta kube.ObjectMeta, attached Entity) (*deploy.Deployment, error) {
	pod, err := NewDefaultPod(meta, attached.Source())
	if err != nil {
		return nil, err
	}

	err = pod.Attach(attached)
	if err != nil {
		return nil, err
	}

	return pod.Deployment()
}
