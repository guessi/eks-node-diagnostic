package types

type AppConfigs struct {
	Region        string   `yaml:"region,omitempty"`
	BucketName    string   `yaml:"bucketName,omitempty"`
	NodeNames     []string `yaml:"nodes,omitempty"`
	ExpireSeconds int      `yaml:"expiredSeconds,omitempty"`
}

type PresignUrlPutObjectInput struct {
	Region        string
	BucketName    string
	NodeName      string
	ExpireSeconds int
}
