package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title     string    `json:"title" gorm:"not null"`
	Body      string    `json:"body"`
	FolderID  uuid.UUID `json:"folder_id" gorm:"type:uuid;not null"`
	OwnerID   uuid.UUID `json:"owner_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Folder *Folder     `json:"folder,omitempty" gorm:"foreignKey:FolderID"`
	Owner  *User       `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Shares []NoteShare `json:"shares,omitempty" gorm:"foreignKey:NoteID"`
}

func (n *Note) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
