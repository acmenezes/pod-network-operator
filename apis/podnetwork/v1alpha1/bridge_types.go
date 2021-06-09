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

// BridgeSpec defines the desired state of Bridge
type BridgeSpec struct {

	// Link layer configs
	Name string `json:"name,omitempty"`

	// possible configs

	// address             -- specify unicast link layer (MAC) address
	// arp                 -- change ARP flag on device
	// broadcast   brd     -- specify broadcast link layer (MAC) address
	// dynamic             -- change DYNAMIC flag on device
	// mtu                 -- specify maximum transmit unit
	// multicast           -- change MULTICAST flag on device
	// peer                -- specify peer link layer (MAC) address
	// promisc             -- set promiscuous mode
	// txqueuelen  txqlen  -- specify length of transmit queue

	// Network layer configs

	// CIDR is a temporary field to hold an IPv4 range
	// while we don't have functions acting on an Ipam plugin
	// Must be in the format "255.255.255.255/32"
	// Otherwise it will fail
	CIDR string `json:"cidr,omitempty"`
}

// BridgePhase for status
type BridgePhase string

// const values for BridgePhase
const (
	BridgePhaseUnset       BridgePhase = "unset"
	BridgePhaseConfiguring BridgePhase = "configuring"
	BridgePhaseConfigured  BridgePhase = "configured"
)

// BridgeStatus defines the observed state of Bridge
type BridgeStatus struct {
	Phase         BridgePhase `json:"phase,omitempty"`
	BridgeConfigs []string    `json:"bridgeConfigs,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Bridge is the Schema for the bridges API
type Bridge struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BridgeSpec   `json:"spec,omitempty"`
	Status BridgeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BridgeList contains a list of Bridge
type BridgeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bridge `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Bridge{}, &BridgeList{})
}
