package repository

import (
	"context"
	"fmt"
	"strings"

	"minerva-web/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	pool *pgxpool.Pool
}

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{pool: pool}
}

func (r *BookRepository) List(ctx context.Context, filter model.BookFilter) ([]model.Book, int, error) {
	var whereClause string
	var args []interface{}

	conditions := []string{"deleted_at IS NULL"}
	if filter.Query != "" {
		conditions = append(conditions, "(title ILIKE $1 OR author ILIKE $1 OR isbn ILIKE $1)")
		args = append(args, "%"+filter.Query+"%")
	}
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM books"
	if whereClause != "" {
		countQuery += " " + whereClause
	}
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting books: %w", err)
	}

	query := "SELECT isbn, title, title_kana, author, author_kana, publisher, pub_date, series_name, description, cover_url, item_price, source_api, location, fetched_at, created_at, updated_at, deleted_at FROM books"
	if whereClause != "" {
		query += " " + whereClause
	}

	sortField := "isbn"
	sortOrder := "ASC"
	allowedFields := map[string]bool{
		"isbn":       true,
		"title":      true,
		"author":     true,
		"publisher":  true,
		"item_price": true,
		"created_at": true,
		"updated_at": true,
	}
	if allowedFields[filter.SortField] {
		sortField = filter.SortField
	}
	if filter.SortOrder == "DESC" {
		sortOrder = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder)

	argIdx := len(args) + 1
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("querying books: %w", err)
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var b model.Book
		if err := rows.Scan(
			&b.ISBN, &b.Title, &b.TitleKana, &b.Author, &b.AuthorKana,
			&b.Publisher, &b.PubDate, &b.SeriesName, &b.Description,
			&b.CoverURL, &b.ItemPrice, &b.SourceAPI, &b.Location, &b.FetchedAt,
			&b.CreatedAt, &b.UpdatedAt, &b.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scanning book: %w", err)
		}
		books = append(books, b)
	}
	return books, total, nil
}

func (r *BookRepository) GetByISBN(ctx context.Context, isbn string) (*model.Book, error) {
	var b model.Book
	err := r.pool.QueryRow(ctx,
		"SELECT isbn, title, title_kana, author, author_kana, publisher, pub_date, series_name, description, cover_url, item_price, source_api, location, fetched_at, created_at, updated_at, deleted_at FROM books WHERE isbn = $1 AND deleted_at IS NULL",
		isbn,
	).Scan(
		&b.ISBN, &b.Title, &b.TitleKana, &b.Author, &b.AuthorKana,
		&b.Publisher, &b.PubDate, &b.SeriesName, &b.Description,
		&b.CoverURL, &b.ItemPrice, &b.SourceAPI, &b.Location, &b.FetchedAt,
		&b.CreatedAt, &b.UpdatedAt, &b.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("getting book by isbn: %w", err)
	}
	return &b, nil
}

func (r *BookRepository) Create(ctx context.Context, b *model.Book) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO books (isbn, title, title_kana, author, author_kana, publisher, pub_date, series_name, description, cover_url, item_price, source_api, location)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		b.ISBN, b.Title, b.TitleKana, b.Author, b.AuthorKana,
		b.Publisher, b.PubDate, b.SeriesName, b.Description,
		b.CoverURL, b.ItemPrice, b.SourceAPI, b.Location,
	)
	if err != nil {
		return fmt.Errorf("creating book: %w", err)
	}
	return nil
}

func (r *BookRepository) Update(ctx context.Context, b *model.Book) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE books SET title=$2, title_kana=$3, author=$4, author_kana=$5, publisher=$6, pub_date=$7, series_name=$8, description=$9, cover_url=$10, item_price=$11, location=$12, updated_at=CURRENT_TIMESTAMP
		 WHERE isbn=$1 AND deleted_at IS NULL`,
		b.ISBN, b.Title, b.TitleKana, b.Author, b.AuthorKana,
		b.Publisher, b.PubDate, b.SeriesName, b.Description,
		b.CoverURL, b.ItemPrice, b.Location,
	)
	if err != nil {
		return fmt.Errorf("updating book: %w", err)
	}
	return nil
}

func (r *BookRepository) Delete(ctx context.Context, isbn string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE books SET deleted_at=CURRENT_TIMESTAMP WHERE isbn=$1 AND deleted_at IS NULL`,
		isbn,
	)
	if err != nil {
		return fmt.Errorf("soft deleting book: %w", err)
	}
	return nil
}
