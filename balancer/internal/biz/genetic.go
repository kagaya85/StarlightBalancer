package biz

import (
	"context"
	"math"
	_ "unsafe"
)

type GeneticConfig struct {
	PopulationSize int
	MutationRate   float32
	CrossoverRate  float32
	MaxGeneration  int
}

type GAInsStatus struct {
	Instance Instance

	OldWeight int
	LinkLoad  int
	Metric    Metric

	IsSameNode bool
	IsSameZone bool
}

type GeneticRunner interface {
	InitPopulation()

	Crossover(target1, target2 []byte) []byte

	Mutation(target []byte) []byte

	Run(context.Context) []weight
}

func NewGARunner(config GeneticConfig, status []GAInsStatus) GeneticRunner {
	initPop := make([]byte, len(status))
	for i, info := range status {
		w := byte(0)
		if info.OldWeight > 255 {
			w = 255
		} else if info.OldWeight < 0 {
			w = 0
		} else {
			w = byte(info.OldWeight)
		}
		initPop[i] = w
	}
	ga := &geneticAlgorithm{
		initPop:    initPop,
		population: make([][]byte, config.PopulationSize),
		generation: 0,
		insStatus:  status,

		mutationThreshold:  uint32(math.MaxUint32 * config.MutationRate),
		crossoverThreshold: uint32(math.MaxUint32 * config.CrossoverRate),
		config:             &config,
	}
	ga.InitPopulation()
	return ga
}

type geneticAlgorithm struct {
	initPop    []byte
	population [][]byte
	generation int
	insStatus  []GAInsStatus

	mutationThreshold  uint32
	crossoverThreshold uint32
	config             *GeneticConfig
}

func (ga *geneticAlgorithm) calcFitness(target []byte) int {
	fitness := 0
	totalWeight := 0
	for _, w := range target {
		totalWeight += int(w)
	}
	// TODO 优化计算方法
	for i, w := range target {
		status := ga.insStatus[i]
		up := (int(w) << 10 / totalWeight) * (50*btoi(status.IsSameNode) + 100*btoi(status.IsSameZone) + int(status.Metric.SuccessRate*100))
		down := status.LinkLoad + status.Metric.Load + status.Metric.ResponseTime + int(status.Metric.CPU*100+status.Metric.Mem*100+math.Sqrt(float64(status.Metric.ConnectionCount)))
		fitness += up << 10 / down
	}
	return fitness
}

func (ga *geneticAlgorithm) InitPopulation() {
	ga.population[0] = ga.initPop
	for i := 1; i < len(ga.population); i++ {
		ga.population[i] = make([]byte, len(ga.initPop))
		for j := 0; j < len(ga.population[0]); j++ {
			ga.population[i][j] = randbyte()
		}
	}
}

func (ga *geneticAlgorithm) Crossover(target1, target2 []byte) []byte {
	if len(target1) != len(target2) {
		return nil
	}
	child := make([]byte, len(target1))
	for i := 0; i < len(target1); i++ {
		if randbyte()&1 == 0 {
			child[i] = target1[i]
		} else {
			child[i] = target2[i]
		}
	}
	return child
}

func (ga *geneticAlgorithm) Mutation(target []byte) []byte {
	child := make([]byte, len(target))
	for i := 0; i < len(target); i++ {
		if randbyte()&1 == 0 {
			child[i] = target[i]
		} else {
			child[i] = randbyte()
		}
	}
	return child
}

func (ga *geneticAlgorithm) Run(ctx context.Context) []weight {
	maxFitness := ga.calcFitness(ga.initPop)
	for ga.generation < ga.config.MaxGeneration {
		if ctx.Err() != nil {
			break
		}
		ga.generation++
		for i := 0; i < len(ga.population); i++ {
			child := ga.population[i]
			if fastrandn(math.MaxUint32) < ga.crossoverThreshold {
				target1 := ga.population[i]
				target2 := ga.population[(i+1)%len(ga.population)]
				child = ga.Crossover(target1, target2)
			}
			if fastrandn(math.MaxUint32) < ga.mutationThreshold {
				child = ga.Mutation(child)
			}
			f := ga.calcFitness(child)
			if f > maxFitness {
				maxFitness = f
			}
			// TODO 优化每一代的选择方法
			ga.population[i] = child
		}
	}
	return ga.bestWeights()
}

func (ga *geneticAlgorithm) bestWeights() []weight {
	best := ga.population[0]
	bestFitness := ga.calcFitness(best)
	for i := 1; i < len(ga.population); i++ {
		fitness := ga.calcFitness(ga.population[i])
		if fitness > bestFitness {
			best = ga.population[i]
			bestFitness = fitness
		}
	}
	weights := make([]weight, len(ga.initPop))
	for i := 0; i < len(best); i++ {
		weights[i] = weight{
			value:    int(best[i]),
			instance: ga.insStatus[i].Instance,
		}
	}
	return weights
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

//go:linkname fastrandn runtime.fastrandn
func fastrandn(max uint32) uint32

func randbyte() byte {
	return byte(fastrandn(255))
}
