package utils

import (
	"fmt"
	"regexp"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
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
