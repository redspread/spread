package component

import (
	"reflect"
	"testing"

	"rsprd.com/spread/pkg/image"

	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
	"rsprd.com/spread/pkg/deploy"
)

func TestNewImage(t *testing.T) {
	imageName := "gcr.io/google_containers/cassandra:v7"
	simple := newImage(t, imageName)

	secret := api.Secret{
		ObjectMeta: api.ObjectMeta{
			Name:      "sosecret",
			Namespace: "default",
		},
		Type: api.SecretTypeOpaque,
		Data: map[string][]byte{"test": []byte("secret")},
	}

	image, err := NewImage(simple, "test", &secret)
	if err != nil {
		t.Fatalf("Could not create Image component: %v", err)
	}

	// check images
	images := image.Images()
	assert.Len(t, images, 1, "supposed to be single image")
	assert.EqualValues(t, simple, images[0], "should return image it was created with")

	expectedPod := api.Pod{
		Spec: api.PodSpec{
			Containers: []api.Container{
				api.Container{
					Image: imageName,
				},
			},
		},
	}

	expected := deploy.NewDeployment()
	expected.Add(&secret)
	expected.Add(&expectedPod)

	actual := image.Deployment()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %#v, saw: %#v", expected, actual)
	}
}

func TestNilImage(t *testing.T) {
	var image *image.Image
	_, err := NewImage(image, "")
	if err == nil {
		t.Errorf("Nil should return an error.")
	}
}

func TestImageType(t *testing.T) {
	image := newImage(t, "ghost:latest")

	component, err := NewImage(image, "")
	if err != nil {
		t.Fatalf("Could not create Image component: %v", err)
	}

	assert.Equal(t, ComponentImage, component.Type(), "incorrect type")
}

func newImage(t *testing.T, imageName string) *image.Image {
	simple, err := image.FromString(imageName)
	if err != nil {
		t.Fatalf("Could not create image.Image: %v", err)
	}
	return simple
}
