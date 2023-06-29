package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sivaprasadreddy/bookmarks-go/internal/domain"
	"github.com/sivaprasadreddy/bookmarks-go/internal/logging"

	"github.com/gin-gonic/gin"
)

type BookmarkController struct {
	repo   domain.BookmarkRepository
	logger *logging.Logger
}

func NewBookmarkController(repository domain.BookmarkRepository, logger *logging.Logger) *BookmarkController {
	return &BookmarkController{repo: repository, logger: logger}
}

func (b BookmarkController) FindAll(c *gin.Context) {
	b.logger.Info("Fetching all bookmarks")
	ctx := c.Request.Context()
	bookmarks, err := b.repo.FindAll(ctx)
	if err != nil {
		b.logger.Errorf("Error while fetching bookmarks")
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

func (b BookmarkController) FindByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		b.logger.Errorf("Error while parsing bookmarkID: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bookmark id",
		})
		return
	}
	b.logger.Infof("Fetching bookmark by id %d", id)
	ctx := c.Request.Context()
	bookmark, err := b.repo.FindByID(ctx, id)
	if err != nil {
		b.logger.Errorf("Error while fetching bookmark by id: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch bookmark by id",
		})
		return
	}
	c.JSON(http.StatusOK, bookmark)
}

func (b BookmarkController) Create(c *gin.Context) {
	b.logger.Info("create bookmark")
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
		URL:         cb.URL,
		CreatedDate: time.Now(),
	}
	bookmark, err := b.repo.Create(ctx, bookmark)
	if err != nil {
		b.logger.Errorf("Error while create bookmark %v", err)
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
		b.logger.Errorf("Error while parsing bookmarkID: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bookmark id",
		})
		return
	}
	b.logger.Infof("update bookmark id=%d", id)
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
		ID:          id,
		Title:       ub.Title,
		URL:         ub.URL,
		UpdatedDate: &now,
	}
	_, err = b.repo.Update(ctx, bookmark)
	if err != nil {
		b.logger.Errorf("Error while update bookmark: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Unable to update bookmark",
		})
		return
	}
	bookmark, _ = b.repo.FindByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, bookmark)
}

func (b BookmarkController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		b.logger.Errorf("Error while parsing bookmarkID: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bookmark id",
		})
		return
	}
	b.logger.Infof("delete bookmark with id=%d", id)
	ctx := c.Request.Context()
	err = b.repo.Delete(ctx, id)
	if err != nil {
		b.logger.Errorf("Error while deleting bookmark: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to delete bookmark",
		})
		return
	}
	c.JSON(http.StatusOK, nil)
}
