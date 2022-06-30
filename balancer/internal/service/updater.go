package service

import (
	"context"
	v1 "starlight/api/balancer/v1"
	"starlight/balancer/internal/biz"
	"starlight/balancer/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// WeightUpdaterService is a weight update service
type WeightUpdaterService struct {
	v1.UnimplementedWeightUpdaterServer

	updater *biz.WeightUpdater

	log log.Helper

	updateInterval time.Duration
}

func NewWeightUpdaterService(c *conf.Updater, updater *biz.WeightUpdater, logger log.Logger) *WeightUpdaterService {
	return &WeightUpdaterService{
		updater:        updater,
		log:            *log.NewHelper(logger),
		updateInterval: c.UpdateInterval.AsDuration(),
	}
}

func (s *WeightUpdaterService) Update(in *v1.UpdateRequeset, stream v1.WeightUpdater_UpdateServer) error {
	// 假设一个实例只有一个服务
	svcInfo := in.Info[0]
	// 更新服务实例信息
	s.updater.UpdateInstance(biz.InstanceInfo{
		ID:      in.Instance,
		Service: svcInfo.Service,
		Port:    svcInfo.Port,
		Pod:     in.Pod,
		Node:    in.Node,
		Zone:    in.Zone,
	})

	operations := make([]biz.Operation, 0, len(svcInfo.Operations))
	svc := svcInfo.Service
	for _, op := range svcInfo.Operations {
		operations = append(operations, biz.NewOperation(svc, op))
	}
	upstreamOperations := []biz.Operation{}
	for _, ops := range in.Upstream {
		svc = ops.Service
		for _, op := range ops.Operations {
			upstreamOperations = append(upstreamOperations, biz.NewOperation(svc, op))
		}
	}
	s.updater.UpdateDependency(svcInfo.Service, operations, upstreamOperations)

	defer s.updater.RemoveInstance(biz.Instance(in.Instance))

	// 设置定时更新权重列表
	ticker := time.NewTicker(s.updateInterval)
	for {
		weightsList := s.updater.UpdateWeights(context.Background(), biz.Instance(in.Instance))
		wl := map[string]*v1.Weight{}
		for op, insWeights := range weightsList {
			iw := map[string]int32{}
			for endpoint, weight := range insWeights {
				iw[string(endpoint)] = int32(weight)
			}
			wl[string(op)] = &v1.Weight{InstanceWeight: iw}
		}
		if err := stream.Send(&v1.UpdateReply{WeightList: wl}); err != nil {
			log.Infof("update stream error: %v", err)
			break
		}
		log.Infof("update instance %s weight list success", in.Instance)
		<-ticker.C
	}
	ticker.Stop()
	return nil
}
