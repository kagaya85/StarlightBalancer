package service

import (
	v1 "starlight/api/balancer/v1"
)

// WeightUpdaterService is a weight update service
type WeightUpdaterService struct {
	v1.UnimplementedWeightUpdaterServer
}

func NewWeightUpdaterService() *WeightUpdaterService {
	return &WeightUpdaterService{}
}

func (u *WeightUpdaterService) Update(in *v1.UpdateRequeset, stream v1.WeightUpdater_UpdateServer) error {
	return nil
}
