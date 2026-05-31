package handler

import (
	"html/template"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB        *pgxpool.Pool
	Templates *template.Template
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	h.Templates.ExecuteTemplate(w, "login.html", nil)
}

func (h *Handler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	h.Templates.ExecuteTemplate(w, "dashboard.html", nil)
}
