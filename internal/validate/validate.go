package validate

import (
	"fmt"
	"regexp"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
	"github.com/guessi/eks-node-diagnostic/internal/types"
)

// Node name: instance ID (i-[a-f0-9]{17}) or private DNS IP name (ip-{IPv4}.{suffix})
// ref: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/hostname-types.html#ec2-instance-private-hostnames
// ref: https://docs.aws.amazon.com/global-infrastructure/latest/regions/aws-regions.html
var nodeNameRegexp = regexp.MustCompile(buildNodeNamePattern())

func buildNodeNamePattern() string {
	octet := `(25[0-5]|2[0-4]\d|1\d\d|[1-9]\d|\d)`
	firstOctet := `(25[0-5]|2[0-4]\d|1\d\d|[1-9]\d|[1-9])`
	ipv4 := firstOctet + `(-` + octet + `){3}`

	regionPrefix := `(af|ap|ca|eu|il|mx|me|sa)`
	regionDirection := `(east|west|north|south|central|northeast|southeast)`
	usRegion := `us-(east-[2-9]|west-[1-9])`
	region := `(` + usRegion + `|` + regionPrefix + `-` + regionDirection + `-[1-9])`
	dnsSuffix := `\.(ec2\.internal|` + region + `\.compute\.internal)`

	instanceID := `i-[a-f0-9]{17}`
	privateDNSIP := `ip-` + ipv4 + dnsSuffix

	return `^(` + instanceID + `|` + privateDNSIP + `)$`
}

func empty(objectType, input string) error {
	if input == "" {
		return fmt.Errorf("%s must be set", objectType)
	}
	return nil
}

func NodeName(nodeName string) error {
	if err := empty("node-name", nodeName); err != nil {
		return err
	}

	if !nodeNameRegexp.MatchString(nodeName) {
		return fmt.Errorf("invalid node-name: %s", nodeName)
	}
	return nil
}

func inRange(fieldName string, input, start, end int) error {
	if input < start || input > end {
		return fmt.Errorf("%s must be between %d and %d", fieldName, start, end)
	}
	return nil
}

func destinationType(dt string) error {
	switch dt {
	case constants.DestinationTypeS3, constants.DestinationTypeNode:
		return nil
	default:
		return fmt.Errorf("destination-type must be %q or %q", constants.DestinationTypeS3, constants.DestinationTypeNode)
	}
}

func AppConfigs(config types.AppConfigs) error {
	if err := empty("region", config.Region); err != nil {
		return err
	}

	if err := destinationType(config.DestinationType); err != nil {
		return err
	}

	if len(config.Nodes) == 0 {
		return fmt.Errorf("nodes must not be empty")
	}

	for _, nodeName := range config.Nodes {
		if err := NodeName(nodeName); err != nil {
			return err
		}
	}

	// S3-specific validations
	if config.DestinationType == constants.DestinationTypeS3 {
		if err := empty("bucket-name", config.BucketName); err != nil {
			return err
		}

		// Allow 0 for expiredSeconds (will use default), otherwise validate range
		if config.ExpiredSeconds != 0 {
			if err := inRange("expire-seconds", config.ExpiredSeconds, constants.MinExpireSeconds, constants.MaxExpireSeconds); err != nil {
				return err
			}
		}
	}

	// Allow 0 for timeout (will use default), otherwise validate range
	if config.Timeout != 0 {
		if err := inRange("timeout", config.Timeout, constants.MinTimeout, constants.MaxTimeout); err != nil {
			return err
		}
	}

	return nil
}
