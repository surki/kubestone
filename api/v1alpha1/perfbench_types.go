package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PerfbenchSpec defines the desired state of Perfbench
type PerfbenchSpec struct {
	// Image defines the perfbench docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// CmdLineArgs are appended to the predefined perfbench parameters
	CmdLineArgs []string `json:"cmdLineArgs"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// Perfbench is the Schema for the perfbenches API
type Perfbench struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PerfbenchSpec   `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PerfbenchList contains a list of Perfbench
type PerfbenchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Perfbench `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Perfbench{}, &PerfbenchList{})
}
