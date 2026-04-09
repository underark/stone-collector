package ticks

import (
	"strconv"
	"time"
)

func parseTime(t string) (time.Time, error) {
	return time.Parse(time.DateTime, t)
}

func calculateTickDiff(lastTick time.Time, now time.Time) int {
	duration := now.Sub(lastTick)
	tick := 5
	return int(duration.Minutes() / float64(tick))
}

func TicksSince(lastTick string) (int, error) {
	t, err := parseTime(lastTick)
	if err != nil {
		return 0, err
	}
	return calculateTickDiff(t, time.Now().UTC()), nil
}

func ConsumeTicks(lastTick string, ticks int) (string, error) {
	t, err := parseTime(lastTick)
	if err != nil {
		return "", err
	}

	tickString := strconv.Itoa(ticks*5) + "m"
	tickDuration, err := time.ParseDuration(tickString)
	if err != nil {
		return "", err
	}

	return t.Add(tickDuration).Format(time.DateTime), nil
}
