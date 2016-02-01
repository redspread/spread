package entity

import (
	"strings"
	"testing"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	"github.com/stretchr/testify/assert"
	kube "k8s.io/kubernetes/pkg/api"
)

func TestImageDeployment(t *testing.T) {
	imageName := "arch"
	simple := newDockerImage(t, imageName)

	image, err := NewImage(simple, kube.ObjectMeta{}, "test")
	assert.NoError(t, err, "valid image")

	containerName := strings.Join([]string{imageName, "container"}, "-")
	expectedPod := kube.Pod{
		ObjectMeta: kube.ObjectMeta{
			GenerateName: imageName,
			Namespace:    "default",
		},
		Spec: kube.PodSpec{
			Containers: []kube.Container{
				kube.Container{
					Name:            containerName,
					Image:           imageName,
					ImagePullPolicy: kube.PullIfNotPresent,
				},
			},
			RestartPolicy: kube.RestartPolicyAlways,
			DNSPolicy:     kube.DNSDefault,
		},
	}

	expected := deploy.Deployment{}
	assert.NoError(t, expected.Add(&expectedPod), "should be able to add pod")

	actual, err := image.Deployment()
	assert.NoError(t, err, "deploy ok")
	if !expected.Equal(actual) {
		t.Errorf("Expected: %#v, saw: %#v", expected, actual)
	}
}

func TestImageImages(t *testing.T) {
	imageName := "gcr.io/google_containers/cassandra:v7"
	simple := newDockerImage(t, imageName)

	image, err := NewImage(simple, kube.ObjectMeta{}, "test")
	if err != nil {
		t.Fatalf("Could not create Image entity: %v", err)
	}

	// check images
	images := image.Images()
	assert.Len(t, images, 1, "supposed to be single image")
	assert.EqualValues(t, simple, images[0], "should return image it was created with")
}

func TestImageNil(t *testing.T) {
	var image *image.Image
	_, err := NewImage(image, kube.ObjectMeta{}, "")
	assert.Error(t, err, "cannot be nil")
}

func TestImageInvalid(t *testing.T) {
	image := new(image.Image)
	_, err := NewImage(image, kube.ObjectMeta{}, "")
	assert.Error(t, err, "not valid")
}

func TestImageAttach(t *testing.T) {
	a := testNewImage(t, "a", kube.ObjectMeta{}, "", testRandomObjects(30))
	b := testNewImage(t, "b", kube.ObjectMeta{}, "", testRandomObjects(30))

	err := a.Attach(b)
	assert.Error(t, err, "Nothing can attach to images")
}

func TestImageType(t *testing.T) {
	image := newDockerImage(t, "ghost:latest")

	entity, err := NewImage(image, kube.ObjectMeta{}, "")
	if err != nil {
		t.Fatalf("Could not create Image entity: %v", err)
	}

	assert.Equal(t, EntityImage, entity.Type(), "incorrect type")
}

func TestImageKube(t *testing.T) {
	imageName := "redis:latest"
	image := newDockerImage(t, imageName)

	entity, err := NewImage(image, kube.ObjectMeta{}, "")
	if err != nil {
		t.Fatalf("Could not create Image entity: %v", err)
	}

	actual := entity.kube()
	assert.Equal(t, imageName, actual, "image names should match")
}

func TestImageBadObject(t *testing.T) {
	imageName := "debian"
	image := newDockerImage(t, imageName)

	service := kube.Service{}

	_, err := NewImage(image, kube.ObjectMeta{}, "", &service)
	assert.Error(t, err, "invalid object, should return error")
}

func testNewImage(t *testing.T, imageName string, defaults kube.ObjectMeta, source string, objects []deploy.KubeObject) *Image {
	image, err := NewImage(newDockerImage(t, imageName), defaults, source, objects...)
	if err != nil {
		t.Fatalf("Could not create Image: %v", err)
	}

	return image
}

func newDockerImage(t *testing.T, imageName string) *image.Image {
	simple, err := image.FromString(imageName)
	if err != nil {
		t.Fatalf("Could not create image.Image: %v", err)
	}
	return simple
}
