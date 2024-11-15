package cfnresource

import (
	"encoding/json"
	"errors"

	"github.com/webdestroya/cfnresource/cfnutils"
	"github.com/webdestroya/cfnresource/encoding"

	cfnTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type Request[Model any, Ctx any] struct {
	Action       string
	AWSAccountId string

	Region string

	StackId   string
	StackName string

	LogicalResourceID          string
	ResourceProperties         *Model
	PreviousResourceProperties *Model

	CallbackContext *Ctx

	StackTags  Tags
	SystemTags Tags

	NextToken string

	TypeConfiguration json.RawMessage

	bearerToken string
	event       *event
}

func (r *Request[Model, Ctx]) UnmarshalJSON(data []byte) error {
	return errors.New("dont marshal the request object directly")
}

func (r *Request[Model, Ctx]) InProgressResponse(model *Model, callbackContext *Ctx) *ProgressEvent[Model, Ctx] {
	return &ProgressEvent[Model, Ctx]{
		OperationStatus: cfnTypes.OperationStatusInProgress,
		CallbackContext: callbackContext,
		ResourceModel:   model,
	}
}

func (r *Request[Model, Ctx]) SuccessResponse(model *Model) *ProgressEvent[Model, Ctx] {
	return &ProgressEvent[Model, Ctx]{
		OperationStatus: cfnTypes.OperationStatusSuccess,
		ResourceModel:   model,
	}
}

func (r *Request[Model, Ctx]) ErrorResponse(err any) *ProgressEvent[Model, Ctx] {
	pe := &ProgressEvent[Model, Ctx]{
		OperationStatus:  cfnTypes.OperationStatusFailed,
		HandlerErrorCode: cfnTypes.HandlerErrorCodeInternalFailure,
		Message:          "",
	}

	switch v := err.(type) {
	case string:
		pe = pe.WithMessage(v)
	case error:
		pe = pe.WithMessage(v.Error())
	}

	return pe
}

func newRequest[Model any, CallbackCtx any](event *event) (*Request[Model, CallbackCtx], error) {
	req := &Request[Model, CallbackCtx]{
		StackId:           event.StackID,
		StackName:         cfnutils.GetStackNameFromArn(event.StackID),
		Action:            event.Action,
		AWSAccountId:      event.AWSAccountID,
		bearerToken:       event.BearerToken,
		Region:            event.Region,
		LogicalResourceID: event.RequestData.LogicalResourceID,
		StackTags:         event.RequestData.StackTags,
		SystemTags:        event.RequestData.SystemTags,
		NextToken:         event.NextToken,
		TypeConfiguration: event.RequestData.TypeConfiguration,
		event:             event,
	}

	if len(event.CallbackContext) > 0 {
		if err := encoding.Unmarshal(event.CallbackContext, req.CallbackContext); err != nil {
			return nil, err
		}
	}

	if len(event.RequestData.ResourceProperties) > 0 {
		if err := encoding.Unmarshal(event.RequestData.ResourceProperties, req.ResourceProperties); err != nil {
			return nil, err
		}
	}

	if len(event.RequestData.PreviousResourceProperties) > 0 {
		if err := encoding.Unmarshal(event.RequestData.PreviousResourceProperties, req.PreviousResourceProperties); err != nil {
			return nil, err
		}
	}

	return req, nil
}
