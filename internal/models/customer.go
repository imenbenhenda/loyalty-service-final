package models

import "time"

// Customer représente un client avec ses points
type Customer struct {
	ID     string `json:"id"`
	Points int    `json:"points"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
}

// Transaction représente une opération sur les points
type Transaction struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Type       string    `json:"type"` // "earn" ou "redeem"
	Points     int       `json:"points"`
	Reason     string    `json:"reason"`
	Timestamp  time.Time `json:"timestamp"`
}

// Request pour ajouter des points
type AddPointsRequest struct {
	Points int    `json:"points"`
	Reason string `json:"reason"`
}

// Request pour utiliser des points
type RedeemPointsRequest struct {
	Points int    `json:"points"`
	Reward string `json:"reward"`
}
