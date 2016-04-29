package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	kube "k8s.io/kubernetes/pkg/api"
)

func TestAppEmpty(t *testing.T) {
	app, err := NewApp(nil, kube.ObjectMeta{}, "valid-app")
	assert.NoError(t, err, "empty apps are allowed")
	assert.Equal(t, 0, len(app.Images()))

	deployment, err := app.Deployment()
	assert.NoError(t, err)
	assert.Equal(t, 0, deployment.Len())
}

func TestAppAttachImage(t *testing.T) {
	numObj := 10

	// create Image
	imageObjects := testRandomObjects(numObj)
	imageEntity := testNewImage(t, "arch", kube.ObjectMeta{}, "", imageObjects)

	app, err := NewApp(nil, kube.ObjectMeta{}, "valid-app")
	assert.NoError(t, err, "empty apps are allowed")

	err = app.Attach(imageEntity)
	assert.NoError(t, err, "any entity can attach to an app")

	// check images
	assert.Equal(t, 1, len(app.Images()))

	// check deployment
	deployment, err := app.Deployment()
	assert.NoError(t, err)
	assert.Equal(t, numObj+1, deployment.Len())
}
