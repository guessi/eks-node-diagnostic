package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/guessi/eks-node-diagnostic/internal/constants"
	k8s "github.com/guessi/eks-node-diagnostic/internal/kubernetes"
	"github.com/guessi/eks-node-diagnostic/internal/s3utils"
	"github.com/guessi/eks-node-diagnostic/internal/types"
	"github.com/guessi/eks-node-diagnostic/internal/utils"

	"github.com/urfave/cli/v3"
	"sigs.k8s.io/yaml"
)

func Entry() *cli.Command {
	return &cli.Command{
		Name:  constants.AppName,
		Usage: constants.AppUsage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config-file",
				Aliases: []string{"c"},
				Value:   "config.yaml",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print version number",
				Action:  utils.Version(),
			},
		},
		Action: func(c context.Context, cmd *cli.Command) error {
			configFile := cmd.String("config-file")

			yamlCfg, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("failed to open %s", configFile)
			}

			cfg := types.AppConfigs{}
			if err = yaml.Unmarshal(yamlCfg, &cfg); err != nil {
				return fmt.Errorf("failed to load %s", configFile)
			}

			if err := utils.ValidateAppConfigs(cfg); err != nil {
				return err
			}

			k8sclient, err := k8s.NewKubeClient()
			if err != nil {
				return fmt.Errorf("failed to create kubernetes client: %w", err)
			}

			for _, nodeName := range cfg.Nodes {
				url, err := s3utils.PresignUrlPutObject(
					types.PresignUrlPutObjectInput{
						Region:         cfg.Region,
						BucketName:     cfg.BucketName,
						NodeName:       nodeName,
						ExpiredSeconds: cfg.ExpiredSeconds,
					},
				)
				if err != nil {
					return err
				}

				err = k8sclient.Apply(nodeName, url)
				if err != nil {
					fmt.Printf("%s\n", err)
				} else {
					fmt.Printf("nodediagnostic.eks.amazonaws.com/%s created\n", nodeName)
				}
			}
			return nil
		},
	}
}
