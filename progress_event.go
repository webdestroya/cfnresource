package cfnresource

import (
	"time"

	cfnTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type ProgressEvent[Model any, CallbackCtx any] struct {
	// OperationStatus indicates whether the handler has reached a terminal state or is
	// still computing and requires more time to complete.
	OperationStatus cfnTypes.OperationStatus `json:"status,omitempty"`

	// HandlerErrorCode should be provided when OperationStatus is FAILED or IN_PROGRESS.
	HandlerErrorCode cfnTypes.HandlerErrorCode `json:"errorCode,omitempty"`

	// Message which can be shown to callers to indicate the
	// nature of a progress transition or callback delay; for example a message
	// indicating "propagating to edge."
	Message string `json:"message,omitempty"`

	// CallbackContext is an arbitrary datum which the handler can return in an
	// IN_PROGRESS event to allow the passing through of additional state or
	// metadata between subsequent retries; for example to pass through a Resource
	// identifier which can be used to continue polling for stabilization
	CallbackContext *CallbackCtx `json:"callbackContext,omitempty"`

	// CallbackDelaySeconds will be scheduled with an initial delay of no less than the number
	// of seconds specified in the progress event. Set this value to <= 0 to
	// indicate no callback should be made.
	CallbackDelaySeconds int64 `json:"callbackDelaySeconds,omitempty"`

	// ResourceModel is the output resource instance populated by a READ/LIST for synchronous results
	// and by CREATE/UPDATE/DELETE for final response validation/confirmation
	ResourceModel *Model `json:"resourceModel,omitempty"`

	// ResourceModels is the output resource instances populated by a LIST for
	// synchronous results. ResourceModels must be returned by LIST so it's
	// always included in the response. When ResourceModels is not set, null is
	// returned.
	ResourceModels []*Model `json:"resourceModels"`

	// NextToken is the token used to request additional pages of resources for a LIST operation
	NextToken string `json:"nextToken,omitempty"`
}

func (pe *ProgressEvent[Model, CallbackCtx]) WithMessage(v string) *ProgressEvent[Model, CallbackCtx] {
	pe.Message = v
	return pe
}

func (pe *ProgressEvent[Model, CallbackCtx]) WithCallbackDelay(v time.Duration) *ProgressEvent[Model, CallbackCtx] {
	pe.CallbackDelaySeconds = int64(v.Seconds())
	return pe
}

func (pe *ProgressEvent[Model, CallbackCtx]) WithModels(models ...*Model) *ProgressEvent[Model, CallbackCtx] {
	pe.ResourceModels = models
	return pe
}

func (pe *ProgressEvent[Model, CallbackCtx]) WithModel(model *Model) *ProgressEvent[Model, CallbackCtx] {
	pe.ResourceModel = model
	return pe
}

func (pe *ProgressEvent[Model, CallbackCtx]) WithErrorCode(code cfnTypes.HandlerErrorCode) *ProgressEvent[Model, CallbackCtx] {
	pe.HandlerErrorCode = code
	return pe
}

func (pe *ProgressEvent[Model, CallbackCtx]) toResponse(bearerToken string) (response, error) {
	return newResponse(pe, bearerToken)
}
