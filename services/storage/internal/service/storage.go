package service

import (
	"context"

	"starlight/balancer/client"

	pb "starlight/api/services/storage/v1"
	"starlight/services/storage/internal/biz"
)

// global weight list for load balance
var GlobalBalancer *client.BalancerClient

type StorageService struct {
	pb.UnimplementedStorageServiceServer

	uc *biz.StorageUsecase
}

func NewStorageService(uc *biz.StorageUsecase) *StorageService {
	return &StorageService{uc: uc}
}

func (s *StorageService) Save(ctx context.Context, req *pb.SaveRequest) (*pb.SaveResponse, error) {
	s.uc.Call(ctx, GlobalBalancer.Default)
	return &pb.SaveResponse{Result: "saved"}, nil
}
