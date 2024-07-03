package provisioner

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
	"sigs.k8s.io/container-object-storage-interface-provisioner-sidecar/pkg/consts"

	log "k8s.io/klog/v2"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

const (
	maxUsernameLength = 64
)

// All errors that can be returned by DriverGrantBucketAccess.
// DriverGrantBucketAccess provides access to Bucket on specific Object Storage Platform.
func (p *Provisioner) DriverGrantBucketAccess(
	ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest,
) (*cosi.DriverGrantBucketAccessResponse, error) {

	// Check if bucketID is not empty.
	if req.GetBucketId() == "" {
		log.Errorf("empty bucket ID", "action", "DriverGrantBucketAccess")
		return nil, fmt.Errorf("empty bucket ID")
	}

	// Check if bucket access name is not empty.
	if req.GetName() == "" {
		log.Errorf("empty bucket access name", "action", "DriverGrantBucketAccess", "bucketID", req.GetBucketId())
		return nil, fmt.Errorf("empty bucket access name")
	}

	// Get bucket name from bucketID.
	bucketName, err := getBucketName(req.GetBucketId())
	if err != nil {
		log.ErrorS(err, "failed to convert bucket name", "action", "DriverGrantBucketAccess", "bucketID", req.GetBucketId())
		return nil, err
	}

	// Equals to "ba-<uid>" with <uid> being the UID of the BucketAccess object.
	userName := req.GetName()
	user, err := p.Powerscale.GetUser(userName)
	if err != nil {
		log.ErrorS(err, "failed to fetch user", "action", "DriverGrantBucketAccess", "bucket", bucketName, "userName", userName)
		return nil, err
	}

	// Create user
	if user == nil {
		if err := p.Powerscale.CreateUser(userName); err != nil {
			log.ErrorS(err, "failed to create user", "action", "DriverGrantBucketAccess", "bucket", bucketName, "userName", userName)
			return nil, fmt.Errorf("failed while creating user %s: %w", userName, err)
		}
	}

	if err := p.Powerscale.EnsureACL(bucketName, userName, "FULL_CONTROL"); err != nil {
		log.ErrorS(err, "failed to add ACL", "action", "DriverGrantBucketAccess", "bucket", bucketName, "userName", userName)
		return nil, err
	}

	// Create Key
	accessKey, err := p.Powerscale.CreateKey(userName)
	if err != nil {
		log.ErrorS(err, "failed to create s3 key", "action", "DriverGrantBucketAccess", "bucket", bucketName, "userName", userName)
		return nil, err
	}

	credentials := assembleCredentials(accessKey, p.Powerscale.S3Endpoint, userName, bucketName)
	return &cosi.DriverGrantBucketAccessResponse{AccountId: userName, Credentials: credentials}, nil
}

// assembleCredentials assembles credentials details and adds them to the credentialRepo.
func assembleCredentials(
	accessKey *iam.CreateAccessKeyOutput,
	s3Endpoint,
	userName,
	bucketName string,
) map[string]*cosi.CredentialDetails {

	secretsMap := make(map[string]string)
	secretsMap[consts.S3SecretAccessKeyID] = *accessKey.AccessKey.AccessKeyId
	secretsMap[consts.S3SecretAccessSecretKey] = *accessKey.AccessKey.SecretAccessKey
	secretsMap[consts.S3Endpoint] = s3Endpoint
	secretsMap["bucketName"] = bucketName

	credentialDetails := cosi.CredentialDetails{Secrets: secretsMap}
	credentials := make(map[string]*cosi.CredentialDetails)
	credentials[consts.S3Key] = &credentialDetails

	return credentials
}
