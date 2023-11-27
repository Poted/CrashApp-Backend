package models

import "github.com/google/uuid"

type File struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Directory string    `json:"directory"`
}

type Directory struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	ParentID uuid.UUID `json:"parentID"`
}
