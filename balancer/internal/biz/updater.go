package biz

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

const svcOpSep = ":" // service:operation

type (
	Instance  string
	Endpoint  string
	Operation string
)

type TraceSource interface {
	ListSpan(context.Context, time.Duration) []Span
}

type MetricSource interface {
	GetByInstanceID(context.Context, Instance) Metric
}

type weight struct {
	instance Instance
	value    int
}

type weightList struct {
	sync.RWMutex

	// instance weight list for each upstream operation
	weights map[Operation][]weight

	total int
}

func (l *weightList) TotalWeight() int {
	l.RLock()
	defer l.RUnlock()
	return l.total
}

func (l *weightList) WeightsOf(op Operation) []weight {
	l.RLock()
	defer l.RUnlock()
	return l.weights[op]
}

func (l *weightList) Update(op Operation, new []weight) {
	l.Lock()
	defer l.Unlock()
	if old, has := l.weights[op]; has {
		for _, w := range old {
			l.total -= w.value
		}
	}
	if l.total < 0 {
		l.total = 0
	}
	l.weights[op] = new
	for _, w := range new {
		l.total += w.value
	}
}

type InstanceInfo struct {
	ID      string
	Service string
	Port    string
	Pod     string
	Node    string
	Zone    string
}

type SerivceInfo struct {
	Service            string
	Instances          []Instance
	Operations         []Operation
	UpstreamOperations []Operation
}

type dependencyGraph struct {
	mu sync.Mutex

	// operation dependency graph
	graph map[Operation][]Operation
}

func (g *dependencyGraph) Update(caller, callee Operation) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if list, has := g.graph[caller]; has {
		for _, op := range list {
			if op == callee {
				return
			}
		}
	}
	g.graph[caller] = append(g.graph[caller], callee)
}

func (g *dependencyGraph) GetDependency(operation Operation) []Operation {
	if res, has := g.graph[operation]; has {
		return res
	}
	return []Operation{}
}

type metricMap struct {
	mu sync.Mutex
	m  map[Instance]Metric
}

func (m *metricMap) GetMetric(id Instance) Metric {
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
	mu         sync.Mutex
	insWeights map[Instance]*weightList // instance id -> weight list
	insInfos   map[Instance]InstanceInfo
	svcInfos   map[string]SerivceInfo // service name -> service info

	// operation dependency depGraph
	depGraph dependencyGraph

	// instance metrics
	insMetrics metricMap

	traceSource  TraceSource  // trace data source
	metricSource MetricSource // metric data source

	log log.Helper
}

func NewWeightUpdater(logger log.Logger, traceSource TraceSource, metricSource MetricSource) *WeightUpdater {
	return &WeightUpdater{
		log:          *log.NewHelper(logger),
		traceSource:  traceSource,
		metricSource: metricSource,
	}
}

func (u *WeightUpdater) UpdateWeights(ctx context.Context, id Instance) map[Operation]map[Endpoint]int {
	// get upstream service
	insInfo := u.insInfos[id]
	upops := u.svcInfos[insInfo.Service].UpstreamOperations
	weightList := u.insWeights[id]
	upstreamSvcs := u.listUpstreamServices(insInfo.Service)
	// update metric for each upstream instance
	for _, svc := range upstreamSvcs {
		for _, ins := range u.svcInfos[svc].Instances {
			observed := u.metricSource.GetByInstanceID(ctx, ins)
			u.insMetrics.Update(id, observed)
		}
	}
	// calculate link load factor for each upstream instance
	for _, upop := range upops {
		new := u.updateWeights(ctx, insInfo, upop.ServiceName(), weightList.WeightsOf(upop))
		if new == nil {
			continue
		}
		weightList.Update(upop, new)
	}

	results := make(map[Operation]map[Endpoint]int, len(upops))

	for _, upop := range upops {
		ws := weightList.WeightsOf(upop)
		opresult := make(map[Endpoint]int, len(ws))
		for _, w := range ws {
			opresult[Endpoint(string(w.instance)+":"+insInfo.Port)] = w.value
		}
		results[upop] = opresult
	}
	return results
}

