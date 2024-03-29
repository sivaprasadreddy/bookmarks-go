package domain

import (
	"time"
)

type Bookmark struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
	CreatedDate time.Time  `json:"created_date"`
	UpdatedDate *time.Time `json:"updated_date"`
}

type CreateBookmarkModel struct {
	Title string `json:"title" binding:"required"`
	URL   string `json:"url" binding:"required,url"`
}

type UpdateBookmarkModel struct {
	Title string `json:"title" binding:"required"`
	URL   string `json:"url" binding:"required,url"`
}
