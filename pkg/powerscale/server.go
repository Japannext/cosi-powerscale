package powerscale

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	log "k8s.io/klog/v2"

	"github.com/japannext/cosi-powerscale/pkg/config"
)

const (
	CreateBucketTraceName       = "CreateBucketRequest"
	DeleteBucketTraceName       = "DeleteBucketRequest"
	GrantBucketAccessTraceName  = "GrantBucketAccessRequest"
	RevokeBucketAccessTraceName = "RevokeBucketAccessRequest"
	splitNumber                 = 2
	allowEffect                 = "Allow"
)

var (
	ErrInvalidBucketID           = errors.New("invalid bucketID")
	ErrFailedToCheckBucketExists = errors.New("failed to check bucket existence")
	ErrFailedToMarshalPolicy     = errors.New("failed to marshal policy into JSON")
	ErrFailedToCheckPolicyExists = errors.New("failed to check bucket policy existence")
	ErrFailedToCheckUserExists   = errors.New("failed to check for user existence")
	ErrInvalidRequest            = errors.New("incoming request invalid")
	defaultTimeout               = time.Second * 20
)

// Server is implementation of driver.Driver interface for ObjectScale platform.
type Server struct {
	Name        string
	S3Endpoint  string
	S3Region    string
	apiEndpoint string
	cacert      string
	zone        string
	username    string
	password    string
	client      *http.Client
	// The path base to use in OneFS, so all buckets are in <basePath>/<bucketName>
	basePath string
}

func (s *Server) basicAuth(req *http.Request) {
	auth := base64.StdEncoding.EncodeToString([]byte(s.username + ":" + s.password))
	req.Header.Add("Authorization", "Basic "+auth)
}

func New(cfg *config.Config) *Server {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.TlsInsecureSkipVerify,
	}
	if cfg.TlsCacert != "" {
		rootCAs := x509.NewCertPool()
		pem, err := ioutil.ReadFile(cfg.TlsCacert)
		if err != nil {
			log.Fatal(err)
		}
		if ok := rootCAs.AppendCertsFromPEM(pem); !ok {
			log.Errorf("No cert found in %s", cfg.TlsCacert)
		}
		tlsConfig.RootCAs = rootCAs
	}
	if cfg.TlsClientCert != "" && cfg.TlsClientKey != "" {
		cert, err := tls.LoadX509KeyPair(cfg.TlsClientCert, cfg.TlsClientKey)
		if err != nil {
			log.Fatal(err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return &Server{
		Name:        cfg.Name,
		apiEndpoint: cfg.ApiEndpoint,
		username:    cfg.ApiUsername,
		password:    cfg.ApiPassword,
		zone:        cfg.Zone,
		S3Endpoint:  cfg.S3Endpoint,
		S3Region:    cfg.S3Region,
		basePath:    cfg.BasePath,
		client: &http.Client{
			Transport: &http.Transport{TLSClientConfig: tlsConfig},
		},
	}
}
