//Filename: internal/data/list.go

package data

import (
	"database/sql"
	"errors"
	"time"

	"todo.joelical.net/internal/validator"
)

type List struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Task      string    `json:"task"`
	Status    string    `json:"status"`
	Version   int32     `json:"version"`
}

func ValidateList(v *validator.Validator, list *List) {
	// use the check() method to execute our validation checks
	v.Check(list.Name != "", "name", "must be provied")
	v.Check(len(list.Name) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(list.Task != "", "task", "must be provied")
	v.Check(len(list.Task) <= 800, "name", "must not be more than 200 bytes long")

	v.Check(list.Status != "", "status", "must be provied")
	v.Check(len(list.Status) <= 300, "status", "must not be more than 200 bytes long")

}

// define a ListModel which wraps a sql.db connection pool
type ListModel struct {
	DB *sql.DB
}

// Insert() allows us to creat a new list
func (m ListModel) Insert(list *List) error {
	query := `
		INSERT INTO lists (name, task, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version
	`
	//collect the data fields into a slice
	args := []interface{}{
		list.Name,
		list.Task,
		list.Status,
	}
	return m.DB.QueryRow(query, args...).Scan(&list.ID, &list.CreatedAt, &list.Version)
}

// Get() allows us to retrieve a specfic list
func (m ListModel) Get(id int64) (*List, error) {
	//ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	//create the query
	query := `
		SELECT id, created_at, name, task, version
		FROM lists
		WHERE id = $1
	`
	//declare a list variable to hold the returned data
	var list List
	//execute the query using QueryRow(
	err := m.DB.QueryRow(query, id).Scan(
		&list.ID,
		&list.CreatedAt,
		&list.Name,
		&list.Task,
		&list.Version,
	)
	//handle any errors
	if err != nil {
		//check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	//success
	return &list, nil
}

// Update() allows us edit a specific list
func (m ListModel) Update(list *List) error {
	return nil
}

// Delete() allows us to remove a specific list
func (m ListModel) Delete(id int64) error {
	return nil
}
