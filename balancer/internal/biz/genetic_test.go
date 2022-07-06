package biz

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkGenetic(b *testing.B) {
	fms := []GAInsStatus{
		{
			Instance:  "instance1",
			OldWeight: 10,
			LinkLoad:  2,
			Metric: Metric{
				CPU:             1.9,
				Mem:             0.6,
				Load:            72,
				ConnectionCount: 50,
				ResponseTime:    3,
				SuccessRate:     0.9342,
			},
		},
		{
			Instance:  "instance2",
			OldWeight: 10,
			LinkLoad:  1,
			Metric: Metric{
				CPU:             0.7,
				Mem:             0.3,
				Load:            23,
				ConnectionCount: 105,
				ResponseTime:    1,
				SuccessRate:     1,
			},
		},
		{
			Instance:  "instance3",
			OldWeight: 10,
			LinkLoad:  15,
			Metric: Metric{
				CPU:             0.5,
				Mem:             0.5,
				Load:            2,
				ConnectionCount: 10,
				ResponseTime:    998,
				SuccessRate:     1,
			},
		},
		{
			Instance:  "instance4",
			OldWeight: 10,
			LinkLoad:  3,
			Metric: Metric{
				CPU:             1.0,
				Mem:             0.2,
				Load:            3,
				ConnectionCount: 17,
				ResponseTime:    34,
				SuccessRate:     0.3,
			},
		},
		{
			Instance:  "instance5",
			OldWeight: 10,
			LinkLoad:  5,
			Metric: Metric{
				CPU:             0.6,
				Mem:             0.3,
				Load:            2,
				ConnectionCount: 5,
				ResponseTime:    3,
				SuccessRate:     1,
			},
		},
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	// defer cancel()

	for i := 0; i < b.N; i++ {
		ga := NewGARunner(GeneticConfig{
			PopulationSize: 100,
			MaxGeneration:  500,
			CrossoverRate:  0.8,
			MutationRate:   0.1,
		}, fms)

		ctx := context.Background()
		fmt.Println(ga.Run(ctx))
	}
}
