package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/Taskon-xyz/kratos-layout/api/_pb/example"
	"github.com/Taskon-xyz/kratos-layout/internal/biz/sample"
)

// SampleService is an sample service.
type SampleService struct {
	example.UnimplementedSampleServer
}

// NewSampleService new an sample service.
func NewSampleService() *SampleService {
	return &SampleService{}
}

// CheckSignupStatus implements example.AuthServer.
func (s *SampleService) CheckSignupStatus(ctx context.Context, in *example.CheckSignupStatusRequest) (*example.CheckSignupStatusResponse, error) {
	log.Context(ctx).Infow("msg", "CheckSignupStatus", "request", in.String())
	registered, err := sample.CheckSignupStatus(ctx, in.Type, in.Value)
	if err != nil {
		log.Context(ctx).Errorf("CheckSignupStatus failed; detail: %+v", err)
		return nil, err
	}

	return &example.CheckSignupStatusResponse{
		Registered: registered,
	}, nil
}
