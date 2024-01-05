package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/Taskon-xyz/kratos-layout/api/_pb/example"
)

type PortalService struct {
	pb.UnimplementedPortalServer
}

func NewPortalService() *PortalService {
	return &PortalService{}
}

func (s *PortalService) Nothing(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {

	return &emptypb.Empty{}, nil
}
