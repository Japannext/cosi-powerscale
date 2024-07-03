package provisioner

import (
	"errors"
)

var (
	ErrInvalidRequest                   = errors.New("incoming request invalid")
	ErrInvalidBucketID                  = errors.New("invalid bucketID")
	ErrEmptyBucketID                    = errors.New("empty bucket ID")
	ErrEmptyBucketAccessName            = errors.New("empty bucket access name")
	ErrInvalidAuthenticationType        = errors.New("invalid authentication type")
	ErrUnknownAuthenticationType        = errors.New("unknown authentication type")
	ErrBucketNotFound                   = errors.New("bucket not found")
	ErrFailedToCreateUser               = errors.New("failed to create user")
	ErrFailedToDecodePolicy             = errors.New("failed to decode bucket policy")
	ErrFailedToUpdatePolicy             = errors.New("failed to update bucket policy")
	ErrFailedToCreateAccessKey          = errors.New("failed to create access key")
	ErrFailedToGeneratePolicyID         = errors.New("failed to generate PolicyID UUID")
	ErrGeneratedPolicyIDIsEmpty         = errors.New("generated PolicyID was empty")
	ErrAuthenticationTypeNotImplemented = errors.New("authentication type IAM not implemented")
)
