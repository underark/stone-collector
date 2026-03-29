// Package user defines a user
package user

import (
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
