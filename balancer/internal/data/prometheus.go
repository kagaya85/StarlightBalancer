package data

import (
	"context"
	"starlight/balancer/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type metricSource struct {
	metricData *MetricData

	log *log.Helper
}

func NewMetricSource(metricData *MetricData, logger log.Logger) biz.MetricSource {
	return &metricSource{
		metricData: metricData,
		log:        log.NewHelper(logger),
	}
}

func (m *metricSource) GetByInstance(ctx context.Context, ins biz.InstanceInfo) biz.Metric {
	cpu := m.metricData.QueryPodCPUUsage(ctx, ins.Pod)
	mem := m.metricData.QueryPodMemUsage(ctx, ins.Pod)
	load := m.metricData.QueryNodeLoad(ctx, ins.Node)
	connenctionCount := m.metricData.QueryNodeConnectionCount(ctx, ins.Node)
	responseTime := m.metricData.QueryAppResponseTime(ctx, ins.Pod)
	successRate := m.metricData.QueryAppSuccessRate(ctx, ins.Pod)
	return biz.Metric{
		CPU:             cpu,
		Mem:             mem,
		Load:            load,
		ConnectionCount: connenctionCount,
		ResponseTime:    responseTime,
		SuccessRate:     successRate,
	}
}
