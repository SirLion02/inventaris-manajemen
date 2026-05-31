package model

type DashboardData struct {
	ProductCount    int
	SupplierCount   int
	RecentProducts  []Product
	RecentMovements []StockMovement
}
