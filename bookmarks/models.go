package bookmarks

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type BookmarkModel struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

type CreateBookmarkModel struct {
	Title string `json:"title" validate:"required"`
	Url   string `json:"url" validate:"required,url"`
}

func (l CreateBookmarkModel) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Title, validation.Required),
		validation.Field(&l.Url, validation.Required, is.URL),
	)
}

type UpdateBookmarkModel struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}
