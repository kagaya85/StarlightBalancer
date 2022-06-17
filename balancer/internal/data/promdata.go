package data

import (
	"context"
	"starlight/balancer/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type MetricData struct {
	conf   *conf.Data_Prometheus
	client api.Client
	api    v1.API
	log    *log.Helper
}

func NewMetricData(c *conf.Data, logger log.Logger) (*MetricData, error) {
	client, err := api.NewClient(api.Config{
		Address: c.Prometheus.GetAddress(),
	})
	if err != nil {
		return nil, err
	}
	return &MetricData{
		conf:   c.Prometheus,
		client: client,
		api:    v1.NewAPI(client),
		log:    log.NewHelper(logger),
	}, nil
}

func (m *MetricData) Query(ctx context.Context, query string) model.Value {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, warnings, err := m.api.Query(ctx, query, time.Now())
	if err != nil {
		log.Errorf("querying Prometheus error: %v", err)
		return nil
	}
	if len(warnings) > 0 {
		log.Warn(warnings)
	}

	return result
}

func (m *MetricData) QueryRange(ctx context.Context, query string) model.Value {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	r := v1.Range{
		Start: time.Now().Add(-5 * time.Second),
		End:   time.Now(),
		Step:  time.Second,
	}

	result, warnings, err := m.api.QueryRange(ctx, query, r)
	if err != nil {
		log.Errorf("querying Prometheus error: %v", err)
		return nil
	}
	if len(warnings) > 0 {
		log.Warn(warnings)
	}

	return result
}

func (m *MetricData) QueryPodCPUUsage(ctx context.Context, target string) int {
	return 0
}

func (m *MetricData) QueryPodMemUsage(ctx context.Context, target string) int {
	return 0
}

func (m *MetricData) QueryNodeLoad(ctx context.Context, target string) int {
	return 0
}

func (m *MetricData) QueryNodeConnectionCount(ctx context.Context, target string) int {
	return 0
}

func (m *MetricData) QueryAppResponseTime(ctx context.Context, target string) int {
	return 0
}

func (m *MetricData) QueryAppCode(ctx context.Context, target string) int {
	return 0
}
