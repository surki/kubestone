package nighthawk

import (
	"path"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

const (
	configsDir = "/configs"
)

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;create;delete;watch

func serverDeploymentName(cr *perfv1alpha1.Nighthawk) string {
	return cr.Name
}

// NewServerDeployment create a nighthawk server deployment from the
// provided Nighthawk server config.
func NewServerDeployment(cr *perfv1alpha1.Nighthawk, configMap *corev1.ConfigMap) *appsv1.Deployment {
	replicas := int32(1)

	labels := map[string]string{
		"kubestone.xridge.io/app":     "nighthawk",
		"kubestone.xridge.io/cr-name": cr.Name,
	}
	// Let's be nice and don't mutate CRs label field
	for k, v := range cr.Spec.ServerConfiguration.PodLabels {
		labels[k] = v
	}

	cmdLineArgs := []string{"-c", path.Join(configsDir, cr.Spec.ServerConfiguration.ConfigFile)}
	cmdLineArgs = append(cmdLineArgs, cr.Spec.ServerConfiguration.CmdLineArgs...)

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        serverDeploymentName(cr),
			Namespace:   cr.Namespace,
			Annotations: cr.Spec.ServerConfiguration.PodConfigurationSpec.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: cr.Spec.ServerConfiguration.PodConfigurationSpec.Annotations,
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: cr.Spec.Image.PullSecret,
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "server",
							Image:           cr.Spec.Image.Name,
							ImagePullPolicy: corev1.PullPolicy(cr.Spec.Image.PullPolicy),
							Command:         []string{"nighthawk_test_server"},
							Args:            cmdLineArgs,
							Ports: []corev1.ContainerPort{
								{
									Name:          "server",
									ContainerPort: cr.Spec.ServerConfiguration.Port,
									Protocol:      corev1.Protocol(corev1.ProtocolTCP),
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{Type: intstr.Int, IntVal: cr.Spec.ServerConfiguration.Port},
									},
								},
								InitialDelaySeconds: 5,
								TimeoutSeconds:      2,
								PeriodSeconds:       2,
							},
							Resources: cr.Spec.ServerConfiguration.Resources,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "configs",
									MountPath: configsDir,
								},
							},
						},
					},
					Affinity:     cr.Spec.ServerConfiguration.PodScheduling.Affinity,
					Tolerations:  cr.Spec.ServerConfiguration.PodScheduling.Tolerations,
					NodeSelector: cr.Spec.ServerConfiguration.PodScheduling.NodeSelector,
					NodeName:     cr.Spec.ServerConfiguration.PodScheduling.NodeName,
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "configs",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMap.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return &deployment
}
