package cfnutils

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

const (
	stackArnPrefix = `stack/`
)

func GetStackNameFromArn(value string) string {
	if !arn.IsARN(value) {
		return ""
	}

	v, _ := arn.Parse(value)

	// "arn:aws:cloudformation:us-east-1:xxxxxxxx:stack/temporal-test/SOME_UUID",
	if !strings.HasPrefix(v.Resource, stackArnPrefix) {
		return ""
	}

	stackName, _, _ := strings.Cut(strings.TrimPrefix(v.Resource, stackArnPrefix), "/")

	return stackName

}
