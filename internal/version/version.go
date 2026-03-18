package version

import (
	"context"
	"fmt"
	"regexp"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
	"github.com/guessi/eks-node-diagnostic/internal/variables"

	"github.com/urfave/cli/v3"
)

var versionRegexp = regexp.MustCompile(`v[0-9]{1,2}\.[0-9]+\.[0-9]+`)

func Print() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println(constants.AppName, versionRegexp.FindString(variables.GitVersion))
		fmt.Println(" Git Commit:", variables.GitVersion)
		fmt.Println(" Build with:", variables.GoVersion)
		fmt.Println(" Build time:", variables.BuildTime)
		return nil
	}
}
