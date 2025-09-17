package models

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Tags      pq.StringArray `json:"tags" gorm:"type:text[]"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
}

func (Post) TableName() string { return "posts" }
