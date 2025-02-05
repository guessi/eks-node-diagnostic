package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
	"github.com/guessi/eks-node-diagnostic/internal/node_diagnostic"
	"github.com/guessi/eks-node-diagnostic/internal/s3utils"
	"github.com/guessi/eks-node-diagnostic/internal/types"
	"github.com/guessi/eks-node-diagnostic/internal/utils"
	"github.com/guessi/eks-node-diagnostic/internal/validators"

	"github.com/urfave/cli/v2"
	"sigs.k8s.io/yaml"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "config-file",
		Aliases: []string{"c"},
	},
	&cli.StringFlag{
		Name:    "region",
		Aliases: []string{"r"},
		Value:   "us-east-1",
		Usage:   "region of the target cluster",
	},
	&cli.StringFlag{
		Name:    "bucket-name",
		Aliases: []string{"b"},
		Usage:   "bucket name for the log files",
	},
	&cli.IntFlag{
		Name:    "expired-seconds",
		Aliases: []string{"t"},
		Value:   300,
		Usage:   "expiration time of the presigned-url in seconds",
	},
	&cli.StringFlag{
		Name:    "node-name",
		Aliases: []string{"n"},
		Usage:   "target node name (to support multiple node names, use \"--config-file\")",
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

func validateAppConfigs(config types.AppConfigs) error {
	if err := validators.ValidateEmpty("region", config.Region); err != nil {
		return err
	}

	for _, nodeName := range config.Nodes {
		if err := validators.ValidateNodeName(nodeName); err != nil {
			return err
		}
	}

	if err := validators.ValidateEmpty("bucket-name", config.BucketName); err != nil {
		return err
	}

	if err := validators.ValidateInRange(config.ExpiredSeconds, constants.MinExpireSeconds, constants.MaxExpireSeconds); err != nil {
		return err
	}
	return nil
}

func Action() cli.ActionFunc {
	return func(c *cli.Context) error {
		configFile := c.String("config-file")
		cfg := types.AppConfigs{}

		if strings.TrimSpace(configFile) != "" {
			yamlCfg, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("failed to open %s", configFile)
			}

			if err = yaml.Unmarshal(yamlCfg, &cfg); err != nil {
				return fmt.Errorf("failed to load %s", configFile)
			}
		} else {
			cfg = types.AppConfigs{
				Region:         c.String("region"),
				BucketName:     c.String("bucket-name"),
				ExpiredSeconds: c.Int("expired-seconds"),
				Nodes: []string{
					c.String("node-name"),
				},
			}
		}

		if err := validateAppConfigs(cfg); err != nil {
			return err
		}

		for _, node := range cfg.Nodes {
			presignUrlPutObjectInput := types.PresignUrlPutObjectInput{
				Region:         cfg.Region,
				BucketName:     cfg.BucketName,
				NodeName:       node,
				ExpiredSeconds: cfg.ExpiredSeconds,
			}
			url, err := s3utils.PresignUrlPutObject(presignUrlPutObjectInput)
			if err != nil {
				return err
			}

			renderErr := node_diagnostic.Render(node, url)
			if renderErr != nil {
				return renderErr
			}
		}
		return nil
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
