package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/xuanviet96/seta-training/internal/cache"
	"github.com/xuanviet96/seta-training/internal/config"
	"github.com/xuanviet96/seta-training/internal/domain/models"
	"github.com/xuanviet96/seta-training/internal/domain/repository"
	search "github.com/xuanviet96/seta-training/internal/search"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostService struct {
	cfg   config.Config
	log   *zap.Logger
	db    *gorm.DB
	cache *redis.Client
	repo  repository.PostRepository
	es    *search.ESClient
}

func NewPostService(cfg config.Config, log *zap.Logger, db *gorm.DB, cache *redis.Client, repo repository.PostRepository, es *search.ESClient) *PostService {
	return &PostService{cfg: cfg, log: log, db: db, cache: cache, repo: repo, es: es}
}

func (s *PostService) Create(ctx context.Context, p *models.Post) (*models.Post, error) {
	var out *models.Post
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// create post + activity log in same transaction
		al := &models.ActivityLog{
			Action:   "new_post",
			LoggedAt: time.Now(),
		}
		if err := s.repo.CreateWithLog(ctx, tx, p, al); err != nil {
			return err
		}
		out = p
		return nil
	})
	if err != nil {
		return nil, err
	}

	// index to ES (best-effort)
	go func(p models.Post) {
		ctx2, cancel := context.WithTimeout(context.Background(), s.cfg.Timeout)
		defer cancel()
		_ = search.IndexPost(ctx2, s.es, s.cfg.ESIndex, search.PostDoc{
			ID:      p.ID,
			Title:   p.Title,
			Content: p.Content,
			Tags:    []string(p.Tags), // <-- cast
		})
	}(*out)

	return out, nil
}

func (s *PostService) GetByID(ctx context.Context, id int) (*models.Post, error) {
	key := fmt.Sprintf("post:%d", id)

	// cache read
	if v, err := s.cache.Get(ctx, key).Result(); err == nil && v != "" {
		var p models.Post
		if json.Unmarshal([]byte(v), &p) == nil {
			return &p, nil
		}
	}

	// miss → db
	p, err := s.repo.GetByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}

	// backfill cache
	if b, err := json.Marshal(p); err == nil {
		_ = s.cache.Set(ctx, key, b, cache.TTL(s.cfg)).Err()
	}

	return p, nil
}

func (s *PostService) Update(ctx context.Context, p *models.Post) (*models.Post, error) {
	if err := s.repo.Update(ctx, s.db, p); err != nil {
		return nil, err
	}
	// invalidate cache
	_ = s.cache.Del(ctx, "post:"+strconv.Itoa(p.ID)).Err()

	// re-index
	go func(p models.Post) {
		ctx2, cancel := context.WithTimeout(context.Background(), s.cfg.Timeout)
		defer cancel()
		_ = search.IndexPost(ctx2, s.es, s.cfg.ESIndex, search.PostDoc{
			ID:      p.ID,
			Title:   p.Title,
			Content: p.Content,
			Tags:    []string(p.Tags), // <-- cast
		})
	}(*p)

	return p, nil
}

func (s *PostService) SearchByTag(ctx context.Context, tag string) ([]models.Post, error) {
	return s.repo.SearchByTag(ctx, s.db, tag)
}

func (s *PostService) SearchES(ctx context.Context, q string) ([]search.PostDoc, int, error) {
	return search.SearchPosts(ctx, s.es, s.cfg.ESIndex, q)
}
