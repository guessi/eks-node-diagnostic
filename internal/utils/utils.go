package utils

import (
	"context"
	"fmt"
	"regexp"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
	"github.com/guessi/eks-node-diagnostic/internal/types"
	"github.com/guessi/eks-node-diagnostic/internal/variables"

	"github.com/urfave/cli/v3"
)

var (
	versionRegexp  = regexp.MustCompile(`v[0-9]{1,2}\.[0-9]+\.[0-9]+`)
	nodeNameRegexp = regexp.MustCompile(constants.NodeNameSuffixPattern)
)

func Version() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println(constants.AppName, versionRegexp.FindString(variables.GitVersion))
		fmt.Println(" Git Commit:", variables.GitVersion)
		fmt.Println(" Build with:", variables.GoVersion)
		fmt.Println(" Build time:", variables.BuildTime)
		return nil
	}
}

func ValidateEmpty(objectType, input string) error {
	if input == "" {
		return fmt.Errorf("%s must be set", objectType)
	}
	return nil
}

func ValidateNodeName(nodeName string) error {
	if err := ValidateEmpty("node-name", nodeName); err != nil {
		return err
	}

	if len(nodeName) != constants.NodeNameLength || nodeName[:2] != constants.NodeNamePrefix || !nodeNameRegexp.MatchString(nodeName[2:]) {
		return fmt.Errorf("invalid node-name, expected to have node name with pattern \"%s[a-f0-9]{17}\"", constants.NodeNamePrefix)
	}
	return nil
}

func ValidateInRange(fieldName string, input, start, end int) error {
	if input < start || input > end {
		return fmt.Errorf("%s must be between %d and %d", fieldName, start, end)
	}
	return nil
}

func ValidateAppConfigs(config types.AppConfigs) error {
	if err := ValidateEmpty("region", config.Region); err != nil {
		return err
	}
	if len(config.Nodes) == 0 {
		return fmt.Errorf("nodes must not be empty")
	}
	for _, nodeName := range config.Nodes {
		if err := ValidateNodeName(nodeName); err != nil {
			return err
		}
	}
	if err := ValidateEmpty("bucket-name", config.BucketName); err != nil {
		return err
	}
	// Allow 0 for expiredSeconds (will use default), otherwise validate range
	if config.ExpiredSeconds != 0 {
		if err := ValidateInRange("expire-seconds", config.ExpiredSeconds, constants.MinExpireSeconds, constants.MaxExpireSeconds); err != nil {
			return err
		}
	}
	// Allow 0 for timeout (will use default), otherwise validate range
	if config.Timeout != 0 {
		if err := ValidateInRange("timeout", config.Timeout, constants.MinTimeout, constants.MaxTimeout); err != nil {
			return err
		}
	}
	return nil
}
