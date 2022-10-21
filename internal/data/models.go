//Filename: internal/data/models.go

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// create a wrapper for our data models
type Models struct {
	List ListModel
}

// NewModels() allows us to create a new models
func NewModels(db *sql.DB) Models {
	return Models{
		List: ListModel{DB: db},
	}
}
