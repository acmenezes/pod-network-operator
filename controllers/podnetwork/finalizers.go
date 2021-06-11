package podnetwork

import (
	"context"
)

func (r *PodNetworkConfigReconciler) Finalizer(finalizer string) (bool, error) {

	// examine DeletionTimestamp to determine if podConfig is under deletion
	if r.podNetworkConfig.ObjectMeta.DeletionTimestamp.IsZero() {

		// podNetworkConfig is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.

		if !containsString(r.podNetworkConfig.GetFinalizers(), finalizer) {
			r.podNetworkConfig.SetFinalizers(append(r.podNetworkConfig.GetFinalizers(), finalizer))
			if err := r.Update(context.Background(), r.podNetworkConfig); err != nil {
				return false, err
			}
		}
		return false, nil
	} else {
		// podNetworkConfig is being deleted
		if containsString(r.podNetworkConfig.GetFinalizers(), finalizer) {
			podList, err := listPodsWithMatchingLabels("podNetworkConfig", r.podNetworkConfig.ObjectMeta.Name)
			if err != nil {
				return true, err
			}
			for _, pod := range podList.Items {

				// Deleting all veth additional networks
				Veth := Link{&Veth{}}
				Veth.Delete(pod, *r.podNetworkConfig)

			}
			// remove our finalizer from the list and update it.
			r.podNetworkConfig.SetFinalizers(removeString(r.podNetworkConfig.GetFinalizers(), finalizer))
			if err := r.Update(context.Background(), r.podNetworkConfig); err != nil {
				return true, err
			}
		}
	}
	return true, nil
}

func (r *BridgeReconciler) Finalizer(finalizer string) (bool, error) {

	// examine DeletionTimestamp to determine if podConfig is under deletion
	if r.bridge.ObjectMeta.DeletionTimestamp.IsZero() {

		// bridge is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.

		if !containsString(r.bridge.GetFinalizers(), finalizer) {
			r.bridge.SetFinalizers(append(r.bridge.GetFinalizers(), finalizer))
			if err := r.Update(context.Background(), &r.bridge); err != nil {
				return false, err
			}
		}
		return false, nil
	} else {
		// bridge is being deleted
		if containsString(r.bridge.GetFinalizers(), finalizer) {

			br := Bridger{}
			err := br.Delete(r.bridge)
			if err != nil {
				return true, err
			}

			// remove our finalizer from the list and update it.
			r.bridge.SetFinalizers(removeString(r.bridge.GetFinalizers(), finalizer))
			if err := r.Update(context.Background(), &r.bridge); err != nil {
				return true, err
			}
		}
	}
	return true, nil

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
