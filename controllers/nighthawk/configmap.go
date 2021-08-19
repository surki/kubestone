package nighthawk

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// NewConfigMap creates a new configmap containing the ConfigsVolume
// for the nighthawk server
func NewConfigMap(cr *perfv1alpha1.Nighthawk) *corev1.ConfigMap {
	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Data: cr.Spec.ServerConfiguration.ConfigsVolume,
	}

	return &configMap
}
