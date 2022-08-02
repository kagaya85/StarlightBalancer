package service

import (
	"context"

	"starlight/balancer/client"

	pb "starlight/api/services/push/v1"
	"starlight/services/push/internal/biz"
)

// global weight list for load balance
var GlobalBalancer *client.BalancerClient

type PushService struct {
	pb.UnimplementedPushServiceServer

	uc *biz.PusherUsecase
}

func NewPusherService(uc *biz.PusherUsecase) *PushService {
	return &PushService{uc: uc}
}

func (s *PushService) PushVideo(ctx context.Context, req *pb.PushRequest) (*pb.PushResponse, error) {
	if err := s.uc.Call(ctx, GlobalBalancer.Default); err != nil {
		return nil, err
	}
	return &pb.PushResponse{Code: 200}, nil
}
