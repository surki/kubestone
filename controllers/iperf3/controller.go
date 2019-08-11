/*
Copyright 2019 The xridge kubestone contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package iperf3

import (
	"context"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/xridge/kubestone/pkg/k8s"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// Reconciler provides fields from manager to reconciler
type Reconciler struct {
	K8S k8s.Access
	Log logr.Logger
}

// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=iperf3s,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=iperf3s/status,verbs=get;update;patch

// Reconcile Iperf3 Benchmark Requests by creating:
//   - iperf3 server deployment
//   - iperf3 server service
//   - iperf3 client pod
// The creation of iperf3 client pod is postponed until the server
// deployment completes. Once the iperf3 client pod is completed,
// the server deployment and service objects are removed from k8s.
func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	var cr perfv1alpha1.Iperf3
	if err := r.K8S.Client.Get(ctx, req.NamespacedName, &cr); err != nil {
		return ctrl.Result{}, k8s.IgnoreNotFound(err)
	}

	// Run to one completion
	if cr.Status.Completed {
		return ctrl.Result{}, nil
	}

	cr.Status.Running = true
	if err := r.K8S.Client.Status().Update(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.K8S.CreateWithReference(ctx, newServerDeployment(&cr), &cr); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.K8S.CreateWithReference(ctx, newServerService(&cr), &cr); err != nil {
		return ctrl.Result{}, err
	}

	serverReady, err := r.serverDeploymentReady(&cr)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !serverReady {
		// Wait for deployment to be ready
		return ctrl.Result{Requeue: true}, nil
	}

	if err := r.K8S.CreateWithReference(ctx, newClientPod(&cr), &cr); err != nil {
		return ctrl.Result{}, err
	}

	clientPodFinished, err := r.clientPodFinished(&cr)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !clientPodFinished {
		// Wait for the client pod to be completed
		return ctrl.Result{Requeue: true}, nil
	}

	if err := r.deleteServerService(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.deleteServerDeployment(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	cr.Status.Running = false
	cr.Status.Completed = true
	if err := r.K8S.Client.Status().Update(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager registers the Iperf3Reconciler with the provided manager
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&perfv1alpha1.Iperf3{}).
		Complete(r)
}
