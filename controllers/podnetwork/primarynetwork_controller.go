/*
Copyright 2021.

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

package podnetwork

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	podnetworkv1alpha1 "github.com/opdev/pod-network-operator/apis/podnetwork/v1alpha1"
)

// PrimaryNetworkReconciler reconciles a PrimaryNetwork object
type PrimaryNetworkReconciler struct {
	client.Client
	Log            logr.Logger
	Scheme         *runtime.Scheme
	PrimaryNetwork *podnetworkv1alpha1.PrimaryNetwork
}

//+kubebuilder:rbac:groups=podnetwork.opdev.io,resources=primarynetworks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=podnetwork.opdev.io,resources=primarynetworks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=podnetwork.opdev.io,resources=primarynetworks/finalizers,verbs=update

func (r *PrimaryNetworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("primarynetwork", req.NamespacedName)

	// Loop through the list of pods with primary newtworks matching labels

	// call finalizer on primary network configuration resource

	// update primary network status

	// log new primary network configuration requested

	// begin configuration task

	// Beginning network configuration task

	// Loop through configuration fields requested

	// log configuration item in progress

	// update status condition InProgress reason item X being configured

	// call appropriate link configuration function for item passing Pod as parameter

	// check error on return
	// log error or
	// log configuration Pod name, field and value - succeeded

	// update status configuration list adding Pod name, field and value configured

	// End configuration task

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PrimaryNetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&podnetworkv1alpha1.PrimaryNetwork{}).
		Complete(r)
}
