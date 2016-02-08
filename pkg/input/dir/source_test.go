package dir

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"rsprd.com/spread/pkg/deploy"

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

func TestSourceNonexistentPath(t *testing.T) {
	doesNotExist := "/dev/null/null"
	_, err := NewFileSource(doesNotExist)
	assert.Error(t, err, "should not create for nonexistent path")
}

func TestSourceValidPath(t *testing.T) {
	exists := "/"
	_, err := NewFileSource(exists)
	assert.NoError(t, err, "valid path")

	relative := "."
	_, err = NewFileSource(relative)
	assert.NoError(t, err, "valid path")
}

func TestSourceObjectsNoKubeDir(t *testing.T) {
	fs := testTempFileSource(t)
	defer os.RemoveAll(string(fs))

	_, err := fs.Objects()
	assert.Error(t, err, "directory doesn't exist")
	if !strings.HasSuffix(err.Error(), "does not exist") {
		t.Error("should have failed from .kube not existing")
	}
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
		testWriteJSONToFile(t, filename, v)
	}

	actual, err := fs.Objects()
	assert.NoError(t, err)
	assert.Len(t, actual, numObjects, "different number of objects than created")
	for _, expectedObj := range expected {
		found := false
		for _, actualObj := range actual {
			if expectedObj.GetObjectMeta().GetName() == actualObj.GetObjectMeta().GetName() {
				found = true
				break
			}
		}
		assert.True(t, found, "should have this object")
	}
}

func testWriteJSONToFile(t *testing.T, filename string, typ interface{}) {
	jsonBytes, err := json.MarshalIndent(typ, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filename, jsonBytes, TestFilePerms)
	if err != nil {
		t.Fatal(err)
	}
}

func testTempFileSource(t *testing.T) fileSource {
	dir := testTempDir(t)
	fs, err := NewFileSource(dir)
	if err != nil {
		t.Error(err)
	}
	return fs
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
		TypeMeta:   unversioned.TypeMeta{Kind: "ReplicationController", APIVersion: "v1"},

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

func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
