package service

import (
	"context"

	kerr "github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/4ont/kratos-layout/api/_pb/example"
	"github.com/4ont/kratos-layout/internal/biz/sample"
)

type AdminService struct {
	pb.UnimplementedAdminServer
}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) Nothing(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	_, ok := sample.AdminFromContext(ctx)
	if !ok {
		return nil, kerr.New(4000, "no privilege", "you are not admin")
	}

	return &emptypb.Empty{}, nil
}
