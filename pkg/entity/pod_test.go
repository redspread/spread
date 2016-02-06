package entity

import (
	"testing"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/image"

	"rsprd.com/spread/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	kube "rsprd.com/spread/Godeps/_workspace/src/k8s.io/kubernetes/pkg/api"
)

func TestPodNil(t *testing.T) {
	_, err := NewPod(nil, kube.ObjectMeta{}, "nilPod")
	assert.Error(t, err, "cannot created entity.Pod from nil Kube Pod")
}

func TestPodInvalid(t *testing.T) {
	kubePod := new(kube.Pod)
	_, err := NewPod(kubePod, kube.ObjectMeta{}, "")
	assert.Error(t, err, "invalid pod")
}

func TestPodNoContainers(t *testing.T) {
	kubePod := testNewKubePod("no-containers")
	pod, err := NewPod(kubePod, kube.ObjectMeta{}, "no-containers")
	assert.NoError(t, err, "should be valid entity.Pod")

	// internals
	assert.Len(t, pod.containers, 0, "no containers should exist")

	// images
	images := pod.Images()
	assert.Len(t, images, 0, "no image should have been created")

	// kube
	_, _, err = pod.data()
	assert.Error(t, err, "doesn't have containers, no kube")

	// deployment
	_, err = pod.Deployment()
	assert.Error(t, err, "doesn't have containers, can't deploy")
}

func TestPodNoImage(t *testing.T) {
	kubePod := testNewKubePod("no-image-pod")
	kubePod.Spec.Containers = []kube.Container{
		kube.Container{
			Name:            "no-image",
			ImagePullPolicy: kube.PullIfNotPresent,
		},
	}

	pod, err := NewPod(kubePod, kube.ObjectMeta{}, "")
	assert.NoError(t, err, "imageless but valid pod")

	_, _, err = pod.data()
	assert.Error(t, err, "kube not possible without image")

	_, err = pod.Deployment()
	assert.Error(t, err, "deployment not possible without image")
}

func TestPodWithContainersImage(t *testing.T) {
	pod, err := NewPod(testCreateKubePodSourcegraph("test"), kube.ObjectMeta{}, "has-containers")
	assert.NoError(t, err, "valid pod")

	images := pod.Images()
	assert.Len(t, images, 2, "there should only be 2 images")

	imageNames := []string{}
	for _, image := range images {
		imageNames = append(imageNames, image.KubeImage())
	}

	assert.Contains(t, imageNames, testKubeContainerPostgres.Image)
	assert.Contains(t, imageNames, testKubeContainerSourcegraph.Image)
}

func TestPodWithContainersKube(t *testing.T) {
	pod, err := NewPod(testCreateKubePodSourcegraph("test"), kube.ObjectMeta{}, "with-containers")
	assert.NoError(t, err, "valid pod")

	// check internals
	assert.Len(t, pod.containers, 2, "should have postgres and sourcegraph containers")

	expected := testCreateKubePodSourcegraph("test")
	expected.Namespace = kube.NamespaceDefault

	actual, _, err := pod.data()
	assert.NoError(t, err, "should generate kube")

	assert.True(t, kube.Semantic.DeepEqual(expected, actual), "Expected: %+v, Actual: %+v", expected, actual)
}

func TestPodWithContainersDeployment(t *testing.T) {
	kubePod := testCreateKubePodSourcegraph("deploy")
	pod, err := NewPod(kubePod, kube.ObjectMeta{}, "containers-deploy")
	assert.NoError(t, err, "valid pod")

	kubePod.Namespace = kube.NamespaceDefault

	expected := new(deploy.Deployment)
	err = expected.Add(kubePod)
	assert.NoError(t, err, "valid pod")

	actual, err := pod.Deployment()
	assert.NoError(t, err, "deployment should be valid")

	testDeploymentEqual(t, expected, actual)
}

func TestPodBadObjects(t *testing.T) {
	objects := []deploy.KubeObject{
		nil, // illegal
	}
	_, err := NewPod(testNewKubePod("bad"), kube.ObjectMeta{}, "", objects...)
	assert.Error(t, err, "bad objects")
}

