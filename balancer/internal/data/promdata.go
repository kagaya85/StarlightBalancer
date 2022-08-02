package data

import (
	"context"
	"fmt"
	"math"
	"starlight/balancer/internal/conf"
	"strconv"
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

func (m *MetricData) QueryPodCPUUsage(ctx context.Context, pod string) float64 {
	query := fmt.Sprintf("rate(pod_cpu_usage_seconds_total{pod=\"%s\"})", pod)
	value := m.Query(ctx, query)
	if res, err := strconv.ParseFloat(value.String(), 64); err == nil {
		return res
	}
	return 0.0
}

func (m *MetricData) QueryPodMemUsage(ctx context.Context, pod string) float64 {
	query := fmt.Sprintf("pod_memory_working_set_bytes{pod=\"%s\"}", pod)
	value := m.Query(ctx, query)
	if used, err := strconv.Atoi(value.String()); err == nil {
		return float64(used) / 128 * 1024 * 1024 // # default mem limit is 128Mi
	}
	return 0.0
}

func (m *MetricData) QueryNodeLoad(ctx context.Context, node string) int {
	query := fmt.Sprintf("node_load1{node=\"%s\"}", node)
	value := m.Query(ctx, query)
	if res, err := strconv.ParseFloat(value.String(), 64); err == nil {
		if res < 1 {
			return 0
		} else {
			return int(math.Floor(res + 0.5))
		}
	}
	return 0
}

func (m *MetricData) QueryNodeConnectionCount(ctx context.Context, node string) int {
	query := fmt.Sprintf("node_sockstat_TCP_alloc{node=\"%s\"}", node)
	value := m.Query(ctx, query)
	if count, err := strconv.Atoi(value.String()); err == nil {
		return count
	}
	return 0
}

func (m *MetricData) QueryAppResponseTime(ctx context.Context, pod string) int {
	query := fmt.Sprintf("rate(server_requests_duration_sec_sum{kubernetes_pod_name=\"%s\"}[1m])/rate(server_requests_duration_sec_count{kubernetes_pod_name=\"%s\"}[1m])", pod, pod)
	value := m.Query(ctx, query)
	if avg, err := strconv.ParseFloat(value.String(), 64); err == nil {
		return int(avg * 1000)
	}
	return 0
}

func (m *MetricData) QueryAppSuccessRate(ctx context.Context, pod string) float64 {
	query := fmt.Sprintf("sum(increase(client_requests_code_total{code=\"0\",kubernetes_pod_name=\"%s\"}[1m]))/sum(increase(client_requests_code_total{kubernetes_pod_name=\"%s\"}[1m]))", pod, pod)
	value := m.Query(ctx, query)
	if rate, err := strconv.ParseFloat(value.String(), 64); err == nil {
		return rate
	}
	return 0
}
