package models

import (
	"fmt"
	"go_app/backend/errorz"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

type File struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"index:idx_name_directory,unique" json:"name"`
	Size      int64     `json:"size"`
	Directory string    `gorm:"index:idx_name_directory,unique" json:"directory_id"`
	Products  []Product `gorm:"many2many:product_files;"`
}

func (f File) NewFile() *File {
	ID, _ := uuid.NewV7()
	return &File{
		ID:        ID,
		Name:      f.Name,
		Size:      f.Size,
		Directory: f.Directory,
		Products:  f.Products,
	}
}

func (f *File) FilePath(returnWithID bool) string {

	path, err := os.Getwd()
	if err != nil {
		errorz.SendError(err)
	}

	path = fmt.Sprint(filepath.Dir(path) + "\\Storage\\")

	for _, v := range f.ID.String() {
		if v == '-' {
			break
		}
		path += fmt.Sprintf("%v\\", string(v))
	}

	if returnWithID {
		path += f.ID.String()
	}

	return path
}

type Directory struct {
	ID       uuid.UUID `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"unique" json:"name"`
	ParentID uuid.UUID `json:"parent_id"`
}

func (d Directory) NewDirectory() *Directory {
	ID, _ := uuid.NewV7()
	return &Directory{
		ID:       ID,
		Name:     d.Name,
		ParentID: d.ParentID,
	}
}
