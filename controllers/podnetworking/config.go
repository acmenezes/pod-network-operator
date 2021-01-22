package controllers

import (
	"fmt"

	podconfigv1alpha1 "github.com/opdev/podconfig-operator/apis/podconfig/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func applyConfig(pod corev1.Pod, podconfig *podconfigv1alpha1.PodConfig) ([]string, error) {

	// Get the first container pid for pod
	pid, err := getPid(pod)
	if err != nil {
		fmt.Printf("Error getting container pid %v", err)
		return []string{}, err
	}

	configList, err := createNetworkAttachments(pid, podconfig.Spec.NetworkAttachments)
	if err != nil {
		fmt.Printf("Error creating network attachments: %v\n", err)
		return configList, err
	}

	return configList, nil
}

func deleteConfig(pod corev1.Pod, podconfig *podconfigv1alpha1.PodConfig) error {
	// Get the first container pid for pod
	pid, err := getPid(pod)
	if err != nil {
		fmt.Printf("Error getting container pid %v", err)
		return err
	}

	err = deleteNetworkAttachments(pid, podconfig.Spec.NetworkAttachments)
	if err != nil {
		fmt.Printf("Error creating network attachments: %v\n", err)
		return err
	}
	return nil
}

func createNetworkAttachments(pid string, networkAttachments []podconfigv1alpha1.Link) ([]string, error) {

	configList := []string{}

	for _, na := range networkAttachments {

		err := getBridgeOnHost(na.Master)

		if err != nil {

			fmt.Printf("%v\n", err)
			fmt.Println("Creating bridge on Host.")

			// Create bridge in host namespace
			err := createBridge(na.Master, ips.getFreeIP(na.CIDR))
			if err != nil {
				fmt.Printf("Error creating bridge device %s: %v\n", na.Master, err)
				return configList, err
			}

		}

		// Create veth pairs for the new networkAttachment
		config, err := createVethForPod(pid, na)
		if err != nil {
			fmt.Printf("Error creating new veth pair for pod: %v\n", err)
			return configList, err
		}
		configList = append(configList, config)
	}

	fmt.Println("New network attachment created successfully.")
	return configList, nil
}

func deleteNetworkAttachments(pid string, networkAttachments []podconfigv1alpha1.Link) error {

	for _, na := range networkAttachments {

		// delete veth pair for pod network attachments
		err := deleteVethForPod(pid, na)
		if err != nil {
			fmt.Printf("Error deleting new veth pair for pod: %v\n", err)
			return err
		}

		// delete remaining bridge
		err = deleteBridge(na.Master)
		if err != nil {
			fmt.Printf("Error creating bridge device %s: %v\n", na.Master, err)
			return err
		}
	}
	return nil
}
