package deploy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

func TestDeploymentSimpleEquals1(t *testing.T) {
	a, b := NewDeployment(), NewDeployment()

	assert.True(t, a.Equals(a), "self")
	assert.True(t, b.Equals(b), "self")
	assert.True(t, a.Equals(b), "both are empty")

	secret := createSecret("test", "a")
	a.Add(secret)
	assert.False(t, a.Equals(b), "a has secret")
}

func TestDeploymentSimpleEquals2(t *testing.T) {
	a, b, c := createSecret("a", "1"), createSecret("b", "2"), createSecret("c", "3")
	deployA, deployB := NewDeployment(), NewDeployment()

	deployA.Add(a)
	deployA.Add(b)
	deployA.Add(c)

	assert.False(t, deployA.Equals(deployB), "deployB is empty")
	assert.False(t, deployB.Equals(deployA), "deployB is empty")

	deployB.Add(c)
	deployB.Add(b)
	deployB.Add(a)

	assert.True(t, deployB.Equals(deployA), "same")
}

func createSecret(name, data string) *api.Secret {
	return &api.Secret{
		ObjectMeta: api.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Data: map[string][]byte{"test": []byte(data)},
		Type: api.SecretTypeOpaque,
	}
}
