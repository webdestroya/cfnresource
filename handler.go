package cfnresource

import (
	"context"
)

type Handler[Model any, CallbackCtx any] interface {
	Create(context.Context, *Request[Model, CallbackCtx]) (*ProgressEvent[Model, CallbackCtx], error)
	Update(context.Context, *Request[Model, CallbackCtx]) (*ProgressEvent[Model, CallbackCtx], error)
	Delete(context.Context, *Request[Model, CallbackCtx]) (*ProgressEvent[Model, CallbackCtx], error)
	Read(context.Context, *Request[Model, CallbackCtx]) (*ProgressEvent[Model, CallbackCtx], error)
	List(context.Context, *Request[Model, CallbackCtx]) (*ProgressEvent[Model, CallbackCtx], error)
}
