package bookmarks

import (
	"context"
	"database/sql"

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
	db *sql.DB
}

func NewBookmarkRepo(db *sql.DB) BookmarkRepository {
	var repo BookmarkRepository = bookmarkRepo{db}
	return repo
}

func (b bookmarkRepo) GetBookmarks() ([]models.Bookmark, error) {
	ctx := context.Background()
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	query := "SELECT id, title, url, created_at FROM bookmarks"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var bookmarks []models.Bookmark

	defer rows.Close()
	for rows.Next() {
		var bookmark = models.Bookmark{}
		err = rows.Scan(&bookmark.Id, &bookmark.Title, &bookmark.Url, &bookmark.CreatedDate)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (b bookmarkRepo) GetBookmarkById(bookmarkId int) (models.Bookmark, error) {
	log.Infof("Fetching bookmark with id=%d", bookmarkId)
	ctx := context.Background()
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Bookmark{}, err
	}
	var bookmark = models.Bookmark{}
	query := "select id, title, url, created_at, updated_at FROM bookmarks where id=$1"
	err = tx.QueryRowContext(ctx, query, bookmarkId).Scan(
		&bookmark.Id, &bookmark.Title, &bookmark.Url, &bookmark.CreatedDate, &bookmark.UpdatedDate)
	if err != nil {
		tx.Rollback()
		return models.Bookmark{}, err
	}
	err = tx.Commit()
	if err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil
}

func (b bookmarkRepo) CreateBookmark(bookmark models.Bookmark) (models.Bookmark, error) {
	ctx := context.Background()
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Bookmark{}, err
	}
	var lastInsertID int
	insertQuery := "insert into bookmarks(title, url, created_at) values($1, $2, $3) RETURNING id"
	err = tx.QueryRowContext(ctx, insertQuery, bookmark.Title, bookmark.Url, bookmark.CreatedDate).Scan(&lastInsertID)
	if err != nil {
		log.Errorf("Error while inserting bookmark row: %v", err)
		tx.Rollback()
		return models.Bookmark{}, err
	}
	bookmark.Id = lastInsertID
	err = tx.Commit()
	if err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil
}

func (b bookmarkRepo) UpdateBookmark(bookmark models.Bookmark) (models.Bookmark, error) {
	ctx := context.Background()
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Bookmark{}, err
	}
	updateQuery := "update bookmarks set title = $1, url=$2, updated_at=$3 where id=$4"
	_, err = tx.ExecContext(ctx, updateQuery, bookmark.Title, bookmark.Url, bookmark.UpdatedDate, bookmark.Id)
	if err != nil {
		tx.Rollback()
		return models.Bookmark{}, err
	}
	err = tx.Commit()
	if err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil
}

func (b bookmarkRepo) DeleteBookmark(bookmarkId int) error {
	ctx := context.Background()
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	deleteStmt := `delete from bookmarks where id=$1`
	_, err = tx.ExecContext(ctx, deleteStmt, bookmarkId)
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}
