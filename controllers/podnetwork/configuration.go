package podnetwork

import (
	podnetworkv1alpha1 "github.com/opdev/pod-network-operator/apis/podnetwork/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// Strategy pattern with the configurator interface
// Every configuration will have its own specific configurator
// Each configurator implements its own logic using the same interface
// with apply and reset methods
type Configurator interface {
	Apply(*corev1.Pod) error
	Delete(*corev1.Pod) error
	Get(*corev1.Pod) error
}

type Configuration struct {
	Configurator Configurator
}

// The method Apply applies a new configuration for pods
func (c *Configuration) Apply(pod *corev1.Pod) error {
	return c.Configurator.Apply(pod)
}

// The method Delete removes configuration applied to pods
func (c *Configuration) Delete(pod *corev1.Pod) error {
	return c.Configurator.Delete(pod)
}

// The method Get brings back the CR/Object instance applied to a pod
func (c *Configuration) Get(pod *corev1.Pod) error {
	return c.Configurator.Delete(pod)
}

type AdditionalNets struct {
	NetworkList *[]podnetworkv1alpha1.AdditionalNetwork
}

func (AdditionalNets) Apply(pod *corev1.Pod) error {

	// 1. Check for existing interfaces with same parameters on pod and return error if so

	// 2. Check for switching / routing master device existence. Create if it doesn't exist.

	// 3. Create the new link.

	// 4. Move one interface to root namespace

	// 5. Setup the master device

	// 6. Set Link up

	// 7. Run connectivity tests

	// 8. Return configuration details (possibly a []map[string]string)

	return nil
}

func (AdditionalNets) Delete(pod *corev1.Pod) error {
	// SecondaryNetwork delete logic here
	return nil
}

func (AdditionalNets) Get(pod *corev1.Pod) error {
	// SecondaryNetwork Get logic here
	return nil
}
