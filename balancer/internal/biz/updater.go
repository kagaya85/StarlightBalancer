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

type weight struct {
	ins   Instance
	value int
}

type weightList struct {
	sync.RWMutex

	// instance weight list for each upstream operation
	weights map[Operation][]weight

	total map[Operation]int
}

func (l *weightList) TotalWeight(op Operation) int {
	l.RLock()
	defer l.RUnlock()
	return l.total[op]
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
			l.total[op] -= w.value
		}
	}
	if l.total[op] < 0 {
		l.total[op] = 0
	}
	l.weights[op] = new
	for _, w := range new {
		l.total[op] += w.value
	}
}

func (l *weightList) RemoveInstance(id Instance, op Operation) {
	l.Lock()
	defer l.Unlock()
	if weights, has := l.weights[op]; has {
		for i, w := range weights {
			if w.ins == id {
				weights = append(weights[:i], weights[i+1:]...)
				l.weights[op] = weights
				l.total[op] -= w.value
				return
			}
		}
	}
}

type InstanceInfo struct {
	ID      string // instance id
	IP      string // ip address
	Service string // service name
	Port    string // port of service
	Pod     string // pod name
	Node    string // node name
	Zone    string // zone name
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

func NewDependencyGraph() *dependencyGraph {
	return &dependencyGraph{
		graph: map[Operation][]Operation{},
	}
}

func (g *dependencyGraph) Update(caller, callee Operation) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if list, has := g.graph[caller]; has {
		for _, op := range list {
			if op == callee {
				return false
			}
		}
	}
	g.graph[caller] = append(g.graph[caller], callee)
	return true
}

