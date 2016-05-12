package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	kube "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	rest "k8s.io/kubernetes/pkg/client/restclient"
	kubecli "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd/config"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/strategicpatch"
)

const DefaultContext = ""

// KubeCluster is able to deploy to Kubernetes clusters. This is a very simple implementation with no error recovery.
type KubeCluster struct {
	client    *kubecli.Client
	context   string
	localkube bool
}

// NewKubeClusterFromContext creates a KubeCluster using a Kubernetes client with the configuration of the given context.
// If the context name is empty, the default context will be used.
func NewKubeClusterFromContext(name string) (*KubeCluster, error) {
	rulesConfig, err := defaultLoadingRules().Load()
	if err != nil {
		return nil, fmt.Errorf("could not load rules: %v", err)
	}

	clientConfig := clientcmd.NewNonInteractiveClientConfig(*rulesConfig, name, &clientcmd.ConfigOverrides{})

	rawConfig, err := clientConfig.RawConfig()
	if err != nil || rawConfig.Contexts == nil {
		return nil, fmt.Errorf("could not access kubectl config: %v", err)
	}

	if name == DefaultContext {
		name = rawConfig.CurrentContext
	}

	if rawConfig.Contexts[name] == nil {
		return nil, fmt.Errorf("context '%s' does not exist", name)
	}

	restClientConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get ClientConfig: %v", err)
	}

	client, err := kubecli.New(restClientConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create Kubernetes client: %v", err)
	}

	return &KubeCluster{
		client:    client,
		context:   name,
		localkube: name == "localkube",
	}, nil
}

// Context returns the kubectl context being used
func (c *KubeCluster) Context() string {
	return c.context
}

// Deploy creates/updates the Deployment's objects on the Kubernetes cluster.
// Currently no error recovery is implemented; if there is an error the deployment process will immediately halt and return the error.
// If update is not set, will error if objects exist. If deleteModifiedPods is set, pods of modified RCs will be deleted.
func (c *KubeCluster) Deploy(dep *Deployment, update, deleteModifiedPods bool) error {
	if c.client == nil {
		return errors.New("client not setup (was nil)")
	}

	// create namespaces before everything else
	for _, nsObj := range dep.ObjectsOfVersionKind("", "Namespace") {
		ns := nsObj.(*kube.Namespace)
		_, err := c.client.Namespaces().Create(ns)
		if err != nil && !alreadyExists(err) {
			return err
		}
	}

	// TODO: add continue on error and error lists
	for _, obj := range dep.Objects() {
		// don't create namespaces again
		if _, isNamespace := obj.(*kube.Namespace); isNamespace {
			continue
		}

		err := c.deploy(obj, update)
		if err != nil {
			return err
		}

		if rc, isRC := obj.(*kube.ReplicationController); isRC && deleteModifiedPods {
			err = c.deletePods(rc)
			if err != nil {
				return fmt.Errorf("could not delete pods for rc `%s/%s`: %v", rc.Namespace, rc.Name, err)
			}
		}
	}

	printLoadBalancers(c.client, dep.ObjectsOfVersionKind("", "Service"), c.localkube)

	// deployed successfully
	return nil
}

