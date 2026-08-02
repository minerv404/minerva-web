package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"minerva-web/internal/repository"
)

type Handler struct {
	repo      *repository.BookRepository
	templates *template.Template
}

func New(repo *repository.BookRepository, templateDir string) (*Handler, error) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"seq": func(start, end int) []int {
			var s []int
			for i := start; i <= end; i++ {
				s = append(s, i)
			}
			return s
		},
		"empty": func(v interface{}) bool {
			if v == nil {
				return true
			}
			switch val := v.(type) {
			case string:
				return val == ""
			case *string:
				return val == nil || *val == ""
			case *int:
				return val == nil
			case *time.Time:
				return val == nil
			}
			return false
		},
		"deref": func(v interface{}) string {
			switch val := v.(type) {
			case string:
				return val
			case *string:
				if val == nil {
					return ""
				}
				return *val
			case *int:
				if val == nil {
					return ""
				}
				return fmt.Sprintf("%d", *val)
			}
			return ""
		},
		"formatPrice": func(v interface{}) string {
			var n int
			switch val := v.(type) {
			case *int:
				if val == nil {
					return ""
				}
				n = *val
			case int:
				n = val
			default:
				return ""
			}
			s := fmt.Sprintf("%d", n)
			if len(s) <= 3 {
				return "¥" + s
			}
			var result strings.Builder
			for i, c := range s {
				if i > 0 && (len(s)-i)%3 == 0 {
					result.WriteByte(',')
				}
				result.WriteRune(c)
			}
			return "¥" + result.String()
		},
		"sortURL": func(field, currentField, currentOrder, query string, perPage int) string {
			newOrder := "ASC"
			if field == currentField && currentOrder == "ASC" {
				newOrder = "DESC"
			}
			if perPage < 1 {
				perPage = 20
			}
			if query != "" {
				return fmt.Sprintf("/books?q=%s&per_page=%d&sort=%s&order=%s", query, perPage, field, newOrder)
			}
			return fmt.Sprintf("/books?per_page=%d&sort=%s&order=%s", perPage, field, newOrder)
		},
		"sortIndicator": func(field, currentField, currentOrder string) string {
			if field != currentField {
				return ""
			}
			if currentOrder == "DESC" {
				return " ▼"
			}
			return " ▲"
		},
		"formatDate": func(v interface{}) string {
			switch val := v.(type) {
			case *time.Time:
				if val == nil {
					return ""
				}
				return val.Format("2006-01-02 15:04")
			case time.Time:
				return val.Format("2006-01-02 15:04")
			}
			return ""
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, err
	}
	return &Handler{repo: repo, templates: tmpl}, nil
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.indexPage)
	mux.HandleFunc("GET /books", h.listBooks)
	mux.HandleFunc("GET /books/new", h.newBookForm)
	mux.HandleFunc("GET /books/{isbn}", h.showBook)
	mux.HandleFunc("POST /books", h.createBook)
	mux.HandleFunc("GET /books/{isbn}/edit", h.editBookForm)
	mux.HandleFunc("PUT /books/{isbn}", h.updateBook)
	mux.HandleFunc("DELETE /books/{isbn}", h.deleteBook)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}
