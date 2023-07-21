package model

import "time"

type Withdraw struct {
	Order       string  `json:"order" binding:"required"`
	Sum         float64 `json:"sum" binding:"required"`
	TimeCreated time.Time
}

type BalanceResp struct {
	Balance float64 `json:"current"`
	Sum     float64 `json:"withdrawn"`
}
