package model

import "time"

type Book struct {
	ISBN        string     `json:"isbn"`
	Title       string     `json:"title"`
	TitleKana   *string    `json:"title_kana"`
	Author      *string    `json:"author"`
	AuthorKana  *string    `json:"author_kana"`
	Publisher   *string    `json:"publisher"`
	PubDate     *string    `json:"pub_date"`
	SeriesName  *string    `json:"series_name"`
	Description *string    `json:"description"`
	CoverURL    *string    `json:"cover_url"`
	ItemPrice   *int       `json:"item_price"`
	SourceAPI   string     `json:"source_api"`
	RawResponse *string    `json:"raw_response"`
	Location    *string    `json:"location"`
	FetchedAt   *time.Time `json:"fetched_at"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type BookFilter struct {
	Query     string
	Offset    int
	Limit     int
	SortField string
	SortOrder string
}
