package k8s

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

type CustomizedClient struct {
	client          *dynamic.DynamicClient
	discoveryClient *discovery.DiscoveryClient
}
