//udah dibikin ke email github
package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"


	"github.com/joho/godotenv"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Data structs
type Category struct {
	ID           string
	Name         string
	ProductCount int
}

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

type Supplier struct {
	ID      string
	Name    string
	Contact string
	Address string
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

type DashboardData struct {
	ProductCount    int
	SupplierCount   int
	RecentProducts  []Product
	RecentMovements []StockMovement
}

//go:embed template/*
var templateFS embed.FS
var templates *template.Template
var db *pgxpool.Pool

func timeAgo(t time.Time) string {
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "dashboard.html", nil)
}

func partialHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		http.Error(w, "Page not specified", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var data any

	switch page {
	case "kategori":
		rows, err := db.Query(ctx, `SELECT id, name, (SELECT count(*) FROM products WHERE category_id = categories.id) as product_count FROM categories`)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var categories []Category
		for rows.Next() {
			var c Category
			if err := rows.Scan(&c.ID, &c.Name, &c.ProductCount); err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			categories = append(categories, c)
		}
		data = categories
	case "produk":
		rows, err := db.Query(ctx, `SELECT p.id, p.category_id, c.name as category_name, p.name, p.sku, p.unit, p.price_buy, p.price_sell, p.stock, p.created_at FROM products p LEFT JOIN categories c ON p.category_id = c.id`)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var products []Product
		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.ID, &p.CategoryID, &p.Category, &p.Name, &p.SKU, &p.Unit, &p.PriceBuy, &p.PriceSell, &p.Stock, &p.CreatedAt); err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}
		data = products
	case "supplier":
		rows, err := db.Query(ctx, `SELECT id, name, contact, address FROM suppliers`)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var suppliers []Supplier
		for rows.Next() {
			var s Supplier
			if err := rows.Scan(&s.ID, &s.Name, &s.Contact, &s.Address); err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			suppliers = append(suppliers, s)
		}
		data = suppliers
	case "dashboard-content":
		var d DashboardData
		err := db.QueryRow(ctx, `SELECT count(*) FROM products`).Scan(&d.ProductCount)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.QueryRow(ctx, `SELECT count(*) FROM suppliers`).Scan(&d.SupplierCount)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Recent Products
		rows, err := db.Query(ctx, `SELECT p.id, p.name, c.name as category, p.stock FROM products p LEFT JOIN categories c ON p.category_id = c.id ORDER BY p.created_at DESC LIMIT 5`)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Stock); err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			d.RecentProducts = append(d.RecentProducts, p)
		}

		// Recent Movements
		rows, err = db.Query(ctx, `SELECT sm.id, p.name as product_name, sm.type, sm.qty, sm.created_at FROM stock_movements sm JOIN products p ON sm.product_id = p.id ORDER BY sm.created_at DESC LIMIT 5`)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var sm StockMovement
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

	err := templates.ExecuteTemplate(w, page+".html", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	var err error
	db, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"formatCurrency": func(n float64) string { return fmt.Sprintf("Rp %.0f", n) },
		"timeAgo":        timeAgo,
	}).ParseFS(templateFS, "template/*.html"))

	fileServer := http.FileServer(http.FS(templateFS))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/partial", partialHandler)

	fmt.Println("Server berjalan di http://localhost:3000")

	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error saat menjalankan server:", err)
	}
}
