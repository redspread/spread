package component

import (
	"math/rand"
	"testing"
	"time"

	"rsprd.com/spread/pkg/deploy"

	"github.com/gh/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestNewBase(t *testing.T) {
	componentType := Type(rand.Intn(5))
	source := randomString(8)
	var objects []deploy.KubeObject

	base, err := newBase(componentType, api.ObjectMeta{}, source, objects)
	assert.NoError(t, err, "valid component")

	assert.Equal(t, componentType, base.Type(), "type cannot change")
	assert.Equal(t, source, base.Source(), "source cannot change")

	emptyDeploy := deploy.Deployment{}
	assert.True(t, emptyDeploy.Equals(base.objects))
}

func TestBaseBadObject(t *testing.T) {
	componentType := ComponentImage
	source := "testSource"
	objects := []deploy.KubeObject{
		&api.Pod{}, // invalid object
	}

	_, err := newBase(componentType, api.ObjectMeta{}, source, objects)
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
