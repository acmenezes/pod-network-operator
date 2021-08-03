# Pod Network Operator
## Open Source Innovative Kubernetes Networking  
---
### Abstract

Computing is being pushed to the edge, internet of things is becoming a reality and containerized multi-tenant platforms are the base for a hybrid cloud. Services are a much stronger business drivers when they are packed together in a solid way in the same point of delivery or platform. Operators demonstrated their strength gaining immense space both on CNCF through the [Operator Framework](https://www.cncf.io/projects/) being accepted as an incubating project and growing number of operators on the operator hubs putting services side by side for kubernetes distributions. Finally the raise of 5G networks, an intricate technology connecting multiple of those layers together come to the table.

With those elements being painted on the scenario the need for more dynamic, flexible and intelligent networks capable of making changes or scheduling resources based on telemetry data becomes something real. Being able to create networks on the fly and delete them as a 5G call terminates, slice networks through multiple hops to separate traffic, classify packets and traffic and treat them accordingly, schedule resources according to network latency are some of the challenges faced by the companies writing the future.

What if we can add a service to our operator hub that can provide all of that? The goal of the pod-network-operator is to experiment, try and test innovative ways to do networking, re-inventing Kubernetes networking for those companies that are creating that new scenario.

### Scope

The scope of this project is network configuration and tuning at the pod level specially designed with CNF multi-tenant environments in mind. For that to be accomplished any <b>host level</b> configuration will be performed on behalf of the pods where it makes sense without disrupting the main CNI plugin used by the cluster, the main routing table on the host and the rules published on iptables by the current kube-proxy workload or other CNI related pieces of software. This behavior on the host of not overriding anything created by the main CNI plugin may be changed as an optional feature but not a default one.

#### That said What this operator is not:

This operator is NOT intended to configure anything cluster wide that affects all pods or an entire node behavior with that intent. Those tasks are already performed by a set of operators such as the ones below:

- `Machine Config Operator` - https://github.com/openshift/machine-config-operator
- `Performance Addons Operator` - https://github.com/openshift-kni/performance-addon-operators
- `Cluster Node Tuning Operator` - https://github.com/openshift/cluster-node-tuning-operator
- `Special Resource Operator` - https://github.com/openshift-psap/special-resource-operator
- `Node Feature Discovery Operator` - https://github.com/openshift/cluster-nfd-operator

And others.

This operator is NOT intended to own specific application deployments, daemonsets, statefulsets or any other application workloads that may result in pods or replicas. It SHOULD own only the configurations it provides gracefully terminating them when a Pod terminates, when the identifying label is deleted or when an entire or partial network configuration object is deleted.

#### Dynamic network use cases:

    Edge computing, IoT and Embedded Systems

    5G networks, network slicing

    Network service chaining

    Complex application meshes

<!-- find solid literature and numbers supporting this -->

#### Security Requirements for Multi-tenant environments

<b>CNFs and the need for privileges</b>

CNF applications often need to deal with additional network configuration and usage in many different ways depending on the application and vendor. That almost always requires privileges on the network level. May it be a Linux capability like CAP_NET_ADMIN or even running some init container or sidecar container under the root user to perform configurations that the application wouldn't be able to do so without the privileges.

When we step in the multi-tenant platforms domain, the privileges associated with those applications are completely undesired and often questioned if they shouldn't be removed from those applications. That brings us to the second key point on the pod-network-operator strategy: delegation of privileges by abstracting away any implementation details and hiding that complexity from the end user.

<!-- put a side note on multi-tenancy here. Important to get the growth of that kind of platform -->

In order to make networking configurations on behalf of a pod some privileges are required to be used by the operator. Those privileges are needed for both pod and host level configurations. Even with linux user namespaces, that will come in at some point in the future to Kubernetes, host level configurations required to complete connectivity to, from or between pods will always demand privileged actions. That puts pod-network-operator on the category of a trusted privileged platform service for all tenants.

<b>A Note on Linux Capabilities and Privilege Delegation</b>

When configuring additional networks with special purposes with an specific feature in mind it often DOESN'T suffice to perform one task alone. It's necessary to run multiple tasks from interface creation, configuration and connection to routing rules or vlan tagging for example. That set of tasks may require CAP_NET_ADMIN, CAP_NET_BIND and CAP_NET_RAW together to be able to accomplish a single goal given that each capability will have its own set of possible privileged tasks for example. Those often need to be combined together to achieve the results required by the business.

Our primary goal is <b>NOT</b> to delegate any specific capability privilege mapping specific tasks or set of tasks to an API with an operator controller behind it. Instead simple abstracted APIs hiding implementation details should be the way to go. We talk a little bit more about that in the next section. Expanding Kubernetes networks to be allocatable resources with `Kubernetes Network Profiles`. An approach to intelligent and resource aware network management.

#### Research and Application Approach 

   Approach: KEP Network Profiles

    Possible outcomes:

    Enabling New protocols (Implemented by Linux but not by CNI or Kubernetes)

      SRv6
      Full IPv6 features
      SCTP

    Dynamic Plugable Routing Protocols (BPG boxes, ISIS, OSPF, Trill etc.)

    Other data planes than the kernel (Rewiring the additional networks)

### Definitions





### General Design