package validators

import (
	"fmt"
	"regexp"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
)

func ValidateNodeName(nodeName string) error {
	if err := ValidateEmpty("node-name", nodeName); err != nil {
		return err
	}

	if nodeName[:2] != constants.NodeNamePrefix {
		return fmt.Errorf("node-name must start with '%s'", constants.NodeNamePrefix)
	}

	expectedNodeNameLength := len(constants.NodeNamePrefix) + constants.NodeNameLength
	if len(nodeName) != expectedNodeNameLength {
		return fmt.Errorf("node-name must be equal to %d characters", expectedNodeNameLength)
	}

	r := regexp.MustCompile(constants.NodeNameSuffixPattern)
	if !r.MatchString(nodeName[2:]) {
		return fmt.Errorf("node-name invalid")
	}

	return nil
}

func ValidateEmpty(objectType, input string) error {
	if input == "" {
		return fmt.Errorf("%s must be set", objectType)
	}
	return nil
}

func ValidateInRange(input, start, end int) error {
	if input < start || input > end {
		return fmt.Errorf("expire-seconds must be between %d and %d", start, end)
	}
	return nil

}
