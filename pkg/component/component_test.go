package component

import (
	"math/rand"
	"testing"

	"rsprd.com/spread/pkg/deploy"

	"github.com/gh/stretchr/testify/assert"
	"github.com/google/gofuzz"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/testapi"
	kubetest "k8s.io/kubernetes/pkg/api/testing"
)

func TestNewBase(t *testing.T) {
	componentType := Type(rand.Intn(5))
	var source string
	var objects []deploy.KubeObject

	fuzzer(t).Fuzz(&source)

	base, err := newBase(componentType, source, objects)
	assert.NoError(t, err, "valid component")

	assert.Equal(t, componentType, base.Type(), "type cannot change")
	assert.Equal(t, source, base.Source(), "source cannot change")

	emptyDeploy := deploy.NewDeployment()
	assert.Equal(t, emptyDeploy, base.objects)
}

func TestBaseBadObject(t *testing.T) {
	componentType := ComponentImage
	source := "testSource"
	objects := []deploy.KubeObject{
		&api.Pod{}, // invalid object
	}

	_, err := newBase(componentType, source, objects)
	assert.Error(t, err, "objects are invalid")
}

func fuzzer(t *testing.T) *fuzz.Fuzzer {
	version := testapi.Default.InternalGroupVersion()
	seed := rand.Int63()
	return kubetest.FuzzerFor(t, version, rand.NewSource(seed))
}
