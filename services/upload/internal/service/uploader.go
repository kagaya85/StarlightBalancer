package service

import (
	"context"

	"starlight/balancer/client"

	pb "starlight/api/services/upload/v1"
	"starlight/services/upload/internal/biz"
)

// global weight list for load balance
var GlobalBalancer *client.BalancerClient

type UploadService struct {
	pb.UnimplementedUploadServiceServer

	uc *biz.UploaderUsecase
}

func NewUploaderService(uc *biz.UploaderUsecase) *UploadService {
	return &UploadService{uc: uc}
}

func (s *UploadService) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	s.uc.CallProcesss(ctx, GlobalBalancer.Default)
	s.uc.CallPush(ctx, GlobalBalancer.Default)
	return &pb.UploadResponse{Result: "hello"}, nil
}
