package powerscale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "k8s.io/klog/v2"
)

func (s *Server) GetUser(userName string) (*User, error) {
	url := fmt.Sprintf("%s/platform/14/auth/users/%s?zone=%s", s.apiEndpoint, userName, s.zone)
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
	var userList UserList
	if err := json.Unmarshal(body, &userList); err != nil {
		return nil, err
	}
	if len(userList.Users) == 0 {
		return nil, nil
	}
	return userList.Users[0], nil
}

func (s *Server) CreateUser(userName string) error {
	data, err := json.Marshal(&User{
		Name:    userName,
		Enabled: true,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/platform/14/auth/users?zone=%s", s.apiEndpoint, s.zone)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
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
	if resp.StatusCode > 299 {
		return fmt.Errorf("Unexpected status code %d: %s", resp.StatusCode, body)
	}

	log.InfoS("CreateUser success", "userName", userName)
	return nil
}

func (s *Server) DeleteUser(userName string) error {
	url := fmt.Sprintf("%s/platform/14/auth/users/%s?zone=%s", s.apiEndpoint, userName, s.zone)
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
		log.InfoS("DeleteUser success (not found)", "userName", userName)
		return nil
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("Unexpected status code %d: %s", resp.StatusCode, body)
	}

	log.InfoS("DeleteUser success", "userName", userName)
	return nil
}
