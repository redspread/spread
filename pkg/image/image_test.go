package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseName(t *testing.T) {
	name := "debian"

	// "debian"
	imageStr := name

	image, err := FromString(imageStr)
	assert.NoError(t, err, "valid image name")

	// Name
	assert.Equal(t, name, image.Name())

	// KubeImage
	assert.Equal(t, imageStr, image.KubeImage())
}

func TestParseNameTag(t *testing.T) {
	name := "debian"
	tag := "jessie"

	// "debian:jessie"
	imageStr := name + ":" + tag

	image, err := FromString(imageStr)
	assert.NoError(t, err, "valid image name")

	// Name
	assert.Equal(t, name, image.Name())

	// KubeImage
	assert.Equal(t, imageStr, image.KubeImage())
}

func TestParseNamespaceName(t *testing.T) {
	ns := "library"
	name := "debian"

	// "library/debian"
	imageStr := ns + "/" + name

	image, err := FromString(imageStr)
	assert.NoError(t, err, "valid image name")

	// Name
	assert.Equal(t, name, image.Name())

	// KubeImage
	assert.Equal(t, imageStr, image.KubeImage())
}

func TestParseNamespaceTag(t *testing.T) {
	ns := "library"
	name := "debian"
	tag := "jessie"

	// "library/debian:jessie"
	imageStr := ns + "/" + name + ":" + tag

	image, err := FromString(imageStr)
	assert.NoError(t, err, "valid image name")

	// Name
	assert.Equal(t, name, image.Name())

	// KubeImage
	assert.Equal(t, imageStr, image.KubeImage())
}

func TestParseRegistryNamespaceName(t *testing.T) {
	registry := "docker.redspread.com:443"
	ns := "library"
	name := "debian"

	// "docker.redspread.com:443/library/debian"
	imageStr := registry + "/" + ns + "/" + name

	image, err := FromString(imageStr)
	assert.NoError(t, err, "valid image name")

	// Name
	assert.Equal(t, name, image.Name())

	// KubeImage
	assert.Equal(t, imageStr, image.KubeImage())
}

func TestParseRegistryNamespaceTag(t *testing.T) {
	registry := "docker.redspread.com:443"
	ns := "library"
	name := "debian"
	tag := "jessie"

	// "docker.redspread.com:443/library/debian:jessie"
	imageStr := registry + "/" + ns + "/" + name + ":" + tag

	image, err := FromString(imageStr)
	assert.NoError(t, err, "valid image name")

	// Name
	assert.Equal(t, name, image.Name())

	// KubeImage
	assert.Equal(t, imageStr, image.KubeImage())
}
