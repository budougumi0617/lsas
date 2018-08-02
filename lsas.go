package lsas

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

// Tag is AMI tag name and value.
type Tag struct {
	Key, Value string
}

// LoadConfig loads AWS setting with option.
func LoadConfig(region string) (aws.Config, error) {
	// FIXME Need to set flexible amount of option
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return aws.Config{}, err
	}
	if len(region) != 0 {
		cfg.Region = region
	}
	return cfg, nil
}
