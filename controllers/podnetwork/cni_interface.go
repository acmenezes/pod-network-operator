package podnetwork

import (
	"fmt"

	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/vishvananda/netlink"
)

// OCP defaults - Those constants will be used to bring the interface back to its original state
// in case of deleting the configurations

const (
	cniMTU int = 8951
)

// ConfigPrimaryIfForPod applies configurations for the CNI provided interface
func (r *PodNetworkConfigReconciler) ConfigPrimaryIfForPod(pid string) (string, error) {

	targetNS, err := ns.GetNS("/tmp/proc/" + pid + "/ns/net")

	if err != nil {
		return "", fmt.Errorf("Error getting Pod network namespace: %v", err)
	}

	// The primary interface is named eth0 by default in Openshift

	cniInterface := "eth0"

	// The Do function takes care of all side effects of switching namespaces
	// and spawning new threads or child processes on the destination namespaces
	// Since targetNS belongs to pod all instructions enclosed by Do() will be run
	// on the pods namespace

	err = targetNS.Do(func(hostNs ns.NetNS) error {

		// verify the existence of the eth0 interface
		eth0, err := netlink.LinkByName(cniInterface)
		if err != nil {
			fmt.Printf("Primary interface %s couldn't be found on the Pod. Skipping configuration ...", cniInterface)
			return nil
		}

		// set the MTU for the eth0 interface - This operation should be idempotent and not cause any harm
		// TODO: needs testing for idempotency
		netlink.LinkSetMTU(eth0, r.podNetworkConfig.Spec.Eth0.MTU)
		if err != nil {
			fmt.Printf("MTU couldn't be set on eth0")
		}
		return nil
	})
	if err != nil {
		fmt.Printf("couldn't configure primary interface %s", cniInterface)
		return "", err
	}

	configStatus := fmt.Sprintf("MTU: %v", r.podNetworkConfig.Spec.Eth0.MTU)
	return configStatus, nil
}

// ResetCNIIfDefaults puts Pod's Eth0 in its initial state
func (r *PodNetworkConfigReconciler) ResetCNIIfDefaults(pid string) error {

	targetNS, err := ns.GetNS("/tmp/proc/" + pid + "/ns/net")

	if err != nil {
		return fmt.Errorf("Error getting Pod network namespace: %v", err)
	}

	// The primary interface is named eth0 by default in OpenShift

	cniInterface := "eth0"

	// The Do function takes care of all side effects of switching namespaces
	// and spawning new threads or child processes on the destination namespaces
	// Since targetNS belongs to pod all instructions enclosed by Do() will be run
	// on the pods namespace

	err = targetNS.Do(func(hostNs ns.NetNS) error {

		// verify the existence of the eth0 interface
		eth0, err := netlink.LinkByName(cniInterface)
		if err != nil {
			fmt.Printf("Primary interface %s couldn't be found on the Pod. Skipping configuration ...", cniInterface)
			return nil
		}

		// Reset the MTU for the eth0 interface - OCP default is 8951

		netlink.LinkSetMTU(eth0, cniMTU)
		if err != nil {
			fmt.Printf("MTU couldn't be set on eth0")
		}
		return nil
	})

	if err != nil {
		fmt.Printf("couldn't configure primary interface %s", cniInterface)
		return err
	}

	return err
}
