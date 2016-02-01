package entity

import (
	"strings"
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

func TestRCBadObjects(t *testing.T) {
	objects := []deploy.KubeObject{
		nil, // illegal
	}

	selector := map[string]string{"required": "value"}
	kubePod := testNewKubePod("testPod")
	kubePod.Labels = selector
	kubeRC := testNewKubeRC(kube.ObjectMeta{}, selector, kubePod)
	_, err := NewReplicationController(kubeRC, kube.ObjectMeta{}, "", objects...)
	assert.Error(t, err, "bad objects")
}

func TestRCAttachImage(t *testing.T) {
	imageName := "rc-attach-image"
	selector := map[string]string{
		"app": "cache",
	}

	// create kube.ReplicationController
	// create ReplicationController
	rcMeta := kube.ObjectMeta{Name: "test-rc"}
	kubeRC := testNewKubeRC(rcMeta, selector, nil)
	rcObjects := testRandomObjects(15)
	rc, err := NewReplicationController(kubeRC, kube.ObjectMeta{}, "", rcObjects...)
	assert.NoError(t, err, "should be valid RC")

	// create Image
	imageObjects := testRandomObjects(10)
	image := testNewImage(t, imageName, kube.ObjectMeta{}, "", imageObjects)

	// Attach image to RC
	// Should assume defaults up tree creating necessary components
	err = rc.Attach(image)
	assert.NoError(t, err, "attachment should be allowed")

	// Compare internal elements
	assert.NotNil(t, rc.rc.Spec.Template, "should of created template")

	// Create struct representation for expected
	rcMeta.Namespace = kube.NamespaceDefault
	containerName := strings.Join([]string{imageName, "container"}, "-")
	expectedRC := &kube.ReplicationController{
		ObjectMeta: rcMeta,
		Spec: kube.ReplicationControllerSpec{
			Selector: selector,
			Template: &kube.PodTemplateSpec{
				ObjectMeta: kube.ObjectMeta{Labels: selector},
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
			},
		},
	}

	// Insert into Deployment
	expected := deploy.Deployment{}
	err = expected.Add(expectedRC)
	assert.NoError(t, err, "should be valid RC")

	// add objects to deployment
	expected.AddDeployment(image.objects)
	expected.AddDeployment(rc.objects)

	// Create Deployment from RC
	actual, err := rc.Deployment()
	assert.NoError(t, err, "should produce valid deployment")

	// Compare deployments
	equal := expected.Equal(actual)
	assert.True(t, equal, "deployments should be same")

	// check images
	images := rc.Images()
	assert.Len(t, images, 1)
	for _, v := range images {
		assert.EqualValues(t, image, v, "should match original image")
	}
}

func TestRCAttachContainer(t *testing.T) {
	// create kube.ReplicationController
	// create ReplicationController

	// create kube.Container
	// create Container from created container

	// Attach container to RC
	// Should assume defaults up tree creating necessary components

	// Compare internal elements

	// Create struct representation for expected
	// Insert into Deployment
	// Create Deployment from RC

	// Compare Len of Deployments
	// Deep-equals deployments
}

func TestRCAttachPod(t *testing.T) {
	// create kube.ReplicationController
	// create ReplicationController

	// create kube.Pod
	// create Pod from created pod

	// Attach pod to RC
	// Should assume defaults up tree creating necessary components

	// Compare internal elements

	// Create struct representation for expected
	// Insert into Deployment
	// Create Deployment from RC

	// Compare Len of Deployments
	// Deep-equals deployments
}

func testNewKubeRC(meta kube.ObjectMeta, selector map[string]string, pod *kube.Pod) *kube.ReplicationController {
	var spec *kube.PodTemplateSpec
	if pod != nil {
		spec = &kube.PodTemplateSpec{
			ObjectMeta: pod.ObjectMeta,
			Spec:       pod.Spec,
		}
	}
	return &kube.ReplicationController{
		ObjectMeta: meta,
		Spec: kube.ReplicationControllerSpec{
			Selector: selector,
			Template: spec,
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
