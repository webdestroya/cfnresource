package cfnresource

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
)

type logSetuper interface {
	SetLogWriter(io.Writer)
}

type logContextSetter interface {
	SetLogContext(context.Context, io.Writer) context.Context
}

type loadOptionsFunc = func(*config.LoadOptions) error

type AwsConfigOptioner interface {
	GetAwsConfigOptions(context.Context) []loadOptionsFunc
}

type defaultCallbackDelayGetter interface {
	DefaultCallbackDelay() time.Duration
}
