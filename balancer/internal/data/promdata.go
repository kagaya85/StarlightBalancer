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
	conf    *conf.Data_Prometheus
	client  api.Client
	timeout time.Duration
	api     v1.API
	log     *log.Helper
}

func NewMetricData(c *conf.Data, logger log.Logger) (*MetricData, error) {
	l := log.NewHelper(logger)
	client, err := api.NewClient(api.Config{
		Address: c.Prometheus.GetAddress(),
	})
	if err != nil {
		l.Errorw("new prometheus client error, error", err)
		return nil, err
	}
	l.Infow("create prometheus client, address", c.Prometheus.GetAddress())
	return &MetricData{
		conf:    c.Prometheus,
		client:  client,
		timeout: c.Prometheus.GetTimeout().AsDuration(),
		api:     v1.NewAPI(client),
		log:     l,
	}, nil
}

func (m *MetricData) Query(ctx context.Context, query string) model.Value {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
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
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
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
	query := fmt.Sprintf("rate(pod_cpu_usage_seconds_total{pod=\"%s\"}[1m])", pod)
	result := m.parseValue(m.Query(ctx, query))
	if result == "" {
		return 0.0
	}
	if res, err := strconv.ParseFloat(result, 64); err == nil {
		return res
	} else {
		m.log.Errorw("query pod cpu usage error, pod", pod, "error", err)
	}
	return 0.0
}

func (m *MetricData) QueryPodMemUsage(ctx context.Context, pod string) float64 {
	query := fmt.Sprintf("pod_memory_working_set_bytes{pod=\"%s\"}", pod)
	result := m.parseValue(m.Query(ctx, query))
	if result == "" {
		return 0.0
	}
	if used, err := strconv.Atoi(result); err == nil {
		return float64(used) / 128 * 1024 * 1024 // # default mem limit is 128Mi
	} else {
		m.log.Errorw("query pod mem usage error, pod", pod, "error", err)
	}
	return 0.0
}

func (m *MetricData) QueryNodeLoad(ctx context.Context, node string) int {
	query := fmt.Sprintf("node_load1{node=\"%s\"}", node)
	result := m.parseValue(m.Query(ctx, query))
	if result == "" {
		return 0.0
	}
	if res, err := strconv.ParseFloat(result, 64); err == nil {
		if res < 1 {
			return 0
		} else {
			return int(math.Floor(res + 0.5))
		}
	} else {
		m.log.Errorw("query node load error, node", node, "error", err)
	}
	return 0
}

func (m *MetricData) QueryNodeConnectionCount(ctx context.Context, node string) int {
	query := fmt.Sprintf("node_sockstat_TCP_alloc{node=\"%s\"}", node)
	result := m.parseValue(m.Query(ctx, query))
	if result == "" {
		return 0.0
	}
	if count, err := strconv.Atoi(result); err == nil {
		return count
	} else {
		m.log.Errorw("query node connection count error, node", node, "error", err)
	}

	return 0
}

func (m *MetricData) QueryAppResponseTime(ctx context.Context, pod string) int {
	query := fmt.Sprintf("rate(server_requests_duration_sec_sum{kubernetes_pod_name=\"%s\"}[5m])/rate(server_requests_duration_sec_count{kubernetes_pod_name=\"%s\"}[5m])", pod, pod)
	result := m.parseValue(m.Query(ctx, query))
	if result == "" {
		return 0.0
	}
	if avg, err := strconv.ParseFloat(result, 64); err == nil {
		return int(avg * 1000)
	} else {
		m.log.Errorw("query app response time error, pod", pod, "error", err)
	}
	return 0
}

func (m *MetricData) QueryAppSuccessRate(ctx context.Context, pod string) float64 {
	query := fmt.Sprintf("sum(increase(client_requests_code_total{code=\"0\",kubernetes_pod_name=\"%s\"}[5m]))/sum(increase(client_requests_code_total{kubernetes_pod_name=\"%s\"}[5m]))", pod, pod)
	result := m.parseValue(m.Query(ctx, query))
	if result == "" {
		return 0.0
	}
	if rate, err := strconv.ParseFloat(result, 64); err == nil {
		return rate
	} else {
		m.log.Errorw("query app success rate error, pod", pod, "error", err)
	}
	return 0
}

func (m *MetricData) parseValue(value model.Value) string {
	if value == nil || value.String() == "" {
		return ""
	}

	switch value.Type() {
	case model.ValScalar:
		return value.(*model.Scalar).Value.String()
	case model.ValVector:
		v := value.(model.Vector)
		if v.Len() > 0 {
			return v[0].Value.String()
		}
		return ""
	default:
		m.log.Errorw("query prometheus error, value", value.String(), "type", value.Type().String())
	}

	return ""
}
