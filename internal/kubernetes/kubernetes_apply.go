package k8s

import (
	"context"
	"fmt"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (k *CustomizedClient) Apply(ctx context.Context, node, url string) error {
	// TODO: apply with structured object
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": fmt.Sprintf("%s/%s", constants.NodeDiagnosticResourceGroup, constants.NodeDiagnosticResourceVersion),
			"kind":       constants.NodeDiagnosticResourceKind,
			"metadata": map[string]interface{}{
				"name": node,
			},
			"spec": map[string]interface{}{
				"logCapture": map[string]interface{}{
					"destination": url,
				},
			},
		},
	}
	gvr := schema.GroupVersionResource{
		Group:    constants.NodeDiagnosticResourceGroup,
		Version:  constants.NodeDiagnosticResourceVersion,
		Resource: constants.NodeDiagnosticResourceName,
	}
	_, err := k.client.Resource(gvr).Create(ctx, obj, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		return fmt.Errorf("nodediagnostic %q already exists, please delete it before retrying", node)
	}
	return err
}
