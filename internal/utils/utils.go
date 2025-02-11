package utils

import (
	"fmt"
	"regexp"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
	"github.com/guessi/eks-node-diagnostic/internal/types"
	"github.com/guessi/eks-node-diagnostic/internal/variables"

	"github.com/urfave/cli/v2"
)

func Version() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		r := regexp.MustCompile(`v[0-9]\.[0-9]+\.[0-9]+`)
		fmt.Println(constants.AppName, r.FindString(variables.GitVersion))
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

	r := regexp.MustCompile(constants.NodeNameSuffixPattern)
	if len(nodeName) != constants.NodeNameLength || !r.MatchString(nodeName[2:]) {
		return fmt.Errorf("invalid node-name, exepected to have node name with pattern \"%s%s\"", constants.NodeNamePrefix, constants.NodeNameSuffixPattern)
	}
	return nil
}

func ValidateInRange(input, start, end int) error {
	if input < start || input > end {
		return fmt.Errorf("expire-seconds must be between %d and %d", start, end)
	}
	return nil
}

func ValidateAppConfigs(config types.AppConfigs) error {
	if err := ValidateEmpty("region", config.Region); err != nil {
		return err
	}
	for _, nodeName := range config.Nodes {
		if err := ValidateNodeName(nodeName); err != nil {
			return err
		}
	}
	if err := ValidateEmpty("bucket-name", config.BucketName); err != nil {
		return err
	}
	if err := ValidateInRange(config.ExpiredSeconds, constants.MinExpireSeconds, constants.MaxExpireSeconds); err != nil {
		return err
	}
	return nil
}
