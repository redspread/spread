package deploy

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	kube "k8s.io/kubernetes/pkg/api"
)

func TestDeploymentSimpleEquals1(t *testing.T) {
	a, b := new(Deployment), new(Deployment)

	assert.True(t, a.Equal(a), "self")
	assert.True(t, b.Equal(b), "self")
	assert.True(t, a.Equal(b), "both are empty")

	secret := createSecret("test", "a")
	a.Add(secret)
	assert.False(t, a.Equal(b), "a has secret")

	b.Add(createSecret("secrete", "data"))
	assert.False(t, a.Equal(b), "still bad")
}

func TestDeploymentSimpleEquals2(t *testing.T) {
	a, b, c := createSecret("a", "1"), createSecret("b", "2"), createSecret("c", "3")
	deployA, deployB := new(Deployment), new(Deployment)

	deployA.Add(a)
	deployA.Add(b)
	deployA.Add(c)

	assert.False(t, deployA.Equal(deployB), "deployB is empty")
	assert.False(t, deployB.Equal(deployA), "deployB is empty")

	deployB.Add(c)
	deployB.Add(b)
	deployB.Add(a)

	assert.True(t, deployB.Equal(deployA), "same")

	c.Name = "new"
	deployB.Add(c)
	assert.False(t, deployB.Equal(deployA), "added another c")
}

func TestNoDuplicateNames(t *testing.T) {
	secretA := createSecret("secret-a", "some data")
	secretB := createSecret("secret-a", "different data")

	deployment := new(Deployment)
	assert.NoError(t, deployment.Add(secretA), "valid add")
	assert.Error(t, deployment.Add(secretA), "duplicate name")
	assert.Error(t, deployment.Add(secretB), "duplicate name")

	// different namespace is okay though
	secretB.Namespace = "somewhere-else"
	assert.NoError(t, deployment.Add(secretB), "same name / different namespace")
}

func TestDeploymentObjects(t *testing.T) {
	secret1, secret2 := createSecret("secret1", "dpn't tell"), createSecret("secret2", "spoilers!")
	pod := kube.Pod{
		ObjectMeta: kube.ObjectMeta{
			Name:      "pod",
			Namespace: "default",
		},
		Spec: kube.PodSpec{
			Containers: []kube.Container{
				kube.Container{
					Name:            "container",
					Image:           "redis",
					ImagePullPolicy: kube.PullAlways,
				},
			},
			RestartPolicy: kube.RestartPolicyAlways,
			DNSPolicy:     kube.DNSClusterFirst,
		},
	}

	deploy := new(Deployment)
	deploy.Add(secret1)
	deploy.Add(secret2)
	deploy.Add(&pod)

	assert.Equal(t, 3, deploy.Len(), "should have 3 items")

	objects := deploy.Objects()

	for i := 0; i < len(objects); i++ {
		if kube.Semantic.DeepEqual(objects[i], secret1) {
			continue
		}

		if kube.Semantic.DeepEqual(objects[i], secret2) {
			continue
		}

		if kube.Semantic.DeepEqual(objects[i], &pod) {
			continue
		}

		t.Errorf("'%s' did not match any, print: %s", objects[i].GetObjectMeta().GetName(), spew.Sdump(objects[i]))
	}

	secrets := deploy.ObjectsOfVersionKind("", "Secret")
	assert.Len(t, secrets, 2)
	for _, v := range secrets {
		_, isSecret := v.(*kube.Secret)
		assert.True(t, isSecret, "must be secret")
	}

	pods := deploy.ObjectsOfVersionKind("", "Pod")
	assert.Len(t, pods, 1)
	for _, v := range pods {
		_, isPod := v.(*kube.Pod)
		assert.True(t, isPod, "must be pod")
	}

	rcs := deploy.ObjectsOfVersionKind("", "ReplicationController")
	assert.Len(t, rcs, 0)
}

func createSecret(name, data string) *kube.Secret {
	return &kube.Secret{
		ObjectMeta: kube.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Data: map[string][]byte{"test": []byte(data)},
		Type: kube.SecretTypeOpaque,
	}
}
