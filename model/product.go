package model

import "time"

type Product struct {
	ID         string
	CategoryID string
	Category   string
	Name       string
	SKU        string
	Unit       string
	PriceBuy   float64
	PriceSell  float64
	Stock      int
	CreatedAt  time.Time
}

type StockMovement struct {
	ID          string
	ProductID   string
	ProductName string
	Type        string
	Qty         int
	Note        string
	CreatedAt   time.Time
}
