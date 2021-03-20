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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	podnetworkv1alpha1 "github.com/opdev/pod-network-operator/apis/podnetwork/v1alpha1"
)

// PodNetworkConfigReconciler reconciles a PodNetworkConfig object
type PodNetworkConfigReconciler struct {
	client.Client
	Log                  logr.Logger
	Scheme               *runtime.Scheme
	podNetworkConfigList *podnetworkv1alpha1.PodNetworkConfigList
	podNetworkConfig     *podnetworkv1alpha1.PodNetworkConfig
}

// controller-gen flags to generate rbac rules

// +kubebuilder:rbac:groups=podnetwork.opdev.io,resources=podnetworkconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=podnetwork.opdev.io,resources=podnetworkconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=podnetwork.opdev.io,resources=podnetworkconfigs/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=daemonsets;deployments;deployments/finalizers;replicasets,verbs=get;list;watch;create;update;patch;delete,namespace=cnf-test

// +kubebuilder:rbac:groups="*",resources="*",verbs="*"

// Reconcile for podnetwork configs
func (r *PodNetworkConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("podnetworkconfig", req.NamespacedName)

	// Get the list of all pod network configurations to be applied
	reqLogger.Info("loading pod network configurations")

	r.podNetworkConfigList = &podnetworkv1alpha1.PodNetworkConfigList{}
	err := r.Client.List(context.TODO(), r.podNetworkConfigList)
	if err != nil {
		return ctrl.Result{}, err
	}

	if len(r.podNetworkConfigList.Items) <= 0 {
		return ctrl.Result{}, nil
	}
	// TODO: Update the status field with conditions while creating the new instance

	for _, podNetworkConfig := range r.podNetworkConfigList.Items {
		// Check the finalizer for all of them and set finalizers for the ones that don't have it
		r.podNetworkConfig = &podNetworkConfig
		finalizer := "podnetworkconfig.finalizers.opdev.io"

		// Check if the item is being deleted

		// examine DeletionTimestamp to determine if podNetworkConfig is under deletion
		if r.podNetworkConfig.ObjectMeta.DeletionTimestamp.IsZero() {

			// podNetworkConfig is not being deleted, so if it does not have our finalizer,
			// then lets add the finalizer and update the object. This is equivalent
			// registering our finalizer.

			if !containsString(r.podNetworkConfig.GetFinalizers(), finalizer) {
				r.podNetworkConfig.SetFinalizers(append(r.podNetworkConfig.GetFinalizers(), finalizer))
				if err := r.Update(context.Background(), r.podNetworkConfig); err != nil {
					return ctrl.Result{}, err
				}
			}
		} else {
			// podNetworkConfig is being deleted
			if containsString(r.podNetworkConfig.GetFinalizers(), finalizer) {

				// finalizer is present, delete configurations

				// Get the pods with matching labels to podConfig
				podList, err := listPodsWithMatchingLabels("podNetworkConfig", r.podNetworkConfig.ObjectMeta.Name)
				if err != nil {
					return ctrl.Result{}, err
				}
				// Delete configuration defined in the podconfig CR from pods with the appropriate label.
				for _, pod := range podList.Items {

					if err := r.deleteConfig(pod); err != nil {
						// if fail to delete the external dependency here, return with error
						// so that it can be retried
						return ctrl.Result{}, err
					}
				}

				// remove our finalizer from the list and update it.
				r.podNetworkConfig.SetFinalizers(removeString(r.podNetworkConfig.GetFinalizers(), finalizer))
				if err := r.Update(context.Background(), r.podNetworkConfig); err != nil {
					return ctrl.Result{}, err
				}
			}

			// Stop reconciliation as the item is being deleted
			return ctrl.Result{}, nil
		}
		// if not being deleted gather the list of pods for each item present on the podnetwork config list by label or annotation
		podList, err := listPodsWithMatchingLabels("podNetworkConfig", r.podNetworkConfig.ObjectMeta.Name)
		if err != nil {
			return ctrl.Result{}, err
		}
		for _, pod := range podList.Items {

			// if pod is not in running phase return
			if pod.Status.Phase != "Running" {
				fmt.Printf("pod %v phase is %v, requeuing... ", pod.ObjectMeta.Name, pod.Status.Phase)
				return ctrl.Result{}, nil
			}

			// begin to reconcile each config element present on the pod network config for each pod with the appropriate annotation or label

			// TODO: to be separated into multiple tasks each on its own file but called by apply config

			// verify the existence of secondary interfaces named in the CR for the pod
			// config_veth.go

			// if they exist reconcile/apply config

			// if they don't exist create and apply config

			// verify vlan IDs to be created on specified secondary interfaces by the CR
			// create a new logic under vlans.go

			// if they exist reconcile/apply config

			// if they don't exist create and apply config

			configList, err := r.applyConfig(pod)
			if err != nil {
				fmt.Printf("%v", err)
				return ctrl.Result{}, nil
			}
			fmt.Printf("%v", configList)

			// Update config status for the actual pod in the list
			configStatus := podnetworkv1alpha1.PodNetworkConfiguration{PodName: pod.ObjectMeta.Name, ConfigList: configList}
			fmt.Printf("%v", r.podNetworkConfig.Status.PodNetworkConfigurations)

			// Refresh cached object to avoid conflicts
			if err := r.Client.Get(context.TODO(), req.NamespacedName, r.podNetworkConfig); err != nil {
				fmt.Printf("%v", err)
				return ctrl.Result{}, err
			}

			// If the pod config didn't reconcile completely update status
			if r.podNetworkConfig.Status.Phase != podnetworkv1alpha1.PodNetworkConfigConfigured {

				isPodNamePresent := false

				for _, p := range r.podNetworkConfig.Status.PodNetworkConfigurations {
					if p.PodName == configStatus.PodName {
						isPodNamePresent = true
					}
				}

				if isPodNamePresent == false {

					r.podNetworkConfig.Status.PodNetworkConfigurations = append(r.podNetworkConfig.Status.PodNetworkConfigurations, configStatus)

					fmt.Printf("%v", r.podNetworkConfig.Status.PodNetworkConfigurations)

					if err := r.Client.Status().Update(context.TODO(), r.podNetworkConfig); err != nil {
						fmt.Printf("%v", err)
						return ctrl.Result{}, err
					}
				}
			}
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodNetworkConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&podnetworkv1alpha1.PodNetworkConfig{}).
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
