package entity

import (
	"testing"

	"github.com/gh/stretchr/testify/assert"
	kube "k8s.io/kubernetes/pkg/api"
)

func TestRCNil(t *testing.T) {
	_, err := NewReplicationController(nil, kube.ObjectMeta{}, "nilRC")
	assert.Error(t, err, "RC cannot be created from ni")
}

func TestRCInvalid(t *testing.T) {
	kubeRC := new(kube.ReplicationController)
	_, err := NewReplicationController(kubeRC, kube.ObjectMeta{}, "")
	assert.Error(t, err, "invalid rc")
}

func TestRCNoPod(t *testing.T) {
	kubeRC := testNewKubeRC("no-pod", nil)
	kubeRC.Spec.Selector = map[string]string{
		"service": "jams",
	}

	rc, err := NewReplicationController(kubeRC, kube.ObjectMeta{}, "")
	assert.NoError(t, err)

	// internals
	assert.Nil(t, rc.pod, "pod should be nil")

	// images
	images := rc.Images()
	assert.Len(t, images, 0, "no image should have been created")

	// deployment
	_, err = rc.Deployment()
	assert.Error(t, err, "does not have pod, cannot deploy")
}

func TestRCInvalidPod(t *testing.T) {
	template := &kube.PodTemplateSpec{
		Spec: kube.PodSpec{},
	}
	kubeRC := testNewKubeRC("invalid-pod", template)
	_, err := NewReplicationController(kubeRC, kube.ObjectMeta{}, "")
	assert.Error(t, err, "should be invalid")
}

func testNewKubeRC(name string, spec *kube.PodTemplateSpec) *kube.ReplicationController {
	return &kube.ReplicationController{
		ObjectMeta: kube.ObjectMeta{
			Name:      name,
			Namespace: kube.NamespaceDefault,
		},
		Spec: kube.ReplicationControllerSpec{
			Template: spec,
		},
	}
}
