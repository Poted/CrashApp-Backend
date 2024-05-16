package models

import "github.com/gofrs/uuid"

type File struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique" json:"name"`
	Size      int64     `json:"size"`
	Directory string    `json:"directory"`
}

type Directory struct {
	ID   uuid.UUID `gorm:"primaryKey" json:"id"`
	Name string    `gorm:"unique" json:"name"`
}