func TestPodFromPodSpec(t *testing.T) {
	spec := kube.PodSpec{
		RestartPolicy: kube.RestartPolicyAlways,
		DNSPolicy:     kube.DNSDefault,
	}
	_, err := NewPodFromPodSpec(kube.ObjectMeta{Name: "no-containers"}, spec, kube.ObjectMeta{}, "no-containers")
	assert.NoError(t, err, "should be valid entity.Pod")
}

func TestPodAttachImage(t *testing.T) {
	podObjects := testRandomObjects(3)
	kubePod := testNewKubePod("containerless")
	pod, err := NewPod(kubePod, kube.ObjectMeta{}, "pod", podObjects...)
	assert.NoError(t, err, "valid")

	imageObjects := testRandomObjects(3)
	kubeImage, err := image.FromString("bprashanth/nginxhttps:1.0")
	assert.NoError(t, err, "image should be valid")

	image, err := NewImage(kubeImage, kube.ObjectMeta{}, "image", imageObjects...)
	assert.NoError(t, err, "valid")

	err = pod.Attach(image)
	assert.NoError(t, err, "should be attachable")

	kubePod.Namespace = kube.NamespaceDefault
	kubePod.Spec.Containers = []kube.Container{
		kube.Container{
			Name:            "nginxhttps",
			Image:           "bprashanth/nginxhttps:1.0",
			ImagePullPolicy: kube.PullAlways,
		},
	}

	objects := append(podObjects, imageObjects...)

	expected := new(deploy.Deployment)
	err = expected.Add(kubePod)
	assert.NoError(t, err, "valid")

	for _, obj := range objects {
		assert.NoError(t, expected.Add(obj))
	}

	actual, err := pod.Deployment()
	assert.NoError(t, err, "deployment should be ok")
	testDeploymentEqual(t, expected, actual)
}

func TestPodAttachContainer(t *testing.T) {
	podObjects := testRandomObjects(60)
	kubePod := testNewKubePod("containerless")
	pod, err := NewPod(kubePod, kube.ObjectMeta{}, "pod", podObjects...)
	assert.NoError(t, err, "valid")

	containerObjects := testRandomObjects(20)

	kubeContainer := testNewKubeContainer("container", "busybox:latest")
	container, err := NewContainer(kubeContainer, kube.ObjectMeta{}, "container", containerObjects...)
	assert.NoError(t, err)

	err = pod.Attach(container)
	assert.NoError(t, err)
	children := pod.children()
	assert.Contains(t, children, container, "should have container as child")

	kubePod.Namespace = kube.NamespaceDefault
	kubePod.Spec.Containers = []kube.Container{
		kubeContainer,
	}

	expected := new(deploy.Deployment)
	err = expected.Add(kubePod)
	assert.NoError(t, err)

	actual, err := pod.Deployment()
	assert.NoError(t, err)

	objects := append(podObjects, containerObjects...)

	for _, obj := range objects {
		assert.NoError(t, expected.Add(obj))
	}

	testDeploymentEqual(t, expected, actual)
}

func testNewKubePod(name string) *kube.Pod {
	return &kube.Pod{
		ObjectMeta: kube.ObjectMeta{
			Name: name,
		},
		Spec: kube.PodSpec{
			RestartPolicy: kube.RestartPolicyAlways,
			DNSPolicy:     kube.DNSDefault,
		},
	}
}

func testCreateKubePodSourcegraph(name string) *kube.Pod {
	return &kube.Pod{
		ObjectMeta: kube.ObjectMeta{Name: name},
		Spec: kube.PodSpec{
			RestartPolicy: kube.RestartPolicyAlways,
			DNSPolicy:     kube.DNSDefault,
			Volumes: []kube.Volume{
				kube.Volume{
					Name: "config",
					VolumeSource: kube.VolumeSource{
						EmptyDir: &kube.EmptyDirVolumeSource{},
					},
				},
				kube.Volume{
					Name: "db",
					VolumeSource: kube.VolumeSource{
						EmptyDir: &kube.EmptyDirVolumeSource{},
					},
				},
			},
			Containers: []kube.Container{
				testKubeContainerSourcegraph,
				testKubeContainerPostgres,
			},
		},
	}
}
