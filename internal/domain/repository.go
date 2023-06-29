package domain

import (
	"context"

	"github.com/sivaprasadreddy/bookmarks-go/internal/logging"

	"github.com/jackc/pgx/v5"
)

type BookmarkRepository interface {
	FindAll(ctx context.Context) ([]Bookmark, error)
	FindByID(ctx context.Context, bookmarkID int) (Bookmark, error)
	Create(ctx context.Context, bookmark Bookmark) (Bookmark, error)
	Update(ctx context.Context, bookmark Bookmark) (Bookmark, error)
	Delete(ctx context.Context, bookmarkID int) error
}

type bookmarkRepo struct {
	db     *pgx.Conn
	logger *logging.Logger
}

func NewBookmarkRepo(db *pgx.Conn, logger *logging.Logger) BookmarkRepository {
	return &bookmarkRepo{db: db, logger: logger}
}

func (repo *bookmarkRepo) FindAll(ctx context.Context) ([]Bookmark, error) {
	sql := "SELECT id, title, url, created_at, updated_at FROM bookmarks"
	rows, err := repo.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	var bookmarks []Bookmark
	defer rows.Close()
	for rows.Next() {
		var b = Bookmark{}
		err = rows.Scan(&b.ID, &b.Title, &b.URL, &b.CreatedDate, &b.UpdatedDate)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, nil
}

func (repo *bookmarkRepo) FindByID(ctx context.Context, id int) (Bookmark, error) {
	repo.logger.Infof("Fetching bookmark with id=%d", id)
	var b = Bookmark{}
	sql := "select id, title, url, created_at, updated_at FROM bookmarks where id=$1"
	err := repo.db.QueryRow(ctx, sql, id).Scan(
		&b.ID, &b.Title, &b.URL, &b.CreatedDate, &b.UpdatedDate)
	if err != nil {
		return Bookmark{}, err
	}
	return b, nil
}

func (repo *bookmarkRepo) Create(ctx context.Context, b Bookmark) (Bookmark, error) {
	var lastInsertID int
	sql := "insert into bookmarks(title, url, created_at, updated_at) values($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(ctx, sql, b.Title, b.URL, b.CreatedDate, b.UpdatedDate).
		Scan(&lastInsertID)
	if err != nil {
		repo.logger.Errorf("Error while inserting bookmark row: %v", err)
		return Bookmark{}, err
	}
	b.ID = lastInsertID
	return b, nil
}

func (repo *bookmarkRepo) Update(ctx context.Context, b Bookmark) (Bookmark, error) {
	sql := "update bookmarks set title = $1, url=$2, updated_at=$3 where id=$4"
	_, err := repo.db.Exec(ctx, sql, b.Title, b.URL, b.UpdatedDate, b.ID)
	if err != nil {
		return Bookmark{}, err
	}
	return b, nil
}

func (repo *bookmarkRepo) Delete(ctx context.Context, id int) error {
	sql := "delete from bookmarks where id=$1"
	_, err := repo.db.Exec(ctx, sql, id)
	return err
}
