package models

import "github.com/gofrs/uuid"

type File struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Directory string    `json:"directory"`
}

type Directory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
