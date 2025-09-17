package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/xuanviet96/seta-training/internal/domain/models"
	service "github.com/xuanviet96/seta-training/internal/domain/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PostHandler struct {
	svc *service.PostService
	val *validator.Validate
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{svc: svc, val: validator.New()}
}

type createPostReq struct {
	Title   string   `json:"title" validate:"required,min=1"`
	Content string   `json:"content" validate:"required,min=1"`
	Tags    []string `json:"tags"`
}

func (h *PostHandler) Create(c *gin.Context) {
	var req createPostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}
	if err := h.val.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "UNPROCESSABLE_ENTITY", "message": err.Error()}})
		return
	}

	p := &models.Post{
		Title:   strings.TrimSpace(req.Title),
		Content: strings.TrimSpace(req.Content),
		Tags:    pq.StringArray(req.Tags), // <-- cast
	}
	out, err := h.svc.Create(c, p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, out)
}

type updatePostReq struct {
	Title   *string   `json:"title,omitempty"`
	Content *string   `json:"content,omitempty"`
	Tags    *[]string `json:"tags,omitempty"`
}

func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid id"}})
		return
	}
	var req updatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Load current
	p, err := h.svc.GetByID(c, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "post not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": err.Error()}})
		return
	}

	if req.Title != nil {
		p.Title = strings.TrimSpace(*req.Title)
	}
	if req.Content != nil {
		p.Content = strings.TrimSpace(*req.Content)
	}
	if req.Tags != nil {
		p.Tags = pq.StringArray(*req.Tags) // <-- cast
	}

	out, err := h.svc.Update(c, p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid id"}})
		return
	}
	p, err := h.svc.GetByID(c, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "post not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *PostHandler) SearchByTag(c *gin.Context) {
	tag := strings.TrimSpace(c.Query("tag"))
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "tag required"}})
		return
	}
	items, err := h.svc.SearchByTag(c, tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})
}

func (h *PostHandler) Search(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "q required"}})
		return
	}
	items, total, err := h.svc.SearchES(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
}

func BuildPostHandler(cfg any, svc *service.PostService) *PostHandler { return NewPostHandler(svc) }
