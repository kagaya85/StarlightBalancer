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
	if err := s.uc.CallProcesss(ctx, GlobalBalancer.Default); err != nil {
		return &pb.UploadResponse{Result: "process error"}, err
	}
	if err := s.uc.CallPush(ctx, GlobalBalancer.Default); err != nil {
		return &pb.UploadResponse{Result: "push error"}, err
	}
	return &pb.UploadResponse{Result: "ok"}, nil
}
