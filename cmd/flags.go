package cmd

import (
	"github.com/guessi/eks-node-diagnostic/internal/constants"
	"github.com/guessi/eks-node-diagnostic/internal/node_diagnostic"
	"github.com/guessi/eks-node-diagnostic/internal/s3utils"
	"github.com/guessi/eks-node-diagnostic/internal/types"
	"github.com/guessi/eks-node-diagnostic/internal/utils"
	"github.com/guessi/eks-node-diagnostic/internal/validators"

	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "region",
		Aliases: []string{"r"},
		Value:   "us-east-1",
		Usage:   "region of the target cluster",
	},
	&cli.StringFlag{
		Name:    "node-name",
		Aliases: []string{"n"},
		Usage:   "target node name",
	},
	&cli.StringFlag{
		Name:    "bucket-name",
		Aliases: []string{"b"},
		Usage:   "bucket name for the log files",
	},
	&cli.IntFlag{
		Name:    "expire-seconds",
		Aliases: []string{"t"},
		Value:   300,
		Usage:   "expiration time of the presigned-url in seconds",
	},
}

var Commands = []*cli.Command{
	{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print version number",
		Action:  utils.Version(),
	},
}

func validateConfig(config types.AppConfig) error {
	if err := validators.ValidateEmpty("region", config.Region); err != nil {
		return err
	}
	if err := validators.ValidateNodeName(config.NodeName); err != nil {
		return err
	}
	if err := validators.ValidateEmpty("bucket-name", config.BucketName); err != nil {
		return err
	}

	if err := validators.ValidateInRange(config.ExpireSeconds, constants.MinExpireSeconds, constants.MaxExpireSeconds); err != nil {
		return err
	}
	return nil
}

func Action() cli.ActionFunc {
	return func(c *cli.Context) error {
		appConfig := types.AppConfig{
			Region:        c.String("region"),
			BucketName:    c.String("bucket-name"),
			NodeName:      c.String("node-name"),
			ExpireSeconds: c.Int("expire-seconds"),
		}

		if err := validateConfig(appConfig); err != nil {
			return err
		}

		presignPutObjectUrl, err := s3utils.PresignUrlPutObject(appConfig)
		if err != nil {
			return err
		}

		return node_diagnostic.Render(appConfig.NodeName, presignPutObjectUrl)
	}
}

func Entry() *cli.App {
	return &cli.App{
		Name:     constants.AppName,
		Usage:    constants.AppUsage,
		Flags:    Flags,
		Commands: Commands,
		Action:   Action(),
	}
}
