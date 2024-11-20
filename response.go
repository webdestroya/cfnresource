package cfnresource

import (
	cfnTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/webdestroya/cfnresource/cfnerr"
	"github.com/webdestroya/cfnresource/encoding"
)

// response represents a response to the
// cloudformation service from a resource handler.
// The zero value is ready to use.
type response struct {
	// Message which can be shown to callers to indicate the nature of a
	// progress transition or callback delay; for example a message
	// indicating "propagating to edge"
	Message string `json:"message,omitempty"`

	// The operationStatus indicates whether the handler has reached a terminal
	// state or is still computing and requires more time to complete
	OperationStatus cfnTypes.OperationStatus `json:"status,omitempty"`

	// ResourceModel it The output resource instance populated by a READ/LIST for
	// synchronous results and by CREATE/UPDATE/DELETE for final response
	// validation/confirmation
	ResourceModel any `json:"resourceModel,omitempty"`

	// ErrorCode is used to report granular failures back to CloudFormation
	ErrorCode string `json:"errorCode,omitempty"`

	// BearerToken is used to report progress back to CloudFormation and is
	// passed back to CloudFormation
	BearerToken string `json:"bearerToken,omitempty"`

	// ResourceModels is the output resource instances populated by a LIST for
	// synchronous results. ResourceModels must be returned by LIST so it's
	// always included in the response. When ResourceModels is not set, null is
	// returned.
	ResourceModels []any `json:"resourceModels"`

	// NextToken the token used to request additional pages of resources for a LIST operation
	NextToken string `json:"nextToken,omitempty"`

	// CallbackContext is an arbitrary datum which the handler can return in an
	// IN_PROGRESS event to allow the passing through of additional state or
	// metadata between subsequent retries; for example to pass through a Resource
	// identifier which can be used to continue polling for stabilization
	CallbackContext any `json:"callbackContext,omitempty"`

	// CallbackDelaySeconds will be scheduled with an initial delay of no less than the number
	// of seconds specified in the progress event. Set this value to <= 0 to
	// indicate no callback should be made.
	CallbackDelaySeconds int `json:"callbackDelaySeconds,omitempty"`
}

// newFailedResponse returns a response pre-filled with the supplied error
func newFailedResponse(err error, bearerToken string) (response, error) {

	if ce, ok := cfnerr.As(err); ok {
		return response{
			OperationStatus: cfnTypes.OperationStatusFailed,
			ErrorCode:       string(ce.Code()),
			Message:         ce.Message(),
			BearerToken:     bearerToken,
		}, nil
	}

	return response{
		OperationStatus: cfnTypes.OperationStatusFailed,
		ErrorCode:       string(cfnTypes.HandlerErrorCodeInternalFailure),
		Message:         err.Error(),
		BearerToken:     bearerToken,
	}, err
}

// newResponse converts a progress event into a useable reponse
// for the CloudFormation Resource Provider service to understand.
func newResponse[Model any, Ctx any](pevt *ProgressEvent[Model, Ctx], bearerToken string) (response, error) {
	model, err := encoding.Stringify(pevt.ResourceModel)
	if err != nil {
		return response{}, err
	}

	var models []any
	if pevt.ResourceModels != nil {
		models = make([]any, len(pevt.ResourceModels))
		for i := range pevt.ResourceModels {
			m, err := encoding.Stringify(pevt.ResourceModels[i])
			if err != nil {
				return response{}, err
			}

			models[i] = m
		}
	}

	cbCtx, err := encoding.Stringify(pevt.CallbackContext)
	if err != nil {
		return response{}, err
	}

	resp := response{
		BearerToken:          bearerToken,
		Message:              pevt.Message,
		OperationStatus:      pevt.OperationStatus,
		ResourceModel:        model,
		ResourceModels:       models,
		NextToken:            pevt.NextToken,
		CallbackContext:      cbCtx,
		CallbackDelaySeconds: pevt.CallbackDelaySeconds,
	}

	if pevt.HandlerErrorCode != "" {
		resp.ErrorCode = string(pevt.HandlerErrorCode)
	}

	return resp, nil
}
