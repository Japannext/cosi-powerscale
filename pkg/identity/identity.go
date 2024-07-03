package identity

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	log "k8s.io/klog/v2"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

type Server struct {
	name string
}

func New(name string) *Server {
	return &Server{name: name}
}

func (srv *Server) DriverGetInfo(_ context.Context,
	_ *cosi.DriverGetInfoRequest,
) (*cosi.DriverGetInfoResponse, error) {

	if srv.name == "" {
		log.Error("driver name is empty")
		return nil, status.Error(codes.InvalidArgument, "DriverName is empty")
	}

	return &cosi.DriverGetInfoResponse{
		Name: srv.name,
	}, nil
}
