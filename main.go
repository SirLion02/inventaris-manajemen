//udah dibikin ke email github
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

func partialHandler(w http.ResponseWriter, r *http.Request) {
    page := r.URL.Query().Get("page")
    if page == "" {
        http.Error(w, "Page not specified", http.StatusBadRequest)
        return
    }

    err := templates.ExecuteTemplate(w, page+".html", nil)
    if err != nil {
        http.Error(w, "Page not found", http.StatusNotFound)
    }
}

func main() {
	templates = template.Must(template.ParseFS(templateFS, "template/*.html"))

	fileServer := http.FileServer(http.FS(templateFS))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)

	http.HandleFunc("/partial", partialHandler)

	fmt.Println("Server berjalan di http://localhost:3000")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error saat menjalankan server:", err)
	}
}
