package encoding_test

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/cfnresource/encoding"
)

func TestMarshaling(t *testing.T) {
	type Nested struct {
		SP *string  `json:",omitempty"`
		BP *bool    `json:",omitempty"`
		IP *int     `json:"intField,omitempty"`
		FP *float64 `json:"floatPointer,omitempty"`

		S string `json:"stringValue,omitempty"`
		B bool   `json:",omitempty"`
		I int
		F float64 `json:",omitempty"`
	}

	type Main struct {
		SP *string
		BP *bool    `json:",omitempty"`
		IP *int     `json:",omitempty"`
		FP *float64 `json:",omitempty"`
		NP *Nested  `json:"nestedPointer,omitempty"`

		S string `json:",omitempty"`
		B bool   `json:"boolValue,omitempty"`
		I int    `json:",omitempty"`
		F float64
		N Nested `json:",omitempty"`
	}

	m := Main{
		SP: aws.String("foo"),
		IP: aws.Int(42),
		NP: &Nested{
			BP: aws.Bool(true),
			FP: aws.Float64(3.14),
		},

		B: true,
		F: 2.72,
		N: Nested{
			S: "bar",
			I: 54,
		},
	}

	stringMap := map[string]interface{}{
		"SP": "foo",
		"IP": "42",
		"nestedPointer": map[string]interface{}{
			"BP":           "true",
			"I":            "0",
			"floatPointer": "3.14",
		},

		"boolValue": "true",
		"F":         "2.72",
		"N": map[string]interface{}{
			"stringValue": "bar",
			"I":           "54",
		},
	}

	var err error

	rep, err := encoding.Marshal(m)
	require.NoError(t, err)

	// Test that rep can be unmarshalled as regular JSON
	var jsonTest map[string]interface{}
	err = json.Unmarshal(rep, &jsonTest)
	require.NoError(t, err)

	// And check it matches the expected form
	require.Empty(t, cmp.Diff(jsonTest, stringMap))

	// Now check we can get the original struct back
	var b Main
	err = encoding.Unmarshal(rep, &b)
	require.NoError(t, err)

	require.Empty(t, cmp.Diff(m, b))
}
