package k8s

import (
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubernetesConfig() *rest.Config {
	cfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{})
	clientConfig, _ := cfg.ClientConfig()
	return clientConfig
}

func NewKubernetesClientSet(cfg *rest.Config) kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	return *clientset
}

func NewDynamicClient(cfg *rest.Config) dynamic.DynamicClient {
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	return *dynamicClient
}

func NewKubeClient() *CustomizedClient {
	cfg := NewKubernetesConfig()
	dynamicClient := NewDynamicClient(cfg)
	k8sClientSet := NewKubernetesClientSet(cfg)

	return &CustomizedClient{
		client:          &dynamicClient,
		discoveryMapper: restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(k8sClientSet.Discovery())),
	}
}
