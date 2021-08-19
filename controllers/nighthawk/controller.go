package nighthawk

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/xridge/kubestone/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// Reconciler provides fields from manager to reconciler
type Reconciler struct {
	K8S k8s.Access
	Log logr.Logger
}

// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=nighthawks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=nighthawks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=nighthawks/finalizers,verbs=update

// Reconcile Nighthawk Benchmark Requests by creating:
//   - nighthawk server deployment
//   - nighthawk server service
//   - nighthawk client pod
// The creation of nighthawk client pod is postponed until the server
// deployment completes. Once the nighthawk client pod is completed,
// the server deployment and service objects are removed from k8s.
func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	var cr perfv1alpha1.Nighthawk
	if err := r.K8S.Client.Get(ctx, req.NamespacedName, &cr); err != nil {
		return ctrl.Result{}, k8s.IgnoreNotFound(err)
	}

	// Run to one completion
	if cr.Status.Completed {
		return ctrl.Result{}, nil
	}

	// Validate on first entry
	if !cr.Status.Completed && !cr.Status.Running {
		if valid, err := IsCrValid(&cr); !valid {
			_ = r.K8S.RecordEventf(&cr, corev1.EventTypeWarning, k8s.CreateFailed,
				"CR validation failed: %v", err)

			// Do not requeue invalid CRs
			return ctrl.Result{}, nil
		}
	}

	cr.Status.Running = true
	if err := r.K8S.Client.Status().Update(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	configMap := NewConfigMap(&cr)
	if err := r.K8S.CreateWithReference(ctx, configMap, &cr); err != nil {
		return ctrl.Result{}, err
	}

	serverDeployment := NewServerDeployment(&cr, configMap)
	if err := r.K8S.CreateWithReference(ctx, serverDeployment, &cr); err != nil {
		return ctrl.Result{}, err
	}

	serverService := NewServerService(&cr)
	if err := r.K8S.CreateWithReference(ctx, serverService, &cr); err != nil {
		return ctrl.Result{}, err
	}

	endpointReady, err := r.K8S.IsEndpointReady(types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Name})
	if err != nil {
		return ctrl.Result{}, err
	}
	if !endpointReady {
		// Wait for deployment to be connected to the service endpoint
		return ctrl.Result{Requeue: true}, nil
	}

	job := NewClientJob(&cr)
	job.Spec.Template.Spec.Containers[0].Command = []string{"nighthawk_client"}

	if err := r.K8S.CreateWithReference(ctx, job, &cr); err != nil {
		return ctrl.Result{}, err
	}

	jobFinished, err := r.K8S.IsJobFinished(types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      clientJobName(&cr),
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	if !jobFinished {
		// Wait for the job to be completed
		return ctrl.Result{Requeue: true}, nil
	}

	if err := r.K8S.DeleteObject(ctx, serverService, &cr); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.K8S.DeleteObject(ctx, serverDeployment, &cr); err != nil {
		return ctrl.Result{}, err
	}

	cr.Status.Running = false
	cr.Status.Completed = true
	if err := r.K8S.Client.Status().Update(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager registers the NighthawkReconciler with the provided manager
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&perfv1alpha1.Nighthawk{}).
		Complete(r)
}
