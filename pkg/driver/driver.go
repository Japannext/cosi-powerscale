package driver

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"

	"google.golang.org/grpc"
	spec "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/japannext/cosi-powerscale/pkg/config"
	"github.com/japannext/cosi-powerscale/pkg/identity"
	"github.com/japannext/cosi-powerscale/pkg/provisioner"
	log "k8s.io/klog/v2"
)

const (
	// COSISocket is a default location of COSI API UNIX socket.
	socket = "/var/lib/cosi/cosi.sock"
)

type Driver struct {
	server *grpc.Server
	lis    net.Listener
}

func New(cfg *config.Config) (*Driver, error) {
	identityName := fmt.Sprintf("%s.powerscale.cosi.japannext.co.jp", cfg.Name)
	identityServer := identity.New(identityName)
	provisionerServer := provisioner.New(cfg)

	options := []grpc.ServerOption{}
	server := grpc.NewServer(options...)
	spec.RegisterIdentityServer(server, identityServer)
	spec.RegisterProvisionerServer(server, provisionerServer)

	if _, err := os.Stat(socket); !errors.Is(err, fs.ErrNotExist) {
		if err := os.RemoveAll(socket); err != nil {
			log.Fatal(err)
		}
	}

	listener, err := net.Listen("unix", socket)
	if err != nil {
		return nil, fmt.Errorf("failed to announce on the local network address: %w", err)
	}

	log.InfoS("Listening on socket", "socket", socket)

	return &Driver{server, listener}, nil
}

func (d *Driver) Run(ctx context.Context) error {

	ready := make(chan struct{})
	go func() {
		close(ready)
		if err := d.server.Serve(d.lis); err != nil {
			log.Fatal(err)
		}
	}()

	<-ready
	log.Info("gRPC server started")
	<-ctx.Done()

	d.server.GracefulStop()

	return nil
}
