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

// LinkTemplateSpec represents the
type LinkTemplateSpec struct {

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

// LinkTemplateStatus defines the observed state of LinkTemplate
type LinkTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// LinkTemplate is the Schema for the linktemplates API
type LinkTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LinkTemplateSpec   `json:"spec,omitempty"`
	Status LinkTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LinkTemplateList contains a list of LinkTemplate
type LinkTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LinkTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LinkTemplate{}, &LinkTemplateList{})
}
