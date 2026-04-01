// Package state defines game state
package state

type State struct {
	Stones int `json:"stones"`
}

type Worker struct {
	LocationID int `json:"locationID"`
}
