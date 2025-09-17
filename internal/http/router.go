package httpserver

import (
	"github.com/xuanviet96/seta-training/internal/config"
	"github.com/xuanviet96/seta-training/internal/domain/repository"
	service "github.com/xuanviet96/seta-training/internal/domain/services"
	"github.com/xuanviet96/seta-training/internal/http/handlers"
	"github.com/xuanviet96/seta-training/internal/http/middleware"
	search "github.com/xuanviet96/seta-training/internal/search"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(cfg config.Config, log *zap.Logger, gdb *gorm.DB, rdb *redis.Client, es *search.ESClient) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery(), middleware.ErrorHandler())

	// health
	health := handlers.NewHealthHandler(gdb, rdb, es)
	r.GET("/health", health.Get)

	// posts
	repo := repository.NewPostRepository()
	svc := service.NewPostService(cfg, log, gdb, rdb, repo, es)
	ph := handlers.NewPostHandler(svc)

	v1 := r.Group("/v1")
	{
		v1.POST("/posts", ph.Create)
		v1.GET("/posts/:id", ph.GetByID)
		v1.PUT("/posts/:id", ph.Update)
		v1.GET("/posts/search-by-tag", ph.SearchByTag)
		v1.GET("/posts/search", ph.Search)
	}

	return r
}
