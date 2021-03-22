package podnetwork

import (
	"context"

	podnetworkv1alpha1 "github.com/opdev/pod-network-operator/apis/podnetwork/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type FinalizerFunc func(obj runtime.Object) error

type Finalizer struct {
	client.Client
}

func (r *PrimaryNetworkReconciler) RegisterFinalizer(finalizer string) error {

	if !containsString(r.PrimaryNetwork.GetFinalizers(), finalizer) {
		r.PrimaryNetwork.SetFinalizers(append(r.PrimaryNetwork.GetFinalizers(), finalizer))
		if err := r.Update(context.Background(), r.PrimaryNetwork); err != nil {
			return err
		}
	}
	return nil
}

func (r *PrimaryNetworkReconciler) ExecuteFinalizer(finalizer string, resetConfigs func(*podnetworkv1alpha1.PrimaryNetwork) error) error {

	if containsString(r.PrimaryNetwork.GetFinalizers(), finalizer) {

		// finalizer is present, undo configurations

		// Delete configuration defined in the podconfig CR from pods with the appropriate label.

		if err := resetConfigs(r.PrimaryNetwork); err != nil {
			// if fail to delete the external dependency here, return with error
			// so that it can be retried
			return err
		}

		// remove our finalizer from the list and update it.
		r.PrimaryNetwork.SetFinalizers(removeString(r.PrimaryNetwork.GetFinalizers(), finalizer))
		if err := r.Update(context.Background(), r.PrimaryNetwork); err != nil {
			return err
		}
	}
	return nil
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
