package deploy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func TestObjectPath(t *testing.T) {
	rc := &kube.ReplicationController{
		TypeMeta: unversioned.TypeMeta{
			APIVersion: "v1",
		},
		ObjectMeta: kube.ObjectMeta{
			Name:      "johnson",
			Namespace: kube.NamespaceDefault,
		},
	}

	expected := "v1/namespaces/default/replicationcontroller/johnson"
	actual, err := ObjectPath(rc)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestObjectPathNoNamespace(t *testing.T) {
	rc := &kube.ReplicationController{
		TypeMeta: unversioned.TypeMeta{
			APIVersion: "v1",
		},
		ObjectMeta: kube.ObjectMeta{
			Name: "johnson",
		},
	}

	expected := "v1/namespaces//replicationcontroller/johnson"
	actual, err := ObjectPath(rc)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
