package nighthawk

import (
	"errors"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
)

// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;create;delete

func clientJobName(cr *perfv1alpha1.Nighthawk) string {
	// Should not match with service name as the pod's
	// hostname is set to it's name. If the two matches
	// the destination ip will resolve to 127.0.0.1 and
	// the server will be unreachable.
	return serverServiceName(cr) + "-client"
}

// NewClientJob creates an Nighthawk Client Job (targeting the
// Server Deployment via the Server Service) from the provided
// Nighthawk config definition.
func NewClientJob(cr *perfv1alpha1.Nighthawk) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      clientJobName(cr),
		Namespace: cr.Namespace,
	}

	cmdLineArgs := append(cr.Spec.ClientConfiguration.CmdLineArgs, serverServiceNamePort(cr))

	job := k8s.NewPerfJob(objectMeta, "nighthawk-client", cr.Spec.Image,
		cr.Spec.ClientConfiguration.PodConfigurationSpec)
	backoffLimit := int32(6)
	job.Spec.BackoffLimit = &backoffLimit
	job.Spec.Template.Spec.Containers[0].Args = cmdLineArgs

	return job
}

// IsCrValid validates the given CR and raises error if semantic errors detected
// For drill it checks that the configFile exists in the ConfigsVolume map
func IsCrValid(cr *perfv1alpha1.Nighthawk) (valid bool, err error) {
	if _, ok := cr.Spec.ServerConfiguration.ConfigsVolume[cr.Spec.ServerConfiguration.ConfigFile]; !ok {
		return false, errors.New("ConfigFile does not exist in ConfigsVolume")
	}

	return true, nil
}
