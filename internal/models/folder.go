package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Folder struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	OwnerID   uuid.UUID `json:"owner_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Owner  *User         `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Notes  []Note        `json:"notes,omitempty" gorm:"foreignKey:FolderID"`
	Shares []FolderShare `json:"shares,omitempty" gorm:"foreignKey:FolderID"`
}

func (f *Folder) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
