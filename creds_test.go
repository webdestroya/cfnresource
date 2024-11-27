package cfnresource

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredentialRedaction(t *testing.T) {
	creds := &credProvider{
		AccessKeyID:     "fake",
		SecretAccessKey: "fake",
		SessionToken:    "fake",
	}

	out, err := json.Marshal(creds)
	require.NoError(t, err)
	require.Equal(t, `"REDACTED"`, string(out))

	rd := requestData{
		CallerCredentials:   creds,
		ProviderCredentials: creds,
		LogicalResourceID:   "dummy",
	}

	out, err = json.Marshal(rd)
	require.NoError(t, err)

	require.Contains(t, string(out), `"callerCredentials":"REDACTED"`)
	require.Contains(t, string(out), `"providerCredentials":"REDACTED"`)

	validJSON := []byte(`{"accessKeyId": "fake", "secretAccessKey": "fake", "sessionToken": "fake"}`)
	newCreds := new(credProvider)
	err = json.Unmarshal(validJSON, newCreds)
	require.NoError(t, err)
	require.Equal(t, "fake", newCreds.AccessKeyID)
	require.Equal(t, "fake", newCreds.SecretAccessKey)
	require.Equal(t, "fake", newCreds.SessionToken)
}
