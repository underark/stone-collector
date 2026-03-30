// Package user defines a user
package user

import (
	"testing"
	"time"
)

func TestParseTimeDummy(t *testing.T) {
	user := User{1, "alex", "2004-10-19 10:23:54"}
	time, err := user.ParseTime()
	if err != nil {
		t.Errorf("Wanted time got %v: %v", time, err)
	}
}

func TestConsumeTicks(t *testing.T) {
	user := User{1, "alex", "0000-01-01 00:00:00"}
	time, err := user.ParseTime()
	if err != nil {
		t.Errorf("Wanted time got %v: %v", time, err)
	}

	newLastTick, err := user.ConsumeTicks(1)
	if err != nil {
		t.Errorf("Wanted time string '0000-01-01 00:05:00' got: %v", err)
	}

	if newLastTick != "0000-01-01 00:05:00" {
		t.Errorf("Wanted time string '0000-01-01 00:05:00' got: %v", err)
	}
}

func TestConsumeTicksNow(t *testing.T) {
	user := User{1, "alex", time.Now().Format(time.DateTime)}
	userTime, err := user.ParseTime()
	if err != nil {
		t.Errorf("Wanted time got %v: %v", userTime, err)
	}

	newLastTick, err := user.ConsumeTicks(10)
	if err != nil {
		t.Errorf("Wanted time string got: %v", err)
	}

	ticksDuration, err := time.ParseDuration("50m")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	nowPlusTen := time.Now().Add(ticksDuration).Format(time.DateTime)
	if newLastTick != nowPlusTen {
		t.Errorf("Wanted value %s got: %s", nowPlusTen, newLastTick)
	}
}
