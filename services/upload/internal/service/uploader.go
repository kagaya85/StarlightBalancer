package service

import (
	"context"

	"starlight/balancer/client"

	pb "starlight/api/services/upload/v1"
	"starlight/services/upload/internal/biz"
)

// global weight list for load balance
var GlobalBalancer *client.BalancerClient

type UploaderService struct {
	pb.UnimplementedUploaderServer

	uc *biz.UploaderUsecase
}

func NewUploaderService(uc *biz.UploaderUsecase) *UploaderService {
	return &UploaderService{uc: uc}
}

func (s *UploaderService) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	s.uc.Call(ctx, GlobalBalancer.Random)
	return &pb.UploadResponse{}, nil
}
