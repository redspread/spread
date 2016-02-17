package deploy

import (
	"encoding/json"
	"errors"
	"fmt"

	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	kubecli "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd/config"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/strategicpatch"
)

// KubeCluster is able to deploy to Kubernetes clusters.
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

// Deploy runs the Deployable remotely.
func (c *KubeCluster) Deploy(dep *Deployment, update bool) error {
	if c.client == nil {
		return errors.New("client not setup (was nil)")
	}

	// create namespaces before everything else
	for _, ns := range dep.namespaces {
		_, err := c.client.Namespaces().Create(ns)
		if err != nil && !alreadyExists(err) {
			return err
		}
	}

	// TODO: add continue on error and error lists
	for _, obj := range dep.Objects() {
		err := c.deploy(obj, update)
		if err != nil {
			return err
		}
	}

	// deployed successfully
	return nil
}

// deploy creates the object on the connected Kubernetes instance. Errors if object exists and not updating.
func (c *KubeCluster) deploy(obj KubeObject, update bool) error {
	if obj == nil {
		return errors.New("tried to deploy nil object")
	}

	if update {
		_, err := c.update(obj, true)
		if err != nil {
			meta := obj.GetObjectMeta()
			err = fmt.Errorf("could not update '%s/%s' (%s): %v", meta.GetNamespace(), meta.GetName(),
				obj.GetObjectKind().GroupVersionKind(), err)
			return err
		}
		return nil
	}

	_, err := c.create(obj)
	return err
}

// update replaces the currently deployed version with a new one. If the objects already match then nothing is done.
func (c *KubeCluster) update(obj KubeObject, create bool) (runtime.Object, error) {
	meta, mapping := obj.GetObjectMeta(), mapping(obj)

	deployedVersion, err := c.get(meta.GetNamespace(), meta.GetName(), true, mapping)
	if doesNotExist(err) && create {
		return c.create(obj)
	} else if err != nil {
		return nil, err
	}

	// if local matches deployed, do nothing
	if kube.Semantic.DeepEqual(obj, deployedVersion) {
		return deployedVersion, nil
	}

	patch, err := diff(deployedVersion, obj)
	if err != nil {
		return nil, fmt.Errorf("could not create diff: %v", err)
	}

	req := c.client.RESTClient.Patch(kube.StrategicMergePatchType).
		Name(meta.GetName()).
		Body(patch)
	setRequestObjectInfo(req, meta.GetNamespace(), mapping)

	return req.Do().Get()
}

// get retrieves the object from the cluster
func (c *KubeCluster) get(namespace, name string, export bool, mapping *meta.RESTMapping) (runtime.Object, error) {
	req := c.client.RESTClient.Get().Name(name)
	setRequestObjectInfo(req, namespace, mapping)

	if export {
		req.Param("export", "true")
	}

	return req.Do().Get()
}

// create adds the object to the cluster
func (c *KubeCluster) create(obj KubeObject) (runtime.Object, error) {
	req := c.client.RESTClient.Post().Body(obj)
	setRequestObjectInfo(req, obj.GetObjectMeta().GetNamespace(), mapping(obj))
	return req.Do().Get()
}

// setRequestObjectInfo adds necessary type information to requests
func setRequestObjectInfo(req *kubecli.Request, namespace string, mapping *meta.RESTMapping) {
	// if namespace scoped resource, set namespace
	req.NamespaceIfScoped(namespace, isNamespaceScoped(mapping))

	// set resource name
	req.Resource(mapping.Resource)
}

// TODO: implement
// alreadyExists checks if the error is for a resource already existing
func alreadyExists(err error) bool {
	if err == nil {
		return false
	}
	return false
}

// TODO: implement
// doesNotExist checks if the error is for a non-existent resource
func doesNotExist(err error) bool {
	if err == nil {
		return false
	}
	return false
}

// mapping returns the appropriate RESTMapping for the object
func mapping(obj KubeObject) *meta.RESTMapping {
	gvk := obj.GetObjectKind().GroupVersionKind()
	mapping, err := kube.RESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		fatalMsg := fmt.Errorf("could not create RESTMapping for %s: %v", gvk, err)
		panic(fatalMsg)
	}
	return mapping
}

// isNamespaceScoped returns if the mapping is scoped by Namespace
func isNamespaceScoped(mapping *meta.RESTMapping) bool {
	return mapping.Scope.Name() == meta.RESTScopeNameNamespace
}

// defaultLoadingRules use the same rules (as of 2/17/16) as kubectl
func defaultLoadingRules() *clientcmd.ClientConfigLoadingRules {
	opts := config.NewDefaultPathOptions()

	loadingRules := opts.LoadingRules
	loadingRules.Precedence = opts.GetLoadingPrecedence()
	return loadingRules
}

// diff creates a patch
func diff(original, modified runtime.Object) (patch []byte, err error) {
	origBytes, err := objToJSONBytes(original)
	if err != nil {
		return nil, err
	}

	modBytes, err := objToJSONBytes(modified)
	if err != nil {
		return nil, err
	}

	return strategicpatch.CreateTwoWayMergePatch(origBytes, modBytes, original)
}

func objToJSONBytes(obj runtime.Object) ([]byte, error) {
	return json.Marshal(obj)
}
