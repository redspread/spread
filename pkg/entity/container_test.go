package entity

import (
	"testing"

	"rsprd.com/spread/pkg/deploy"

	"github.com/gh/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

func TestContainerSimpleDeployment(t *testing.T) {
	kubeContainer := api.Container{
		Name:            "simple-container",
		Image:           "busybox:latest",
		Command:         []string{"/bin/busybox", "ls"},
		ImagePullPolicy: api.PullAlways,
	}

	ctr, err := NewContainer(kubeContainer, api.ObjectMeta{}, "simpleTest")
	assert.NoError(t, err, "should be able to create container")

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
	actual := ctr.Deployment()

	assert.True(t, expected.Equals(actual), "should be equivlant")
}

func newKubeContainer(name, imageName string) api.Container {
	return api.Container{
		Name:            name,
		Image:           imageName,
		Command:         []string{"/bin/busybox", "ls"},
		ImagePullPolicy: api.PullAlways,
	}
}
