package validate

import (
	"testing"

	"github.com/guessi/eks-node-diagnostic/internal/types"
)

func TestAppConfigs(t *testing.T) {
	base := types.AppConfigs{
		Region:          "us-west-2",
		DestinationType: "s3",
		BucketName:      "my-bucket",
		Nodes:           []string{"i-0123456789abcdef0"},
	}

	tests := []struct {
		name    string
		modify  func(types.AppConfigs) types.AppConfigs
		wantErr bool
	}{
		// S3-specific validations
		{
			"s3 missing bucket",
			func(c types.AppConfigs) types.AppConfigs { c.BucketName = ""; return c },
			true,
		},
		{
			"s3 expire out of range",
			func(c types.AppConfigs) types.AppConfigs { c.ExpiredSeconds = 1; return c },
			true,
		},
		{
			"s3 expire at min boundary",
			func(c types.AppConfigs) types.AppConfigs { c.ExpiredSeconds = 120; return c },
			false,
		},
		{
			"s3 expire at max boundary",
			func(c types.AppConfigs) types.AppConfigs { c.ExpiredSeconds = 86400; return c },
			false,
		},
		{
			"s3 expire above max",
			func(c types.AppConfigs) types.AppConfigs { c.ExpiredSeconds = 86401; return c },
			true,
		},

		// node destination type skips S3 validations
		{
			"node skips bucket check",
			func(c types.AppConfigs) types.AppConfigs {
				c.DestinationType = "node"
				c.BucketName = ""
				return c
			},
			false,
		},
		{
			"node skips expire check",
			func(c types.AppConfigs) types.AppConfigs {
				c.DestinationType = "node"
				c.ExpiredSeconds = 1
				return c
			},
			false,
		},

		// timeout applies to both destination types
		{
			"timeout above max",
			func(c types.AppConfigs) types.AppConfigs { c.Timeout = 301; return c },
			true,
		},
		{
			"node still validates timeout",
			func(c types.AppConfigs) types.AppConfigs {
				c.DestinationType = "node"
				c.Timeout = 301
				return c
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.modify(base)
			err := AppConfigs(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppConfigs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
