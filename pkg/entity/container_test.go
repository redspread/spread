package entity

import (
	"testing"

	"rsprd.com/spread/pkg/deploy"

	"github.com/gh/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

func TestContainerWithImageDeployment(t *testing.T) {
	const imageName = "busybox:latest"
	kubeContainer := api.Container{
		Name:            "simple-container",
		Image:           imageName,
		Command:         []string{"/bin/busybox", "ls"},
		ImagePullPolicy: api.PullAlways,
	}

	ctr, err := NewContainer(kubeContainer, api.ObjectMeta{}, "simpleTest")
	assert.NoError(t, err, "should be able to create container")
	assert.NotNil(t, ctr.image, "an image should have been created")

	// check images
	images := ctr.Images()
	assert.Len(t, images, 1, "should have single image")

	expectedImage := newDockerImage(t, imageName)
	actualImage := images[0]
	assert.Equal(t, expectedImage.DockerName(), actualImage.DockerName(), "image should not have changed")

	// check kube
	kube, err := ctr.kube()
	assert.NoError(t, err, "should be able to produce kube")
	assert.True(t, api.Semantic.DeepEqual(&kube, &kubeContainer), "kube should be same as container")

	pod := api.Pod{
		ObjectMeta: api.ObjectMeta{
			GenerateName: kubeContainer.Name,
			Namespace:    api.NamespaceDefault,
		},
		Spec: api.PodSpec{
			Containers: []api.Container{
				kubeContainer,
			},
		},
	}

	expected := deploy.Deployment{}
	assert.NoError(t, expected.Add(&pod), "valid pod")
	actual, err := ctr.Deployment()
	assert.NoError(t, err, "deploy ok")

	assert.True(t, expected.Equal(actual), "should be equivlant")
}

func TestContainerNoImageDeployment(t *testing.T) {
	kubeContainer := api.Container{
		Name: "no-image-container",
		// no image
		Command:         []string{"/bin/busybox", "ls"},
		ImagePullPolicy: api.PullAlways,
	}

	ctr, err := NewContainer(kubeContainer, api.ObjectMeta{}, "noImage")
	assert.NoError(t, err, "should be able to create container")
	assert.Nil(t, ctr.image, "no image should exist")

	images := ctr.Images()
	assert.Len(t, images, 0, "no image should have been created")

	_, err = ctr.kube()
	assert.Error(t, err, "container is not ready")

	_, err = ctr.Deployment()
	assert.Error(t, err, "cannot be deployed without image")
}

func TestContainerAttach(t *testing.T) {
	imageName := "to-be-attached"
	image := testNewImage(t, imageName, api.ObjectMeta{}, "test", []deploy.KubeObject{})

	kubeContainer := newKubeContainer("test-container", "") // no image
	container, err := NewContainer(kubeContainer, api.ObjectMeta{}, "attach")
	assert.NoError(t, err, "valid container")

	_, err = container.Deployment()
	assert.Error(t, err, "cannot be deployed without image")

	err = container.Attach(image)
	assert.NoError(t, err, "attach should be allowed")

	kubeContainer.Image = imageName
	pod := api.Pod{
		ObjectMeta: api.ObjectMeta{
			GenerateName: kubeContainer.Name,
			Namespace:    api.NamespaceDefault,
		},
		Spec: api.PodSpec{
			Containers: []api.Container{
				kubeContainer,
			},
		},
	}

	expected := deploy.Deployment{}
	assert.NoError(t, expected.Add(&pod), "valid pod")
	actual, err := container.Deployment()
	assert.NoError(t, err, "deploy ok")

	assert.True(t, expected.Equal(actual), "should be equivlant")
}

func TestContainerBadObject(t *testing.T) {
	kubeContainer := newKubeContainer("test-container", "test-image")
	objects := []deploy.KubeObject{
		createSecret(""), // invalid - must have name
	}

	_, err := NewContainer(kubeContainer, api.ObjectMeta{}, "invalidobjects", objects...)
	assert.Error(t, err, "container should not be created with invalid objects")
}

func TestContainerInvalidContainer(t *testing.T) {
	kubeContainer := api.Container{
		// invalid - no name
		Image:           "invalid-container",
		Command:         []string{"/bin/busybox", "ls"},
		ImagePullPolicy: api.PullAlways,
	}
	_, err := NewContainer(kubeContainer, api.ObjectMeta{}, "invalidcontainer")
	assert.Error(t, err, "name is missing, container is invalid")
}

func TestContainerInvalidImage(t *testing.T) {
	imageName := "*T*H*I*S* IS ILLEGAL"
	kubeContainer := newKubeContainer("invalid-image", imageName)
	_, err := NewContainer(kubeContainer, api.ObjectMeta{}, "invalidimage")
	assert.Error(t, err, "image was invalid")
}

func newKubeContainer(name, imageName string) api.Container {
	return api.Container{
		Name:            name,
		Image:           imageName,
		Command:         []string{"/bin/busybox", "ls"},
		ImagePullPolicy: api.PullAlways,
	}
}
