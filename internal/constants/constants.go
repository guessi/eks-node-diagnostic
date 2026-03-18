package constants

const (
	// Application basic info
	AppName  string = "eks-node-diagnostic"
	AppUsage string = "friendly NodeDiagnostic generator"

	// NodeDiagnostic
	MinExpireSeconds     = 120
	MaxExpireSeconds     = 86400
	DefaultExpireSeconds = 300 // seconds

	// Timeout
	MinTimeout     = 10  // seconds
	MaxTimeout     = 300 // seconds
	DefaultTimeout = 30  // seconds

	// Destination types
	DestinationTypeS3   = "s3"
	DestinationTypeNode = "node"

	// Log generated pattern
	LogfileNamePattern = "node-diagnostic/log__%s__%s__%s.tar.gz"

	NodeDiagnosticResourceGroup   = "eks.amazonaws.com"
	NodeDiagnosticResourceVersion = "v1alpha1"
	NodeDiagnosticResourceKind    = "NodeDiagnostic"
	NodeDiagnosticResourceName    = "nodediagnostics"
)
