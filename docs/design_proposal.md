### Design Proposal
#### Motivation

Multiple networking challenges faced by Telco companies and Internet Service Providers working on 5G networks as well as other CNF driven solutions are not addressed by current working Kubernetes CNI standard. Some work is being done to address some of those challenges to be generally available in the long run. This operator aims at short-term solutions for a few use cases that can be solved by a kubernetes operator quickly given the cloud native nature of those while waiting for the official upstream solutions in the long run. It's NOT intended to replace the CNI plugins in any way. Its intention is to enable businesses as fast as we can while we don't have an official long-term solution.

#### Some of those challenges are:

- Runtime network configurations dynamically triggered by CNF applications
- Use of not so common protocols such as SCTP and its configuration on host
- Use of protocols or configurations not implemented by CNI plugins but offered by Linux
- Separation of those configurations and respective permissions per tenant
- A centralized point to monitor and configure pod networking without tenant access and completely secured at a host's administrative level
- Protect the actual CNI plugin in place on the host as well as the Linux host networking from possible disruptive proprietary tenant's software at host level.
- Have a common open source repository that can bring all partners to contribute together on a common solution for all those problems

#### Scope

The scope of this project is network configuration, monitoring and tuning at the pod level specially designed with CNF multi-tenant environments in mind. For that to be accomplished any host level configuration will be performed on behalf of the pods where it makes sense without disrupting the main CNI plugin used by the cluster, the main routing table on the host and the rules published on iptables by the current kube-proxy workload or other CNI related pieces of software.

#### That said What this operator is not:

This operator is NOT intended to configure anything cluster wide that affects all pods or an entire node behavior with that intent. Those tasks are already performed by a set of operators such as the ones below:

- `Machine Config Operator` - https://github.com/openshift/machine-config-operator
- `Performance Addons Operator` - https://github.com/openshift-kni/performance-addon-operators
- `Cluster Node Tuning Operator` - https://github.com/openshift/cluster-node-tuning-operator
- `Special Resource Operator` - https://github.com/openshift-psap/special-resource-operator
- `Node Feature Discovery Operator` - https://github.com/openshift/cluster-nfd-operator

This operator is NOT intended to own specific application deployments, daemonsets, statefulsets or any other application workloads that may result in pods or replicas. It SHOULD own only the configurations it provides gracefully terminating them when a Pod terminates, when the identifying label is deleted or when an entire or partial network configuration object is deleted.
Security Concerns

In order to make networking configurations on behalf of a pod some privileges must be in place available to the operator making those. Those privileges are needed for both pod and host level configs at this moment. Even with linux user namespaces, that will come in the future, it won't be possible to contain workload privileges if host level configurations are needed. And this imposes a privileged workload to take care of those tasks on behalf of pods. That may perfectly be an operator.

That said, the pod network operator should remain in a separate kubernetes namespace and all its RBAC resources should be separated on that kubernetes namespace as well to prevent any tenant from using them in the first place.

Network configurations will be defined in CRDs and can be seen by all tenants, but can only be used by those who have the proper RBAC permissions for that. No tenant will be able to change them. Only cluster admins.

#### Why an Operator?

OpenShift nodes are managed by Kubernetes operators. The operators take care of multiple disciplines on OCP nodes. That means that all those configurations are repeatable, reconcilable and visible in a cloud native way via the Kubernetes API extensions.

Any architecture that doesn't comply with that standard will have a hard time adapting to OCP's premises. Normally it won't be seen as a cloud native/container native way of managing the platform and will most probably bring a lot more work with it since it doesn't inherit all the already proven goods that come from the Kubernetes API itself.
Here below I put an extract with common features that we gain from using CRDs within kubernetes when we develop using the operator pattern.
Extracted from https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/.
Beyond all those good features from the extending the kubernetes API we still have the advantage of using the Operator Framework that will give us right away an embedded metrics endpoint to expose the metrics retrieved by the reconciler process as well as all the tooling to seamlessly publish this operator using OLM (Operator Lifecycle Manager) that is already native in OpenShift.

#### Architecture

#### Overview:

<img src='docs/img/pod_network_operator.png'></img>

1. Ideally a CNF operator spins up pods to run CNF applications in the tenants namespace with the proper permissions and running with the restricted SCC. It could be deployed by helm or yaml manifests as well;
2. The CNF requires privileged pod network configurations. Those could be at initialization time or runtime. The operator or the application itself may request those configurations as part of the deployment using basic yaml manifests or dynamically at CNFs lifetime on demand according to its own conditions by sending an update or create request to the kubeapi-server with the network configuration CRD.
3. An admission process will take place and both authentication and authorization will be checked before granting any access to the that specific network configuration CRD;
4. After authorizing, the kubeapi-server allows access to the network configuration;
5. Any configuration changes, creation or deletion will be watched and received by the pod network operator triggering the reconciler;
6. It may make any configuration on network interfaces provided by Linux or other available packages on behalf of the pod;
7. It may pass on sockets to the pods as needed without granting privileged access to that pod if it make sense to do so;
8. It change routes in the pods context to complete its configuration if needed;
9. Any kind of tunneling techniques available on Linux and other possible custom packages may be used to build on demand tunnels and connect the CNF pod with other workloads.
10. All those possible configurations will be completed with whatever is needed at host level as long as it doesn't disrupt the CNI plugin being used, the main host routing table, iptables and other main network stack components. Ideally all configurations will be using separate tables and configurations.

#### Fine Grained Permission Control

<img src='docs/img/multi-crd-network-operator.png'></img>

Every network configuration in the system will be available through an individual custom resource definition. That allows for fine grained RBAC rules to be implemented on top of them giving the administrators the freedom to grant or deny access to very specific actions on pod networking creating plenty of possibilities for multi-tenant environments setups.

1. The pod network operator may hold multiple controllers and have access to all CRDs in order to reconcile different networking configurations and components on multiple levels;

2. The CNF tenant will have only the permissions needed for the configuration subset it requires;

3. Each CRD as an API endpoint can be used as any other kubernetes object. They can be referenced in roles with the well known rules and verbs, combined with the proper role bindings having the proper service accounts as subjects.

The only caveat here is from the management perspective. A new CRD may be created to track down all the information from all the separate pieces of configuration into a single object for visibility and administration.

#### The Controller Workflow

<img src='docs/img/pod-network-controller.png'></img>

The controller workflow represented in the diagram above shows simplified steps on how the reconciliation process occurs. A few steps before actually running the configuration functions it's necessary to find out what pods need new configurations, grab the first container ID from the Pod resource object and pass it as a parameter with a ContainerStatusRequest to CRI-O. From the ContainerStatusResponse we can get the process ID for that container. It's the same process that crictl inspect does.

Here we have an important observation. If the configuration is to be available to a Pod (a.k.a. shared linux namespaces between containers) then the first container ID is fine. If it's container or even process specific (for application with more than one process "inside" a container) then the procedure is a little bit more complex than the one represented on the diagram above.

After that we can read the correct path /proc/<PID>/ns/<desired namespace> that has a symbolic link with the namespace type and inode for the namespace we want to jump in from the operator. Then we use the ns package from github.com/containernetworking/plugins. Within the proper namespace all changes will affect the desired container or pod.

When it comes to the host it's the same process. We move to PID number 1 and the desired namespace and that's all.

For the reconciler function itself it will depend on the logic and configurations that we're trying to achieve. Each different configuration object may require different libraries to perform the configurations.

#### Conclusion

We hope this design covers the lack of attention that very special applications in the Telecom and Edge computing domain may require for their next generation containerized cloud native initiatives.

