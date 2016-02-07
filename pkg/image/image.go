package image

import (
	"io"
	"strings"

	"github.com/docker/distribution/reference"
	docker "github.com/fsouza/go-dockerclient"
)

// Image contains configuration necessary to deploy an image or if necessary, built it.
type Image struct {
	tag   bool
	image reference.NamedTagged
	Build *Build
}

// KubeImage returns a reference to the Image for use with the Image field of kube.Container.
func (i Image) KubeImage() (out string) {
	if i.image == nil {
		return ""
	}
	imageStr := i.image.String()
	if !i.tag {
		// remove default latest tag
		return strings.TrimSuffix(imageStr, ":latest")
	}
	return imageStr
}

// Name returns a human readable identifier for the Image. Should be a DNS Label.
func (i Image) Name() string {
	if i.image == nil {
		return ""
	}
	userName := strings.Split(remoteName(i.image), "/")
	return userName[1]
}

// PushOptions returns the parameters needed to push an image.
func (i Image) PushOptions(out io.Writer, json bool) docker.PushImageOptions {
	if i.image == nil {
		return docker.PushImageOptions{}
	}
	name := strings.TrimPrefix(fullName(i.image), "docker.io/")
	opts := docker.PushImageOptions{
		Name:          name,
		Registry:      hostname(i.image),
		OutputStream:  out,
		RawJSONStream: json,
	}

	if i.tag {
		opts.Tag = i.image.Tag()
	}
	return opts
}

// FromString creates an Image using a string representation
func FromString(str string) (*Image, error) {
	named, err := reference.ParseNamed(str)
	if err != nil {
		return nil, err
	}

	// check for tag
	hasTag := !isNameOnly(named)

	// implicit default tag
	tagged := withDefaultTag(named)

	return &Image{
		image: tagged.(reference.NamedTagged),
		tag:   hasTag,
	}, nil
}

// Implementations below from docker/docker upstream
// TODO: once kubernetes updates Docker version these should be removed and use "github.com/docker/docker/reference"

const (
	// DefaultTag defines the default tag used when performing images related actions and no tag or digest is specified
	DefaultTag = "latest"
	// DefaultHostname is the default built-in hostname
	DefaultHostname = "docker.io"
	// LegacyDefaultHostname is automatically converted to DefaultHostname
	LegacyDefaultHostname = "index.docker.io"
	// DefaultRepoPrefix is the prefix used for default repositories in default host
	DefaultRepoPrefix = "library/"
)

// isNameOnly returns true if reference only contains a repo name.
func isNameOnly(ref reference.Named) bool {
	if _, ok := ref.(reference.NamedTagged); ok {
		return false
	}
	if _, ok := ref.(reference.Canonical); ok {
		return false
	}
	return true
}

// withDefaultTag adds a default tag to a reference if it only has a repo name.
func withDefaultTag(ref reference.Named) reference.Named {
	if isNameOnly(ref) {
		ref, _ = reference.WithTag(ref, DefaultTag)
	}
	return ref
}

// fullName returns full repository name with hostname, like "docker.io/library/ubuntu"
func fullName(r reference.Named) string {
	hostname, remoteName := splitHostname(r.Name())
	return hostname + "/" + remoteName
}

// hostname returns hostname for the reference, like "docker.io"
func hostname(r reference.Named) string {
	hostname, _ := splitHostname(r.Name())
	return hostname
}

// remoteName returns the repository component of the full name, like "library/ubuntu"
func remoteName(r reference.Named) string {
	_, remoteName := splitHostname(r.Name())
	return remoteName
}

// splitHostname splits a repository name to hostname and remotename string.
// If no valid hostname is found, the default hostname is used. Repository name
// needs to be already validated before.
func splitHostname(name string) (hostname, remoteName string) {
	i := strings.IndexRune(name, '/')
	if i == -1 || (!strings.ContainsAny(name[:i], ".:") && name[:i] != "localhost") {
		hostname, remoteName = DefaultHostname, name
	} else {
		hostname, remoteName = name[:i], name[i+1:]
	}
	if hostname == LegacyDefaultHostname {
		hostname = DefaultHostname
	}
	if hostname == DefaultHostname && !strings.ContainsRune(remoteName, '/') {
		remoteName = DefaultRepoPrefix + remoteName
	}
	return
}
