package biz

import (
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

type InstanceID string

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

	log log.Helper
}

func NewWeightUpdater(logger log.Logger) *WeightUpdater {
	return &WeightUpdater{log: *log.NewHelper(logger)}
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
