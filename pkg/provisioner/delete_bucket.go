package provisioner

import (
	"context"
	"errors"
	"fmt"

	log "k8s.io/klog/v2"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// All errors that can be returned by DriverDeleteBucket.
var ErrFailedToDeleteBucket = errors.New("bucket was not successfully deleted")

// DriverDeleteBucket deletes Bucket on specific Object Storage Platform.
func (p *Provisioner) DriverDeleteBucket(ctx context.Context,
	req *cosi.DriverDeleteBucketRequest,
) (*cosi.DriverDeleteBucketResponse, error) {

	// Check if bucketID is not empty.
	if req.GetBucketId() == "" {
		log.Errorf("empty bucket ID", "action", "DriverDeleteBucket", "bucketID", req.BucketId)
		return nil, fmt.Errorf("Empty bucket ID")
	}

	// Extract bucket name from bucketID.
	bucketName, err := getBucketName(req.BucketId)
	if err != nil {
		log.ErrorS(err, "error extracting bucket name", "action", "DriverDeleteBucket", "bucketID", req.BucketId)
		return nil, err
	}

	// Delete the directory
	if err := p.Powerscale.DeleteDirectoryForBucket(bucketName); err != nil {
		log.ErrorS(err, "error deleting directory", "action", "DriverDeleteBucket", "bucketID", req.BucketId)
		return &cosi.DriverDeleteBucketResponse{}, err
	}

	// Delete bucket.
	if err := p.Powerscale.DeleteBucket(bucketName); err != nil {
		log.ErrorS(err, "error deleting bucket", "action", "DriverDeleteBucket", "bucketID", req.BucketId)
		return &cosi.DriverDeleteBucketResponse{}, err
	}

	return &cosi.DriverDeleteBucketResponse{}, nil
}
