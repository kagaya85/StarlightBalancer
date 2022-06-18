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
	ticker := time.NewTicker(s.updateInterval)
	for {
		weightsList := s.updater.UpdateInstance(context.Background(), in.InstanceID)
		wl := map[string]*v1.Weight{}
		for op, insWeights := range weightsList {
			iw := map[string]int32{}
			for ins, weight := range insWeights {
				iw[string(ins)] = int32(weight)
			}
			wl[op] = &v1.Weight{InstanceWeight: iw}
		}
		if err := stream.Send(&v1.UpdateReply{WeightList: wl}); err != nil {
			log.Infof("update stream error: %v", err)
			break
		}
		log.Infof("update instance %s weight list success", in.InstanceID)
		<-ticker.C
	}
	ticker.Stop()
	return nil
}
