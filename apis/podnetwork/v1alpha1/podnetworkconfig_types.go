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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Link type for new Pod Interfaces
type Link struct {
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

// CNIInterface is a type to make adjustments on the primary interface
// provided by the CNI plugin
type CNIInterface struct {

	// for applications that need a specific MTU
	// setting other then the CNI configured MTU
	// this will work on a per pod basis
	// not cluster wide
	MTU int `json:"mtu,omitempty"`
}

// PodNetworkConfigSpec defines the desired state of PodNetworkConfig
type PodNetworkConfigSpec struct {

	// List of new links to be configured on Pod
	Links []Link `json:"links,omitempty"`
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
	Phase                    PodNetworkConfigPhase `json:"phase,omitempty"`
	PodNetworkConfigurations []PodNetworkConfig    `json:"PodNetworkConfigs,omitempty"`
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
