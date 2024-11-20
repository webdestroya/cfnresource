package cfnerr

import (
	"errors"

	cfnTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type Error interface {
	error

	// Returns an error code
	Code() cfnTypes.HandlerErrorCode

	// Returns the error message
	Message() string

	// Returns the original error
	OrigErr() error
}

type errChecker interface {
	isCfnError() bool
}

type cfnErr struct {
	// Classification of error
	code cfnTypes.HandlerErrorCode

	// Detailed information about error
	message string

	// Optional original error this error is based off of. Allows building
	// chained errors.
	err error
}

var _ error = (*cfnErr)(nil)
var _ Error = (*cfnErr)(nil)

func (cfnErr) isCfnError() bool {
	return true
}

func (b cfnErr) Error() string {
	return string(b.code) + ":" + b.err.Error()
}

func (b cfnErr) String() string {
	return b.Error()
}

// Code returns the short phrase depicting the classification of the error.
func (b cfnErr) Code() cfnTypes.HandlerErrorCode {
	return b.code
}

// Message returns the error details message.
func (b cfnErr) Message() string {
	return b.message
}

// Message returns the error details message.
func (b cfnErr) OrigErr() error {
	return b.err
}

func (b cfnErr) Is(err error) bool {
	if err == nil {
		return false
	}

	if Is(err) {
		return true
	}

	return errors.Is(b.err, err)
}

func (b cfnErr) Unwrap() error {
	return b.err
}
