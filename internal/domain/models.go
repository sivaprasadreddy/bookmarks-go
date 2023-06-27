package domain

import (
	"time"
)

type Bookmark struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

type CreateBookmarkModel struct {
	Title string `json:"title" binding:"required"`
	Url   string `json:"url" binding:"required,url"`
}

type UpdateBookmarkModel struct {
	Id    int    `json:"id"`
	Title string `json:"title" binding:"required"`
	Url   string `json:"url" binding:"required,url"`
}
