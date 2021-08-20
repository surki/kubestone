package perfbench

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
)

// NewJob creates a perfbench benchmark job
func NewJob(cr *perfv1alpha1.Perfbench) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	job := k8s.NewPerfJob(objectMeta, "perfbench", cr.Spec.Image, cr.Spec.PodConfig)
	job.Spec.Template.Spec.Containers[0].Command = []string{"perf", "bench"}
	job.Spec.Template.Spec.Containers[0].Args = cr.Spec.CmdLineArgs
	privileged := true
	job.Spec.Template.Spec.Containers[0].SecurityContext = &corev1.SecurityContext{
		Privileged: &privileged,
	}
	return job
}
