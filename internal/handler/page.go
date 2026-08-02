package handler

import (
	"net/http"
	"strconv"

	"minerva-web/internal/model"
)

func (h *Handler) indexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	query := r.URL.Query().Get("q")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 20
	}

	sortField := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")

	offset := (page - 1) * perPage

	books, total, err := h.repo.List(r.Context(), model.BookFilter{
		Query:     query,
		Offset:    offset,
		Limit:     perPage,
		SortField: sortField,
		SortOrder: sortOrder,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (total + perPage - 1) / perPage

	data := map[string]interface{}{
		"Books":      books,
		"Query":      query,
		"Page":       page,
		"PerPage":    perPage,
		"TotalPages": totalPages,
		"Total":      total,
		"SortField":  sortField,
		"SortOrder":  sortOrder,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
