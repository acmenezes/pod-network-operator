package podnetwork

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

func (r *PodNetworkConfigReconciler) applyConfig(pod corev1.Pod) ([]string, error) {

	var configList []string
	// Get the first container pid for pod
	pid, err := getPid(pod)
	if err != nil {
		fmt.Printf("Error getting container pid %v", err)
		return []string{}, err
	}

	// Use the pid as input to configure Pod's CNI primary interface
	configStatus, err := r.ConfigPrimaryIfForPod(pid)

	configList = append(configList, configStatus)

	return configList, nil
}

func (r *PodNetworkConfigReconciler) deleteConfig(pod corev1.Pod) error {
	// Get the first container pid for pod
	pid, err := getPid(pod)
	if err != nil {
		fmt.Printf("Error getting container pid %v", err)
		return err
	}

	// Brings primary CNI interface to its defaults
	r.ResetCNIIfDefaults(pid)

	return nil
}
