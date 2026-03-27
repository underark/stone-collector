package game

import (
	"testing"
	"time"
)

func TestTicksSince1(t *testing.T) {
	earlier, err := time.Parse(time.DateTime, "2000-01-01 00:00:00")
	if err != nil {
		t.Errorf("Wanted parsed string got %v: %v", earlier, err)
	}

	later, err := time.Parse(time.DateTime, "2000-01-01 00:05:00")
	if err != nil {
		t.Errorf("Wanted parsed string got %v: %v", earlier, err)
	}

	want := 1.0
	ticks := TicksSince(earlier, later)
	if ticks != want {
		t.Errorf("wanted %f got %f", want, ticks)
	}
}

func TestTicksSince2(t *testing.T) {
	earlier, err := time.Parse(time.DateTime, "2000-01-01 00:00:00")
	if err != nil {
		t.Errorf("Wanted parsed string got %v: %v", earlier, err)
	}

	later, err := time.Parse(time.DateTime, "2000-01-01 00:10:00")
	if err != nil {
		t.Errorf("Wanted parsed string got %v: %v", earlier, err)
	}

	want := 2.0
	ticks := TicksSince(earlier, later)
	if ticks != want {
		t.Errorf("wanted %f got %f", want, ticks)
	}
}

func TestTicksSince2Floor(t *testing.T) {
	earlier, err := time.Parse(time.DateTime, "2000-01-01 00:00:00")
	if err != nil {
		t.Errorf("Wanted parsed string got %v: %v", earlier, err)
	}

	later, err := time.Parse(time.DateTime, "2000-01-01 00:10:20")
	if err != nil {
		t.Errorf("Wanted parsed string got %v: %v", earlier, err)
	}

	want := 2.0
	ticks := TicksSince(earlier, later)
	if ticks != want {
		t.Errorf("wanted %f got %f", want, ticks)
	}
}
