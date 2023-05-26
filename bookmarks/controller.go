package bookmarks

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/sivaprasadreddy/bookmarks-go/helpers"
)

type BookmarkController struct {
	repository *BookmarkRepository
}

func NewBookmarkController(repository *BookmarkRepository) *BookmarkController {
	return &BookmarkController{repository}
}

func (b *BookmarkController) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Info("Fetching all bookmarks")
	bookmarks, err := b.repository.GetBookmarks()
	if err != nil {
		log.Errorf("Error while fetching bookmarks")
		helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to fetch bookmarks")
		return
	}
	if bookmarks == nil {
		bookmarks = []Bookmark{}
	}
	helpers.RespondWithJSON(w, http.StatusOK, bookmarks)
}

func (b *BookmarkController) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	log.Infof("Fetching bookmark by id %d", id)
	bookmark, err := b.repository.GetBookmarkById(id)
	if err != nil {
		log.Errorf("Error while fetching bookmark by id")
		helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to fetch bookmark by id")
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, bookmark)
}

func (b *BookmarkController) Create(w http.ResponseWriter, r *http.Request) {
	log.Info("create bookmark")
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && contentType != "application/json" {
		helpers.RespondWithError(w, http.StatusUnsupportedMediaType, "Content-Type header is not application/json")
		return
	}
	var createBookmark CreateBookmarkModel
	err := json.NewDecoder(r.Body).Decode(&createBookmark)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Unable to parse request body. Error: "+err.Error())
		return
	}
	err = createBookmark.Validate()
	if err != nil {
		log.Errorf("Error while create bookmark %v", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to create bookmark")
		return
	}
	bookmark := Bookmark{
		Title:       createBookmark.Title,
		Url:         createBookmark.Url,
		CreatedDate: time.Time{},
	}
	bookmark, err = b.repository.CreateBookmark(bookmark)
	if err != nil {
		log.Errorf("Error while create bookmark %v", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to create bookmark")
		return
	}
	helpers.RespondWithJSON(w, http.StatusCreated, bookmark)
}

func (b *BookmarkController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	log.Infof("update bookmark id=%d", id)
	var bookmark Bookmark
	err := json.NewDecoder(r.Body).Decode(&bookmark)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Unable to parse request body. Error: "+err.Error())
		return
	}
	bookmark.Id = id
	bookmark.UpdatedDate = time.Now()
	bookmark, err = b.repository.UpdateBookmark(bookmark)
	if err != nil {
		log.Errorf("Error while update bookmark")
		helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to update bookmark")
		return
	}
	bookmark, _ = b.repository.GetBookmarkById(id)
	helpers.RespondWithJSON(w, http.StatusOK, bookmark)
}

func (b *BookmarkController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	log.Infof("delete bookmark with id=%d", id)
	err := b.repository.DeleteBookmark(id)
	if err != nil {
		log.Errorf("Error while deleting bookmark")
		helpers.RespondWithError(w, http.StatusInternalServerError, "Unable to delete bookmark")
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, nil)
}
