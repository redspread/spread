package component

import (
	"reflect"
	"testing"

	"rsprd.com/spread/pkg/image"

	"k8s.io/kubernetes/pkg/api"
	"rsprd.com/spread/pkg/deploy"
)

func TestNewImage(t *testing.T) {
	imageName := "gcr.io/google_containers/cassandra:v7"
	simple, err := image.FromString(imageName)
	if err != nil {
		t.Fatalf("Could not create image.Image: %v", err)
	}

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
