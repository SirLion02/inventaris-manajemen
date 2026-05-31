package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"inventaris/handler"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed template/*
var templateFS embed.FS

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	templates := template.Must(template.New("").Funcs(template.FuncMap{
		"formatCurrency": func(n float64) string { return fmt.Sprintf("Rp %.0f", n) },
		"timeAgo":        handler.TimeAgo,
	}).ParseFS(templateFS, "template/*.html"))

	h := &handler.Handler{
		DB:        db,
		Templates: templates,
	}

	fileServer := http.FileServer(http.FS(templateFS))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/login", h.LoginHandler)
	http.HandleFunc("/dashboard", h.DashboardHandler)
	http.HandleFunc("/partial", h.PartialHandler)

	fmt.Println("Server berjalan di http://localhost:3000")

	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error saat menjalankan server:", err)
	}
}
