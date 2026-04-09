// Package drops defines locations
package drops

import (
	"math"
	"math/rand/v2"

	"github.com/underark/stone-collector/internal/models/types"
)

type Drop struct {
	Material string
	Amount   int
}

type drop struct {
	material string
	rate     float64
}

var table = []drop{
	{types.Limestone, 0.6},
	{types.Granite, 0.3},
	{types.Basalt, 0.1},
}

// Drops estimates the binomial distribution of the drop table using the normal distribution.
// Estimating the binomial distribution is reserved for large tick numbers as an optimization
// Normal distributions work best with large sample sizes and average drop rates.
func Drops(ticks int) (drops []Drop) {
	drops = make([]Drop, 0)
	remaining := ticks
	if ticks > 100 {
		for i, d := range table {
			if i == len(table)-1 {
				drops = append(drops, Drop{d.material, remaining})
			} else {
				mean := int(d.rate * float64(ticks))
				variance := int(float64(mean) * (1 - d.rate))
				stddev := int(math.Sqrt(float64(variance)))
				amount := rand.IntN((mean+stddev)+1-(mean-stddev)) + (mean - stddev)
				drops = append(drops, Drop{d.material, amount})
				remaining = remaining - amount
			}
		}
	} else {
		for range ticks {
			material := table[rand.IntN(len(table))].material
			drops = append(drops, Drop{material, 1})
		}
	}
	return
}

func Droppable() []string {
	s := make([]string, 0)
	for _, d := range table {
		s = append(s, d.material)
	}
	return s
}
