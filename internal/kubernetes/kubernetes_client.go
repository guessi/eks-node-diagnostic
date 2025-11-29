package k8s

import (
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubernetesConfig() (*rest.Config, error) {
	cfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{})
	clientConfig, err := cfg.ClientConfig()
	if err != nil {
		return nil, err
	}
	return clientConfig, nil
}

func NewKubernetesClientSet(cfg *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func NewDynamicClient(cfg *rest.Config) (*dynamic.DynamicClient, error) {
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return dynamicClient, nil
}

func NewKubeClient() (*CustomizedClient, error) {
	cfg, err := NewKubernetesConfig()
	if err != nil {
		return nil, err
	}

	dynamicClient, err := NewDynamicClient(cfg)
	if err != nil {
		return nil, err
	}

	k8sClientSet, err := NewKubernetesClientSet(cfg)
	if err != nil {
		return nil, err
	}

	return &CustomizedClient{
		client:          dynamicClient,
		discoveryMapper: restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(k8sClientSet.Discovery())),
	}, nil
}
