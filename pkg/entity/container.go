package entity

import (
	"fmt"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	kube "rsprd.com/spread/Godeps/_workspace/src/k8s.io/kubernetes/pkg/api"
	"rsprd.com/spread/Godeps/_workspace/src/k8s.io/kubernetes/pkg/api/validation"
)

// DefaultContainer is the default for all new Containers
var DefaultContainer = kube.Container{
	ImagePullPolicy: kube.PullAlways,
}

// NoImageErrStr is the error message returned by k8s when the image field has been left blank.
const NoImageErrStr = "spec.containers[0].image: Required value"

// Container represents kube.Container in the Redspread hierarchy.
type Container struct {
	base
	container kube.Container
	image     *Image
}

// NewContainer creates a new Entity for the provided kube.Container. Container must be valid.
func NewContainer(container kube.Container, defaults kube.ObjectMeta, source string, objects ...deploy.KubeObject) (*Container, error) {
	err := validateContainer(container)
	if err != nil {
		return nil, fmt.Errorf("could not create Container from `%s`: %v", source, err)
	}

	base, err := newBase(EntityContainer, defaults, source, objects)
	if err != nil {
		return nil, err
	}

	newContainer := Container{base: base}
	if len(container.Image) != 0 {
		image, err := image.FromString(container.Image)
		if err != nil {
			return nil, err
		}

		newContainer.image, err = NewImage(image, defaults, source)
		if err != nil {
			return nil, err
		}
		container.Image = "placeholder"
	}

	newContainer.container = container
	return &newContainer, nil
}

func newDefaultContainer(name, source string) (*Container, error) {
	kubeContainer := DefaultContainer
	kubeContainer.Name = name
	return NewContainer(kubeContainer, kube.ObjectMeta{}, source)
}

// Deployment is created with Container attached to Pod. The pod is named after the kube.Container.
func (c *Container) Deployment() (*deploy.Deployment, error) {
	meta := kube.ObjectMeta{
		GenerateName: c.name(),
	}

	return deployWithPod(meta, c)
}

// Images returns the Container's image
func (c *Container) Images() []*image.Image {
	var images []*image.Image
	if c.image != nil {
		images = c.image.Images()
	}
	return images
}

// Attach is only possible for images.
func (c *Container) Attach(e Entity) error {
	if c.image != nil {
		return ErrorMaxAttached
	}

	if err := c.validAttach(e); err != nil {
		return err
	}

	// entity must be image
	c.image = e.(*Image)
	return nil
}

func (c *Container) name() string {
	return c.container.Name
}

func (c *Container) children() []Entity {
	return []Entity{
		c.image,
	}
}

func (c *Container) data() (container kube.Container, objects deploy.Deployment, err error) {
	if c.image == nil {
		return kube.Container{}, deploy.Deployment{}, ErrorEntityNotReady
	}

	// if image exists should always return valid result
	container = c.container
	image, objects, err := c.image.data()
	if err != nil {
		return container, deploy.Deployment{}, err
	}

	err = objects.AddDeployment(c.objects)
	if err != nil {
		return container, deploy.Deployment{}, err
	}
	container.Image = image
	return container, objects, nil
}

func validateContainer(c kube.Container) error {
	validMeta := kube.ObjectMeta{
		Name:      "valid",
		Namespace: "object",
	}

	pod := kube.Pod{
		ObjectMeta: validMeta,
		Spec: kube.PodSpec{
			Containers:    []kube.Container{c},
			RestartPolicy: kube.RestartPolicyAlways,
			DNSPolicy:     kube.DNSClusterFirst,
		},
	}

	// fake volumes to allow validation
	for _, mount := range c.VolumeMounts {
		volume := kube.Volume{
			Name:         mount.Name,
			VolumeSource: kube.VolumeSource{EmptyDir: &kube.EmptyDirVolumeSource{}},
		}
		pod.Spec.Volumes = append(pod.Spec.Volumes, volume)
	}

	errList := validation.ValidatePod(&pod)

	// Remove error for missing image field
	return errList.Filter(func(e error) bool {
		return e.Error() == NoImageErrStr
	}).ToAggregate()
}
