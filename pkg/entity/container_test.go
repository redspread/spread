package entity

import (
	"testing"

	"rsprd.com/spread/pkg/deploy"

	"github.com/gh/stretchr/testify/assert"
	kube "k8s.io/kubernetes/pkg/api"
)

func TestContainerWithImageDeployment(t *testing.T) {
	const imageName = "busybox:latest"
	kubeContainer := kube.Container{
		Name:            "simple-container",
		Image:           imageName,
		Command:         []string{"/bin/busybox", "ls"},
		ImagePullPolicy: kube.PullAlways,
	}

	ctr, err := NewContainer(kubeContainer, kube.ObjectMeta{}, "simpleTest")
	assert.NoError(t, err, "should be able to create container")
	assert.NotNil(t, ctr.image, "an image should have been created")

	// check images
	images := ctr.Images()
	assert.Len(t, images, 1, "should have single image")

	expectedImage := newDockerImage(t, imageName)
	actualImage := images[0]
	assert.Equal(t, expectedImage.DockerName(), actualImage.DockerName(), "image should not have changed")

	// check kube
	kubectr, err := ctr.data()
	assert.NoError(t, err, "should be able to produce kube")
	assert.True(t, kube.Semantic.DeepEqual(&kubectr, &kubeContainer), "kube should be same as container")

	pod := kube.Pod{
		ObjectMeta: kube.ObjectMeta{
			GenerateName: kubeContainer.Name,
			Namespace:    kube.NamespaceDefault,
		},
		Spec: kube.PodSpec{
			Containers: []kube.Container{
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
	kubeContainer := kube.Container{
		Name: "no-image-container",
		// no image
		Command:         []string{"/bin/busybox", "ls"},
		ImagePullPolicy: kube.PullAlways,
	}

	ctr, err := NewContainer(kubeContainer, kube.ObjectMeta{}, "noImage")
	assert.NoError(t, err, "should be able to create container")
	assert.Nil(t, ctr.image, "no image should exist")

	images := ctr.Images()
	assert.Len(t, images, 0, "no image should have been created")

	_, err = ctr.data()
	assert.Error(t, err, "container is not ready")

	_, err = ctr.Deployment()
	assert.Error(t, err, "cannot be deployed without image")
}

func TestContainerAttach(t *testing.T) {
	imageName := "to-be-attached"
	image := testNewImage(t, imageName, kube.ObjectMeta{}, "test", []deploy.KubeObject{})

	kubeContainer := testNewKubeContainer("test-container", "") // no image
	container, err := NewContainer(kubeContainer, kube.ObjectMeta{}, "attach")
	assert.NoError(t, err, "valid container")

	_, err = container.Deployment()
	assert.Error(t, err, "cannot be deployed without image")

	err = container.Attach(image)
	assert.NoError(t, err, "attach should be allowed")

	images := container.Images()
	assert.Len(t, images, 1, "only one image should exist")

	kubeContainer.Image = imageName
	pod := kube.Pod{
		ObjectMeta: kube.ObjectMeta{
			GenerateName: kubeContainer.Name,
			Namespace:    kube.NamespaceDefault,
		},
		Spec: kube.PodSpec{
			Containers: []kube.Container{
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
	kubeContainer := testNewKubeContainer("test-container", "test-image")
	objects := []deploy.KubeObject{
		createSecret(""), // invalid - must have name
	}

	_, err := NewContainer(kubeContainer, kube.ObjectMeta{}, "invalidobjects", objects...)
	assert.Error(t, err, "container should not be created with invalid objects")
}

func TestContainerInvalidContainer(t *testing.T) {
	kubeContainer := kube.Container{
		// invalid - no name
		Image:           "invalid-container",
		Command:         []string{"/bin/busybox", "ls"},
		ImagePullPolicy: kube.PullAlways,
	}
	_, err := NewContainer(kubeContainer, kube.ObjectMeta{}, "invalidcontainer")
	assert.Error(t, err, "name is missing, container is invalid")
}

func TestContainerInvalidImage(t *testing.T) {
	imageName := "*T*H*I*S* IS ILLEGAL"
	kubeContainer := testNewKubeContainer("invalid-image", imageName)
	_, err := NewContainer(kubeContainer, kube.ObjectMeta{}, "invalidimage")
	assert.Error(t, err, "image was invalid")
}

func testNewKubeContainer(name, imageName string) kube.Container {
	return kube.Container{
		Name:            name,
		Image:           imageName,
		Command:         []string{"/bin/busybox", "ls"},
		ImagePullPolicy: kube.PullAlways,
	}
}

func testRandomContainer(defaults kube.ObjectMeta, source string, objects []deploy.KubeObject) *Container {
	kubeContainer := testNewKubeContainer(randomString(10), randomString(15))
	container, _ := NewContainer(kubeContainer, defaults, source, objects...)
	return container
}

var (
	testKubeContainerSourcegraph = kube.Container{
		Name:  "src",
		Image: "sourcegraph/sourcegraph:latest",
		VolumeMounts: []kube.VolumeMount{
			kube.VolumeMount{Name: "config", MountPath: "/home/sourcegraph/.sourcegraph"},
		},
		Ports: []kube.ContainerPort{
			kube.ContainerPort{ContainerPort: 80, Protocol: kube.ProtocolTCP},
			kube.ContainerPort{ContainerPort: 443, Protocol: kube.ProtocolTCP},
		},
		ImagePullPolicy: kube.PullAlways,
	}

	testKubeContainerPostgres = kube.Container{
		Name:  "postgres",
		Image: "postgres:9.5",
		VolumeMounts: []kube.VolumeMount{
			kube.VolumeMount{Name: "db", MountPath: "/var/lib/postgresql/data/pgdata"},
		},
		Ports: []kube.ContainerPort{
			kube.ContainerPort{ContainerPort: 5432, Protocol: kube.ProtocolTCP},
		},
		ImagePullPolicy: kube.PullAlways,
	}
)
