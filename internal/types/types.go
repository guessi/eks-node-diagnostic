package types

type AppConfigs struct {
	Region         string   `yaml:"region,omitempty"`
	BucketName     string   `yaml:"bucketName,omitempty"`
	Nodes          []string `yaml:"nodes,omitempty"`
	ExpiredSeconds int      `yaml:"expiredSeconds,omitempty"`
}

type PresignUrlPutObjectInput struct {
	Region         string
	BucketName     string
	NodeName       string
	ExpiredSeconds int
}
