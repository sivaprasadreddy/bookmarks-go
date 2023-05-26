package bookmarks

import (
	"context"

	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"github.com/sivaprasadreddy/bookmarks-go/models"
)

type BookmarkRepository interface {
	GetBookmarks() ([]models.Bookmark, error)
	GetBookmarkById(bookmarkId int) (models.Bookmark, error)
	CreateBookmark(bookmark models.Bookmark) (models.Bookmark, error)
	UpdateBookmark(bookmark models.Bookmark) (models.Bookmark, error)
	DeleteBookmark(bookmarkId int) error
}

type bookmarkRepo struct {
	db *pgx.Conn
}

func NewBookmarkRepo(db *pgx.Conn) BookmarkRepository {
	var repo BookmarkRepository = bookmarkRepo{db}
	return repo
}

func (b bookmarkRepo) GetBookmarks() ([]models.Bookmark, error) {
	ctx := context.Background()
	query := "SELECT id, title, url, created_at FROM bookmarks"
	rows, err := b.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var bookmarks []models.Bookmark
	defer rows.Close()
	for rows.Next() {
		var bookmark = models.Bookmark{}
		err = rows.Scan(&bookmark.Id, &bookmark.Title, &bookmark.Url, &bookmark.CreatedDate)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}
	return bookmarks, nil
}

func (b bookmarkRepo) GetBookmarkById(bookmarkId int) (models.Bookmark, error) {
	log.Infof("Fetching bookmark with id=%d", bookmarkId)
	ctx := context.Background()
	var bookmark = models.Bookmark{}
	query := "select id, title, url, created_at, updated_at FROM bookmarks where id=$1"
	err := b.db.QueryRow(ctx, query, bookmarkId).Scan(
		&bookmark.Id, &bookmark.Title, &bookmark.Url, &bookmark.CreatedDate, &bookmark.UpdatedDate)
	if err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil
}

func (b bookmarkRepo) CreateBookmark(bookmark models.Bookmark) (models.Bookmark, error) {
	ctx := context.Background()
	var lastInsertID int
	insertQuery := "insert into bookmarks(title, url, created_at) values($1, $2, $3) RETURNING id"
	err := b.db.QueryRow(ctx, insertQuery, bookmark.Title, bookmark.Url, bookmark.CreatedDate).Scan(&lastInsertID)
	if err != nil {
		log.Errorf("Error while inserting bookmark row: %v", err)
		return models.Bookmark{}, err
	}
	bookmark.Id = lastInsertID
	return bookmark, nil
}

func (b bookmarkRepo) UpdateBookmark(bookmark models.Bookmark) (models.Bookmark, error) {
	ctx := context.Background()
	updateQuery := "update bookmarks set title = $1, url=$2, updated_at=$3 where id=$4"
	_, err := b.db.Exec(ctx, updateQuery, bookmark.Title, bookmark.Url, bookmark.UpdatedDate, bookmark.Id)
	if err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil
}

func (b bookmarkRepo) DeleteBookmark(bookmarkId int) error {
	ctx := context.Background()
	deleteStmt := `delete from bookmarks where id=$1`
	_, err := b.db.Exec(ctx, deleteStmt, bookmarkId)
	return err
}
