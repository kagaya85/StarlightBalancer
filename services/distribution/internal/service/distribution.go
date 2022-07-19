package service

import (
	"context"

	"starlight/balancer/client"

	pb "starlight/api/services/distribution/v1"
	"starlight/services/distribution/internal/biz"
)

// global weight list for load balance
var GlobalBalancer *client.BalancerClient

type DistributionService struct {
	pb.UnimplementedDistributionServiceServer

	uc *biz.DistributionUsecase
}

func NewDistributionService(uc *biz.DistributionUsecase) *DistributionService {
	return &DistributionService{uc: uc}
}

func (s *DistributionService) Distribution(ctx context.Context, req *pb.DistributeRequest) (*pb.DistributeResponse, error) {
	return &pb.DistributeResponse{Code: 200}, nil
}
