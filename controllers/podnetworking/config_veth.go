package controllers

import (
	"fmt"

	"github.com/containernetworking/plugins/pkg/ns"
	podconfigv1alpha1 "github.com/opdev/pod-network-operator/apis/podconfig/v1alpha1"
	"github.com/vishvananda/netlink"
)

func (pid string, networkAttachment podconfigv1alpha1.Link) (string, error) {

	type vethPodConfig struct {
		podVethName  string
		podIPAddr    string
		peerVethName string
		bridge       string
	}

	var vethConfig = vethPodConfig{}

	// Get the pods namespace object
	targetNS, err := ns.GetNS("/tmp/proc/" + pid + "/ns/net")

	if err != nil {
		return "", fmt.Errorf("Error getting Pod network namespace: %v", err)
	}

	// Appending the process id number to the names to identify the links
	// with the container processes

	podVethName := networkAttachment.Name + pid
	hostVethName := "h" + networkAttachment.Name + pid

	// The Do function takes care of all side effects of switching namespaces
	// and spawning new threads or child processes on the destination namespaces
	// Since targetNS belongs to pod all instructions enclosed by Do() will be run
	// on the pods namespace

	err = targetNS.Do(func(hostNs ns.NetNS) error {

		// Attempt to check the existence of the pod veth
		// If if already exists it skips creation and configuration
		_, err := netlink.LinkByName(podVethName)
		if err == nil {
			fmt.Printf("Veth link %s already exists on the Pod. Skipping creation ...", podVethName)
			return nil
		}

		veth := &netlink.Veth{
			LinkAttrs: netlink.LinkAttrs{
				Name: podVethName,
			},
			PeerName: hostVethName,
		}
		err = netlink.LinkAdd(veth)
		if err != nil {
			return fmt.Errorf("failed to set %q up: %w", podVethName, err)
		}

		// Get newly created pod link by name
		podVeth, err := netlink.LinkByName(podVethName)

		if err != nil {
			return fmt.Errorf("failed to lookup %q: %v", podVethName, err)
		}

		// Add ip address to pod veth
		addr := ips.getFreeIP(networkAttachment.CIDR)
		err = netlink.AddrAdd(podVeth, addr)
		if err != nil {
			return fmt.Errorf("failed to add IP addr to %q: %v", podVeth, err)
		}
		vethConfig.podIPAddr = fmt.Sprintf("%v", addr)

		// Set pod veth link up
		err = netlink.LinkSetUp(podVeth)
		if err != nil {
			return fmt.Errorf("failed to set %q up: %w", podVethName, err)
		}

		// Move host end of the link to the host and continue
		// the configuration from the host network namespace

		targetNS, err := ns.GetNS("/tmp/proc/1/ns/net")

		hostVeth, _ := netlink.LinkByName(hostVethName)

		err = netlink.LinkSetNsFd(hostVeth, int(targetNS.Fd()))
		if err != nil {
			return fmt.Errorf("failed to move veth to host netns: %v", err)
		}
		// }
		return nil
	})
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	targetNS, err = ns.GetNS("/tmp/proc/1/ns/net")

	err = targetNS.Do(func(hostNs ns.NetNS) error {

		// Get host veth link by name
		hostVeth, err := netlink.LinkByName(hostVethName)
		if err != nil {
			return fmt.Errorf("failed to lookup %q: %v", hostVethName, err)
		}

		if hostVeth.Attrs().OperState != netlink.OperUp {
			// Set host veth link up ( for PoC purposes it's only layer 2 on bridge)
			if err = netlink.LinkSetUp(hostVeth); err != nil {
				return fmt.Errorf("failed to set %q up: %w", hostVethName, err)
			}
		}

		// Set host veth link master bridge
		br, err := netlink.LinkByName(networkAttachment.Master)

		if hostVeth.Attrs().MasterIndex != br.Attrs().Index {

			err = netlink.LinkSetMaster(hostVeth, br)
			if err != nil {
				return fmt.Errorf("Error setting master device to %s: %v", hostVethName, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// Setup config information for pod
	vethConfig.podVethName = podVethName
	vethConfig.bridge = networkAttachment.Master
	vethConfig.peerVethName = hostVethName

	config := fmt.Sprintf("%+v", vethConfig)

	fmt.Println("Veth pair created successfully")
	return config, nil
}

func deleteVethForPod(pid string, networkAttachment podconfigv1alpha1.Link) error {

	targetNS, err := ns.GetNS("/tmp/proc/" + pid + "/ns/net")

	if err != nil {
		return fmt.Errorf("Error getting Pod network namespace: %v", err)
	}

	// Appending the process id number to the names to identify the links
	// with the container processes

	podVethName := networkAttachment.Name + pid

	// The Do function takes care of all side effects of switching namespaces
	// and spawning new threads or child processes on the destination namespaces
	// Since targetNS belongs to pod all instructions enclosed by Do() will be run
	// on the pods namespace

	err = targetNS.Do(func(hostNs ns.NetNS) error {
		// Get newly created pod link by name
		podVeth, err := netlink.LinkByName(podVethName)
		if err != nil {
			return fmt.Errorf("failed to lookup %q: %v", podVethName, err)
		}

		err = netlink.LinkDel(podVeth)
		if err != nil {
			return fmt.Errorf("failed to delete link %q: %v", podVethName, err)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	return nil
}
