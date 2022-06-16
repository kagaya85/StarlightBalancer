package biz

import (
	"context"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type InstanceID string

type TraceSource interface {
	ListSpanFrom(context.Context, time.Time) []Span
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
	spanID       string
	parentSpanID string
	traceID      string

	start    time.Time
	duration time.Duration

	service   string
	operation string
}

type Metric struct {
	cpu          float64
	mem          float64
	load         int
	connecnt     int
	responseTime int // ms
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
