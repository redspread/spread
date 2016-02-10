package dir

import (
	"os"
	"path"
	"testing"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/entity"

	"github.com/stretchr/testify/assert"
	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func TestInputContainersOnly(t *testing.T) {
	input := testTempFileInput(t)
	defer os.RemoveAll(input.Path())

	objects := testWriteRandomObjects(t, input.Path(), 5)

	expected, err := entity.NewDefaultPod(kube.ObjectMeta{GenerateName: "spread"}, input.Path(), objects...)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		name := randomString(8)

		kubeContainer := kube.Container{
			Name:            name,
			Image:           randomString(6),
			ImagePullPolicy: kube.PullAlways,
		}
		filename := path.Join(input.Path(), name+ContainerExtension)

		container, err := entity.NewContainer(kubeContainer, kube.ObjectMeta{}, filename)
		assert.NoError(t, err)

		testWriteYAMLToFile(t, filename, &kubeContainer)

		assert.NoError(t, expected.Attach(container), "should be able to attach container to pod")
	}

	actual, err := input.Build()
	assert.NoError(t, err, "should have built entity successfully")

	testCompareEntity(t, expected, actual)
}

func TestInputPodwithContainers(t *testing.T) {
	input := testTempFileInput(t)
	defer os.RemoveAll(input.Path())

	terminationPeriod := int64(30)
	kubePod := &kube.Pod{
		ObjectMeta: kube.ObjectMeta{
			GenerateName: "spread",
			Namespace:    kube.NamespaceDefault,
		},
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		Spec: kube.PodSpec{
			Containers: []kube.Container{
				kube.Container{
					Name:                   "wiki",
					Image:                  "mediawiki",
					ImagePullPolicy:        kube.PullAlways,
					TerminationMessagePath: kube.TerminationMessagePathDefault,
				},
			},
			RestartPolicy:                 kube.RestartPolicyAlways,
			DNSPolicy:                     kube.DNSDefault,
			TerminationGracePeriodSeconds: &terminationPeriod,
			SecurityContext:               &kube.PodSecurityContext{},
		},
	}

	podFile := path.Join(input.Path(), PodFile)
	testWriteYAMLToFile(t, podFile, kubePod)

	testClearTypeInfo(kubePod)

	objects := testWriteRandomObjects(t, input.Path(), 5)

	expected, err := entity.NewPod(kubePod, kube.ObjectMeta{}, podFile, objects...)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		name := randomString(8)

		kubeContainer := kube.Container{
			Name:                   name,
			Image:                  randomString(6),
			ImagePullPolicy:        kube.PullAlways,
			TerminationMessagePath: kube.TerminationMessagePathDefault,
		}
		filename := path.Join(input.Path(), name+ContainerExtension)

		container, err := entity.NewContainer(kubeContainer, kube.ObjectMeta{}, filename)
		assert.NoError(t, err)

		testWriteYAMLToFile(t, filename, &kubeContainer)

		assert.NoError(t, expected.Attach(container), "should be able to attach container to pod")
	}

	actual, err := input.Build()
	assert.NoError(t, err, "should have built entity successfully")

	testCompareEntity(t, expected, actual)
}

func testTempFileInput(t *testing.T) *fileInput {
	dir := testTempDir(t)
	input, err := NewFileInput(dir)
	if err != nil {
		t.Error(err)
	}
	return input
}

// testWriteRandomObjects randomly generates Kubernetes objects and writes them in the specified path.
// The created objects are returned with Type information clean. If object count is < 0, a random number will be used.
func testWriteRandomObjects(t *testing.T, dir string, count int) []deploy.KubeObject {
	objects := testRandomObjects(count)

	kubeDir := path.Join(dir, ObjectsDir)
	err := os.Mkdir(kubeDir, TestFilePerms)
	if err != nil {
		t.Fatal(err)
	}

	for _, obj := range objects {
		filename := path.Join(kubeDir, obj.GetObjectMeta().GetName()+".yml")
		testWriteYAMLToFile(t, filename, obj)

		// cleanup type information which is removed from decoded objects
		testClearTypeInfo(obj)
	}

	return objects
}
