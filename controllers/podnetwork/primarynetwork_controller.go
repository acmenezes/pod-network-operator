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
	"fmt"

	"github.com/go-logr/logr"
	podnetworkv1alpha1 "github.com/opdev/pod-network-operator/apis/podnetwork/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	reqLogger := r.Log.WithValues("primarynetwork", req.NamespacedName)

	reqLogger.Info("Starting Reconciliation...")

	r.PrimaryNetwork = &podnetworkv1alpha1.PrimaryNetwork{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, r.PrimaryNetwork)
	if err != nil {
		return ctrl.Result{}, err
	}

	// call finalizer on primary network configuration resource

	// update primary network status condition unknown

	condition := podnetworkv1alpha1.Condition{}

	condition.Type = podnetworkv1alpha1.ConditionTypeUnknown
	condition.Status = true
	condition.Reason = "BeginningReconcileFunction"
	condition.LastHeartbeatTime = ""
	condition.LastTransitionTime = ""

	r.PrimaryNetwork.Status.Conditions = append(r.PrimaryNetwork.Status.Conditions, condition)

	err = r.Client.Status().Update(context.TODO(), r.PrimaryNetwork)
	if err != nil {
		reqLogger.Error(err, "Couldn't update Primary Network", "PrimaryNetwork", r.PrimaryNetwork.ObjectMeta.Name)
		return ctrl.Result{}, err
	}

	// Loop through the list of pods with primary newtworks matching labels
	podList, err := listPodsWithMatchingLabels("PrimaryNetworkConfiguration", r.PrimaryNetwork.ObjectMeta.Name)
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, pod := range podList.Items {

		// if pod is not in running phase return
		if pod.Status.Phase != "Running" {
			fmt.Printf("pod %v phase is %v, requeuing... ", pod.ObjectMeta.Name, pod.Status.Phase)
			return ctrl.Result{}, nil
		}

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
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PrimaryNetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&podnetworkv1alpha1.PrimaryNetwork{}).
		Complete(r)
}

func (r *PrimaryNetworkReconciler) updateConditions(Type podnetworkv1alpha1.ConditionType, Status bool, Reason string) error {

	condition := podnetworkv1alpha1.Condition{}

	condition.Type = Type
	condition.Status = true
	condition.Reason = "BeginningReconcileFunction"

	condition.LastHeartbeatTime = ""
	condition.LastTransitionTime = ""

	r.PrimaryNetwork.Status.Conditions = append(r.PrimaryNetwork.Status.Conditions, condition)

	err := r.Client.Status().Update(context.TODO(), r.PrimaryNetwork)
	if err != nil {
		return err
	}
	return nil
}
