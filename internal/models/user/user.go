// Package user defines a user
package user

import (
	"strconv"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LastTick string `json:"lastTick"`
}

func (u *User) ParseTime() (time.Time, error) {
	return time.Parse(time.DateTime, u.LastTick)
}

func (u *User) ConsumeTicks(ticks int) (string, error) {
	lastTick, err := u.ParseTime()
	if err != nil {
		return "", err
	}

	tickString := strconv.Itoa(ticks*5) + "m"
	tickDuration, err := time.ParseDuration(tickString)
	if err != nil {
		return "", err
	}

	return lastTick.Add(tickDuration).Format(time.DateTime), err
}
