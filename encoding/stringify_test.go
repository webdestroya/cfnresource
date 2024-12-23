package encoding_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/cfnresource/encoding"
)

func TestStringifyTypes(t *testing.T) {
	type Struct struct {
		S string
	}

	s := "foo"
	b := true
	i := 42
	f := 3.14
	l := []interface{}{s, b, i, f}
	m := map[string]interface{}{
		"l": l,
	}
	o := Struct{S: s}
	var nilPointer *Struct

	for _, testCase := range []struct {
		data     interface{}
		expected interface{}
	}{
		// Basic types
		{s, "foo"},
		{b, "true"},
		{i, "42"},
		{f, "3.14"},
		{l, []interface{}{"foo", "true", "42", "3.14"}},
		{m, map[string]interface{}{"l": []interface{}{"foo", "true", "42", "3.14"}}},
		{o, struct{ S string }{S: "foo"}},

		// Pointers
		{&s, "foo"},
		{&b, "true"},
		{&i, "42"},
		{&f, "3.14"},
		{&l, []interface{}{"foo", "true", "42", "3.14"}},
		{&m, map[string]interface{}{"l": []interface{}{"foo", "true", "42", "3.14"}}},
		{&o, struct{ S string }{S: "foo"}},

		// Nils are stripped
		{map[string]interface{}{"foo": nil}, map[string]interface{}{}},

		// Nil pointers are nil
		{nilPointer, nil},

		// Nils are nil
		{nil, nil},
	} {
		actual, err := encoding.Stringify(testCase.data)
		require.NoError(t, err)

		require.Empty(t, cmp.Diff(actual, testCase.expected))
	}
}

func TestStringifyModel(t *testing.T) {
	type Model struct {
		BucketName      *string
		Key             *string
		Body            *string
		IsBase64Encoded *bool
		ContentType     *string
		ContentLength   *int
		ACL             *string
		Grants          map[string][]string
	}

	m := Model{
		BucketName:  aws.String("foo"),
		Key:         aws.String("bar"),
		Body:        aws.String("baz"),
		ContentType: aws.String("quux"),
		ACL:         aws.String("mooz"),
	}

	expected := struct {
		BucketName      string
		Key             string
		Body            string
		IsBase64Encoded string
		ContentType     string
		ContentLength   string
		ACL             string
		Grants          map[string]interface{}
	}{
		BucketName:  "foo",
		Key:         "bar",
		Body:        "baz",
		ContentType: "quux",
		ACL:         "mooz",
		Grants:      map[string]interface{}{},
	}

	actual, err := encoding.Stringify(m)
	require.NoError(t, err)

	require.Empty(t, cmp.Diff(actual, expected))
}
