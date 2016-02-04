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

func TestParseTag(t *testing.T) {
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

func TestParseUserName(t *testing.T) {
	user := "base"
	name := "debian"

	// "base/debian"
	imageStr := user + "/" + name

	image, err := FromString(imageStr)
	assert.NoError(t, err, "valid image name")

	// Name
	assert.Equal(t, name, image.Name())

	// KubeImage
	assert.Equal(t, imageStr, image.KubeImage())
}

func TestParseUserTag(t *testing.T) {
	user := "base"
	name := "debian"
	tag := "jessie"

	// "base/debian:jessie"
	imageStr := user + "/" + name + ":" + tag

	image, err := FromString(imageStr)
	assert.NoError(t, err, "valid image name")

	// Name
	assert.Equal(t, name, image.Name())

	// KubeImage
	assert.Equal(t, imageStr, image.KubeImage())
}

func TestParseRegistryUserName(t *testing.T) {
	registry := "docker.redspread.com:443"
	user := "base"
	name := "debian"

	// "docker.redspread.com:443/base/debian"
	imageStr := registry + "/" + user + "/" + name

	image, err := FromString(imageStr)
	assert.NoError(t, err, "valid image name")

	// Name
	assert.Equal(t, name, image.Name())

	// KubeImage
	assert.Equal(t, imageStr, image.KubeImage())
}

func TestParseRegistryUserTag(t *testing.T) {
	registry := "docker.redspread.com:443"
	user := "base"
	name := "debian"
	tag := "jessie"

	// "docker.redspread.com:443/base/debian:jessie"
	imageStr := registry + "/" + user + "/" + name + ":" + tag

	image, err := FromString(imageStr)
	assert.NoError(t, err, "valid image name")

	// Name
	assert.Equal(t, name, image.Name())

	// KubeImage
	assert.Equal(t, imageStr, image.KubeImage())
}

func TestParseInvalidImage(t *testing.T) {
	imageName := "H * A * P * P * Y"
	_, err := FromString(imageName)
	assert.Error(t, err, "invalid image name")
}
