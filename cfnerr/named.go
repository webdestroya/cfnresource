package cfnerr

import (
	cfnTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

const (
	AccessDenied            = cfnTypes.HandlerErrorCodeAccessDenied
	AlreadyExists           = cfnTypes.HandlerErrorCodeAlreadyExists
	GeneralServiceException = cfnTypes.HandlerErrorCodeGeneralServiceException
	InternalFailure         = cfnTypes.HandlerErrorCodeInternalFailure
	InvalidCredentials      = cfnTypes.HandlerErrorCodeInvalidCredentials
	InvalidRequest          = cfnTypes.HandlerErrorCodeInvalidRequest
	NetworkFailure          = cfnTypes.HandlerErrorCodeNetworkFailure
	NotFound                = cfnTypes.HandlerErrorCodeNotFound
	NotStabilized           = cfnTypes.HandlerErrorCodeServiceTimeout
	NotUpdatable            = cfnTypes.HandlerErrorCodeNotUpdatable
	ResourceConflict        = cfnTypes.HandlerErrorCodeResourceConflict
	ServiceInternalError    = cfnTypes.HandlerErrorCodeServiceInternalError
	ServiceLimitExceeded    = cfnTypes.HandlerErrorCodeServiceLimitExceeded
	Throttling              = cfnTypes.HandlerErrorCodeThrottling
)
