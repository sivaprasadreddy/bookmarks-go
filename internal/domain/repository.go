package domain

import (
	"context"

	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type BookmarkRepository interface {
	GetAll(ctx context.Context) ([]Bookmark, error)
	GetById(ctx context.Context, bookmarkId int) (Bookmark, error)
	Create(ctx context.Context, bookmark Bookmark) (Bookmark, error)
	Update(ctx context.Context, bookmark Bookmark) (Bookmark, error)
	Delete(ctx context.Context, bookmarkId int) error
}

type bookmarkRepo struct {
	db *pgx.Conn
}

func NewBookmarkRepo(db *pgx.Conn) BookmarkRepository {
	return bookmarkRepo{db}
}

func (b bookmarkRepo) GetAll(ctx context.Context) ([]Bookmark, error) {
	query := "SELECT id, title, url, created_at FROM bookmarks"
	rows, err := b.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var bookmarks []Bookmark
	defer rows.Close()
	for rows.Next() {
		var bookmark = Bookmark{}
		err = rows.Scan(&bookmark.Id, &bookmark.Title, &bookmark.Url, &bookmark.CreatedDate)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}
	return bookmarks, nil
}

func (b bookmarkRepo) GetById(ctx context.Context, bookmarkId int) (Bookmark, error) {
	log.Infof("Fetching bookmark with id=%d", bookmarkId)
	var bookmark = Bookmark{}
	query := "select id, title, url, created_at, updated_at FROM bookmarks where id=$1"
	err := b.db.QueryRow(ctx, query, bookmarkId).Scan(
		&bookmark.Id, &bookmark.Title, &bookmark.Url, &bookmark.CreatedDate, &bookmark.UpdatedDate)
	if err != nil {
		return Bookmark{}, err
	}
	return bookmark, nil
}

func (b bookmarkRepo) Create(ctx context.Context, bookmark Bookmark) (Bookmark, error) {
	var lastInsertID int
	insertQuery := "insert into bookmarks(title, url, created_at) values($1, $2, $3) RETURNING id"
	err := b.db.QueryRow(ctx, insertQuery, bookmark.Title, bookmark.Url, bookmark.CreatedDate).Scan(&lastInsertID)
	if err != nil {
		log.Errorf("Error while inserting bookmark row: %v", err)
		return Bookmark{}, err
	}
	bookmark.Id = lastInsertID
	return bookmark, nil
}

func (b bookmarkRepo) Update(ctx context.Context, bookmark Bookmark) (Bookmark, error) {
	updateQuery := "update bookmarks set title = $1, url=$2, updated_at=$3 where id=$4"
	_, err := b.db.Exec(ctx, updateQuery, bookmark.Title, bookmark.Url, bookmark.UpdatedDate, bookmark.Id)
	if err != nil {
		return Bookmark{}, err
	}
	return bookmark, nil
}

func (b bookmarkRepo) Delete(ctx context.Context, bookmarkId int) error {
	deleteStmt := `delete from bookmarks where id=$1`
	_, err := b.db.Exec(ctx, deleteStmt, bookmarkId)
	return err
}
