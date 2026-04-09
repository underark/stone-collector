package models

type Drop struct {
	Material string
	Amount   int
}

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

const (
	Limestone = "Limestone"
	Basalt    = "Basalt"
	Sand      = "Sand"
	Shell     = "Shell"
	Granite   = "Granite"
)
