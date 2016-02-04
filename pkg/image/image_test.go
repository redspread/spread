package image

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
)

const DefaultDockerRegistry = "docker.io"

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

	// PushOptions
	out := testSampleWriter(2)
	json := false
	expected := docker.PushImageOptions{
		Name:          "library/" + name,
		Registry:      DefaultDockerRegistry,
		OutputStream:  out,
		RawJSONStream: json,
	}
	assert.Equal(t, expected, image.PushOptions(out, json))
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

	// PushOptions
	out := testSampleWriter(3)
	json := true
	expected := docker.PushImageOptions{
		Name:          "library/" + name,
		Tag:           tag,
		Registry:      DefaultDockerRegistry,
		OutputStream:  out,
		RawJSONStream: json,
	}
	assert.Equal(t, expected, image.PushOptions(out, json))
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

	// PushOptions
	out := testSampleWriter(3)
	json := true
	expected := docker.PushImageOptions{
		Name:          imageStr,
		Registry:      DefaultDockerRegistry,
		OutputStream:  out,
		RawJSONStream: json,
	}
	assert.Equal(t, expected, image.PushOptions(out, json))
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

	// PushOptions
	out := testSampleWriter(4)
	json := false
	expected := docker.PushImageOptions{
		Name:          user + "/" + name,
		Tag:           tag,
		Registry:      DefaultDockerRegistry,
		OutputStream:  out,
		RawJSONStream: json,
	}
	assert.Equal(t, expected, image.PushOptions(out, json))
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

	// PushOptions
	out := testSampleWriter(5)
	json := true
	expected := docker.PushImageOptions{
		Name:          registry + "/" + user + "/" + name,
		Registry:      registry,
		OutputStream:  out,
		RawJSONStream: json,
	}
	assert.Equal(t, expected, image.PushOptions(out, json))
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

	// PushOptions
	out := testSampleWriter(6)
	json := false
	expected := docker.PushImageOptions{
		Name:          registry + "/" + user + "/" + name,
		Tag:           tag,
		Registry:      registry,
		OutputStream:  out,
		RawJSONStream: json,
	}
	assert.Equal(t, expected, image.PushOptions(out, json))
}

func TestParseInvalidImage(t *testing.T) {
	imageName := "H * A * P * P * Y"
	_, err := FromString(imageName)
	assert.Error(t, err, "invalid image name")
}

func TestNonInitImage(t *testing.T) {
	image := new(Image)
	assert.Len(t, image.Name(), 0, "not setup")
	assert.Len(t, image.KubeImage(), 0, "not setup")
	assert.EqualValues(t, docker.PushImageOptions{}, image.PushOptions(nil, false), "not setup")
}

type testSampleWriter int

func (testSampleWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
