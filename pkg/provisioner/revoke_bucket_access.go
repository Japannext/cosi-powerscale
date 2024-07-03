package provisioner

import (
	"context"
	"errors"
	"fmt"

	log "k8s.io/klog/v2"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// All errors that can be returned by DriverRevokeBucketAccess.
var (
	ErrEmptyAccountID             = errors.New("empty accountID")
	ErrExistingPolicyIsEmpty      = errors.New("existing policy is empty")
	ErrFailedToUpdateBucketPolicy = errors.New("failed to update bucket policy")
	ErrFailedToListAccessKeys     = errors.New("failed to list access keys")
	ErrFailedToDeleteAccessKey    = errors.New("failed to delete access key")
	ErrFailedToDeleteUser         = errors.New("failed to delete user")
)

// All warnings that can be returned by DriverRevokeBucketAccess.
var (
	WarnBucketNotFound = "Bucket not found."
	WarnUserNotFound   = "User not found."
)

// DriverRevokeBucketAccess revokes access from Bucket on specific Object Storage Platform.
func (p *Provisioner) DriverRevokeBucketAccess(ctx context.Context,
	req *cosi.DriverRevokeBucketAccessRequest,
) (*cosi.DriverRevokeBucketAccessResponse, error) {

	// Check if bucketID is not empty.
	if req.GetBucketId() == "" {
		log.Errorf("empty bucket ID", "action", "DriverRevokeBucketAccess")
		return nil, fmt.Errorf("empty bucket ID")
	}

	// Check if bucket access name is not empty.
	if req.GetAccountId() == "" {
		log.Errorf("empty account ID", "action", "DriverRevokeBucketAccess", "bucketID", req.GetBucketId())
		return nil, fmt.Errorf("empty bucket access name")
	}

	// Get bucket name from bucketID.

	userName := req.AccountId
	bucketName, err := getBucketName(req.GetBucketId())
	if err != nil {
		log.ErrorS(err, "failed to convert bucket name", "action", "DriverRevokeBucketAccess", "bucketID", req.GetBucketId())
		return nil, err
	}

	// Check if bucket for revoking access exists.
	bucket, err := p.Powerscale.GetBucket(bucketName)
	if err != nil {
		log.ErrorS(err, "error fetching bucket", "action", "DriverRevokeBucketAccess", "bucket", bucketName)
		return nil, err
	}
	if bucket != nil {
		if err := p.Powerscale.DeleteACL(bucketName, userName); err != nil {
			log.ErrorS(err, "error removing acl", "action", "DriverRevokeBucketAccess", "bucket", bucketName, "userName", userName)
			return nil, err
		}
	}

	key, err := p.Powerscale.GetKey(userName)
	if err != nil {
		log.ErrorS(err, "error fetching key", "action", "DriverRevokeBucketAccess", "bucket", bucketName, "userName", userName)
		return nil, err
	}
	if key != nil {
		if err := p.Powerscale.DeleteKey(userName); err != nil {
			log.ErrorS(err, "error deleting key", "action", "DriverRevokeBucketAccess", "bucket", bucketName, "userName", userName)
			return nil, err
		}
	}

	user, err := p.Powerscale.GetUser(userName)
	if err != nil {
		log.ErrorS(err, "error fetching user", "action", "DriverRevokeBucketAccess", "bucket", bucketName, "userName", userName)
		return nil, err
	}
	if user != nil {
		if err := p.Powerscale.DeleteUser(userName); err != nil {
			log.ErrorS(err, "error deleting user", "action", "DriverRevokeBucketAccess", "bucket", bucketName, "userName", userName)
			return nil, err
		}
	}

	return &cosi.DriverRevokeBucketAccessResponse{}, nil
}
