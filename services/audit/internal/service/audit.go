package service

import (
	"context"

	"starlight/balancer/client"
	"starlight/services/audit/internal/biz"

	pb "starlight/api/services/audit/v1"
)

// global weight list for load balance
var GlobalBalancer *client.BalancerClient

type AuditService struct {
	pb.UnimplementedAuditServiceServer

	uc *biz.AuditUsecase
}

func NewAuditService(uc *biz.AuditUsecase) *AuditService {
	return &AuditService{uc: uc}
}

func (s *AuditService) Audit(ctx context.Context, req *pb.AuditRequest) (*pb.AuditResponse, error) {
	return &pb.AuditResponse{Result: "OK"}, nil
}
