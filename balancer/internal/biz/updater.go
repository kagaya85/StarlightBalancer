package biz

import (
	"context"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type InstanceID string

type TraceSource interface {
	ListSpan(context.Context, time.Duration) []Span
}

type MetricSource interface {
	GetByInstanceID(context.Context, InstanceID) Metric
}

type weightList struct {
	sync.Mutex

	// operation name -> upstream instance id -> weight
	weights map[string]map[string]int
}

type dependencyGraph struct {
	sync.Mutex

	// operation dependency graph
	graph map[string][]string
}

type Span struct {
	SpanID       string
	ParentSpanID string
	TraceID      string

	Start    time.Time
	Duration time.Duration

	Service   string
	Instance  string
	Operation string
}

type Metric struct {
	CPU float64 // pod cpu usage rate
	Mem float64 // pod memory usage rate

	Load            int // node exporter
	ConnectionCount int // node_netstat_Tcp_ActiveOpens

	ResponseTime int     // app range response time
	SuccessRate  float64 // app range success rate
}

type WeightUpdater struct {
	// instance id -> weight list
	instances map[InstanceID]weightList

	// dependency graph
	graph dependencyGraph

	// instance metrics
	histroyMetrics map[InstanceID]Metric

	// trace data source
	traceSource TraceSource

	// metric data source
	metricSource MetricSource

	log log.Helper
}

func NewWeightUpdater(logger log.Logger, traceSource TraceSource, metricSource MetricSource) *WeightUpdater {
	return &WeightUpdater{
		log:          *log.NewHelper(logger),
		traceSource:  traceSource,
		metricSource: metricSource,
	}
}

func (u *WeightUpdater) UpdateInstance(instanceID string) map[InstanceID]map[string]int {
	// update dependency from skywalking

	// get upstream dependency

	// get current metric of upstream

	// calculate link load factor

	// genetci algorithm

	// update instance weight list
	return nil
}
