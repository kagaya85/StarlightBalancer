package biz

import (
	"context"
	"sync"
)

type MetricSource interface {
	GetByInstance(context.Context, InstanceInfo) Metric
}

type Metric struct {
	CPU float64 // pod cpu usage rate (0~1)
	Mem float64 // pod memory usage rate (0~1)

	Load            int // node exporter
	ConnectionCount int // node_netstat_Tcp_ActiveOpens

	ResponseTime int     // app range response time (ms)
	SuccessRate  float64 // app range success rate (0~1)
}

type metricMap struct {
	mu sync.RWMutex
	m  map[Instance]Metric
}

func NewMetricMap() *metricMap {
	return &metricMap{
		m: map[Instance]Metric{},
	}
}

func (m *metricMap) GetMetric(id Instance) Metric {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m[id]
}

func (m *metricMap) Update(id Instance, observed Metric) {
	ewma := func(old, new float64) float64 {
		beta := 0.8
		return beta*old + (1-beta)*new
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if history, has := m.m[id]; has {
		m.m[id] = Metric{
			CPU:             ewma(history.CPU, observed.CPU),
			Mem:             ewma(history.Mem, observed.Mem),
			Load:            int(ewma(float64(history.Load), float64(observed.Load))),
			ConnectionCount: int(ewma(float64(history.ConnectionCount), float64(observed.ConnectionCount))),
			ResponseTime:    int(ewma(float64(history.ResponseTime), float64(observed.ResponseTime))),
			SuccessRate:     ewma(history.SuccessRate, observed.SuccessRate),
		}
	} else {
		m.m[id] = observed
	}
}

func (m *metricMap) Clear(id Instance) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, id)
}
