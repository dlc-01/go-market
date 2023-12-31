package model

import "time"

type Order struct {
	ID          string    `json:"number"`
	Status      string    `json:"status"`
	Accrual     float64   `json:"accrual,omitempty"`
	TimeCreated time.Time `json:"uploaded_at,omitempty"`
}

const (
	NEW        = "NEW"
	PROCESSING = "PROCESSING"
	INVALID    = "INVALID"
	PROCESSED  = "PROCESSED"
)
