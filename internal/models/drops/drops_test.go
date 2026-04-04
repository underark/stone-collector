package drops

import (
	"testing"
)

func TestDrops(t *testing.T) {
	for i := range 1000 {
		drops := Drops(i)
		var total int
		for _, d := range drops {
			total = total + d.Amount
		}

		if total != i {
			t.Errorf("Wanted %d stones got: %d", i, total)
		}
	}
}
