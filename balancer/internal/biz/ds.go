package biz

func fusion(m [6]float64) bool {
	a, b := 0.0, 0.0
	for _, v := range m {
		a *= v
		b *= (1 - v)
	}

	k := 1 - (a + b)

	if (k-1) < 1e-6 || k > 1 {
		return false
	}

	c_1 := 1 / (1 - k)
	return c_1*a > c_1*b
}

func confidence(metric Metric, weights [6][2]int) [6]float64 {
	m := [6]float64{}
	if metric.Load > 90 {
		m[0] = (1 + (float64(weights[0][0]) / float64(weights[0][1]))) / 2
	} else {
		m[0] = (1 - (float64(weights[0][0]) / float64(weights[0][1]))) / 2
	}

	if metric.ConnectionCount > 100 {
		m[1] = (1 + (float64(weights[1][0]) / float64(weights[1][1]))) / 2
	} else {
		m[1] = (1 - (float64(weights[1][0]) / float64(weights[1][1]))) / 2
	}

	if metric.CPU > 0.8 {
		m[2] = (1 + (float64(weights[2][0]) / float64(weights[2][1]))) / 2
	} else {
		m[2] = (1 - (float64(weights[2][0]) / float64(weights[2][1]))) / 2
	}

	if metric.Mem > 0.8 {
		m[3] = (1 + (float64(weights[3][0]) / float64(weights[3][1]))) / 2
	} else {
		m[3] = (1 - (float64(weights[3][0]) / float64(weights[3][1]))) / 2
	}

	if metric.ResponseTime > 1000 {
		m[4] = (1 + (float64(weights[4][0]) / float64(weights[4][1]))) / 2
	} else {
		m[4] = (1 - (float64(weights[4][0]) / float64(weights[4][1]))) / 2
	}

	if metric.SuccessRate < 0.8 {
		m[5] = (1 + (float64(weights[5][0]) / float64(weights[5][1]))) / 2
	} else {
		m[5] = (1 - (float64(weights[5][0]) / float64(weights[5][1]))) / 2
	}

	return m
}

func (u *WeightUpdater) updateDSWeights(id Instance, results []bool) {
	count := 0
	for _, res := range results {
		if res {
			count++
		}
	}
	u.dsmu.Lock()
	defer u.dsmu.Unlock()
	ws := u.dsWeights[id]
	for i, isOverload := range results {
		if isOverload {
			ws[i][0]++
		}
		ws[i][1] += count
	}
}

func (u *WeightUpdater) IsOverload(id Instance) int {
	m := u.insMetrics.GetMetric(id)
	ws, has := u.dsWeights[id]
	if !has {
		ws = [6][2]int{{1, 6}, {1, 6}, {1, 6}, {1, 6}, {1, 6}, {1, 6}}
	}

	fusion(confidence(m, ws))

	if m.CPU > 0.8 || m.Mem > 0.8 || m.Load > 90 || m.ConnectionCount > 100 || m.ResponseTime > 1000 || m.SuccessRate < 0.8 {
		return 1
	}

	return 0
}
