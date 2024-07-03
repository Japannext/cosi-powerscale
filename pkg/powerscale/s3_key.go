package powerscale

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/service/iam"

	log "k8s.io/klog/v2"
)

func (s *Server) GetKey(userName string) (*Key, error) {
	url := fmt.Sprintf("%s/platform/14/protocols/s3/keys/%s?zone=%s", s.apiEndpoint, userName, s.zone)
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
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("Unexpected status code %d: %s", resp.StatusCode, body)
	}
	var keys Keys
	if err := json.Unmarshal(body, &keys); err != nil {
		return nil, err
	}
	if keys.Keys.AccessID == "" {
		return nil, nil
	}
	return &keys.Keys, nil
}

func (s *Server) CreateKey(userName string) (*iam.CreateAccessKeyOutput, error) {

	url := fmt.Sprintf("%s/platform/14/protocols/s3/keys/%s?zone=%s", s.apiEndpoint, userName, s.zone)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	s.basicAuth(req)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("Unexpected status code %d: %s", resp.StatusCode, body)
	}

	var keys Keys
	if err := json.Unmarshal(body, &keys); err != nil {
		return nil, err
	}

	accessKey := &iam.CreateAccessKeyOutput{AccessKey: &iam.AccessKey{
		AccessKeyId:     &keys.Keys.AccessID,
		SecretAccessKey: &keys.Keys.SecretKey,
	}}
	log.InfoS("CreateKey success", "userName", userName)
	return accessKey, nil
}

func (s *Server) DeleteKey(userName string) error {
	url := fmt.Sprintf("%s/platform/14/protocols/s3/keys/%s?zone=%s", s.apiEndpoint, userName, s.zone)
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
		log.InfoS("DeleteKey success (not found)", "user", userName)
		return nil
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("Unexpected status code %d: %s", resp.StatusCode, body)
	}
	log.InfoS("DeleteKey success", "user", userName)

	return nil
}
