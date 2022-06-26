package biz

func (u *WeightUpdater) IsOverload(id Instance) int {
	m := u.insMetrics.GetMetric(id)
	if m.CPU > 0.8 || m.Mem > 0.8 || m.Load > 90 || m.ConnectionCount > 100 || m.ResponseTime > 1000 || m.SuccessRate < 0.8 {
		return 1
	}

	return 0
}
