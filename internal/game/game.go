// Package game defines logic for the game's execution
package game

import (
	"math"
	"time"
)

func TicksSince(lastTick time.Time, now time.Time) float64 {
	duration := now.Sub(lastTick)
	tick := 5.0
	return math.Floor(duration.Minutes() / tick)
}
