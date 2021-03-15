#### Changing Eth0 Primary Interface's Configuration

Once we have our operator in place, supposedly deployed by the cluster's administrator or the platform owner, and also the CNF deployment running on a tenant namespace properly setup with tenant's permissions we may try requesting new network configurations for the Pods.

The first use case explored here is changing the Pod's MTU size in order to avoid heavy packet fragmentation on third party networks connecting clusters across different regions.

For that we may create a podNetworkConfig that looks like below:
```
apiVersion: podnetwork.opdev.io/v1alpha1
kind: PodNetworkConfig
metadata:
  name: podnetwork-sample-a
  namespace: cnf-telco
spec:
  Name: podnetwork-sample-a
  eth0:
      mtu: 1500

```
First let's check one of our sample deployment pods and take a look at it's primary interface the eth0:

`oc exec -it cnf-example-a-7cdb5b9fff-cr2pv -- /bin/bash`
```
bash-5.1$ ip link
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN mode DEFAULT group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
3: eth0@if432: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 8951 qdisc noqueue state UP mode DEFAULT group default
    link/ether 0a:58:0a:80:03:ab brd ff:ff:ff:ff:ff:ff link-netnsid 0
bash-5.1$
```
Notice that the OpenShift default's MTU for the eth0 interface is 8951 which may too big if the network traffic is being sent through non managed networks until it reaches it destiny. That's normally OK for general applications. But for applications dealing with jitter or latency sensitive data the fragmentation that will occur at some points on the path may impact performance on delivering Voice or Video streams for example.

In order to simulate a tenant with low privileges and request that default to be replaced by 1500, which is a pretty standard number, on a pod basis, we may run:

`oc apply -f config/samples/podnetwork_v1alpha1_podnetworkconfig.yaml --as=system:serviceaccount:cnf-telco:cnf-telco-sa`

The `--as` parameter allow us to impersonate the service account created for the tenant and run with it's privileges.

After applying that podNetworkConfig we should see that the interface eth0 on that pod we logged in now is 1500:

```
ip link
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN mode DEFAULT group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
3: eth0@if432: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP mode DEFAULT group default
    link/ether 0a:58:0a:80:03:ab brd ff:ff:ff:ff:ff:ff link-netnsid 0
```
And we may also see what pods requested that configuration by running `oc describe podnetworkconfig podnetwork-sample-a` that will present a status field like the one below:

```
Status:
  Pod Network Configs:
    Config List:
      MTU: 1500
    Pod Name:  cnf-example-a-7cdb5b9fff-tjrkk
    Config List:
      MTU: 1500
    Pod Name:  cnf-example-a-7cdb5b9fff-cr2pv
```

Finally, since the operator is constantly reconciling the target configurations with the podNetworkConfigs, if we delete we should see it coming back to its default values without disrupting the Pod and/or restarting the Pod at all.

For that it's enough to run:
`oc delete -f config/samples/podnetwork_v1alpha1_podnetworkconfig.yaml --as=system:serviceaccount:cnf-telco:cnf-telco-sa`

And we should see the MTU 8951 back on the Pod:

```
bash-5.1$ ip link
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN mode DEFAULT group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
3: eth0@if432: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 8951 qdisc noqueue state UP mode DEFAULT group default
    link/ether 0a:58:0a:80:03:ab brd ff:ff:ff:ff:ff:ff link-netnsid 0
bash-5.1$
```