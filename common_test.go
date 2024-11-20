package cfnresource_test

import (
	"github.com/webdestroya/cfnresource"
)

type model struct {
	Name    string  `json:",omitempty"`
	BoolVal bool    `json:",omitempty"`
	IntVal  int     `json:",omitempty"`
	Floaty  float64 `json:",omitempty"`
}

type callbackCtx struct {
}

type requestType = *cfnresource.Request[model, callbackCtx]
type progEventType = *cfnresource.ProgressEvent[model, callbackCtx]
