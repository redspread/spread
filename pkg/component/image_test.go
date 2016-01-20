package component

import (
	"reflect"
	"testing"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

func TestImageImages(t *testing.T) {
	imageName := "gcr.io/google_containers/cassandra:v7"
	simple := newImage(t, imageName)

	image, err := NewImage(simple, "test")
	if err != nil {
		t.Fatalf("Could not create Image component: %v", err)
	}

	// check images
	images := image.Images()
	assert.Len(t, images, 1, "supposed to be single image")
	assert.EqualValues(t, simple, images[0], "should return image it was created with")
}

func TestImageDeployment(t *testing.T) {
	imageName := "arch"
	simple := newImage(t, imageName)

	image, err := NewImage(simple, "test")
	if err != nil {
		t.Fatalf("Could not create Image component: %v", err)
	}

	deployName := "image"

	expectedPod := api.Pod{
		ObjectMeta: api.ObjectMeta{
			Name:      deployName,
			Namespace: "default",
		},
		Spec: api.PodSpec{
			Containers: []api.Container{
				api.Container{
					Name:            "container",
					Image:           imageName,
					ImagePullPolicy: api.PullAlways,
				},
			},
			RestartPolicy: api.RestartPolicyAlways,
			DNSPolicy:     api.DNSClusterFirst,
		},
	}

	expected := deploy.NewDeployment()
	assert.NoError(t, expected.Add(&expectedPod), "should be able to add pod")

	actual := image.Deployment()
	t.Logf("\n%+v\n%+v", actual, expected)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %#v, saw: %#v", expected, actual)
	}
}

func TestNilImage(t *testing.T) {
	var image *image.Image
	_, err := NewImage(image, "")
	assert.Error(t, err, "cannot be nil")
}

func TestImageType(t *testing.T) {
	image := newImage(t, "ghost:latest")

	component, err := NewImage(image, "")
	if err != nil {
		t.Fatalf("Could not create Image component: %v", err)
	}

	assert.Equal(t, ComponentImage, component.Type(), "incorrect type")
}

func TestImageKube(t *testing.T) {
	imageName := "redis:latest"
	image := newImage(t, imageName)

	component, err := NewImage(image, "")
	if err != nil {
		t.Fatalf("Could not create Image component: %v", err)
	}

	actual := component.kube()
	assert.Equal(t, imageName, actual, "image names should match")
}

func TestImageBadObject(t *testing.T) {
	imageName := "debian"
	image := newImage(t, imageName)

	service := api.Service{}

	_, err := NewImage(image, "", &service)
	assert.Error(t, err, "invalid object, should return error")
}

func newImage(t *testing.T, imageName string) *image.Image {
	simple, err := image.FromString(imageName)
	if err != nil {
		t.Fatalf("Could not create image.Image: %v", err)
	}
	return simple
}
