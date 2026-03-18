package k8s

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubeClient() (*CustomizedClient, error) {
	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &CustomizedClient{
		client: dynamicClient,
	}, nil
}
