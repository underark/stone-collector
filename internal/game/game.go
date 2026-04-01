// Package game defines logic for the game's execution
package game

import (
	"fmt"
	"time"

	"github.com/underark/stone-collector/internal/models/locations"
	"github.com/underark/stone-collector/internal/models/user"
)

func TicksSince(user user.User) (int, error) {
	lastTick, err := user.ParseTime()
	if err != nil {
		return 0, err
	}
	return calculateTickDiff(lastTick, time.Now().UTC()), nil
}

func calculateTickDiff(lastTick time.Time, now time.Time) int {
	duration := now.Sub(lastTick)
	tick := 5
	fmt.Printf("lastTick in calculateTickDiff is %s\n", lastTick.String())
	fmt.Printf("Minutes in calculateTickDiff are %f\n", duration.Minutes())
	return int(duration.Minutes() / float64(tick))
}

func DropsFromLocation(locationID int, ticks int) (map[string]int, error) {
	location, err := locations.IDToLocation(locationID)
	if err != nil {
		return nil, err
	}

	return location.Drops(ticks), nil
}
