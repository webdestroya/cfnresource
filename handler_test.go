package cfnresource

import (
	"context"
	"encoding/json"
	"io"
	"testing"

	cfnTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/cfnresource/cfncontext"
)

type fancyStr string

type model struct {
	Name       string   `json:",omitempty"`
	BoolVal    bool     `json:",omitempty"`
	IntVal     int      `json:",omitempty"`
	Floaty     float64  `json:",omitempty"`
	FancyPants fancyStr `json:",omitempty"`
}

type callbackCtx struct {
	Step *int `json:",omitempty"`
}

type ctxKeyType string

var dummyCtx = ctxKeyType("dummy")

type requestType = *Request[model, callbackCtx]
type progEventType = *ProgressEvent[model, callbackCtx]

type basicHandler struct {
	t *testing.T
}

func (basicHandler) SetLogContext(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, dummyCtx, "some test val")
}

func (h basicHandler) Create(ctx context.Context, req requestType) (progEventType, error) {

	require.Equal(h.t, "Test Thing", req.ResourceProperties.Name)
	require.Equal(h.t, fancyStr("yaryar"), req.ResourceProperties.FancyPants)

	val, ok := ctx.Value(dummyCtx).(string)
	require.True(h.t, ok)
	require.Equal(h.t, "some test val", val)

	cfg, err := cfncontext.GetAwsConfig(ctx)
	require.NoError(h.t, err)

	creds, err := cfg.Credentials.Retrieve(ctx)
	require.NoError(h.t, err)

	require.Equal(h.t, "fake", creds.AccessKeyID)
	require.Equal(h.t, "SampleStack", req.StackName)

	require.Equal(h.t, "logically", req.LogicalResourceID)

	// return req.ErrorResponse("oops"), nil
	return req.SuccessResponse(&model{
		Name:    req.ResourceProperties.Name,
		BoolVal: true,
	}), nil
}

func (basicHandler) Update(ctx context.Context, req requestType) (progEventType, error) {
	val := 1234
	return req.InProgressResponse(req.ResourceProperties, &callbackCtx{
		Step: &val,
	}), nil
}

func (basicHandler) Delete(ctx context.Context, req requestType) (progEventType, error) {

	_ = req.CallbackContext.Step

	return nil, nil
}

func (basicHandler) Read(ctx context.Context, req requestType) (progEventType, error) {
	return nil, nil
}

func (basicHandler) List(ctx context.Context, req requestType) (progEventType, error) {
	return nil, nil
}

func TestHandlerBasic(t *testing.T) {
	fn := makeEventFunc(basicHandler{t})

	creds := &credProvider{
		AccessKeyID:     "fake",
		SecretAccessKey: "fake",
		SessionToken:    "fake",
	}

	ev := &event{
		AWSAccountID:        "000000000000",
		BearerToken:         "xxxbearerxxx",
		Region:              "us-east-1",
		Action:              createAction,
		ResourceType:        "Dummy::Thing::Basic",
		ResourceTypeVersion: "1.0",
		CallbackContext:     nil,
		RequestData: requestData{
			CallerCredentials:          creds,
			LogicalResourceID:          "logically",
			ResourceProperties:         json.RawMessage(`{"Name": "Test Thing", "IntVal": 1234, "FancyPants": "yaryar"}`),
			PreviousResourceProperties: json.RawMessage(`{}`),
			ProviderCredentials:        creds,
			ProviderLogGroupName:       "loggroup",
			StackTags:                  map[string]string{},
			SystemTags:                 map[string]string{},
			TypeConfiguration:          nil,
		},
		StackID:   "arn:aws:cloudformation:us-east-1:123456789012:stack/SampleStack/e722ae60-fe62-11e8-9a0e-0ae8cc519968",
		NextToken: "",
	}

	resp, err := fn(context.Background(), ev)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// respData, err := json.Marshal(resp)
	// require.NoError(t, err)
	// t.Logf("STUFF: %v", string(respData))

	ev.Action = updateAction
	resp, err = fn(context.Background(), ev)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, cfnTypes.OperationStatusInProgress, resp.OperationStatus)

	ev.Action = updateAction
	ev.CallbackContext = json.RawMessage(`{"Step": "1234"}`)
	resp, err = fn(context.Background(), ev)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, cfnTypes.OperationStatusInProgress, resp.OperationStatus)

}

func TestPanicker(t *testing.T) {
	fn := makeEventFunc(basicHandler{t})

	creds := &credProvider{
		AccessKeyID:     "fake",
		SecretAccessKey: "fake",
		SessionToken:    "fake",
	}

	ev := &event{
		AWSAccountID:        "000000000000",
		BearerToken:         "xxxbearerxxx",
		Region:              "us-east-1",
		Action:              deleteAction,
		ResourceType:        "Dummy::Thing::Basic",
		ResourceTypeVersion: "1.0",
		CallbackContext:     nil,
		RequestData: requestData{
			CallerCredentials:          creds,
			LogicalResourceID:          "logically",
			ResourceProperties:         json.RawMessage(`{"Name": "Test Thing", "IntVal": 1234}`),
			PreviousResourceProperties: json.RawMessage(`{}`),
			ProviderCredentials:        creds,
			ProviderLogGroupName:       "loggroup",
			StackTags:                  map[string]string{},
			SystemTags:                 map[string]string{},
			TypeConfiguration:          nil,
		},
		StackID:   "arn:aws:cloudformation:us-east-1:123456789012:stack/SampleStack/e722ae60-fe62-11e8-9a0e-0ae8cc519968",
		NextToken: "",
	}

	resp, err := fn(context.Background(), ev)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, cfnTypes.OperationStatusFailed, resp.OperationStatus)
	require.EqualValues(t, cfnTypes.HandlerErrorCodeInternalFailure, resp.ErrorCode)
	t.Logf("msg=%s", resp.Message)
}
