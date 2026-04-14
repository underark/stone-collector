package models

type Drop struct {
	Material string `json:"material"`
	Amount   int    `json:"amount"`
}

type State struct {
	Total  int         `json:"total"`
	Stones []Inventory `json:"stones"`
	Trades []Trade
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
	Material string
	Amount   int64
}

const (
	Limestone = "Limestone"
	Basalt    = "Basalt"
	Sand      = "Sand"
	Shell     = "Shell"
	Granite   = "Granite"
)
