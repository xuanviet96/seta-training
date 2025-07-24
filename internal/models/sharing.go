package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FolderShare struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FolderID  uuid.UUID `json:"folder_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Access    string    `json:"access" gorm:"not null;check:access IN ('read', 'write')"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Folder *Folder `json:"folder,omitempty" gorm:"foreignKey:FolderID"`
	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type NoteShare struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NoteID    uuid.UUID `json:"note_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Access    string    `json:"access" gorm:"not null;check:access IN ('read', 'write')"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Note *Note `json:"note,omitempty" gorm:"foreignKey:NoteID"`
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (fs *FolderShare) BeforeCreate(tx *gorm.DB) error {
	if fs.ID == uuid.Nil {
		fs.ID = uuid.New()
	}
	return nil
}

func (ns *NoteShare) BeforeCreate(tx *gorm.DB) error {
	if ns.ID == uuid.Nil {
		ns.ID = uuid.New()
	}
	return nil
}
