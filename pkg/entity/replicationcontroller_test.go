package entity

import (
	"testing"

	"rsprd.com/spread/pkg/deploy"

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
	kubeRC := kube.ReplicationController{
		ObjectMeta: kube.ObjectMeta{
			Name: "no-pod",
		},
		Spec: kube.ReplicationControllerSpec{
			Selector: map[string]string{"service": "jams"},
			Template: nil,
		},
	}

	rc, err := NewReplicationController(&kubeRC, kube.ObjectMeta{}, "")
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
	kubeRC := testNewKubeRC(kube.ObjectMeta{Name: "invalid-pod"}, nil, &kube.Pod{})
	_, err := NewReplicationController(kubeRC, kube.ObjectMeta{}, "")
	assert.Error(t, err, "should be invalid")
}

func TestRCValidPodImages(t *testing.T) {
	selector := map[string]string{
		"service": "hotline",
	}
	kubeRC := testNewRC(t, selector)

	images := kubeRC.Images()
	assert.Len(t, images, 2, "RC should have 2 images")
}

func TestRCValidPodDeployment(t *testing.T) {
	selector := map[string]string{
		"service": "postgres",
	}
	kubePod := testCreateKubePodSourcegraph("sourcegraph")
	kubePod.Labels = selector
	kubeRC := testNewKubeRC(kube.ObjectMeta{Name: "sourcegraph-rc"}, selector, kubePod)

	rc, err := NewReplicationController(kubeRC, kube.ObjectMeta{}, "")
	assert.NoError(t, err)
	assert.NotNil(t, rc.pod, "a pod should have be created")
	assert.Len(t, rc.pod.containers, 2, "two containers should have been created")

	expected := deploy.Deployment{}
	expected.Add(kubeRC)

	actual, err := rc.Deployment()
	assert.NoError(t, err)
	assert.True(t, expected.Equal(actual))
}

func testNewKubeRC(meta kube.ObjectMeta, selector map[string]string, pod *kube.Pod) *kube.ReplicationController {
	return &kube.ReplicationController{
		ObjectMeta: meta,
		Spec: kube.ReplicationControllerSpec{
			Selector: selector,
			Template: &kube.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			},
		},
	}
}

func testNewRC(t *testing.T, selector map[string]string) *ReplicationController {
	return testNewRCWithOpts(t, selector, kube.ObjectMeta{}, []deploy.KubeObject{})
}

func testNewRCWithOpts(t *testing.T, selector map[string]string, defaults kube.ObjectMeta, objects []deploy.KubeObject) *ReplicationController {
	pod := testCreateKubePodSourcegraph("sourcegraph")
	meta := kube.ObjectMeta{
		Name: "sourcegraph-rc",
	}
	kubeRC := testNewKubeRC(meta, selector, pod)
	kubeRC.Spec.Template.Labels = selector

	rc, err := NewReplicationController(kubeRC, kube.ObjectMeta{}, "valid-rc", objects...)
	assert.NoError(t, err)
	return rc
}
