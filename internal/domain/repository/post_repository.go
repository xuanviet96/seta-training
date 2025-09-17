package repository

import (
	"context"

	"github.com/xuanviet96/seta-training/internal/domain/models"

	"gorm.io/gorm"
)

type PostRepository interface {
	CreateWithLog(ctx context.Context, tx *gorm.DB, p *models.Post, log *models.ActivityLog) error
	GetByID(ctx context.Context, db *gorm.DB, id int) (*models.Post, error)
	Update(ctx context.Context, db *gorm.DB, p *models.Post) error
	SearchByTag(ctx context.Context, db *gorm.DB, tag string) ([]models.Post, error)
}

type postRepo struct{}

func NewPostRepository() PostRepository { return &postRepo{} }

func (r *postRepo) CreateWithLog(ctx context.Context, tx *gorm.DB, p *models.Post, al *models.ActivityLog) error {
	if err := tx.WithContext(ctx).Create(p).Error; err != nil {
		return err
	}
	al.PostID = p.ID
	if err := tx.WithContext(ctx).Create(al).Error; err != nil {
		return err
	}
	return nil
}

func (r *postRepo) GetByID(ctx context.Context, db *gorm.DB, id int) (*models.Post, error) {
	var p models.Post
	if err := db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *postRepo) Update(ctx context.Context, db *gorm.DB, p *models.Post) error {
	return db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", p.ID).
		Updates(map[string]any{
			"title":   p.Title,
			"content": p.Content,
			"tags":    p.Tags,
		}).Error
}

func (r *postRepo) SearchByTag(ctx context.Context, db *gorm.DB, tag string) ([]models.Post, error) {
	var posts []models.Post
	// Use GIN index: WHERE tags @> ARRAY[$1]::text[]
	err := db.WithContext(ctx).
		Where("tags @> ARRAY[?]::text[]", tag).
		Order("id DESC").
		Find(&posts).Error
	return posts, err
}
