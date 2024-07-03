package powerscale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "k8s.io/klog/v2"
)

func aclInsert(acls []ACL, newAcl ACL) []ACL {
	for i, acl := range acls {
		if acl.Grantee.Name == newAcl.Grantee.Name {
			acls[i] = newAcl
			return acls
		}
	}
	acls = append(acls, newAcl)
	return acls
}

func (s *Server) EnsureACL(bucketName, userName, permission string) error {
	bucket, err := s.GetBucket(bucketName)
	if err != nil {
		return err
	}

	newAcl := ACL{
		Grantee:    &AclUser{Type: "user", Name: userName},
		Permission: permission,
	}
	bucketUpdate := &PartialBucket{
		Acl: aclInsert(bucket.Acl, newAcl),
	}

	data, err := json.Marshal(&bucketUpdate)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/platform/14/protocols/s3/buckets/%s?zone=%s", s.apiEndpoint, bucketName, s.zone)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	s.basicAuth(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("Unexpected status code %d: %s", resp.StatusCode, body)
	}

	log.InfoS("EnsureACL success", "bucket", bucket.Name, "userName", userName)
	return nil
}

func (s *Server) DeleteACL(bucketName, userName string) error {
	bucket, err := s.GetBucket(bucketName)
	if err != nil {
		return err
	}
	if bucket == nil {
		log.InfoS("DeleteACL success (bucket not found)", "bucket", bucketName)
		return nil
	}

	newAcls := []ACL{}
	if len(bucket.Acl) > 0 {
		for _, acl := range bucket.Acl {
			if acl.Grantee.Name != userName {
				newAcls = append(newAcls, acl)
			}
		}
	}
	bucketUpdate := &PartialBucket{
		Acl: newAcls,
	}

	data, err := json.Marshal(&bucketUpdate)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/platform/14/protocols/s3/buckets/%s?zone=%s", s.apiEndpoint, bucketName, s.zone)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	s.basicAuth(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		log.InfoS("DeleteACL success (not found)", "bucket", bucketName, "userName", userName)
		return nil
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("Unexpected status code %d: %s", resp.StatusCode, body)
	}

	log.InfoS("DeleteACL success", "bucket", bucket.Name, "userName", userName)
	return nil
}
