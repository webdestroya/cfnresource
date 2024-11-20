package cfnerr

import (
	"errors"

	cfnTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func As(err error) (Error, bool) {
	var ce *cfnErr
	if errors.As(err, &ce) {
		return ce, true
	}
	return nil, false
}

func Is(err error) bool {
	_, ok := err.(errChecker)
	return ok
}

func New(code cfnTypes.HandlerErrorCode, msg string, err error) Error {

	if code == "" {
		code = InternalFailure
	}

	return &cfnErr{
		code:    code,
		message: msg,
		err:     err,
	}
}

func NewMessage(code cfnTypes.HandlerErrorCode, msg string) Error {
	return New(code, msg, nil)
}

func Wrap(code cfnTypes.HandlerErrorCode, err error) Error {
	return New(code, err.Error(), err)
}