func (u *WeightUpdater) updateWeights(ctx context.Context, insInfo InstanceInfo, upstreamService string, old []weight) []weight {
	svcInfo, has := u.svcInfos[upstreamService]
	if !has {
		log.Errorf("upstream service %s infomation not found", upstreamService)
		return nil
	}

	fms := make([]GAInsStatus, 0, len(svcInfo.Instances))
	for _, ins := range svcInfo.Instances {
		var weight int
		for _, w := range old {
			if w.instance == ins {
				weight = w.value
				break
			}
		}
		// prepare metric for each instance
		fms = append(fms, GAInsStatus{
			Instance:  ins,
			OldWeight: weight,
			Metric:    u.insMetrics.GetMetric(ins),
			LinkLoad:  u.calcLinkLoad(ins),

			IsSameNode: insInfo.Node == u.insInfos[ins].Node,
			IsSameZone: insInfo.Zone == u.insInfos[ins].Zone,
		})
	}

	ga := NewGARunner(GeneticConfig{
		PopulationSize: 100,
		MaxGeneration:  100,
		CrossoverRate:  0.8,
		MutationRate:   0.1,
	}, fms)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return ga.Run(ctx)
}

func (u *WeightUpdater) calcLinkLoad(id Instance) int {
	info := u.insInfos[id]
	svc := info.Service
	upops := u.svcInfos[svc].UpstreamOperations

	// check load status for the last instance
	if len(upops) == 0 {
		return u.IsOverload(id)
	}

	// check link load status for the next
	weights := u.insWeights[id]
	var load int
	for _, upop := range upops {
		for _, w := range weights.WeightsOf(upop) {
			load += (w.value << 10) / weights.TotalWeight() * u.calcLinkLoad(w.instance)
		}
	}

	return load >> 10
}

func (u *WeightUpdater) listUpstreamServices(service string) []string {
	ops := u.svcInfos[service].Operations
	svcs := map[string]struct{}{}

	for _, op := range ops {
		stack := []Operation{op}
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			upops := u.depGraph.GetDependency(cur)
			if len(upops) == 0 {
				continue
			}

			for _, upop := range upops {
				svc := upop.ServiceName()
				if _, has := svcs[svc]; has {
					continue
				}
				svcs[svc] = struct{}{}
				stack = append(stack, upop)
			}
		}
	}

	result := make([]string, 0, len(svcs))
	for k := range svcs {
		result = append(result, k)
	}
	return result
}

func (u *WeightUpdater) UpdateDependency(service string, operations []Operation, upstreams []Operation) {
	// TODO operation 分组
	for _, caller := range operations {
		for _, callee := range upstreams {
			u.depGraph.Update(caller, callee)
		}
	}
	if info, has := u.svcInfos[service]; has {
		info.Operations = operations
		info.UpstreamOperations = upstreams
		u.svcInfos[service] = info
	}
}

func (u *WeightUpdater) UpdateInstance(ins InstanceInfo) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.insInfos[Instance(ins.ID)] = ins
	if info, has := u.svcInfos[ins.Service]; has {
		info.Instances = append(info.Instances, Instance(ins.ID))
		u.svcInfos[info.Service] = info
	} else {
		u.svcInfos[info.Service] = SerivceInfo{
			Service:   ins.Service,
			Instances: []Instance{Instance(ins.ID)},
		}
	}
}

func (u *WeightUpdater) RemoveInstance(id Instance) {
	u.mu.Lock()
	defer u.mu.Unlock()
	service := u.insInfos[id].Service
	delete(u.insInfos, id)
	if info, has := u.svcInfos[service]; has {
		for i, ins := range info.Instances {
			if ins == id {
				info.Instances = append(info.Instances[:i], info.Instances[i+1:]...)
			}
		}
		if len(info.Instances) <= 0 {
			delete(u.svcInfos, service)
		}
	}
	u.insMetrics.Clear(id)
}

func NewOperation(service, operation string) Operation {
	return Operation(service + svcOpSep + operation)
}

func (o *Operation) ServiceName() string {
	return strings.SplitN(string(*o), svcOpSep, 2)[0]
}