// deploy creates the object on the connected Kubernetes instance. Errors if object exists and not updating.
func (c *KubeCluster) deploy(obj KubeObject, update bool) error {
	if obj == nil {
		return errors.New("tried to deploy nil object")
	}

	mapping, err := mapping(obj)
	if err != nil {
		return err
	}

	if update {
		_, err := c.update(obj, true, mapping)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = c.create(obj, mapping)
	return err
}

// update replaces the currently deployed version with a new one. If the objects already match then nothing is done.
func (c *KubeCluster) update(obj KubeObject, create bool, mapping *meta.RESTMapping) (KubeObject, error) {
	meta := obj.GetObjectMeta()

	deployed, err := c.get(meta.GetNamespace(), meta.GetName(), true, mapping)
	if doesNotExist(err) && create {
		return c.create(obj, mapping)
	} else if err != nil {
		return nil, err
	}

	// TODO: need a better way to handle resource versioning
	// set resource version on local to same as remote
	deployedVersion := deployed.GetObjectMeta().GetResourceVersion()
	meta.SetResourceVersion(deployedVersion)

	copyImmutables(deployed, obj)

	// if local matches deployed, do nothing
	if kube.Semantic.DeepEqual(obj, deployed) {
		return deployed, nil
	}

	patch, err := diff(deployed, obj)
	if err != nil {
		return nil, fmt.Errorf("could not create diff: %v", err)
	}

	req := c.client.RESTClient.Patch(kube.StrategicMergePatchType).
		Name(meta.GetName()).
		Body(patch)

	setRequestObjectInfo(req, meta.GetNamespace(), mapping)

	runtimeObj, err := req.Do().Get()
	if err != nil {
		return nil, resourceError("update", meta.GetNamespace(), meta.GetName(), mapping, err)
	}

	return asKubeObject(runtimeObj)
}

// get retrieves the object from the cluster.
func (c *KubeCluster) get(namespace, name string, export bool, mapping *meta.RESTMapping) (KubeObject, error) {
	req := c.client.RESTClient.Get().Name(name)
	setRequestObjectInfo(req, namespace, mapping)

	if export {
		req.Param("export", "true")
	}

	runtimeObj, err := req.Do().Get()
	if err != nil {
		return nil, resourceError("get", namespace, name, mapping, err)
	}

	return asKubeObject(runtimeObj)
}

// create adds the object to the cluster.
func (c *KubeCluster) create(obj KubeObject, mapping *meta.RESTMapping) (KubeObject, error) {
	meta := obj.GetObjectMeta()
	req := c.client.RESTClient.Post().Body(obj)

	setRequestObjectInfo(req, meta.GetNamespace(), mapping)

	runtimeObj, err := req.Do().Get()
	if err != nil {
		return nil, resourceError("create", meta.GetName(), meta.GetNamespace(), mapping, err)
	}

	return asKubeObject(runtimeObj)
}

func (c *KubeCluster) deletePods(rc *kube.ReplicationController) error {
	if rc == nil {
		return errors.New("rc was nil")
	}

	// list pods
	opts := kube.ListOptions{
		LabelSelector: labels.Set(rc.Spec.Selector).AsSelector(),
	}
	podList, err := c.client.Pods(rc.Namespace).List(opts)
	if err != nil {
		return err
	}

	// delete pods
	for _, pod := range podList.Items {
		err := c.client.Pods(pod.Namespace).Delete(pod.Name, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// setRequestObjectInfo adds necessary type information to requests.
func setRequestObjectInfo(req *rest.Request, namespace string, mapping *meta.RESTMapping) {
	// if namespace scoped resource, set namespace
	req.NamespaceIfScoped(namespace, isNamespaceScoped(mapping))

	// set resource name
	req.Resource(mapping.Resource)
}

// alreadyExists checks if the error is for a resource already existing.
func alreadyExists(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasSuffix(err.Error(), "already exists")
}

// doesNotExist checks if the error is for a non-existent resource.
func doesNotExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasSuffix(err.Error(), "not found")
}

// mapping returns the appropriate RESTMapping for the object.
func mapping(obj KubeObject) (*meta.RESTMapping, error) {
	gvk, err := kube.Scheme.ObjectKind(obj)
	if err != nil {
		return nil, err
	}

	mapping, err := kube.RESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, fmt.Errorf("could not create RESTMapping for %s: %v", gvk, err)
	}
	return mapping, nil
}

// isNamespaceScoped returns if the mapping is scoped by Namespace.
func isNamespaceScoped(mapping *meta.RESTMapping) bool {
	return mapping.Scope.Name() == meta.RESTScopeNameNamespace
}

// defaultLoadingRules use the same rules (as of 2/17/16) as kubectl.
func defaultLoadingRules() *clientcmd.ClientConfigLoadingRules {
	opts := config.NewDefaultPathOptions()

	loadingRules := opts.LoadingRules
	loadingRules.Precedence = opts.GetLoadingPrecedence()
	return loadingRules
}

// diff creates a patch.
func diff(original, modified runtime.Object) (patch []byte, err error) {
	origBytes, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}

	modBytes, err := json.Marshal(modified)
	if err != nil {
		return nil, err
	}

	return strategicpatch.CreateTwoWayMergePatch(origBytes, modBytes, original)
}

// asKubeObject attempts use the object as a KubeObject. It will return an error if not possible.
func asKubeObject(runtimeObj runtime.Object) (KubeObject, error) {
	kubeObj, ok := runtimeObj.(KubeObject)
	if !ok {
		return nil, errors.New("was unable to use runtime.Object as deploy.KubeObject")
	}
	return kubeObj, nil
}

func resourceError(action, namespace, name string, mapping *meta.RESTMapping, err error) error {
	if mapping == nil || mapping.GroupVersionKind.IsEmpty() {
		return fmt.Errorf("could not %s '%s/%s': %v", action, namespace, name, err)
	}
	gvk := mapping.GroupVersionKind
	return fmt.Errorf("could not %s '%s/%s' (%s): %v", action, namespace, name, gvk.Kind, err)
}

// copyImmutables sets any immutable fields from src on dst. Will panic if objects not of same type.
func copyImmutables(src, dst KubeObject) {
	if src == nil || dst == nil {
		return
	}

	// each type has specific fields that must be copied
	switch src := src.(type) {
	case *kube.Service:
		dst := dst.(*kube.Service)
		dst.Spec.ClusterIP = src.Spec.ClusterIP
	}
}

// printLoadBalancers blocks until all Services of type LoadBalancer have been deployed, printing it's details as it becomes available.
// Will panic if given something other than services
func printLoadBalancers(client *kubecli.Client, services []KubeObject, localkube bool) {
	if len(services) == 0 {
		return
	}

	first := true
	completed := map[string]bool{}

	// checks when we've seen every service
	done := func() bool {
		for _, svcObj := range services {
			s := svcObj.(*kube.Service)
			if s.Spec.Type == kube.ServiceTypeLoadBalancer && !completed[s.Name] {
				return false
			}
		}
		return true
	}

	for {
		if done() {
			return
		}

		if first {
			fmt.Println("Waiting for load balancer deployment...")
			first = false
		}

		for _, svcObj := range services {
			s := svcObj.(*kube.Service)
			if s.Spec.Type == kube.ServiceTypeLoadBalancer && !completed[s.Name] {
				clusterVers, err := client.Services(s.Namespace).Get(s.Name)
				if err != nil {
					fmt.Printf("Error getting service `%s`: %v\n", s.Name, err)
				}

				if localkube {
					completed[s.Name] = true
					for _, port := range clusterVers.Spec.Ports {
						fmt.Printf("'%s/%s' - %s available on localkube host port:\t %d\n", s.Namespace, s.Name, port.Name, port.NodePort)
					}
				}

				loadBalancers := clusterVers.Status.LoadBalancer.Ingress
				if len(loadBalancers) == 1 {
					completed[s.Name] = true
					fmt.Printf("Service '%s/%s' available at: \t%s\n", s.Namespace, s.Name, loadBalancers[0].IP)
				}
			}
		}

		// prevents warning about throttling
		time.Sleep(250 * time.Millisecond)
	}
}
