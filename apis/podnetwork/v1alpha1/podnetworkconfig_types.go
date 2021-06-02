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

package v1alpha1

import (
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Link type for new Pod Interfaces
type LinkAttributes struct {
	// Device Info configurations
	NamePrefix string `json:"namePrefix,omitempty"`

	// Link Layer configurations
	LinkType string `json:"linkType,omitemtpy"` // temporarily used for veth pair
	Parent   string `json:"parent,omitemtpy"`   // name for the parent interface
	Master   string `json:"master,omitempty"`   // name for the master bridge
	MTU      int    `json:"mtu,omitempty"`

	// For use with the netlink package  may access all types on the ip stack
	// Index        int                     `json:"index,omitempty"`
	// MTU          int                     `json:"mtu,omitempty"`
	// TxQLen       int                     `json:"txqlen,omitempty"` // Transmit Queue Length
	// HardwareAddr net.HardwareAddr        `json:"hardwareAddr,omitempty"`
	// Flags        net.Flags               `json:"flags,omitempty"`
	// RawFlags     uint32                  `json:"rawFlags,omitempty"`
	// ParentIndex  int                     `json:"parentIndex,omitempty"` // index of the parent link device
	// MasterIndex  int                     `json:"masterIndex,omitempty"` // must be the index of a bridge
	// Alias        string                  `json:"alias,omitempty"`
	// Statistics   *netlink.LinkStatistics `json:"statistics,omitempty"`
	// Promisc      int                     `json:"promisc,omitempty"`
	// Xdp          *netlink.LinkXdp        `json:"xdp,omitempty"`
	// EncapType    string                  `json:"encapType,omitempty"`
	// Protinfo     *netlink.Protinfo       `json:"protinfo,omitempty"`
	// OperState    netlink.LinkOperState   `json:"operState,omitempty"`
	// NumTxQueues  int                     `json:"numTxQueues,omitempty"`
	// NumRxQueues  int                     `json:"numRxQueues,omitempty"`
	// GSOMaxSize   uint32                  `json:"gsoMaxSize,omitempty"`
	// GSOMaxSegs   uint32                  `json:"gsoMaxSegs,omitempty"`
	// Vfs          []netlink.VfInfo        `json:"vfs,omitempty"` // virtual functions available on link
	// Group        uint32                  `json:"group,omitempty"`
	// Slave        netlink.LinkSlave       `json:"slave,omitempty"`

	// Network Layer configurations
	CIDR string `json:"cidr,omitempty"`
}

// AdditionalNetwork for Pod configuration
type AdditionalNetwork struct {
	// Just a name to identify the network
	// Must be a short name (15 char) with no special characters
	NetworkName string `json:"networkName,omitempty"`
	// NetworkDescription should shortly describe the use for this network
	NetworkDescription string `json:"networkDescription,omitempty"`

	// Layer 2 configurations:

	// Intention is to support all Linux Link types
	// First one being implemented is Veth
	// Available link types can be found here:
	// https://man7.org/linux/man-pages/man8/ip-link.8.html
	LinkType string `json:"linkType,omitempty"`

	// Link Attribute (L2) Configurations below is a subset of netlink.LinkAttrs

	// Master device to attach the new network
	// If it doesn't existe it will be created with default options
	// If set empty it will set the link as no master link
	Master string `json:"master,omitempty"`
	// MTU if set empty it will be set to the default value of the
	// underlying OS.
	MTU int `json:"mtu,omitempty"`
	// Layer 2 address for the link being created
	HardwareAddr net.HardwareAddr `json:"hardwareAddr,omitempty"`
	// Alias for in system symbolic link description
	Alias string `json:"alias,omitempty"`

	// Ip or Layer 3 configurations:

	// CIDR is a temporary field to hold an IPv4 range
	// while we don't have functions acting on an Ipam plugin
	// Must be in the format "255.255.255.255/32"
	// Otherwise it will fail
	CIDR string `json:"cidr,omitempty"`
}

// PodNetworkConfigSpec defines the desired state of PodNetworkConfig
type PodNetworkConfigSpec struct {
	// This name represents the network profile desired for a set of pods
	// Pods containing the label PodNetworkConfig: with this name will trigger
	// the controller to add this additional network interface to the pod
	// Must be a short name with no special characters
	Name string `json:"name,omitempty"`

	// List of new links to be configured on Pod
	AdditionalNetworks []AdditionalNetwork `json:"additionalNetworks,omitempty"`
}

// PodNetworkConfiguration verified configs for status
type PodNetworkConfiguration struct {
	PodName    string   `json:"podName,omitempty"`
	ConfigList []string `json:"configList,omitempty"`
}

// PodNetworkConfigPhase type for status
type PodNetworkConfigPhase string

// pod network config phases
const (
	PodNetworkConfigUnset       PodNetworkConfigPhase = "unset"
	PodNetworkConfigConfiguring PodNetworkConfigPhase = "configuring"
	PodNetworkConfigConfigured  PodNetworkConfigPhase = "configured"
)

// PodNetworkConfigStatus defines the observed state of PodNetworkConfig
type PodNetworkConfigStatus struct {
	Phase                    PodNetworkConfigPhase     `json:"phase,omitempty"`
	PodNetworkConfigurations []PodNetworkConfiguration `json:"PodNetworkConfigs,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PodNetworkConfig is the Schema for the podnetworkconfigs API
type PodNetworkConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodNetworkConfigSpec   `json:"spec,omitempty"`
	Status PodNetworkConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PodNetworkConfigList contains a list of PodNetworkConfig
type PodNetworkConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodNetworkConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodNetworkConfig{}, &PodNetworkConfigList{})
}
