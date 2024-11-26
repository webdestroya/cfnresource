package cfnresource

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
)

type loadOptionsFunc = func(*config.LoadOptions) error

type AwsConfigOptioner interface {
	GetAwsConfigOptions(context.Context) []loadOptionsFunc
}

type PostInitializer interface {
	PostInitialize(context.Context, io.Writer) (context.Context, error)
}

type eventLogger interface {
	LogEvent(context.Context, any)
}

type defaultCallbackDelayGetter interface {
	DefaultCallbackDelay() time.Duration
}
