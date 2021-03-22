package podnetwork

import (
	corev1 "k8s.io/api/core/v1"
)

// Strategy pattern with the configurator interface
// Every configuration will have its own specific configurator
// Each configurator implements its own logic using the same interface
// with apply and reset methods
type Configurator interface {
	Apply(*corev1.Pod) error
	Reset(*corev1.Pod) error
}

type Configuration struct {
	Configurator Configurator
}

// The method configure applies a new configuration for the pod
func (c *Configuration) Configure(pod *corev1.Pod) error {
	return c.Configurator.Apply(pod)
}

// The method reset may return configuration to last previous value
// or delete the configuration entirely depending on its specific use case
func (c *Configuration) Reset(pod *corev1.Pod) error {
	return c.Configurator.Reset(pod)
}

// MTU configurator implementation
type MTUConfig struct {
	vNIC string
	MTU  int
}

func (MTUConfig) Apply(pod *corev1.Pod) error {
	// mtu configuration logic here
	return nil
}

func (MTUConfig) Reset(pod *corev1.Pod) error {
	// mtu reset logic here
	return nil
}
