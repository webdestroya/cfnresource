package cfnresource

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	cfnTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/webdestroya/cfnresource/cfncontext"
	"github.com/webdestroya/cfnresource/cfnerr"
	"github.com/webdestroya/cfnresource/cfnutils"
	"github.com/webdestroya/cfnresource/cloudwatchwriter"
)

var (
	logSetup  sync.Once
	logWriter io.Writer
)

const (
	stackNameSystemTag = `aws:cloudformation:stack-name`
)

const (
	unknownAction = "UNKNOWN"
	createAction  = "CREATE"
	readAction    = "READ"
	updateAction  = "UPDATE"
	deleteAction  = "DELETE"
	listAction    = "LIST"
)

const (
	invalidRequestError  = "InvalidRequest"
	serviceInternalError = "ServiceInternal"
	unmarshalingError    = "UnmarshalingError"
	marshalingError      = "MarshalingError"
	validationError      = "Validation"
	timeoutError         = "Timeout"
	sessionNotFoundError = "SessionNotFound"
)

func Start[Model any, Ctx any](handler Handler[Model, Ctx]) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Handler panicked: %s", r)
			panic(r) // Continue the panic
		}
	}()

	log.Printf("Handler starting")
	lambda.Start(makeEventFunc(handler))
}

func makeEventFunc[Model any, Ctx any](handler Handler[Model, Ctx]) func(context.Context, *event) (response, error) {
	return func(ctx context.Context, event *event) (response, error) {

		providerCfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(event.RequestData.ProviderCredentials))
		if err != nil {
			return newFailedResponse(err, event.BearerToken)
		}

		// setup the caller aws config
		callerCfgOptions := []loadOptionsFunc{config.WithCredentialsProvider(event.RequestData.CallerCredentials)}
		if haws, ok := handler.(AwsConfigOptioner); ok {
			callerCfgOptions = append(callerCfgOptions, haws.GetAwsConfigOptions(ctx)...)
		}
		callerCfg, err := config.LoadDefaultConfig(ctx, callerCfgOptions...)
		if err != nil {
			return newFailedResponse(err, event.BearerToken)
		}
		ctx = cfncontext.SetAwsConfig(ctx, callerCfg)

		logicalId := event.RequestData.LogicalResourceID

		// logging setup
		// logSetup.Do(func() {
		logStreamName := fmt.Sprintf("%s/%s", cfnutils.GetStackNameFromArn(event.StackID), logicalId)
		cwClient := cloudwatchlogs.NewFromConfig(providerCfg)
		logWriter = cloudwatchwriter.NewSync(ctx, cwClient, event.RequestData.ProviderLogGroupName, logStreamName)

		log.SetOutput(logWriter)
		log.SetFlags(0)
		log.SetPrefix("")

		// })
		if hlog, ok := handler.(PostInitializer); ok {
			ctx, err = hlog.PostInitialize(ctx, logWriter)
			if err != nil {
				return newFailedResponse(err, event.BearerToken)
			}
		}

		if hlog, ok := handler.(eventLogger); ok {
			hlog.LogEvent(ctx, event)
		}

		handlerFn, err := router(event.Action, handler)
		if err != nil {
			return newFailedResponse(err, event.BearerToken)
		}

		req, err := newRequest[Model, Ctx](event)
		if err != nil {
			return newFailedResponse(err, event.BearerToken)
		}

		pe := invoke(handlerFn, ctx, req)
		resp, err := newResponse(pe, event.BearerToken)
		if err != nil {
			return newFailedResponse(err, event.BearerToken)
		}
		return resp, nil
	}
}

func router[Model any, Ctx any](action string, handler Handler[Model, Ctx]) (func(context.Context, *Request[Model, Ctx]) (*ProgressEvent[Model, Ctx], error), error) {
	switch action {
	case createAction:
		return handler.Create, nil
	case readAction:
		return handler.Read, nil
	case updateAction:
		return handler.Update, nil
	case deleteAction:
		return handler.Delete, nil
	case listAction:
		return handler.List, nil
	default:
		return nil, cfnerr.New(cfnerr.InvalidRequest, "No action/invalid action specified", nil)
	}
}

func invoke[Model any, Ctx any](handlerFn func(context.Context, *Request[Model, Ctx]) (*ProgressEvent[Model, Ctx], error), ctx context.Context, request *Request[Model, Ctx]) *ProgressEvent[Model, Ctx] {

	ch := make(chan *ProgressEvent[Model, Ctx], 1)

	go func() {
		ch <- invokeWrap(handlerFn, ctx, request)
	}()

	return <-ch
}

func invokeWrap[Model any, Ctx any](handlerFn func(context.Context, *Request[Model, Ctx]) (*ProgressEvent[Model, Ctx], error), ctx context.Context, request *Request[Model, Ctx]) (respPE *ProgressEvent[Model, Ctx]) {
	defer func() {
		// Catch any panics and return a failed ProgressEvent
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = errors.New(fmt.Sprint(r))
			}

			log.Printf("Trapped error in handler: %v", err)

			respPE = request.ErrorResponse(err).WithErrorCode(cfnTypes.HandlerErrorCodeInternalFailure)
		}
	}()

	pe, err := handlerFn(ctx, request)
	if err != nil {
		pe = request.ErrorResponse(err)
	}
	return pe
}
