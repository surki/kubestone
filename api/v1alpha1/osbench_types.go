package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OsbenchSpec contains the configuration parameters
// with scheduling options for the osbench benchmark.
type OsbenchSpec struct {
	// Image defines the osbench docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,inline"`

	// Options is a list of zero or more command line options.
	// +optional
	Options string `json:"options,omitempty"`

	// TestName is the name of a built-in test (e.g. `create_threads`, `create_files`, etc.)
	TestName string `json:"testName"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// Osbench is the Schema for the osbenches API
type Osbench struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OsbenchSpec     `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OsbenchList contains a list of Osbench
type OsbenchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Osbench `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Osbench{}, &OsbenchList{})
}
