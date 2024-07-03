package provisioner

import (
	"context"
	"errors"
	"fmt"

	log "k8s.io/klog/v2"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// All errors that can be returned by DriverCreateBucket.
var (
	ErrEmptyBucketName      = errors.New("empty bucket name")
	ErrFailedToCreateBucket = errors.New("failed to create bucket")
)

// DriverCreateBucket creates Bucket on specific Object Storage Platform.
func (p *Provisioner) DriverCreateBucket(
	ctx context.Context,
	req *cosi.DriverCreateBucketRequest,
) (*cosi.DriverCreateBucketResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	bucketName := req.GetName()

	// Check if bucket name is not empty.
	if bucketName == "" {
		log.Errorf("empty bucket name", "action", "DriverCreateBucket", "bucket", bucketName)
		return nil, fmt.Errorf("Empty bucket name")
	}

	// Check if bucket exist
	bucket, err := p.Powerscale.GetBucket(bucketName)
	if err != nil {
		log.ErrorS(err, "error attempting to fetch bucket", "action", "DriverCreateBucket", "bucket", bucketName)
		return nil, err
	}
	if bucket != nil {
		return nil, nil
	}

	// Create bucket.
	err = p.Powerscale.CreateBucket(bucketName)
	if err != nil {
		log.ErrorS(err, "error creating bucket", "action", "DriverCreateBucket", "bucket", bucketName)
		return nil, err
	}

	// Return response.
	return &cosi.DriverCreateBucketResponse{
		BucketId: fmt.Sprintf("%s-%s", p.ID(), bucketName),
	}, nil
}
