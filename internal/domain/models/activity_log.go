package models

import "time"

type ActivityLog struct {
	ID       int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Action   string    `json:"action"`
	PostID   int       `json:"post_id" gorm:"index"`
	LoggedAt time.Time `json:"logged_at" gorm:"autoCreateTime"`
}

func (ActivityLog) TableName() string { return "activity_logs" }
