package k8s

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

type CustomizedClient struct {
	client          *dynamic.DynamicClient
	discoveryMapper *restmapper.DeferredDiscoveryRESTMapper
}
