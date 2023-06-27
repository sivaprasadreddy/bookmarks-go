package api

import (
	"github.com/sivaprasadreddy/bookmarks-go/internal/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type BookmarkController struct {
	repo domain.BookmarkRepository
}

func NewBookmarkController(repository domain.BookmarkRepository) *BookmarkController {
	return &BookmarkController{repository}
}

func (b BookmarkController) FindAll(c *gin.Context) {
	log.Info("Fetching all bookmarks")
	ctx := c.Request.Context()
	bookmarks, err := b.repo.FindAll(ctx)
	if err != nil {
		log.Errorf("Error while fetching bookmarks")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch bookmarks",
		})
		return
	}
	if bookmarks == nil {
		bookmarks = []domain.Bookmark{}
	}
	c.JSON(http.StatusOK, bookmarks)
}

func (b BookmarkController) FindById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Error while parsing bookmarkId: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bookmark id",
		})
		return
	}
	log.Infof("Fetching bookmark by id %d", id)
	ctx := c.Request.Context()
	bookmark, err := b.repo.FindById(ctx, id)
	if err != nil {
		log.Errorf("Error while fetching bookmark by id: %v", err)
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
	var cb domain.CreateBookmarkModel
	if err := c.ShouldBindJSON(&cb); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Unable to parse request body. Error: " + err.Error(),
		})
		return
	}
	bookmark := domain.Bookmark{
		Title:       cb.Title,
		Url:         cb.Url,
		CreatedDate: time.Now(),
	}
	bookmark, err := b.repo.Create(ctx, bookmark)
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Error while parsing bookmarkId: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bookmark id",
		})
		return
	}
	log.Infof("update bookmark id=%d", id)
	ctx := c.Request.Context()
	var ub domain.UpdateBookmarkModel
	if err := c.ShouldBindJSON(&ub); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Unable to parse request body. Error: " + err.Error(),
		})
		return
	}
	now := time.Now()
	bookmark := domain.Bookmark{
		Id:          id,
		Title:       ub.Title,
		Url:         ub.Url,
		UpdatedDate: &now,
	}
	bookmark, err = b.repo.Update(ctx, bookmark)
	if err != nil {
		log.Errorf("Error while update bookmark: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Unable to update bookmark",
		})
		return
	}
	bookmark, _ = b.repo.FindById(c.Request.Context(), id)
	c.JSON(http.StatusOK, bookmark)
}

func (b BookmarkController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Error while parsing bookmarkId: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bookmark id",
		})
		return
	}
	log.Infof("delete bookmark with id=%d", id)
	ctx := c.Request.Context()
	err = b.repo.Delete(ctx, id)
	if err != nil {
		log.Errorf("Error while deleting bookmark: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to delete bookmark",
		})
		return
	}
	c.JSON(http.StatusOK, nil)
}
