package podnetwork

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// Temporary IP logic to facilitate the PoC
// Needs a real IPAM software

type ipsInUse struct {
	ipList []string
}

var ips = &ipsInUse{ipList: []string{}}

func (ips *ipsInUse) Contains(ip net.IP) bool {

	for _, element := range ips.ipList {
		if element == fmt.Sprintf("%v", ip) {
			return true
		}
	}
	return false
}

func (ips *ipsInUse) AllocateIP(ip net.IP) {

	ips.ipList = append(ips.ipList, fmt.Sprintf("%v", ip))

}

func (ips *ipsInUse) getFreeIP(network string) *netlink.Addr {

	ipv4Addr, _, _ := net.ParseCIDR(network)
	ip := net.ParseIP(fmt.Sprintf("%v", ipv4Addr))
	ip = ip.To4()

	for i := 1; i <= 254; i++ {
		ip[3]++
		if !ips.Contains(ip) {

			ips.AllocateIP(ip)
			addr, _ := netlink.ParseAddr(fmt.Sprintf("%v/24", ip))
			return addr

		}

	}
	return nil
}
