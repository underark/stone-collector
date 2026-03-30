// Package locations defines locations
package locations

import "math/rand"

type Location struct {
	stoneTypes []string
}

var Park = Location{
	[]string{
		"Limestone",
		"Granite",
		"Basalt",
	},
}

func (l *Location) Material() string {
	return l.stoneTypes[rand.Intn(len(l.stoneTypes))]
}
