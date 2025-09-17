package handlers

import (
	"context"
	"net/http"
	"time"

	httpserversearch "github.com/xuanviet96/seta-training/internal/search"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	DB *gorm.DB
	RD *redis.Client
	ES *httpserversearch.ESClient
}

func NewHealthHandler(db *gorm.DB, rd *redis.Client, es *httpserversearch.ESClient) *HealthHandler {
	return &HealthHandler{DB: db, RD: rd, ES: es}
}

func (h *HealthHandler) Get(c *gin.Context) {
	dbOK := "ok"
	if sqlDB, err := h.DB.DB(); err != nil || sqlDB.Ping() != nil {
		dbOK = "down"
	}

	// Derive a request-scoped context with timeout for external checks
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	redisOK := "ok"
	if err := h.RD.Ping(ctx).Err(); err != nil {
		redisOK = "down"
	}

	esOK := "ok"
	if resp, err := h.ES.Client.Info(h.ES.Client.Info.WithContext(ctx)); err != nil {
		esOK = "down"
	} else {
		_ = resp.Body.Close()
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"db":     dbOK,
		"redis":  redisOK,
		"es":     esOK,
	})
}
