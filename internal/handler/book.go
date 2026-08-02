package handler

import (
	"net/http"
	"strconv"

	"minerva-web/internal/model"
)

func (h *Handler) listBooks(w http.ResponseWriter, r *http.Request) {
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
	if r.Header.Get("HX-Request") == "true" {
		if err := h.templates.ExecuteTemplate(w, "_book_list.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		if err := h.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *Handler) newBookForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "_book_form.html", map[string]interface{}{
		"Action": "create",
		"Book":   &model.Book{SourceAPI: "manual"},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isbn := r.FormValue("isbn")
	if isbn == "" {
		http.Error(w, "ISBN is required", http.StatusBadRequest)
		return
	}

	sourceAPI := r.FormValue("source_api")
	if sourceAPI == "" {
		sourceAPI = "manual"
	}

	book := &model.Book{
		ISBN:      isbn,
		Title:     r.FormValue("title"),
		SourceAPI: sourceAPI,
	}

	if v := r.FormValue("title_kana"); v != "" {
		book.TitleKana = &v
	}
	if v := r.FormValue("author"); v != "" {
		book.Author = &v
	}
	if v := r.FormValue("author_kana"); v != "" {
		book.AuthorKana = &v
	}
	if v := r.FormValue("publisher"); v != "" {
		book.Publisher = &v
	}
	if v := r.FormValue("pub_date"); v != "" {
		book.PubDate = &v
	}
	if v := r.FormValue("series_name"); v != "" {
		book.SeriesName = &v
	}
	if v := r.FormValue("description"); v != "" {
		book.Description = &v
	}
	if v := r.FormValue("cover_url"); v != "" {
		book.CoverURL = &v
	}
	if v := r.FormValue("location"); v != "" {
		book.Location = &v
	}
	if v := r.FormValue("item_price"); v != "" {
		if price, err := strconv.Atoi(v); err == nil {
			book.ItemPrice = &price
		}
	}

	if err := h.repo.Create(r.Context(), book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func (h *Handler) showBook(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")
	book, err := h.repo.GetByISBN(r.Context(), isbn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Book":      book,
		"IsDetail":  true,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) editBookForm(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")
	book, err := h.repo.GetByISBN(r.Context(), isbn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fromDetail := r.URL.Query().Get("from") == "detail"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "_book_form.html", map[string]interface{}{
		"Action":     "edit",
		"Book":       book,
		"FromDetail": fromDetail,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book := &model.Book{
		ISBN:  isbn,
		Title: r.FormValue("title"),
	}

	if v := r.FormValue("title_kana"); v != "" {
		book.TitleKana = &v
	}
	if v := r.FormValue("author"); v != "" {
		book.Author = &v
	}
	if v := r.FormValue("author_kana"); v != "" {
		book.AuthorKana = &v
	}
	if v := r.FormValue("publisher"); v != "" {
		book.Publisher = &v
	}
	if v := r.FormValue("pub_date"); v != "" {
		book.PubDate = &v
	}
	if v := r.FormValue("series_name"); v != "" {
		book.SeriesName = &v
	}
	if v := r.FormValue("description"); v != "" {
		book.Description = &v
	}
	if v := r.FormValue("cover_url"); v != "" {
		book.CoverURL = &v
	}
	if v := r.FormValue("location"); v != "" {
		book.Location = &v
	}
	if v := r.FormValue("item_price"); v != "" {
		if price, err := strconv.Atoi(v); err == nil {
			book.ItemPrice = &price
		}
	}

	if err := h.repo.Update(r.Context(), book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := "/"
	if r.FormValue("from_detail") == "true" {
		redirectURL = "/books/" + isbn
	}
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusSeeOther)
}

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")

	if err := h.repo.Delete(r.Context(), isbn); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
