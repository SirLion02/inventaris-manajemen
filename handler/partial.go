package handler

import (
	"fmt"
	"net/http"
	"time"

	"inventaris/model"
)

func (h *Handler) PartialHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		http.Error(w, "Page not specified", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var data any

	switch page {
	case "kategori":
		rows, err := h.DB.Query(ctx, `SELECT id, name, (SELECT count(*) FROM products WHERE category_id = categories.id) as product_count FROM categories`)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var categories []model.Category
		for rows.Next() {
			var c model.Category
			if err := rows.Scan(&c.ID, &c.Name, &c.ProductCount); err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			categories = append(categories, c)
		}
		data = categories
	case "produk":
		rows, err := h.DB.Query(ctx, `SELECT p.id, p.category_id, c.name as category_name, p.name, p.sku, p.unit, p.price_buy, p.price_sell, p.stock, p.created_at FROM products p LEFT JOIN categories c ON p.category_id = c.id`)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var products []model.Product
		for rows.Next() {
			var p model.Product
			if err := rows.Scan(&p.ID, &p.CategoryID, &p.Category, &p.Name, &p.SKU, &p.Unit, &p.PriceBuy, &p.PriceSell, &p.Stock, &p.CreatedAt); err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}
		data = products
	case "supplier":
		rows, err := h.DB.Query(ctx, `SELECT id, name, contact, address FROM suppliers`)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var suppliers []model.Supplier
		for rows.Next() {
			var s model.Supplier
			if err := rows.Scan(&s.ID, &s.Name, &s.Contact, &s.Address); err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			suppliers = append(suppliers, s)
		}
		data = suppliers
	case "dashboard-content":
		var d model.DashboardData
		err := h.DB.QueryRow(ctx, `SELECT count(*) FROM products`).Scan(&d.ProductCount)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = h.DB.QueryRow(ctx, `SELECT count(*) FROM suppliers`).Scan(&d.SupplierCount)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Recent Products
		rows, err := h.DB.Query(ctx, `SELECT p.id, p.name, c.name as category, p.stock FROM products p LEFT JOIN categories c ON p.category_id = c.id ORDER BY p.created_at DESC LIMIT 5`)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var p model.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Stock); err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			d.RecentProducts = append(d.RecentProducts, p)
		}

		// Recent Movements
		rows, err = h.DB.Query(ctx, `SELECT sm.id, p.name as product_name, sm.type, sm.qty, sm.created_at FROM stock_movements sm JOIN products p ON sm.product_id = p.id ORDER BY sm.created_at DESC LIMIT 5`)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var sm model.StockMovement
			if err := rows.Scan(&sm.ID, &sm.ProductName, &sm.Type, &sm.Qty, &sm.CreatedAt); err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			d.RecentMovements = append(d.RecentMovements, sm)
		}
		data = d
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	err := h.Templates.ExecuteTemplate(w, page+".html", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	if diff < time.Minute {
		return fmt.Sprintf("%d detik yang lalu", int(diff.Seconds()))
	}
	if diff < time.Hour {
		return fmt.Sprintf("%d menit yang lalu", int(diff.Minutes()))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%d jam yang lalu", int(diff.Hours()))
	}
	return t.Format("02 Jan 2006")
}
