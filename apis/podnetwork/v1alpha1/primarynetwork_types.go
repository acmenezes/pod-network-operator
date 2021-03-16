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
	net
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PrimaryNetworkSpec for extra configurations on the primary network
// provided by the CNI plugin
type PrimaryNetworkSpec struct {

	MTU int `json:"mtu,omitempty"`
	TxQLen int `json:"txQLen,omitempty"`
	HardwareAddr net.HardwareAddr `json:"HardwareAddr,omitempty"`

}

type ConditionType string

const (
	ConditionTypeReady	ConditionType = "Ready"
	ConditionTypeInProgress ConditionType = "InProgress"
	ConditionTypeFailed ConditionType = "Failed"
	ConditionTypeUnknown ConditionType = "Unknown"
)

type Condition struct {
	Type ConditionType `json:"type,omitempty"`
	Status bool `json:"status,omitemtpy"`
	Reason string `json:"reason,omitempty"`
	LastHeartbeatTime string `json:"lastHeartbeatTime,omitempty"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
}

// PrimaryNetworkStatus defines the observed state of PrimaryNetwork
type PrimaryNetworkStatus struct {

	Conditions []Condition `json:"condition,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PrimaryNetwork is the Schema for the primarynetworks API
type PrimaryNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PrimaryNetworkSpec   `json:"spec,omitempty"`
	Status PrimaryNetworkStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PrimaryNetworkList contains a list of PrimaryNetwork
type PrimaryNetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PrimaryNetwork `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PrimaryNetwork{}, &PrimaryNetworkList{})
}
