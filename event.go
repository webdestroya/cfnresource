package cfnresource

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// Tags are stored as key/value paired strings
type Tags map[string]string

type event struct {
	AWSAccountID        string          `json:"awsAccountId"`
	BearerToken         string          `json:"bearerToken" validate:"nonzero"`
	Region              string          `json:"region" validate:"nonzero"`
	Action              string          `json:"action"`
	ResourceType        string          `json:"resourceType"`
	ResourceTypeVersion string          `json:"resourceTypeVersion"` // encoding.Float
	CallbackContext     json.RawMessage `json:"callbackContext,omitempty"`
	RequestData         requestData     `json:"requestData"`
	StackID             string          `json:"stackId"`
	NextToken           string
}

// requestData is internal to the RPDK. It contains a number of fields that are for
// internal use only.
type requestData struct {
	// This would be xxxxx-us-east-1-cfn-role (the role assigned to the Resource)
	CallerCredentials *credProvider `json:"callerCredentials"`

	LogicalResourceID          string          `json:"logicalResourceId"`
	ResourceProperties         json.RawMessage `json:"resourceProperties"`
	PreviousResourceProperties json.RawMessage `json:"previousResourceProperties"`

	// this would be the tf-cloudformation role
	ProviderCredentials *credProvider `json:"providerCredentials"`

	ProviderLogGroupName string          `json:"providerLogGroupName"`
	StackTags            Tags            `json:"stackTags"`
	SystemTags           Tags            `json:"systemTags"`
	TypeConfiguration    json.RawMessage `json:"typeConfiguration"`
}

type credProvider struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
}

var _ aws.CredentialsProvider = (*credProvider)(nil)

func (c *credProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return credentials.NewStaticCredentialsProvider(c.AccessKeyID, c.SecretAccessKey, c.SessionToken).Retrieve(ctx)
}

func (credProvider) MarshalJSON() ([]byte, error) {
	return []byte(`"REDACTED"`), nil
}
