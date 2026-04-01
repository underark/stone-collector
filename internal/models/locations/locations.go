// Package locations defines locations
package locations

import (
	"math"
	"math/rand/v2"

	"github.com/underark/stone-collector/internal/models/types"
)

type Location struct {
	drops []drop
}

type Drop struct {
	material string
	amount   int
}

type drop struct {
	material string
	rate     float64
}

var Park = Location{
	[]drop{
		{types.Limestone, 0.6},
		{types.Granite, 0.3},
		{types.Basalt, 0.1},
	},
}

var Beach = Location{
	[]drop{
		{types.Sand, 0.5},
		{types.Shell, 0.25},
		{types.Basalt, 0.25},
	},
}

func (l *Location) Drops(ticks int) (drops map[string]int) {
	// TODO: find a way to make a better data structure (slice?) work with the loop approach below
	drops = make(map[string]int)
	remaining := ticks
	if ticks > 100 {
		for i, d := range l.drops {
			if i == len(l.drops)-1 {
				drops[d.material] = remaining
			} else {
				mean := int(d.rate * float64(ticks))
				variance := int(float64(mean) * (1 - d.rate))
				stddev := int(math.Sqrt(float64(variance)))
				amount := rand.IntN((mean+stddev)+1-(mean-stddev)) + (mean - stddev)
				drops[d.material] = amount
				remaining = remaining - amount
			}
		}
	} else {
		for range ticks {
			material := l.drops[rand.IntN(len(l.drops))].material
			drops[material] = drops[material] + 1
		}
	}
	return
}

func IDToLocation(id int) (Location, error) {
	switch id {
	case 0:
		return Park, nil
	case 1:
		return Beach, nil
	default:
		return Location{}, nil
	}
}
