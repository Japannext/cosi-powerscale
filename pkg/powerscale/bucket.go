package powerscale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	log "k8s.io/klog/v2"
)

// getBucket is used to obtain bucket info from the Provisioner.
func (s *Server) GetBucket(bucketName string) (*Bucket, error) {
	url := fmt.Sprintf("%s/platform/14/protocols/s3/buckets/%s?zone=%s", s.apiEndpoint, bucketName, s.zone)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	s.basicAuth(req)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("Unexpected status code %d: %s", resp.StatusCode, body)
	}

	var bucketList BucketList
	if err := json.Unmarshal(body, &bucketList); err != nil {
		return nil, err
	}
	if len(bucketList.Buckets) == 0 {
		return nil, nil
	}

	return bucketList.Buckets[0], nil
}

// createBucket is used to create bucket on the Provisioner.
func (s *Server) CreateBucket(bucketName string) error {
	bucket := &Bucket{
		Name:            bucketName,
		Path:            s.getBucketPath(bucketName),
		CreatePath:      true,
		ObjectACLPolicy: "replace",
		Acl:             []ACL{},
		Owner:           "root",
		Description:     "Created by cosi-powerscale",
	}

	data, err := json.Marshal(&bucket)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/platform/14/protocols/s3/buckets?zone=%s", s.apiEndpoint, s.zone)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	s.basicAuth(req)
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

	log.InfoS("CreateBucket success", "bucket", bucket.Name)
	return nil
}

func (s *Server) DeleteBucket(bucketName string) error {
	url := fmt.Sprintf("%s/platform/14/protocols/s3/buckets/%s?zone=%s", s.apiEndpoint, bucketName, s.zone)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	s.basicAuth(req)
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
		log.InfoS("DeleteBucket success (not found)", "bucket", bucketName)
		return nil
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("Unexpected status code %d: %s", resp.StatusCode, body)
	}

	log.InfoS("DeleteBucket success", "bucket", bucketName)

	return nil
}

func (s *Server) getBucketPath(bucketName string) string {
	return strings.Join([]string{s.basePath, bucketName}, "/")
}

func (s *Server) DeleteDirectoryForBucket(bucketName string) error {
	path := s.getBucketPath(bucketName)
	url := fmt.Sprintf("%s/namespace/%s?recursive=true", s.apiEndpoint, path)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	s.basicAuth(req)
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
		log.InfoS("DeleteDirectory success (not found)", "directory", path)
		return nil
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("Unexpected status code %d: %s", resp.StatusCode, body)
	}

	log.InfoS("DeleteDirectory success", "directory", path)
	return nil
}