func (g *dependencyGraph) GetDependency(operation Operation) []Operation {
	if res, has := g.graph[operation]; has {
		return res
	}
	return []Operation{}
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

type WeightUpdater struct {
	mu         sync.Mutex
	insWeights map[Instance]*weightList // instance id -> weight list
	insInfos   map[Instance]InstanceInfo
	svcInfos   map[string]SerivceInfo // service name -> service info

	// operation dependency depGraph
	depGraph *dependencyGraph

	// instance metrics
	insMetrics *metricMap

	traceSource  TraceSource  // trace data source
	metricSource MetricSource // metric data source

	log log.Helper
}

func NewWeightUpdater(logger log.Logger, traceSource TraceSource, metricSource MetricSource) *WeightUpdater {
	return &WeightUpdater{
		insWeights: map[Instance]*weightList{},
		insInfos:   map[Instance]InstanceInfo{},
		svcInfos:   map[string]SerivceInfo{},

		depGraph:   NewDependencyGraph(),
		insMetrics: NewMetricMap(),

		traceSource:  traceSource,
		metricSource: metricSource,

		log: *log.NewHelper(logger),
	}
}

// UpdateWeights returns the updated Endpoint weights of the each Service.
func (u *WeightUpdater) UpdateWeights(ctx context.Context, id Instance) map[string]map[Endpoint]int {
	// get upstream service
	insInfo := u.insInfos[id]
	upops := u.svcInfos[insInfo.Service].UpstreamOperations
	weightList := u.insWeights[id]
	upstreamSvcs := u.listUpstreamServices(insInfo.Service)
	// update metric for each upstream instance
	for _, svc := range upstreamSvcs {
		for _, ins := range u.svcInfos[svc].Instances {
			observed := u.metricSource.GetByInstance(ctx, insInfo)
			u.insMetrics.Update(id, observed)
			u.log.Debugf("update upstream %s metric %v", ins, u.insMetrics.GetMetric(id))
		}
	}
	// calculate link load factor for each upstream instance
	for _, upop := range upops {
		new := u.updateWeights(ctx, insInfo, upop.ServiceName(), weightList.WeightsOf(upop))
		if new == nil {
			continue
		}
		weightList.Update(upop, new)
		u.log.Debugf("update %s upstream weights %s: %v", id, upop, new)
	}

	results := make(map[string]map[Endpoint]int, len(upops))

	for _, upop := range upops {
		ws := weightList.WeightsOf(upop)
		opresult := make(map[Endpoint]int, len(ws))
		for _, w := range ws {
			if insInfo, has := u.insInfos[w.ins]; has {
				opresult[Endpoint(insInfo.IP+":"+insInfo.Port)] = w.value
			} else {
				// 该实例信息已经被删除，从权重list中删除
				weightList.RemoveInstance(w.ins, upop)
			}
		}
		results[upop.ServiceName()] = opresult
	}
	return results
}

func (u *WeightUpdater) updateWeights(ctx context.Context, insInfo InstanceInfo, upstreamService string, old []weight) []weight {
	initWeight := 100
	svcInfo, has := u.svcInfos[upstreamService]
	if !has {
		log.Debugf("upstream service %s infomation not found", upstreamService)
		return nil
	}
	insNum := len(svcInfo.Instances)
	if insNum < 1 {
		log.Debugf("no avaliable instance for service %s", upstreamService)
		return []weight{}
	}
	if insNum == 1 {
		return []weight{
			{
				ins:   svcInfo.Instances[0],
				value: initWeight,
			},
		}
	}
	fms := make([]GAInsStatus, 0, insNum)
	for _, ins := range svcInfo.Instances {
		var weight int = initWeight
		for _, w := range old {
			if w.ins == ins {
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
		MaxGeneration:  1000,
		CrossoverRate:  0.8,
		MutationRate:   0.01,
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
		total := weights.TotalWeight(upop)
		if total == 0 {
			total = 1
		}
		for _, w := range weights.WeightsOf(upop) {
			load += (w.value << 10) / total * u.calcLinkLoad(w.ins)
		}
	}

	return load >> 10
}

func (u *WeightUpdater) listUpstreamServices(service string) []string {
	ops := u.svcInfos[service].Operations
	svcs := map[string]struct{}{}

	stack := make([]Operation, 0, len(ops))
	stack = append(stack, ops...)
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

	result := make([]string, 0, len(svcs))
	for k := range svcs {
		result = append(result, k)
	}
	return result
}

func (u *WeightUpdater) initWeightList(id Instance) *weightList {
	svc := u.insInfos[id].Service
	ops := u.svcInfos[svc].Operations
	weights := make(map[Operation][]weight, len(ops))
	for _, op := range ops {
		weights[op] = []weight{}
	}
	return &weightList{
		weights: weights,
		total:   make(map[Operation]int),
	}
}

func (u *WeightUpdater) UpdateDependency(operations []Operation, upstreams []Operation) {
	// TODO operation 分组
	for _, caller := range operations {
		for _, callee := range upstreams {
			if ok := u.depGraph.Update(caller, callee); ok {
				u.log.Debugf("add dependency %s -> %s", caller, callee)
			}
		}
	}
}

func (u *WeightUpdater) UpdateInstance(ins InstanceInfo, operations []Operation, upstreams []Operation) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if info, has := u.svcInfos[ins.Service]; has {
		info.Instances = append(info.Instances, Instance(ins.ID))
		info.Operations = operations
		info.UpstreamOperations = upstreams
		u.svcInfos[ins.Service] = info
		u.log.Infow("new instance of service", info.Service, "instances count", len(info.Instances))
	} else {
		u.svcInfos[ins.Service] = SerivceInfo{
			Service:            ins.Service,
			Instances:          []Instance{Instance(ins.ID)},
			Operations:         operations,
			UpstreamOperations: upstreams,
		}
		u.log.Infow("add new service", ins.Service)
	}
	u.insInfos[Instance(ins.ID)] = ins
	u.insWeights[Instance(ins.ID)] = u.initWeightList(Instance(ins.ID))
}

func (u *WeightUpdater) RemoveInstance(target Instance) {
	u.mu.Lock()
	service := u.insInfos[target].Service
	if info, has := u.svcInfos[service]; has {
		for i, ins := range info.Instances {
			if ins == target {
				info.Instances = append(info.Instances[:i], info.Instances[i+1:]...)
				u.log.Debugf("remove %s from instance list", target)
				break
			}
		}
		if len(info.Instances) <= 0 {
			delete(u.svcInfos, service)
			u.log.Infow("no avaliable instance, remove service", service)
		}
	}
	delete(u.insInfos, target)
	delete(u.insWeights, target)
	u.log.Infof("instance %s removed", target)
	u.mu.Unlock()
	u.insMetrics.Clear(target)
}

func NewOperation(service, operation string) Operation {
	return Operation(service + svcOpSep + operation)
}

func (o *Operation) ServiceName() string {
	return strings.SplitN(string(*o), svcOpSep, 2)[0]
}
