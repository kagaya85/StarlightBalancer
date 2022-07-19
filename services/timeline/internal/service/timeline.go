package service

import (
	"context"

	"starlight/balancer/client"

	pb "starlight/api/services/timeline/v1"
	"starlight/services/timeline/internal/biz"
)

// global weight list for load balance
var GlobalBalancer *client.BalancerClient

type TimelineService struct {
	pb.UnimplementedTimelineServiceServer

	uc *biz.TimelineUsecase
}

func NewTimelineService(uc *biz.TimelineUsecase) *TimelineService {
	return &TimelineService{uc: uc}
}

func (s *TimelineService) PushTimeline(ctx context.Context, req *pb.PushTimelineRequest) (*pb.PushTimelineResponse, error) {
	return &pb.PushTimelineResponse{Status: 200}, nil
}
