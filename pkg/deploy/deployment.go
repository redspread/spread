package deploy

import (
	"k8s.io/kubernetes/pkg/api/meta"
)

// A Deployment is a collection of Kubernetes deployed. Deployment stores a slice of deployable Kubernetes objects.
// It can be used to create deployments deployments and is how the current state of a deployment is returned.
type Deployment struct {
	Name    string
	Objects []meta.ObjectMetaAccessor
}

// WithOptions customizes deployment with values set in options. The rules this follows can be found
// in DeploymentOptions.
func (d Deployment) WithOptions(opts DeploymentOptions) {
	// setting namespace if not set
	if len(opts.Namespace) > 0 {
		for _, v := range d.Objects {
			m := v.GetObjectMeta()
			if len(m.GetNamespace()) == 0 {
				m.SetNamespace(opts.Namespace)
			}
		}
	}

	d.ApplyAnnotations(opts.Annotations)
	d.ApplyLabels(opts.Labels)
}

// ApplyLabels applies a provided map on top of the current set of Labels of each object.
func (d *Deployment) ApplyLabels(labels map[string]string) {
	d.apply(func(m meta.Object) {
		current := m.GetLabels()
		m.SetLabels(applyMap(labels, current))
	})
}

// ApplyAnnotations applies a provided map on top of the current set of Annotations of each object.
func (d *Deployment) ApplyAnnotations(annotations map[string]string) {
	d.apply(func(m meta.Object) {
		current := m.GetAnnotations()
		m.SetAnnotations(applyMap(annotations, current))
	})
}

func (d *Deployment) apply(f func(meta.Object)) {
	for _, object := range d.Objects {
		meta := object.GetObjectMeta()
		f(meta)
	}
}

func applyMap(src, dst map[string]string) map[string]string {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// DeploymentOptions specify deployment time parameters.
//
// If the Namespace field is set, then all objects without a namespace explicitly set will default to it’s value.
//
// Deployments can have a sets of Annotations and Labels which will be applied to each deployed object. If an Annotation
// or Label already exists it will be overwritten.
type DeploymentOptions struct {
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string
}
