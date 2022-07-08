package service

import (
	"context"

	"starlight/balancer/client"

	pb "starlight/api/services/process/v1"
	"starlight/services/process/internal/biz"
)

// global weight list for load balance
var GlobalBalancer *client.BalancerClient

type ProcessService struct {
	pb.UnimplementedProcessServiceServer

	uc *biz.ProcesserUsecase
}

func NewProcessService(uc *biz.ProcesserUsecase) *ProcessService {
	return &ProcessService{uc: uc}
}

func (s *ProcessService) Upload(ctx context.Context, req *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	s.uc.Call(ctx, GlobalBalancer.Random)
	return &pb.ProcessResponse{Result: "ok"}, nil
}
