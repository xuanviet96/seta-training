package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Members  []TeamMember  `json:"members" gorm:"foreignKey:TeamID"`
	Managers []TeamManager `json:"managers" gorm:"foreignKey:TeamID"`
}

type TeamMember struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TeamID uuid.UUID `json:"team_id" gorm:"type:uuid;not null"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`

	// Relationships
	Team *Team `json:"team,omitempty" gorm:"foreignKey:TeamID"`
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
}

type TeamManager struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TeamID uuid.UUID `json:"team_id" gorm:"type:uuid;not null"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`

	// Relationships
	Team *Team `json:"team,omitempty" gorm:"foreignKey:TeamID"`
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (tm *TeamMember) BeforeCreate(tx *gorm.DB) error {
	if tm.ID == uuid.Nil {
		tm.ID = uuid.New()
	}
	return nil
}

func (tm *TeamManager) BeforeCreate(tx *gorm.DB) error {
	if tm.ID == uuid.Nil {
		tm.ID = uuid.New()
	}
	return nil
}
