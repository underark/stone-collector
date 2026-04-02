package locations

import (
	"testing"
)

func TestDrops(t *testing.T) {
	for i := range 1000 {
		drops := Park.Drops(i)
		var total int
		for _, d := range drops {
			total = total + d.Amount
		}

		if total != i {
			t.Errorf("Wanted %d stones got: %d", i, total)
		}
	}
}
