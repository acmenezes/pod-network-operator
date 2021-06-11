package podnetwork

import (
	podnetworkv1alpha1 "github.com/opdev/pod-network-operator/apis/podnetwork/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// Strategy pattern with the linker interface
// Every link type will have its own specific linker
// Each linker implements its own logic using the same interface
// with apply and reset methods
type Linker interface {
	Apply(corev1.Pod, podnetworkv1alpha1.PodNetworkConfig) error
	Delete(corev1.Pod, podnetworkv1alpha1.PodNetworkConfig) error
	Get(corev1.Pod, podnetworkv1alpha1.PodNetworkConfig) error
}

type Link struct {
	Linker Linker
}

// The method Apply applies a new configuration for pods
func (c *Link) Apply(pod corev1.Pod, podNetworkConfig podnetworkv1alpha1.PodNetworkConfig) error {
	return c.Linker.Apply(pod, podNetworkConfig)
}

// The method Delete removes configuration applied to pods
func (c *Link) Delete(pod corev1.Pod, podNetworkConfig podnetworkv1alpha1.PodNetworkConfig) error {
	return c.Linker.Delete(pod, podNetworkConfig)
}

// The method Get brings back the CR/Object instance applied to a pod
func (c *Link) Get(pod corev1.Pod, podNetworkConfig podnetworkv1alpha1.PodNetworkConfig) error {
	return c.Linker.Get(pod, podNetworkConfig)
}
