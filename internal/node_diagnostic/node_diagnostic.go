package node_diagnostic

import (
	"os"
	"text/template"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
)

const NodeDiagnosticObjectTemplate = `---
apiVersion: {{ .ApiVersion }}
kind: NodeDiagnostic
metadata:
    name: {{ .NodeName }}
spec:
    logCapture:
        destination: {{ .Destination }}
`

type NodeTemplateData struct {
	ApiVersion  string
	NodeName    string
	Destination string
}

func Render(nodeName, presignPutObjectUrl string) error {
	t := template.Must(template.New("node-diagnostic").Parse(NodeDiagnosticObjectTemplate))
	data := NodeTemplateData{
		ApiVersion:  constants.NodeDiagnosticApiVersion,
		NodeName:    nodeName,
		Destination: presignPutObjectUrl,
	}
	return t.Execute(os.Stdout, data)
}
