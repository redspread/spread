package deploy

var shortForms = map[string]string{
	"cs":     "componentstatuses",
	"ds":     "daemonsets",
	"ep":     "endpoints",
	"ev":     "events",
	"hpa":    "horizontalpodautoscalers",
	"ing":    "ingresses",
	"limits": "limitranges",
	"no":     "nodes",
	"ns":     "namespaces",
	"po":     "pods",
	"psp":    "podSecurityPolicies",
	"pvc":    "persistentvolumeclaims",
	"pv":     "persistentvolumes",
	"quota":  "resourcequotas",
	"rc":     "replicationcontrollers",
	"rs":     "replicasets",
	"svc":    "services",
}

func KubeShortForm(resource string) string {
	if long, ok := shortForms[resource]; ok {
		return long
	}
	return resource
}
