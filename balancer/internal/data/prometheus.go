package data

import (
	"context"
	"starlight/balancer/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type metricSource struct {
	metricData *MetricData

	instanceInfo map[biz.Instance]info
	log          *log.Helper
}

type info struct {
	Service string
	PodID   string
	NodeID  string
	Zone    string
}

func NewMetricSource(metricData *MetricData, logger log.Logger) biz.MetricSource {
	return &metricSource{
		metricData:   metricData,
		instanceInfo: make(map[biz.Instance]info),
		log:          log.NewHelper(logger),
	}
}

func (m *metricSource) GetByInstanceID(ctx context.Context, id biz.Instance) biz.Metric {
	info, ok := m.instanceInfo[id]
	if !ok {
		var err error
		info, err = m.updateInstanceInfo(ctx, id)
		if err != nil {
			// TODO
		}
		m.instanceInfo[id] = info
	}
	return biz.Metric{}
}

func (m *metricSource) updateInstanceInfo(ctx context.Context, id biz.Instance) (info, error) {
	return info{}, nil
}
