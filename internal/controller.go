package bookmarks

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type BookmarkController struct {
	repository BookmarkRepository
}

func NewBookmarkController(repository BookmarkRepository) *BookmarkController {
	return &BookmarkController{repository}
}

func (b BookmarkController) GetAll(c *gin.Context) {
	log.Info("Fetching all bookmarks")
	ctx := c.Request.Context()
	bookmarks, err := b.repository.GetBookmarks(ctx)
	if err != nil {
		log.Errorf("Error while fetching bookmarks")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch bookmarks",
		})
		return
	}
	if bookmarks == nil {
		bookmarks = []Bookmark{}
	}
	c.JSON(http.StatusOK, bookmarks)
}

func (b BookmarkController) GetById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	log.Infof("Fetching bookmark by id %d", id)
	ctx := c.Request.Context()
	bookmark, err := b.repository.GetBookmarkById(ctx, id)
	if err != nil {
		log.Errorf("Error while fetching bookmark by id")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch bookmark by id",
		})
		return
	}
	c.JSON(http.StatusOK, bookmark)
}

func (b BookmarkController) Create(c *gin.Context) {
	log.Info("create bookmark")
	ctx := c.Request.Context()
	var createBookmark CreateBookmarkModel
	if err := c.BindJSON(&createBookmark); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Unable to parse request body. Error: " + err.Error(),
		})
		return
	}
	err := createBookmark.Validate()
	if err != nil {
		log.Errorf("Error while create bookmark %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to create bookmark",
		})
		return
	}
	bookmark := Bookmark{
		Title:       createBookmark.Title,
		Url:         createBookmark.Url,
		CreatedDate: time.Time{},
	}
	bookmark, err = b.repository.CreateBookmark(ctx, bookmark)
	if err != nil {
		log.Errorf("Error while create bookmark %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to create bookmark",
		})
		return
	}
	c.JSON(http.StatusCreated, bookmark)
}

func (b BookmarkController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	log.Infof("update bookmark id=%d", id)
	ctx := c.Request.Context()
	var bookmark Bookmark
	if err := c.BindJSON(&bookmark); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Unable to parse request body. Error: " + err.Error(),
		})
		return
	}
	bookmark.Id = id
	bookmark.UpdatedDate = time.Now()
	bookmark, err := b.repository.UpdateBookmark(ctx, bookmark)
	if err != nil {
		log.Errorf("Error while update bookmark")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Unable to update bookmark",
		})
		return
	}
	bookmark, _ = b.repository.GetBookmarkById(c.Request.Context(), id)
	c.JSON(http.StatusOK, bookmark)
}

func (b BookmarkController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	log.Infof("delete bookmark with id=%d", id)
	ctx := c.Request.Context()
	err := b.repository.DeleteBookmark(ctx, id)
	if err != nil {
		log.Errorf("Error while deleting bookmark")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to delete bookmark",
		})
		return
	}
	c.JSON(http.StatusOK, nil)
}
