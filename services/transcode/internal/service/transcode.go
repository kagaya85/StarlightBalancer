package service

import (
	"context"

	"starlight/balancer/client"

	pb "starlight/api/services/transcode/v1"
	"starlight/services/transcode/internal/biz"
)

// global weight list for load balance
var GlobalBalancer *client.BalancerClient

type TranscodeService struct {
	pb.UnimplementedTranscodeServiceServer

	uc *biz.TranscodeUsecase
}

func NewTranscodeService(uc *biz.TranscodeUsecase) *TranscodeService {
	return &TranscodeService{uc: uc}
}

func (s *TranscodeService) Transcode(ctx context.Context, req *pb.TranscodeRequest) (*pb.TranscodeResponse, error) {
	return &pb.TranscodeResponse{Target: "foobar"}, nil
}
