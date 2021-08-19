package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NighthawkConfigurationSpec contains configuration parameters
// with scheduling options for the both the nighthawk client
// and server instances.
type NighthawkClientConfigurationSpec struct {
	PodConfigurationSpec `json:",inline"`

	// CmdLineArgs are appended to the predefined nighthawk parameters
	CmdLineArgs []string `json:"cmdLineArgs"`
}

// NighthawkConfigurationSpec contains configuration parameters
// with scheduling options for the both the nighthawk client
// and server instances.
type NighthawkServerConfigurationSpec struct {
	PodConfigurationSpec `json:",inline"`

	// CmdLineArgs are appended to the predefined nighthawk parameters
	CmdLineArgs []string `json:"cmdLineArgs,omitempty"`

	// ConfigsVolume holds the content of config files.  The key of the
	// map specifies the filename and the value is the content of the
	// file. ConfigMap is created from the map which is mounted as
	// config directory to the server/client pods.
	ConfigsVolume map[string]string `json:"configsVolume"`

	// Config file to use
	ConfigFile string `json:"configFile"`

	// Port the server listens on, used for readiness check and service
	// creation
	Port int32 `json:"port"`
}

// NighthawkSpec defines the Nighthawk Benchmark Stone which
// consist of server deployment with service definition
// and client pod.
type NighthawkSpec struct {
	// Image defines the nighthawk docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// ServerConfiguration contains the configuration of the nighthawk server
	ServerConfiguration NighthawkServerConfigurationSpec `json:"serverConfiguration,omitempty"`

	// ClientConfiguration contains the configuration of the nighthawk client
	ClientConfiguration NighthawkClientConfigurationSpec `json:"clientConfiguration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// Nighthawk is the Schema for the nighthawks API
type Nighthawk struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NighthawkSpec   `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NighthawkList contains a list of Nighthawk
type NighthawkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Nighthawk `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Nighthawk{}, &NighthawkList{})
}
