package domain

import (
	"context"

	"github.com/sivaprasadreddy/bookmarks-go/internal/logging"
	"gorm.io/gorm"
)

type BookmarkRepository interface {
	FindAll(ctx context.Context) ([]Bookmark, error)
	FindByID(ctx context.Context, bookmarkID int) (Bookmark, error)
	Create(ctx context.Context, bookmark Bookmark) (Bookmark, error)
	Update(ctx context.Context, bookmark Bookmark) (Bookmark, error)
	Delete(ctx context.Context, bookmarkID int) error
}

type bookmarkRepo struct {
	db     *gorm.DB
	logger *logging.Logger
}

func NewBookmarkRepo(db *gorm.DB, logger *logging.Logger) BookmarkRepository {
	return &bookmarkRepo{db: db, logger: logger}
}

func (repo *bookmarkRepo) FindAll(ctx context.Context) ([]Bookmark, error) {
	var bookmarks []Bookmark
	result := repo.db.Find(&bookmarks)
	return bookmarks, result.Error
}

func (repo *bookmarkRepo) FindByID(ctx context.Context, id int) (Bookmark, error) {
	var bookmark Bookmark
	result := repo.db.First(&bookmark, id)
	return bookmark, result.Error
}

func (repo *bookmarkRepo) Create(ctx context.Context, b Bookmark) (Bookmark, error) {
	result := repo.db.Create(&b)
	return b, result.Error
}

func (repo *bookmarkRepo) Update(ctx context.Context, b Bookmark) (Bookmark, error) {
	var bookmark Bookmark
	result := repo.db.First(&bookmark, b.ID)
	if result.Error != nil {
		return Bookmark{}, result.Error
	}
	bookmark.Title = b.Title
	bookmark.URL = b.URL
	bookmark.UpdatedDate = b.UpdatedDate
	result = repo.db.Save(&bookmark)
	if result.Error != nil {
		return Bookmark{}, result.Error
	}
	return bookmark, nil
}

func (repo *bookmarkRepo) Delete(ctx context.Context, id int) error {
	result := repo.db.Delete(&Bookmark{}, id)
	return result.Error
}
