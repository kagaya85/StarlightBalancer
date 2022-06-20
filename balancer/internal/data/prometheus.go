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

func (m *metricSource) GetByInstanceID(ctx context.Context, id biz.Instance) biz.Metric {
	return biz.Metric{}
}
