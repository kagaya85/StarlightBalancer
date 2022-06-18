package biz

import (
	"context"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type (
	Instance  string
	Operation string
	Service   string
)
type TraceSource interface {
	ListSpan(context.Context, time.Duration) []Span
}

type MetricSource interface {
	GetByInstanceID(context.Context, Instance) Metric
}

type weightList struct {
	sync.Mutex

	// operation name -> upstream instance id -> weight
	weights map[string]map[string]int
}

type dependencyGraph struct {
	sync.Mutex

	// operation dependency graph
	graph map[Operation][]Operation

	serviceMap map[Operation]Service
}

func (g *dependencyGraph) Update(callerOp Operation, calleeOp Operation, callerSvc Service, calleeSvc Service) {
	g.graph[callerOp] = append(g.graph[callerOp], calleeOp)
	if _, has := g.serviceMap[callerOp]; !has {
		g.serviceMap[callerOp] = callerSvc
	}
	if _, has := g.serviceMap[calleeOp]; !has {
		g.serviceMap[calleeOp] = calleeSvc
	}
}

type Span struct {
	TraceID  string
	Start    time.Time
	Duration time.Duration

	CallerIns string
	CalleeIns string

	CallerSvc string
	CalleeSvc string

	CallerOp string
	CalleeOp string
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
	instances map[Instance]weightList

	// operation dependency depGraph
	depGraph dependencyGraph

	// instance metrics
	histroyMetrics map[Instance]Metric

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

func (u *WeightUpdater) UpdateInstance(ctx context.Context, id string) map[string]map[Instance]int {
	// try to update dependency from skywalking
	if ok := u.depGraph.TryLock(); ok {
		log.Info("start update dependency graph")
		spans := u.traceSource.ListSpan(ctx, 1*time.Minute)
		for _, span := range spans {
			callerOp := span.CallerOp
			calleeOp := span.CalleeOp
			if callerOp == "" || calleeOp == "" || callerOp == calleeOp {
				continue
			}
			u.depGraph.Update(Operation(callerOp), Operation(calleeOp), Service(span.CallerSvc), Service(span.CalleeSvc))
		}
		u.depGraph.Unlock()
	}

	// get upstream service
	upstreamSvcs := u.listUpstreamSvc(id)
	// get upstream service instance
	for _, svc := range upstreamSvcs {

	}
	// get current metric of upstream

	// calculate link load factor

	// genetci algorithm

	// update instance weight list
	return nil
}

func (u *WeightUpdater) listUpstreamSvc(id string) []Service {
	// TODO
	return []Service{}
}
