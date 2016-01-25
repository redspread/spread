package entity

import (
	"testing"

	"rsprd.com/spread/pkg/deploy"

	"github.com/gh/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

func TestPodNil(t *testing.T) {
	var pod *api.Pod
	_, err := NewPod(pod, api.ObjectMeta{}, "nilPod")
	assert.Error(t, err, "cannot created entity.Pod from nil Kube Pod")
}

func TestPodNoContainers(t *testing.T) {
	kubePod := testNewKubePod("no-containers")
	pod, err := NewPod(kubePod, api.ObjectMeta{}, "no-containers")
	assert.NoError(t, err, "should be valid entity.Pod")

	// internals
	assert.Len(t, pod.containers, 0, "no containers should exist")

	// images
	images := pod.Images()
	assert.Len(t, images, 0, "no image should have been created")

	// kube
	_, err = pod.kube()
	assert.Error(t, err, "doesn't have containers, no kube")

	// deployment
	_, err = pod.Deployment()
	assert.Error(t, err, "doesn't have containers, can't deploy")
}

func TestPodWithContainersImage(t *testing.T) {
	pod, err := NewPod(testCreateKubePodSourcegraph("test"), api.ObjectMeta{}, "has-containers")
	assert.NoError(t, err, "valid pod")

	images := pod.Images()
	assert.Len(t, images, 2, "there should only be 2 images")

	imageNames := []string{}
	for _, image := range images {
		imageNames = append(imageNames, image.DockerName())
	}

	assert.Contains(t, imageNames, testKubeContainerPostgres.Image)
	assert.Contains(t, imageNames, testKubeContainerSourcegraph.Image)
}

func TestPodWithContainersKube(t *testing.T) {
	pod, err := NewPod(testCreateKubePodSourcegraph("test"), api.ObjectMeta{}, "with-containers")
	assert.NoError(t, err, "valid pod")

	// check internals
	assert.Len(t, pod.containers, 2, "should have postgres and sourcegraph containers")

	expected := testCreateKubePodSourcegraph("test")
	expected.Namespace = api.NamespaceDefault

	actual, err := pod.kube()
	assert.NoError(t, err, "should generate kube")

	assert.True(t, api.Semantic.DeepEqual(expected, actual), "Expected: %+v, Actual: %+v", expected, actual)
}

func TestPodWithContainersDeployment(t *testing.T) {
	kubePod := testCreateKubePodSourcegraph("deploy")
	pod, err := NewPod(kubePod, api.ObjectMeta{}, "containers-deploy")
	assert.NoError(t, err, "valid pod")

	kubePod.Namespace = api.NamespaceDefault

	expected := deploy.Deployment{}
	err = expected.Add(kubePod)
	assert.NoError(t, err, "valid pod")

	actual, err := pod.Deployment()
	assert.NoError(t, err, "deployment should be valid")

	assert.True(t, expected.Equal(actual), "deployments should be the same")
}

func TestPodBadObjects(t *testing.T) {
	objects := []deploy.KubeObject{
		nil, // illegal
	}
	_, err := NewPod(testNewKubePod("bad"), api.ObjectMeta{}, "", objects...)
	assert.Error(t, err, "bad objects")
}

func TestPodFromPodSpec(t *testing.T) {
	spec := api.PodSpec{
		RestartPolicy: api.RestartPolicyAlways,
		DNSPolicy:     api.DNSDefault,
	}
	_, err := NewPodFromPodSpec("no-containers", spec, api.ObjectMeta{}, "no-containers")
	assert.NoError(t, err, "should be valid entity.Pod")
}

func testNewKubePod(name string) *api.Pod {
	return &api.Pod{
		ObjectMeta: api.ObjectMeta{
			Name: name,
		},
		Spec: api.PodSpec{
			RestartPolicy: api.RestartPolicyAlways,
			DNSPolicy:     api.DNSDefault,
		},
	}
}

func testCreateKubePodSourcegraph(name string) *api.Pod {
	return &api.Pod{
		ObjectMeta: api.ObjectMeta{Name: name},
		Spec: api.PodSpec{
			RestartPolicy: api.RestartPolicyAlways,
			DNSPolicy:     api.DNSDefault,
			Volumes: []api.Volume{
				api.Volume{
					Name: "config",
					VolumeSource: api.VolumeSource{
						EmptyDir: &api.EmptyDirVolumeSource{},
					},
				},
				api.Volume{
					Name: "db",
					VolumeSource: api.VolumeSource{
						EmptyDir: &api.EmptyDirVolumeSource{},
					},
				},
			},
			Containers: []api.Container{
				testKubeContainerSourcegraph,
				testKubeContainerPostgres,
			},
		},
	}
}
