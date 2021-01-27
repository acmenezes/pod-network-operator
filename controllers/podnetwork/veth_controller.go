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

// VethReconciler reconciles a Veth object
type VethReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	VethList *podnetworkv1alpha1.VethList
}

// +kubebuilder:rbac:groups=podnetwork.opdev.io,resources=veths,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=podnetwork.opdev.io,resources=veths/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=podnetwork.opdev.io,resources=veths/finalizers,verbs=update

// Reconcile function for Veth CRD instances
func (r *VethReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger = r.Log.WithName("podnetwork").WithValues("veth", req.NamespacedName)

	// get the list of Veths CR for pods managed by PNO
	// (Gotta put a knob to take control of existing network interfaces for small configs like MTU)
	r.VethList = &podnetworkv1alpha1.VethList{}
	err := r.Client.List(context.TODO(), r.VethList)
	if err != nil {
		return ctrl.Result{}, err
	}

	if len(r.VethList.Items) <= 0 {

		return ctrl.Result{}, nil

	}

	// setup the finalizers for veth controller

	for _, veth := range r.VethList.Items {

		finalizer := "finalizers.veth.podnetwork.opdev.io"
		// examine DeletionTimestamp to determine if veth is under deletion
		if veth.ObjectMeta.DeletionTimestamp.IsZero() {

			// veth is not being deleted, so if it does not have our finalizer,
			// then lets add the finalizer and update the object. This is equivalent
			// to registering our finalizer.

			if !containsString(veth.GetFinalizers(), finalizer) {
				veth.SetFinalizers(append(veth.GetFinalizers(), finalizer))
				if err := r.Update(context.Background(), &veth); err != nil {
					return ctrl.Result{}, err
				}
			}

		} else {

			// veth is being deleted
			if containsString(veth.GetFinalizers(), finalizer) {

				// finalizer is present, delete configurations

				// Get the pods with matching labels to veth

			}

		}

	}

	// check the state of each one and reconcile configurations

	// apply configs - may include many

	// write to status

	// Requeue

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VethReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&podnetworkv1alpha1.Veth{}).
		Complete(r)
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
