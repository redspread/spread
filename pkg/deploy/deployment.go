package deploy

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pmezard/go-difflib/difflib"
	kube "k8s.io/kubernetes/pkg/api"

	pb "rsprd.com/spread/pkg/spreadproto"
)

// DeploymentFromDocMap produces a new Deployment from a map of Documents.
func DeploymentFromDocMap(docs map[string]*pb.Document) (deploy *Deployment, err error) {
	var obj KubeObject
	for path, doc := range docs {
		doc.GetInfo().Path = path
		obj, err = KubeObjectFromDocument(path, doc)
		if err != nil {
			return
		}

		err = deploy.Add(obj)
		if err != nil {
			return
		}
	}
	return
}

// A Deployment is a representation of a Kubernetes cluster's object registry.
// It can be used to specify objects to be deployed and is how the current state of a deployment is returned.
// The key of objects is formed by "<apiVersion>/namespaces/<namespace>/<kind>/<name>"
type Deployment struct {
	objects map[string]KubeObject
}

// Add inserts an object into a deployment. The object must be a valid Kubernetes object or it will fail.
// There can only be a single object of the same name, namespace, and type. Objects are deep-copied into the Deployment.
func (d *Deployment) Add(obj KubeObject) error {
	// create objects if doesn't exist
	if d.objects == nil {
		d.objects = make(map[string]KubeObject, 1)
	}

	// generate path using meta and type info
	path, err := ObjectPath(obj)
	if err != nil {
		return fmt.Errorf("could not add object: %v", err)
	}

	// there can only be one object per path
	if _, exists := d.objects[path]; exists {
		return ErrorConflict
	}

	// create a copy of values
	copy, err := deepCopy(obj)
	if err != nil {
		return err
	}

	// add object to
	d.objects[path] = copy
	return nil
}

// AddDeployment inserts the contents of one Deployment into another.
func (d *Deployment) AddDeployment(deployment Deployment) (err error) {
	// this is inefficient-it results in two deep copies being made, ones that's thrown out
	// if this becomes frequently used it should be reimplemented

	// TODO: perform check for collisions before mutation to prevent incomplete additions
	for _, obj := range deployment.Objects() {
		err = d.Add(obj)
		if err != nil {
			return fmt.Errorf("could not add `%s`: %v", obj.GetObjectMeta().GetName(), err)
		}
	}
	return nil
}

// Get returns the object with the given path from the Deployment. Error is returned if object does not exist.
func (d *Deployment) Get(name string) (KubeObject, error) {
	obj, ok := d.objects[name]
	if !ok {
		return nil, fmt.Errorf("no object '%s' exists", name)
	}
	return obj, nil
}

// Equal performs a deep equality check between Deployments. Internal ordering is ignored.
func (d *Deployment) Equal(other *Deployment) bool {
	if other == nil {
		return false
	}

	if len(d.objects) != len(other.objects) {
		return false
	}

	for k, thisVal := range d.objects {
		if otherVal, contains := other.objects[k]; contains {
			// check if object in same key matches
			if !kube.Semantic.DeepEqual(thisVal, otherVal) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// Objects returns the contents of a Deployment. No ordering guarantees are given.
func (d Deployment) Objects() []KubeObject {
	num := 0
	objs := make([]KubeObject, len(d.objects))
	for _, val := range d.objects {
		objs[num] = val
		num++
	}
	return objs
}

// ObjectsOfVersionKind returns all objects matching the given Version and Kind. No ordering guarantees are given.
// If an empty string is given for either selector, all options for that selector will be returned.
func (d Deployment) ObjectsOfVersionKind(version, kind string) (objs []KubeObject) {
	// checkVersion returns true if the version is matched or empty.
	checkVersion := func(objVersion string) bool {
		if len(version) == 0 {
			return true
		}

		return version == objVersion
	}

	// checkKind returns true if the kind is matched or empty.
	checkKind := func(objKind string) bool {
		if len(kind) == 0 {
			return true
		}

		return kind == objKind
	}

	for _, val := range d.objects {
		gvk, err := objectKind(val)
		if err != nil {
			continue
		}
		if checkVersion(gvk.Version) && checkKind(gvk.Kind) {
			objs = append(objs, val)
		}
	}

	return objs
}

// Len returns the number of objects in a Deployment.
func (d Deployment) Len() int {
	return len(d.objects)
}

// String returns a JSON representation of a Deployment
func (d *Deployment) String() string {
	output, err := json.MarshalIndent(d.objects, "", "\t")
	if err != nil {
		panic(err)
	}

	return string(output)
}

// Diff returns the difference between the textual representation of two deployments
func (d *Deployment) Diff(other *Deployment) string {
	if other == nil {
		return "other was nil"
	}
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(d.String()),
		B:        difflib.SplitLines(other.String()),
		FromFile: "ThisDeployment",
		ToFile:   "OtherDeployment",
	}

	out, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		panic(err)
	}

	return out
}

// PathDiff returns the list of the paths of objects.
// Currently doesn't detect modifications
func (d *Deployment) PathDiff(other *Deployment) (added, removed, modified []string) {
	for path, obj := range d.objects {
		if oObj, has := other.objects[path]; has {
			if !kube.Semantic.DeepEqual(obj, oObj) {
				//modified = append(modified, path)
			}

		} else {
			added = append(added, path)
		}
	}

	for path := range other.objects {
		if _, has := d.objects[path]; !has {
			removed = append(removed, path)
		}
	}
	return
}

// Stat returns change information about a deployment.
func Stat(index, head, cluster *Deployment) DiffStat {
	stat := DiffStat{}
	stat.IndexNew, stat.IndexDeleted, stat.IndexModified = index.PathDiff(head)
	stat.ClusterNew, stat.ClusterDeleted, stat.ClusterModified = cluster.PathDiff(index)
	return stat
}

type DiffStat struct {
	IndexNew        []string
	IndexModified   []string
	IndexDeleted    []string
	ClusterNew      []string
	ClusterModified []string
	ClusterDeleted  []string
}

// deepCopy creates a deep copy of the Kubernetes object given.
func deepCopy(obj KubeObject) (KubeObject, error) {
	copy, err := kube.Scheme.DeepCopy(obj)
	if err != nil {
		return nil, err
	}
	return copy.(KubeObject), nil
}

var (
	// ErrorConflict is returned when an object with an identical path already exists in the Deployment.
	ErrorConflict = errors.New("name/namespace combination already exists for type")
)
