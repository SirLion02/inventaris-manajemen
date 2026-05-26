package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed template/*
var templateFS embed.FS
var templates *template.Template

func loginHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "dashboard.html", nil)
}

func main() {
	templates = template.Must(template.ParseFS(templateFS, "template/*.html"))

	fileServer := http.FileServer(http.FS(templateFS))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)

	fmt.Println("Server berjalan di http://localhost:3000")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error saat menjalankan server:", err)
	}
}
