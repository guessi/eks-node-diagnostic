package k8s

import (
	"fmt"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
	"k8s.io/client-go/discovery"
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

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &CustomizedClient{
		client:          dynamicClient,
		discoveryClient: discoveryClient,
	}, nil
}

func (k *CustomizedClient) ValidateCRD() error {
	resources, err := k.discoveryClient.ServerResourcesForGroupVersion(
		fmt.Sprintf("%s/%s", constants.NodeDiagnosticResourceGroup, constants.NodeDiagnosticResourceVersion),
	)
	if err != nil {
		return fmt.Errorf("NodeDiagnostic CRD not found, please ensure the cluster has EKS Auto Mode enabled or Node Monitoring Agent installed: %w", err)
	}
	for _, r := range resources.APIResources {
		if r.Name == constants.NodeDiagnosticResourceName {
			return nil
		}
	}
	return fmt.Errorf("NodeDiagnostic CRD not found, please ensure the cluster has EKS Auto Mode enabled or Node Monitoring Agent installed")
}
