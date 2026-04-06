// Package state defines game state
package state

type State struct {
	Stones int `json:"stones"`
}

type Trade struct {
	ID          int    `json:"id"`
	OwnerID     int    `json:"ownerID"`
	Material    string `json:"material"`
	Amount      int    `json:"amount"`
	MaterialReq string `json:"materialReq"`
	AmountReq   int    `json:"amountReq"`
}

type Inventory struct {
	ID       int    `json:"id"`
	OwnerID  int    `json:"ownerID"`
	Material string `json:"material"`
	Amount   int    `json:"amount"`
}
