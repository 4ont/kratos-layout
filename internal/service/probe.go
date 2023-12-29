package service

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"

	probeapi "github.com/4ont/kratos-layout/api/_pb/probe"
	"github.com/4ont/kratos-layout/internal/biz/probe"
)

// ProbeService is a probe service.
type ProbeService struct {
	probeapi.UnimplementedProbeServer
}

// NewProbeService new a probe service.
func NewProbeService() *ProbeService {
	return &ProbeService{}
}

// Healthy implements probe.ProbeServer.
func (s *ProbeService) Healthy(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// Ready implements probe.ProbeServer.
func (s *ProbeService) Ready(ctx context.Context, in *structpb.Struct) (*probeapi.ReadinessProbeResponse, error) {
	// 通过将structpb转成json，再从json重新解析到go struct达到复用已有的go struct定义的目的。
	// 这个方法的缺点是多了一轮json的转换，牺牲性能且不利于API的后续维护，并不推荐这样做
	bys, err := in.MarshalJSON()
	if err != nil {
		log.Context(ctx).Errorf("marshal to json failed; detail: %+v", err)
		return nil, err
	}

	var anything map[string]interface{}
	if err = json.Unmarshal(bys, &anything); err != nil {
		log.Context(ctx).Error("parse argument failed", zap.Error(err))
		return nil, err
	}

	err = probe.Ready(ctx, anything)
	if err != nil {
		log.Context(ctx).Error("Readiness probe failed", zap.Error(err))
		return nil, err
	}

	return &probeapi.ReadinessProbeResponse{Status: "success"}, nil
}
