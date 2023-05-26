package bookmarks

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sivaprasadreddy/bookmarks-go/models"
)

type bookmarkService struct {
	repo BookmarkRepository
}

func NewBookmarkService(repo BookmarkRepository) *bookmarkService {
	return &bookmarkService{repo: repo}
}

func (b *bookmarkService) GetBookmarks() ([]models.Bookmark, error) {
	return b.repo.GetBookmarks()
}

func (b *bookmarkService) GetBookmarkById(bookmarkId int) (models.Bookmark, error) {
	return b.repo.GetBookmarkById(bookmarkId)
}

func (b *bookmarkService) CreateBookmark(createBookmark CreateBookmarkModel) (models.Bookmark, error) {
	err := createBookmark.Validate()
	if err != nil {
		log.Error(err)
		return models.Bookmark{}, err
	}

	bookmark := models.Bookmark{
		Title:       createBookmark.Title,
		Url:         createBookmark.Url,
		CreatedDate: time.Time{},
	}
	return b.repo.CreateBookmark(bookmark)
}

func (b *bookmarkService) UpdateBookmark(bookmark models.Bookmark) (models.Bookmark, error) {
	return b.repo.UpdateBookmark(bookmark)
}

func (b *bookmarkService) DeleteBookmark(bookmarkId int) error {
	return b.repo.DeleteBookmark(bookmarkId)
}
