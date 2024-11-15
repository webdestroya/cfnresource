package cfncontext

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type ctxKey string

var ErrContextValueMissingError = errors.New("Config missing")

const (
	awsCfgKey = ctxKey(`awscfg`)
)

func SetAwsConfig(ctx context.Context, cfg aws.Config) context.Context {
	return context.WithValue(ctx, awsCfgKey, cfg)
}

func GetAwsConfig(ctx context.Context) (aws.Config, error) {
	val, ok := ctx.Value(awsCfgKey).(aws.Config)
	if !ok {
		return aws.Config{}, ErrContextValueMissingError
	}
	return val, nil
}
