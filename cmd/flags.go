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

			s3client, err := s3utils.NewS3Client(cfg.Region)
			if err != nil {
				return fmt.Errorf("failed to create S3 client: %w", err)
			}

			var hasErrors bool
			for _, nodeName := range cfg.Nodes {
				url, err := s3utils.PresignUrlPutObject(
					s3client,
					types.PresignUrlPutObjectInput{
						Region:         cfg.Region,
						BucketName:     cfg.BucketName,
						NodeName:       nodeName,
						ExpiredSeconds: cfg.ExpiredSeconds,
					},
				)
				if err != nil {
					fmt.Printf("failed to generate presigned URL for %s: %s\n", nodeName, err)
					hasErrors = true
					continue
				}

				err = k8sclient.Apply(nodeName, url)
				if err != nil {
					fmt.Printf("failed to apply nodediagnostic for %s: %s\n", nodeName, err)
					hasErrors = true
				} else {
					fmt.Printf("nodediagnostic.eks.amazonaws.com/%s created\n", nodeName)
				}
			}

			if hasErrors {
				return fmt.Errorf("one or more nodes failed to process")
			}
			return nil
		},
	}
}
