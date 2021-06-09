package podnetwork

import (
	"fmt"

	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/opdev/pod-network-operator/apis/podnetwork/v1alpha1"
	"github.com/vishvananda/netlink"
)

type Bridger struct{}

func (b *Bridger) Apply(bridge v1alpha1.Bridge) error {

	targetNS, err := ns.GetNS("/tmp/proc/1/ns/net")
	if err != nil {
		return fmt.Errorf("error getting host network namespace: %v", err)
	}

	err = targetNS.Do(func(hostNs ns.NetNS) error {

		// Attempt to check the existence of the bridge
		// If if already exists it skips creation and configuration
		_, err := netlink.LinkByName(bridge.Spec.Name)
		if err == nil {
			fmt.Printf("Bridge %s already exists on the Pod. Skipping creation ...", bridge.Spec.Name)
			return nil
		}

		br := &netlink.Bridge{
			LinkAttrs: netlink.LinkAttrs{
				Name: bridge.Spec.Name,
			},
		}
		// Creating bridge
		err = netlink.LinkAdd(br)
		if err != nil {
			return fmt.Errorf("failed to create bridge %v: %v", bridge.Spec.Name, err)
		}

		// Setting bridge ip address
		err = netlink.AddrAdd(br, ips.getFreeIP(bridge.Spec.CIDR))
		if err != nil {
			return fmt.Errorf("failed to set bridge ip address: %v", err)
		}

		// Setting bridge up
		err = netlink.LinkSetUp(br)
		if err != nil {
			return fmt.Errorf("failed to set bridge up: %v", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Println("Bridge created successfully on the host.")
	return nil

}

func (b *Bridger) Delete(bridge v1alpha1.Bridge) error {
	targetNS, err := ns.GetNS("/tmp/proc/1/ns/net")
	if err != nil {
		return fmt.Errorf("error getting host network namespace: %v", err)
	}

	err = targetNS.Do(func(hostNs ns.NetNS) error {

		br, err := netlink.LinkByName(bridge.Spec.Name)
		if err != nil {
			return fmt.Errorf("error looking up for bridge %v %v", bridge.Spec.Name, err)
		}

		err = netlink.LinkDel(br)
		if err != nil {
			return fmt.Errorf("failed to delete bridge %q: %v", bridge.Spec.Name, err)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("%v\n", err)
	}

	return nil
}
