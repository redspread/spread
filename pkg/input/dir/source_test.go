package dir

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"rsprd.com/spread/pkg/deploy"
	"rsprd.com/spread/pkg/entity"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

var TestUsedNames = map[string]bool{}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

const (
	SpreadTestDir = "spread-test"
	TestFilePerms = 0777
)

func TestSourceObjectsNoKubeDir(t *testing.T) {
	fs := testTempFileSource(t)
	defer os.RemoveAll(string(fs))

	_, err := fs.Objects()
	assert.NoError(t, err, "ok")
}

func TestSourceObjectsEmptyKubeDir(t *testing.T) {
	fs := testTempFileSource(t)
	defer os.RemoveAll(string(fs))

	kubeDir := path.Join(string(fs), ObjectsDir)
	err := os.Mkdir(kubeDir, 0777)
	if err != nil {
		t.Fatal(err)
	}

	objects, err := fs.Objects()
	assert.NoError(t, err, "should be okay")
	assert.Len(t, objects, 0, "should not have any objects")
}

func TestSourceObjectsKubeDir(t *testing.T) {
	fs := testTempFileSource(t)
	defer os.RemoveAll(string(fs))

	kubeDir := path.Join(string(fs), ObjectsDir)
	err := os.Mkdir(kubeDir, TestFilePerms)
	if err != nil {
		t.Fatal(err)
	}

	numObjects := 5
	expected := testRandomObjects(numObjects)
	for _, v := range expected {
		filename := path.Join(string(fs), ObjectsDir, v.GetObjectMeta().GetName()+".json")
		testWriteYAMLToFile(t, filename, v)
	}

	actual, err := fs.Objects()
	assert.NoError(t, err)
	assert.Len(t, actual, numObjects, "different number of objects than created")
	for _, expectedObj := range expected {
		found := false
		for _, actualObj := range actual {
			if expectedObj.GetObjectMeta().GetName() == actualObj.GetObjectMeta().GetName() {
				testClearTypeInfo(expectedObj)
				found = kube.Semantic.DeepEqual(expectedObj, actualObj)
				break
			}
		}
		assert.True(t, found, "should have this object")
	}
}

// TODO: Add tests for entities in wrong folders

func TestSourceInvalidEntity(t *testing.T) {
	fs := testTempFileSource(t)
	defer os.RemoveAll(string(fs))

	// create an invalid entity.Type
	invalidEntityType := entity.Type(42)
	_, err := fs.Entities(invalidEntityType)
	assert.EqualError(t, err, ErrInvalidType.Error(), "should error for invalid entity types")
}

func TestSourceEntitiesNoFile(t *testing.T) {
	fs := testTempFileSource(t)
	defer os.RemoveAll(string(fs))

	rcs, err := fs.Entities(entity.EntityReplicationController)
	assert.NoError(t, err)
	assert.Len(t, rcs, 0)

	pods, err := fs.Entities(entity.EntityPod)
	assert.NoError(t, err)
	assert.Len(t, pods, 0)

	containers, err := fs.Entities(entity.EntityContainer)
	assert.NoError(t, err)
	assert.Len(t, containers, 0)

	images, err := fs.Entities(entity.EntityImage)
	assert.NoError(t, err)
	assert.Len(t, images, 0)
}

func TestSourceEntitiesEmptyFile(t *testing.T) {
	fs := testTempFileSource(t)
	defer os.RemoveAll(string(fs))

	entityFiles := []string{
		path.Join(string(fs), RCFile),
		path.Join(string(fs), PodFile),
	}

	// create files
	for _, file := range entityFiles {
		_, err := os.Create(file)
		if err != nil {
			t.Fatal(err)
		}
	}

	rcs, err := fs.Entities(entity.EntityReplicationController)
	assert.NoError(t, err, "should be okay")
	assert.Len(t, rcs, 0, "should not have any rcs")

	pods, err := fs.Entities(entity.EntityPod)
	assert.NoError(t, err, "should be okay")
	assert.Len(t, pods, 0, "should not have any pods")

	containers, err := fs.Entities(entity.EntityContainer)
	assert.NoError(t, err, "should be okay")
	assert.Len(t, containers, 0, "should not have any containers")

	images, err := fs.Entities(entity.EntityImage)
	assert.NoError(t, err, "should be okay")
	assert.Len(t, images, 0, "should not have any images")
}

func TestSourceRCs(t *testing.T) {
	fs := testTempFileSource(t)
	defer os.RemoveAll(string(fs))

	rcFile := path.Join(string(fs), RCFile)

	selector := map[string]string{"app": "example"}

	terminationPeriod := int64(30)
	kubeRC := &kube.ReplicationController{
		ObjectMeta: kube.ObjectMeta{
			Name:      "example-rc",
			Namespace: kube.NamespaceDefault,
		},
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ReplicationController",
			APIVersion: "v1",
		},
		Spec: kube.ReplicationControllerSpec{
			Selector: selector,
			Replicas: 2,
			Template: &kube.PodTemplateSpec{
				ObjectMeta: kube.ObjectMeta{
					Labels: selector,
				},
				Spec: kube.PodSpec{
					Containers: []kube.Container{
						kube.Container{
							Name:                   "example",
							Image:                  "hello-world",
							ImagePullPolicy:        kube.PullAlways,
							TerminationMessagePath: kube.TerminationMessagePathDefault,
						},
						kube.Container{
							Name:                   "cache",
							Image:                  "redis",
							ImagePullPolicy:        kube.PullAlways,
							TerminationMessagePath: kube.TerminationMessagePathDefault,
						},
					},
					SecurityContext:               &kube.PodSecurityContext{},
					RestartPolicy:                 kube.RestartPolicyAlways,
					DNSPolicy:                     kube.DNSDefault,
					TerminationGracePeriodSeconds: &terminationPeriod,
				},
			},
		},
	}

	testWriteYAMLToFile(t, rcFile, kubeRC)

	rcs, err := fs.Entities(entity.EntityReplicationController)
	assert.NoError(t, err, "should be okay")
	assert.Len(t, rcs, 1, "should have single rc")

	expectedKubeRC := kubeRC
	testClearTypeInfo(expectedKubeRC)
	expectedKubeRC.Labels = selector
	expected, err := entity.NewReplicationController(expectedKubeRC, kube.ObjectMeta{}, rcFile)

	actual := rcs[0]

	testCompareEntity(t, expected, actual)
}

