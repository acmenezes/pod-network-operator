package controllers

import (
	"fmt"

	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/vishvananda/netlink"
)

// Bridge creation logic
// TODO: verify existence first and only creates if doesn't exist
// TODO: do the clean up on podconfig deletion - probable use for a finalizer
func getBridgeOnHost(bridge string) error {

	targetNS, err := ns.GetNS("/tmp/proc/1/ns/net")
	if err != nil {
		return fmt.Errorf("error getting host network namespace: %v", err)
	}

	err = targetNS.Do(func(hostNs ns.NetNS) error {

		_, err := netlink.LinkByName(bridge)
		if err != nil {
			return fmt.Errorf("error looking up for bridge %v %v", bridge, err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
func createBridge(bridge string, ipAddr *netlink.Addr) error {

	targetNS, err := ns.GetNS("/tmp/proc/1/ns/net")
	if err != nil {
		return fmt.Errorf("error getting host network namespace: %v", err)
	}

	err = targetNS.Do(func(hostNs ns.NetNS) error {
		br := &netlink.Bridge{
			LinkAttrs: netlink.LinkAttrs{
				Name: bridge,
			},
		}
		// Creating bridge
		err := netlink.LinkAdd(br)
		if err != nil {
			return fmt.Errorf("failed to create bridge %v: %v", bridge, err)
		}

		// Setting bridge ip address
		err = netlink.AddrAdd(br, ipAddr)
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

func deleteBridge(bridge string) error {

	targetNS, err := ns.GetNS("/tmp/proc/1/ns/net")
	if err != nil {
		return fmt.Errorf("error getting host network namespace: %v", err)
	}

	err = targetNS.Do(func(hostNs ns.NetNS) error {

		br, err := netlink.LinkByName(bridge)
		if err != nil {
			return fmt.Errorf("error looking up for bridge %v %v", bridge, err)
		}

		err = netlink.LinkDel(br)
		if err != nil {
			return fmt.Errorf("failed to delete bridge %q: %v", bridge, err)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("%v\n", err)
	}

	return nil
}
