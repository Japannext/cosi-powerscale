package provisioner

import (
	"strings"
	"time"

	"github.com/japannext/cosi-powerscale/pkg/config"
	"github.com/japannext/cosi-powerscale/pkg/powerscale"
)

const (
	defaultTimeout = 20 * time.Second
)

type Provisioner struct {
	Powerscale *powerscale.Server
}

func New(cfg *config.Config) *Provisioner {
	return &Provisioner{
		Powerscale: powerscale.New(cfg),
	}
}

func (p *Provisioner) ID() string {
	return p.Powerscale.Name
}

// getBucketName splits BucketID by -, the first element is backendID, the second element is bucketName.
func getBucketName(bucketID string) (string, error) {
	list := strings.SplitN(bucketID, "-", 2)

	if len(list) != 2 || list[1] == "" { //nolint:gomnd
		return "", ErrInvalidBucketID
	}

	return list[1], nil
}