func TestSourcePods(t *testing.T) {
	fs := testTempFileSource(t)
	defer os.RemoveAll(string(fs))

	podFile := path.Join(string(fs), PodFile)

	terminationPeriod := int64(30)
	kubePod := &kube.Pod{
		ObjectMeta: kube.ObjectMeta{
			Name:      "example-pod",
			Namespace: kube.NamespaceDefault,
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
				kube.Container{
					Name:                   "db",
					Image:                  "postgres",
					ImagePullPolicy:        kube.PullAlways,
					TerminationMessagePath: kube.TerminationMessagePathDefault,
				},
			},
			SecurityContext:               &kube.PodSecurityContext{},
			RestartPolicy:                 kube.RestartPolicyAlways,
			DNSPolicy:                     kube.DNSDefault,
			TerminationGracePeriodSeconds: &terminationPeriod,
		},
	}

	testWriteYAMLToFile(t, podFile, kubePod)

	pods, err := fs.Entities(entity.EntityPod)
	assert.NoError(t, err, "should be okay")
	assert.Len(t, pods, 1, "should have single pod")

	expectedKubePod := kubePod
	testClearTypeInfo(expectedKubePod)

	expected, err := entity.NewPod(expectedKubePod, kube.ObjectMeta{}, podFile)

	actual := pods[0]

	testCompareEntity(t, expected, actual)
}

func TestSourceContainersDir(t *testing.T) {
	fs := testTempFileSource(t)
	defer os.RemoveAll(string(fs))

	expected := []*entity.Container{}

	for i := 0; i < 15+rand.Intn(100); i++ {
		name := randomString(8)

		kubeContainer := kube.Container{
			Name:            name,
			Image:           randomString(6),
			ImagePullPolicy: kube.PullAlways,
		}
		filename := path.Join(string(fs), name+ContainerExtension)

		container, err := entity.NewContainer(kubeContainer, kube.ObjectMeta{}, filename)
		assert.NoError(t, err)

		expected = append(expected, container)

		testWriteYAMLToFile(t, filename, &kubeContainer)
	}

	actual, err := fs.Entities(entity.EntityContainer)
	assert.NoError(t, err, "should be okay")
	assert.Len(t, actual, len(expected), "should have same number of containers as start")

	for _, expectedContainer := range expected {
		found := false
		for _, actualContainer := range actual {
			if expectedContainer.Source() == actualContainer.Source() {
				found = testCompareEntity(t, expectedContainer, actualContainer)
				if found {
					break
				}
			}
		}
		assert.True(t, found, "should have found container")
	}
}

// TODO: Add test for file

func testWriteYAMLToFile(t *testing.T, filename string, typ interface{}) {
	jsonBytes, err := yaml.Marshal(typ)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filename, jsonBytes, TestFilePerms)
	if err != nil {
		t.Fatal(err)
	}
}

func testTempFileSource(t *testing.T) FileSource {
	dir := testTempDir(t)
	return FileSource(dir)
}

func testTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", SpreadTestDir)
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func createSecret(name string) *kube.Secret {
	return &kube.Secret{
		ObjectMeta: kube.ObjectMeta{Name: name},
		TypeMeta:   unversioned.TypeMeta{Kind: "Secret", APIVersion: "v1"},

		Type: kube.SecretTypeOpaque,
		Data: map[string][]byte{randomString(10): []byte(randomString(80))},
	}
}

// testRandomObjects returns a slice of randomly generated objects. If it is called with an object
// count of < 0 a random number of slices (with an upper bound of 100) will be generated.
func testRandomObjects(num int) (objects []deploy.KubeObject) {
	if num < 0 {
		num = rand.Intn(100)
	}

	for i := 0; i < num; i++ {
		// TODO: create different types of objects
		name := ""
		for {
			name = randomString(10)
			if !TestUsedNames[name] {
				break
			}
		}
		TestUsedNames[name] = true
		objects = append(objects, createSecret(name))
	}
	return
}

func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func testClearTypeInfo(obj deploy.KubeObject) {
	obj.GetObjectKind().SetGroupVersionKind(nil)
}

func testCompareEntity(t *testing.T, expected, actual entity.Entity) bool {
	assert.Equal(t, expected.Source(), actual.Source(), "spurces should match")
	assert.Equal(t, expected.DefaultMeta(), actual.DefaultMeta())

	actualImages := actual.Images()
	expectedImages := expected.Images()
	assert.Len(t, actualImages, len(expectedImages))
	for _, expectImage := range expectedImages {
		found := false
		for _, actualImage := range actualImages {
			found = expectImage.Equal(actualImage)
			if found {
				break
			}
		}
		assert.True(t, found, "should have found image")
	}

	expectedDeploy, err := expected.Deployment()
	assert.NoError(t, err)
	actualDeploy, err := actual.Deployment()
	assert.NoError(t, err)

	if !assert.True(t, expectedDeploy.Equal(actualDeploy)) {
		t.Log(expectedDeploy.Diff(actualDeploy))
		return false
	}
	return true
}
