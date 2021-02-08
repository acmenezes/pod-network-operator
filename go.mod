module github.com/opdev/pod-network-operator

go 1.15

require (
	github.com/containernetworking/plugins v0.9.0
	github.com/go-logr/logr v0.3.0
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.3
	github.com/vishvananda/netlink v1.1.1-0.20201029203352-d40f9887b852
	golang.org/x/mod v0.4.0 // indirect
	golang.org/x/tools v0.1.1-0.20210201201750-4d4ee958a9b7 // indirect
	golang.org/x/tools/gopls v0.6.5 // indirect
	google.golang.org/grpc v1.27.1
	honnef.co/go/tools v0.0.1-2020.1.6 // indirect
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/cri-api v0.20.2
	sigs.k8s.io/controller-runtime v0.7.0
)
