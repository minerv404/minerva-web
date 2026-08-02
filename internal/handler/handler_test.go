package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"minerva-web/internal/config"
	"minerva-web/internal/db"
	"minerva-web/internal/handler"
	"minerva-web/internal/repository"
)

func setupTestHandler(t *testing.T) *handler.Handler {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DSN())
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	t.Cleanup(func() { pool.Close() })

	repo := repository.NewBookRepository(pool)
	h, err := handler.New(repo, "../../templates")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}
	return h
}

func TestIndexPage(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("Response should be a full HTML page")
	}
	if !strings.Contains(body, "water.css") {
		t.Error("Response should contain water.css link")
	}
	if !strings.Contains(body, "htmx.min.js") {
		t.Error("Response should contain htmx.min.js script")
	}
	if !strings.Contains(body, "table-container") {
		t.Error("Response should contain table-container class")
	}
}

func TestListBooks(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<table>") {
		t.Error("Response should contain table element")
	}
	if !strings.Contains(body, "title-link") {
		t.Error("Response should contain title-link class for clickable titles")
	}
	if !strings.Contains(body, "water.css") {
		t.Error("Direct browser access to /books should return full page with water.css")
	}
}

func TestListBooksHTMX(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/books", nil)
	req.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<table>") {
		t.Error("Response should contain table element")
	}
	if strings.Contains(body, "water.css") {
		t.Error("HTMX request to /books should return fragment without water.css")
	}
}

func TestListBooksWithSearch(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/books?q=講談社", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<table>") {
		t.Error("Response should contain table element")
	}
	if !strings.Contains(body, "water.css") {
		t.Error("Direct browser access to /books?q= should return full page with water.css")
	}
}

func TestShowBookDetail(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/books/9784062575690", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "water.css") {
		t.Error("Detail page should contain water.css link")
	}
	if !strings.Contains(body, "htmx.min.js") {
		t.Error("Detail page should contain htmx.min.js script")
	}
	if !strings.Contains(body, "9784062575690") {
		t.Error("Response should contain ISBN")
	}
	if !strings.Contains(body, "電磁気学") {
		t.Error("Response should contain book title")
	}
	if !strings.Contains(body, "一覧に戻る") {
		t.Error("Response should contain back link")
	}
}

func TestShowBookNotFound(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/books/0000000000000", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestNewBookForm(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/books/new", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "新規書籍登録") {
		t.Error("Response should contain form title")
	}
	if !strings.Contains(body, `name="isbn"`) {
		t.Error("Response should contain ISBN input")
	}
}

func TestEditBookForm(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/books/9784062575690/edit", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "書籍編集") {
		t.Error("Response should contain form title")
	}
	if !strings.Contains(body, "電磁気学") {
		t.Error("Response should contain book title")
	}
}

func TestCreateBook(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	testISBN := fmt.Sprintf("9999999%06d", time.Now().UnixNano()%1000000)

	formData := "isbn=" + testISBN + "&title=テスト書籍&author=テスト著者&publisher=テスト出版社&item_price=1000&source_api=test"
	req := httptest.NewRequest("POST", "/books", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status 303, got %d. Body: %s", w.Code, w.Body.String())
	}

	if w.Header().Get("HX-Redirect") != "/" {
		t.Error("Response should contain HX-Redirect header")
	}

	defer func() {
		req := httptest.NewRequest("DELETE", "/books/"+testISBN, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}()
}

func TestUpdateBook(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	formData := "title=更新テスト書籍&author=更新著者&publisher=更新出版社&item_price=2000"
	req := httptest.NewRequest("PUT", "/books/9784062575690", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status 303, got %d", w.Code)
	}

	formData = "title=新装版 電磁気学のABC&author=福島 肇&publisher=講談社&item_price=1100"
	req = httptest.NewRequest("PUT", "/books/9784062575690", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
}

func TestDeleteBook(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	testISBN := fmt.Sprintf("8888888%06d", time.Now().UnixNano()%1000000)
	createForm := "isbn=" + testISBN + "&title=削除テスト書籍&source_api=test"
	createReq := httptest.NewRequest("POST", "/books", strings.NewReader(createForm))
	createReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createW := httptest.NewRecorder()
	mux.ServeHTTP(createW, createReq)

	req := httptest.NewRequest("DELETE", "/books/"+testISBN, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	listReq := httptest.NewRequest("GET", "/books?q=削除テスト書籍", nil)
	listW := httptest.NewRecorder()
	mux.ServeHTTP(listW, listReq)

	if strings.Contains(listW.Body.String(), testISBN) {
		t.Error("Deleted book should not appear in list (soft delete)")
	}
}

func TestSoftDelete(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	testISBN := fmt.Sprintf("7777777%06d", time.Now().UnixNano()%1000000)
	createForm := "isbn=" + testISBN + "&title=ソフトデリートテスト&source_api=test"
	createReq := httptest.NewRequest("POST", "/books", strings.NewReader(createForm))
	createReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createW := httptest.NewRecorder()
	mux.ServeHTTP(createW, createReq)

	detailReq := httptest.NewRequest("GET", "/books/"+testISBN, nil)
	detailW := httptest.NewRecorder()
	mux.ServeHTTP(detailW, detailReq)
	if detailW.Code != http.StatusOK {
		t.Errorf("Book should be accessible before delete, got status %d", detailW.Code)
	}

	deleteReq := httptest.NewRequest("DELETE", "/books/"+testISBN, nil)
	deleteW := httptest.NewRecorder()
	mux.ServeHTTP(deleteW, deleteReq)
	if deleteW.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", deleteW.Code)
	}

	detailReq2 := httptest.NewRequest("GET", "/books/"+testISBN, nil)
	detailW2 := httptest.NewRecorder()
	mux.ServeHTTP(detailW2, detailReq2)
	if detailW2.Code != http.StatusNotFound {
		t.Error("Soft deleted book should return 404 on detail page")
	}
}

func TestPriceFormat(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "¥1,100") && !strings.Contains(body, "¥") {
		t.Error("Response should contain formatted price with yen symbol")
	}
}

func TestSortBooks(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/books?sort=title&order=DESC", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "sort=title") {
		t.Error("Response should contain sort parameter in URLs")
	}
	if !strings.Contains(body, "▼") {
		t.Error("Response should contain sort indicator for descending order")
	}
}

func TestPerPage(t *testing.T) {
	h := setupTestHandler(t)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/books?per_page=50", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "per_page=50") {
		t.Error("Response should contain per_page parameter in URLs")
	}
}
