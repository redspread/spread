package deploy

import (
	"fmt"

	kubecli "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd/config"
)

// KubeCluster provides implements Deployer for Kubernetes clusters.
type KubeCluster struct {
	client *kubecli.Client
}

// NewKubeClusterFromContext creates a KubeCluster using a Kubernetes client with the configuration of the given context.
// If the context name is empty, the default context will be used
func NewKubeClusterFromContext(name string) (*KubeCluster, error) {
	rules := defaultLoadingRules()
	overrides := &clientcmd.ConfigOverrides{
		CurrentContext: name,
	}

	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)

	clientConfig, err := config.ClientConfig()
	if err != nil {
		if len(name) == 0 {
			return nil, fmt.Errorf("could not use default context: %v", err)
		}
		return nil, fmt.Errorf("could not use context `%s`: %v", name, err)
	}

	client, err := kubecli.New(clientConfig)
	if err != nil {
		return nil, err
	}

	return &KubeCluster{
		client: client,
	}, nil
}

func defaultLoadingRules() *clientcmd.ClientConfigLoadingRules {
	opts := config.NewDefaultPathOptions()

	loadingRules := opts.LoadingRules
	loadingRules.Precedence = opts.GetLoadingPrecedence()
	return loadingRules
}
