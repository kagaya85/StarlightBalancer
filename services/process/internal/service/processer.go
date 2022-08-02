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

func (s *ProcessService) Process(ctx context.Context, req *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	result := ""
	res, err := s.uc.CallAudit(ctx, GlobalBalancer.Default)
	result = "audit:" + res
	if err != nil {
		return nil, err
	}

	res, err = s.uc.CallTranscode(ctx, GlobalBalancer.Default)
	result += "/transcode:" + res
	if err != nil {
		return nil, err
	}

	res, err = s.uc.CallStorage(ctx, GlobalBalancer.Default)
	result += "/storage:" + res
	if err != nil {
		return nil, err
	}

	return &pb.ProcessResponse{Result: result}, nil
}
