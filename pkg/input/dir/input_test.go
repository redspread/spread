package dir

import (
	"fmt"
	"os"
	"path"
	"testing"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/entity"

	"github.com/stretchr/testify/assert"
	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func TestInputPod(t *testing.T) {
	input := testTempFileInput(t)
	defer os.RemoveAll(input.Path())

	kubePod := testKubePod()

	podFile := fmt.Sprintf("test.%s.%s", PodExtension, testRandomExtension())
	podPath := path.Join(input.Path(), podFile)
	testWriteYAMLToFile(t, podPath, kubePod)

	testClearTypeInfo(kubePod)

	objects := testWriteRandomObjects(t, input.Path(), 5)

	expectedPod, err := entity.NewPod(kubePod, kube.ObjectMeta{}, podPath)
	assert.NoError(t, err)
	expected, err := entity.NewApp([]entity.Entity{expectedPod}, kube.ObjectMeta{}, input.Path(), objects...)
	assert.NoError(t, err)

	actual, err := input.Build()
	assert.NoError(t, err, "should have built entity successfully")

	testCompareEntity(t, expected, actual)
}

func TestInputRCAndPod(t *testing.T) {
	input := testTempFileInput(t)
	defer os.RemoveAll(input.Path())

	// setup pod
	kubePod := testKubePod()

	podFile := fmt.Sprintf("test.%s.%s", PodExtension, testRandomExtension())
	podPath := path.Join(input.Path(), podFile)
	testWriteYAMLToFile(t, podPath, kubePod)
	testClearTypeInfo(kubePod)
	pod, err := entity.NewPod(kubePod, kube.ObjectMeta{}, podPath)
	assert.NoError(t, err)

	// setup rc
	objects := testWriteRandomObjects(t, input.Path(), 5)
	kubeRC := &kube.ReplicationController{
		ObjectMeta: kube.ObjectMeta{
			Name:      "spread",
			Namespace: kube.NamespaceDefault,
		},
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ReplicationController",
			APIVersion: "v1",
		},
		Spec: kube.ReplicationControllerSpec{
			Selector: map[string]string{"valid": "selector"},
			Template: &kube.PodTemplateSpec{
				Spec: kubePod.Spec,
			},
		},
	}

	rcFile := fmt.Sprintf("test.%s.%s", RCExtension, testRandomExtension())
	rcPath := path.Join(input.Path(), rcFile)
	testWriteYAMLToFile(t, rcPath, kubeRC)
	testClearTypeInfo(kubeRC)
	rc, err := entity.NewReplicationController(kubeRC, kube.ObjectMeta{}, rcPath)
	assert.NoError(t, err)

	expected, err := entity.NewApp(nil, kube.ObjectMeta{}, input.Path(), objects...)
	assert.NoError(t, err)

	// attach rc and pod
	err = expected.Attach(rc)
	assert.NoError(t, err)
	err = expected.Attach(pod)
	assert.NoError(t, err)

	// generate entity
	actual, err := input.Build()
	assert.NoError(t, err, "should have built entity successfully")

	testCompareEntity(t, expected, actual)
}

func testTempFileInput(t *testing.T) *FileInput {
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

func testKubePod() *kube.Pod {
	terminationPeriod := int64(30)
	return &kube.Pod{
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
				{
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
}

func testWriteAndAttachRandomContainers(t *testing.T, num int, dir string, parent entity.Entity) {
	for i := 0; i < num; i++ {
		name := randomString(8)

		kubeContainer := kube.Container{
			Name:                   name,
			Image:                  randomString(6),
			ImagePullPolicy:        kube.PullAlways,
			TerminationMessagePath: kube.TerminationMessagePathDefault,
		}
		filename := path.Join(dir, name+ContainerExtension)

		container, err := entity.NewContainer(kubeContainer, kube.ObjectMeta{}, filename)
		assert.NoError(t, err)

		testWriteYAMLToFile(t, filename, &kubeContainer)

		assert.NoError(t, parent.Attach(container), "should be able to attach container to pod")
	}
}
