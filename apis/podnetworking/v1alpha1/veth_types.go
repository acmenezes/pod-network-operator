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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VethSpec defines the desired state of Veth
type VethSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Name     string `json:"name,omitempty"`
	LinkType string `json:"linkType,omitemtpy"` // temporarily used for veth pair
	Parent   string `json:"parent,omitemtpy"`   // name for the parent interface
	Master   string `json:"master,omitempty"`   // name for the master bridge
	CIDR     string `json:"cidr,omitempty"`

	// For use with the netlink package  may access all types on the ip stack
	// Index        int                     `json:"index,omitempty"`
	MTU int `json:"mtu,omitempty"`
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
}

// VethStatus defines the observed state of Veth
type VethStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Veth is the Schema for the veths API
type Veth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VethSpec   `json:"spec,omitempty"`
	Status VethStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VethList contains a list of Veth
type VethList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Veth `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Veth{}, &VethList{})
}
