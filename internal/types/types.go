package types

type AppConfigs struct {
	BucketRegion    string   `yaml:"region,omitempty"`
	DestinationType string   `yaml:"destinationType,omitempty"`
	BucketName      string   `yaml:"bucketName,omitempty"`
	Nodes           []string `yaml:"nodes,omitempty"`
	ExpiredSeconds  int      `yaml:"expiredSeconds,omitempty"`
	Timeout         int      `yaml:"timeout,omitempty"` // in seconds
}

type PresignUrlPutObjectInput struct {
	BucketRegion   string
	BucketName     string
	NodeName       string
	ExpiredSeconds int
}
