// Package stones defines stone generation logic
package stones

import "github.com/underark/stone-collector/internal/models/locations"

type Stone struct {
	Material string
}

// New returns a new Stone
func New(location locations.Location) Stone {
	return Stone{
		location.Material(),
	}
}
