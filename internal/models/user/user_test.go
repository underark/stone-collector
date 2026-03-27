// Package user defines a user
package user

import (
	"fmt"
	"testing"
)

func TestParseTimeDummy(t *testing.T) {
	user := User{1, "alex", "2004-10-19 10:23:54"}
	time, err := user.parseTime()
	if err != nil {
		t.Errorf("Wanted time got %v: %v", time, err)
	}
	fmt.Println(time)
}
