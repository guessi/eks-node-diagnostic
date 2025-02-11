package constants

const (
	// Application basic info
	AppName  string = "eks-node-diagnostic"
	AppUsage string = "friendly NodeDiagnostic generator"

	// Node under Auto Mode would be EC2 instance id, where would be prefixed with "i-" follow by 17-character IDs
	// - https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeIdFormat.html
	// - https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_ModifyIdFormat.html
	// - https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeAggregateIdFormat.html
	NodeNameLength        = 19
	NodeNamePrefix        = "i-"
	NodeNameSuffixPattern = "[a-f0-9]{17}"

	// NodeDiagnostic
	NodeDiagnosticApiVersion = "eks.amazonaws.com/v1alpha1"
	MinExpireSeconds         = 120
	MaxExpireSeconds         = 86400

	// Log generated pattern
	LogfileNamePattern = "node-diagnostic/log__%s__%s__%s.tar.gz"

	NodeDiagnosticResourceGroup   = "eks.amazonaws.com"
	NodeDiagnosticResourceVersion = "v1alpha1"
	NodeDiagnosticResourceKind    = "NodeDiagnostic"
	NodeDiagnosticResourceName    = "nodediagnostics"
)
