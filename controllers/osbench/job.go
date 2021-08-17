package osbench

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/firepear/qsplit"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
)

// NewJob creates a osbench benchmark job
func NewJob(cr *perfv1alpha1.Osbench) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	osbenchCmdLineArgs := []string{}
	osbenchCmdLineArgs = append(osbenchCmdLineArgs, cr.Spec.TestName)
	osbenchCmdLineArgs = append(osbenchCmdLineArgs, qsplit.ToStrings([]byte(cr.Spec.Options))...)

	job := k8s.NewPerfJob(objectMeta, "osbench", cr.Spec.Image, cr.Spec.PodConfig)
	job.Spec.Template.Spec.Containers[0].Args = osbenchCmdLineArgs
	return job
}
