package nighthawk

import (
	"fmt"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;create;delete;watch

func serverServiceName(cr *perfv1alpha1.Nighthawk) string {
	return cr.Name
}

func serverServicePort(cr *perfv1alpha1.Nighthawk) int32 {
	return cr.Spec.ServerConfiguration.Port
}

func serverServiceNamePort(cr *perfv1alpha1.Nighthawk) string {
	return fmt.Sprintf("%s:%d", serverServiceName(cr), serverServicePort(cr))
}

// NewServerService creates k8s headless service (which targets the server deployment)
// from the Nighthawk server definition
func NewServerService(cr *perfv1alpha1.Nighthawk) *corev1.Service {
	labels := map[string]string{
		"kubestone.xridge.io/app":     "nighthawk",
		"kubestone.xridge.io/cr-name": cr.Name,
	}
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        serverServiceName(cr),
			Namespace:   cr.Namespace,
			Annotations: cr.Spec.ServerConfiguration.PodConfigurationSpec.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "nighthawk",
					Protocol: corev1.Protocol(corev1.ProtocolTCP),
					Port:     serverServicePort(cr),
				},
			},
			Selector:  labels,
			ClusterIP: "None", // Headless service!
		},
	}

	return &service
}
