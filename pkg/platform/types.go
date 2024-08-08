package platform

import "k8s.io/apimachinery/pkg/runtime/schema"

// ObjectReference is a reference to a Kubernetes resource which Platform uses to enable certain capabilities.
// These custom resources serve as single point of configuration for enabling given capability for the component.
type ObjectReference struct {
	// GroupVersionKind specifies the group, version, and kind of the resource.
	schema.GroupVersionKind `json:"gvk,omitempty"`
	// Resources is the type of resource being protected in a plural form, e.g., "pods", "services".
	Resources string `json:"resources,omitempty"`
}

// RoutingTarget represents a target object that routing capability
// will watch to ensure proper routing configuration.
type RoutingTarget struct {
	ObjectReference `json:"ref,omitempty"`
}

// ProtectedResource defines a custom resource type that the component requires capability for.
type ProtectedResource struct {
	ObjectReference `json:"ref,omitempty"`
	// WorkloadSelector is a map of labels used to select the workload.
	WorkloadSelector map[string]string `json:"workloadSelector,omitempty"`
	// HostPaths is a list of host paths associated with the resource.
	HostPaths []string `json:"hostPaths,omitempty"`
	// Ports is a list of ports associated with the resource.
	Ports []string `json:"ports,omitempty"`
}
