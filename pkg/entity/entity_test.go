package entity

import (
	"math/rand"
	"testing"
	"time"

	"rsprd.com/spread/pkg/deploy"

	"github.com/stretchr/testify/assert"
	kube "k8s.io/kubernetes/pkg/api"
)

var TestUsedNames = map[string]bool{}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestBaseNew(t *testing.T) {
	entityType := Type(rand.Intn(5))
	source := randomString(8)
	var objects []deploy.KubeObject

	base, err := newBase(entityType, kube.ObjectMeta{}, source, objects)
	assert.NoError(t, err, "valid entity")

	assert.Equal(t, entityType, base.Type(), "type cannot change")
	assert.Equal(t, source, base.Source(), "source cannot change")

	emptyDeploy := deploy.Deployment{}
	assert.True(t, emptyDeploy.Equal(&base.objects))
}

func TestBaseNilObjects(t *testing.T) {
	objects := []deploy.KubeObject{
		createSecret("test"),
		nil, // illegal
	}
	_, err := newBase(EntityPod, kube.ObjectMeta{}, "source", objects)
	assert.Error(t, err, "should not be able to create base with nil components")
}

func TestBaseNoDefaults(t *testing.T) {
	defaults := kube.ObjectMeta{}
	obj := createSecret(randomString(8))

	base, err := newBase(EntityApplication, defaults, "src", []deploy.KubeObject{obj})
	assert.NoError(t, err, "valid base")
	assert.True(t, kube.Semantic.DeepEqual(defaults, base.DefaultMeta()), "defaults should have not changed")

	objects := base.Objects()
	assert.Len(t, objects, 1, "should only have secret")

	obj.GetObjectMeta().SetNamespace(kube.NamespaceDefault)

	actual := objects[0]
	assert.True(t, kube.Semantic.DeepEqual(obj, actual), "secrets should be same")
}

func TestBaseNamespaceDefaults(t *testing.T) {
	defaults := kube.ObjectMeta{
		Namespace: "set-by-defaults",
	}

	nsSetName := "namespace-set"
	nsSet := createSecret(nsSetName)
	nsSet.Namespace = "set-on-object"

	nsUnsetName := "namespace-unset"
	nsUnset := createSecret(nsUnsetName)

	objects := []deploy.KubeObject{nsSet, nsUnset}
	base, err := newBase(EntityReplicationController, defaults, "src", objects)
	assert.NoError(t, err, "valid base")

	for _, obj := range base.Objects() {
		meta := obj.GetObjectMeta()
		switch meta.GetName() {
		case nsSetName:
			assert.Equal(t, "set-on-object", meta.GetNamespace(), "object namespace should override defaults")
		case nsUnsetName:
			assert.Equal(t, "set-by-defaults", meta.GetNamespace(), "should use defaults for namespace")
		default:
			t.Errorf("unexpected object `%s`", meta.GetName())
		}
	}
}

func TestBaseDefaultGenerateName(t *testing.T) {
	defaults := kube.ObjectMeta{
		GenerateName: "inventory",
	}

	objects := []deploy.KubeObject{
		createSecret(""), // empty name set
	}

	base, err := newBase(EntityApplication, defaults, "src", objects)
	assert.NoError(t, err, "valid base")

	for _, obj := range base.Objects() {
		assert.Equal(t, defaults.GenerateName, obj.GetObjectMeta().GetGenerateName(), "generate name should have been set")
	}
}

func TestBaseDefaultAnnotationsAndLabels(t *testing.T) {
	initial := map[string]string{
		"overwritten":     "no",
		"not-overwritten": "yes",
	}

	defaults := kube.ObjectMeta{
		Labels:      initial,
		Annotations: initial,
	}

	override := map[string]string{
		"overwritten": "yes",
	}

	obj := createSecret(randomString(8))
	obj.GetObjectMeta().SetLabels(override)
	obj.GetObjectMeta().SetAnnotations(override)

	base, err := newBase(EntityContainer, defaults, "src", []deploy.KubeObject{obj})
	assert.NoError(t, err, "valid base")
	assert.True(t, kube.Semantic.DeepEqual(defaults, base.DefaultMeta()), "defaults should have not changed")

	objects := base.Objects()
	assert.Len(t, objects, 1)

	expected := initial
	expected["overwritten"] = "yes"

	output := objects[0]
	meta := output.GetObjectMeta()
	assert.Equal(t, expected, meta.GetLabels(), "labels should match")
	assert.Equal(t, expected, meta.GetAnnotations(), "annotations should match")
}

func TestBaseNoDefaultAnnotationsAndLabels(t *testing.T) {
	defaults := kube.ObjectMeta{}
	obj := createSecret("postgres")
	obj.Annotations = map[string]string{"object": "data"}
	obj.Labels = map[string]string{"data": "object"}

	base, err := newBase(EntityContainer, defaults, "src", []deploy.KubeObject{obj})
	assert.NoError(t, err)

	objects := base.Objects()
	assert.Len(t, objects, 1)

	expected := obj
	actual := objects[0]

	assert.Equal(t, expected.Annotations, actual.GetObjectMeta().GetAnnotations())
	assert.Equal(t, expected.Labels, actual.GetObjectMeta().GetLabels())
}

func TestBaseCheckAttach(t *testing.T) {
	baseImage := newDockerImage(t, "sample-image")
	image, err := NewImage(baseImage, kube.ObjectMeta{}, "")
	assert.NoError(t, err, "valid image")

	kubeContainer := testNewKubeContainer("sample-container", "golang")
	container, err := NewContainer(kubeContainer, kube.ObjectMeta{}, "")
	assert.NoError(t, err, "valid container")

	assert.Error(t, image.validAttach(container), "containers should not be allowed to attach to images")
	assert.NoError(t, container.validAttach(image), "images should be allowed to attach to containers")
}

func TestBaseBadObject(t *testing.T) {
	entityType := EntityImage
	source := "testSource"
	objects := []deploy.KubeObject{
		&kube.Pod{}, // invalid object
	}

	_, err := newBase(entityType, kube.ObjectMeta{}, source, objects)
	assert.Error(t, err, "objects are invalid")
}

func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func createSecret(name string) *kube.Secret {
	return &kube.Secret{
		ObjectMeta: kube.ObjectMeta{Name: name},

		Type: kube.SecretTypeOpaque,
		Data: map[string][]byte{randomString(10): []byte(randomString(80))},
	}
}

// testRandomObjects returns a slice of randomly generated objects. If it is called with an object
// count of 0, a random number of slices (with an upper bound of 100) will be generated.
func testRandomObjects(num int) (objects []deploy.KubeObject) {
	if num == 0 {
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

func testDeploymentEqual(t *testing.T, expected, actual *deploy.Deployment) bool {
	equal := expected.Equal(actual)
	return assert.True(t, equal, "expected: %s,\n actual: %s,\n diff:\n%s", expected, actual, expected.Diff(actual))
}
