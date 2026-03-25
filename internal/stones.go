// Package stones defines stone generation logic
package stones

type Stone struct {
	Material string
}

// New returns a new Stone
func New() Stone {
	return Stone{"Limestone"}
}
